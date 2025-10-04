# tldr++ Makefile

.PHONY: help build-go build-python test-go test-python test clean install-go install-python dev-go dev-python

# Default target
help:
	@echo "tldr++ - Interactive cheat-sheets"
	@echo ""
	@echo "Available targets:"
	@echo "  build-go      Build Go binary"
	@echo "  build-python  Build Python package"
	@echo "  test-go       Run Go tests"
	@echo "  test-python   Run Python tests"
	@echo "  test          Run all tests"
	@echo "  clean         Clean build artifacts"
	@echo "  install-go    Install Go binary"
	@echo "  install-python Install Python package"
	@echo "  dev-go        Run Go version in dev mode"
	@echo "  dev-python    Run Python version in dev mode"

# Build targets
build-go:
	@echo "Building Go binary..."
	@cd cmd/tldrpp && go build -o ../../bin/tldrpp-go -ldflags "-X main.version=$(shell git describe --tags --always --dirty) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" .

build-python:
	@echo "Building Python package..."
	@python -m build

# Test targets
test-go:
	@echo "Running Go tests..."
	@go test -v ./...

test-python:
	@echo "Running Python tests..."
	@python -m pytest tests/ -v

test: test-go test-python

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@rm -rf build/
	@rm -rf *.egg-info/
	@go clean

# Install targets
install-go: build-go
	@echo "Installing Go binary..."
	@mkdir -p ~/.local/bin
	@cp bin/tldrpp-go ~/.local/bin/tldrpp-go
	@echo "Go binary installed to ~/.local/bin/tldrpp-go"

install-python: build-python
	@echo "Installing Python package..."
	@pip install dist/*.whl

# Development targets
dev-go:
	@echo "Running Go version in dev mode..."
	@cd cmd/tldrpp && go run . --dev

dev-python:
	@echo "Running Python version in dev mode..."
	@python -m tldrpp --dev

# Setup targets
setup-go:
	@echo "Setting up Go development environment..."
	@go mod download
	@go mod tidy

setup-python:
	@echo "Setting up Python development environment..."
	@pip install -e ".[dev]"

# Lint targets
lint-go:
	@echo "Linting Go code..."
	@golangci-lint run

lint-python:
	@echo "Linting Python code..."
	@black --check .
	@isort --check-only .
	@flake8 .
	@mypy .

lint: lint-go lint-python

# Format targets
format-go:
	@echo "Formatting Go code..."
	@go fmt ./...

format-python:
	@echo "Formatting Python code..."
	@black .
	@isort .

format: format-go format-python

# Release targets
release-go:
	@echo "Building Go release binaries..."
	@mkdir -p dist/go
	@cd cmd/tldrpp && \
		GOOS=linux GOARCH=amd64 go build -o ../../dist/go/tldrpp-linux-amd64 -ldflags "-X main.version=$(shell git describe --tags --always --dirty) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" . && \
		GOOS=darwin GOARCH=amd64 go build -o ../../dist/go/tldrpp-darwin-amd64 -ldflags "-X main.version=$(shell git describe --tags --always --dirty) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" . && \
		GOOS=darwin GOARCH=arm64 go build -o ../../dist/go/tldrpp-darwin-arm64 -ldflags "-X main.version=$(shell git describe --tags --always --dirty) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" . && \
		GOOS=windows GOARCH=amd64 go build -o ../../dist/go/tldrpp-windows-amd64.exe -ldflags "-X main.version=$(shell git describe --tags --always --dirty) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" .

release-python:
	@echo "Building Python release package..."
	@python -m build

release: release-go release-python
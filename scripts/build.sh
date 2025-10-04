#!/bin/bash

# tldr++ build script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to build Go version
build_go() {
    print_status "Building Go version..."
    
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.22+ to build the Go version."
        return 1
    fi
    
    # Check Go version
    go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
    required_version="1.22"
    
    if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_error "Go version $go_version is too old. Please install Go 1.22+ to build the Go version."
        return 1
    fi
    
    # Create bin directory
    mkdir -p bin
    
    # Build binary
    cd cmd/tldrpp
    go build -o ../../bin/tldrpp-go \
        -ldflags "-X main.version=$(git describe --tags --always --dirty) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" .
    cd ../..
    
    print_status "Go binary built successfully: bin/tldrpp-go"
}

# Function to build Python version
build_python() {
    print_status "Building Python version..."
    
    if ! command_exists python3; then
        print_error "Python 3 is not installed. Please install Python 3.11+ to build the Python version."
        return 1
    fi
    
    # Check Python version
    python_version=$(python3 -c "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')")
    required_version="3.11"
    
    if [ "$(printf '%s\n' "$required_version" "$python_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_error "Python version $python_version is too old. Please install Python 3.11+ to build the Python version."
        return 1
    fi
    
    # Install build dependencies
    if ! command_exists pip; then
        print_error "pip is not installed. Please install pip to build the Python version."
        return 1
    fi
    
    # Install build tools
    pip install --user build wheel
    
    # Build package
    python3 -m build
    
    print_status "Python package built successfully: dist/"
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    
    # Run Go tests
    if command_exists go; then
        print_status "Running Go tests..."
        go test -v ./...
    else
        print_warning "Skipping Go tests (Go not installed)"
    fi
    
    # Run Python tests
    if command_exists python3; then
        print_status "Running Python tests..."
        python3 -m pytest tests/ -v
    else
        print_warning "Skipping Python tests (Python 3 not installed)"
    fi
}

# Function to clean build artifacts
clean() {
    print_status "Cleaning build artifacts..."
    rm -rf bin/
    rm -rf dist/
    rm -rf build/
    rm -rf *.egg-info/
    if command_exists go; then
        go clean
    fi
    print_status "Clean completed"
}

# Function to show help
show_help() {
    echo "tldr++ build script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --go        Build Go version only"
    echo "  --python    Build Python version only"
    echo "  --test      Run tests"
    echo "  --clean     Clean build artifacts"
    echo "  --all       Build both versions (default)"
    echo "  --help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                # Build both versions"
    echo "  $0 --go           # Build Go version only"
    echo "  $0 --python       # Build Python version only"
    echo "  $0 --test         # Run tests"
    echo "  $0 --clean        # Clean build artifacts"
}

# Main function
main() {
    local build_go_flag=false
    local build_python_flag=false
    local test_flag=false
    local clean_flag=false
    local all_flag=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --go)
                build_go_flag=true
                all_flag=false
                shift
                ;;
            --python)
                build_python_flag=true
                all_flag=false
                shift
                ;;
            --test)
                test_flag=true
                shift
                ;;
            --clean)
                clean_flag=true
                shift
                ;;
            --all)
                all_flag=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Execute actions
    if [ "$clean_flag" = true ]; then
        clean
    fi
    
    if [ "$test_flag" = true ]; then
        run_tests
    fi
    
    if [ "$all_flag" = true ]; then
        build_go
        build_python
    else
        if [ "$build_go_flag" = true ]; then
            build_go
        fi
        
        if [ "$build_python_flag" = true ]; then
            build_python
        fi
    fi
    
    print_status "Build completed successfully!"
}

# Run main function with all arguments
main "$@"
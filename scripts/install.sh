#!/bin/bash

# tldr++ installation script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_header() {
    echo -e "${BLUE}[tldr++]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "darwin"
    elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

# Function to detect architecture
detect_arch() {
    case $(uname -m) in
        x86_64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        armv7l)
            echo "armv7"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Function to install Go version
install_go() {
    print_status "Installing Go version..."
    
    local os=$(detect_os)
    local arch=$(detect_arch)
    
    if [ "$os" = "unknown" ] || [ "$arch" = "unknown" ]; then
        print_error "Unsupported OS/architecture combination: $os/$arch"
        return 1
    fi
    
    # Create bin directory
    mkdir -p bin
    
    # Build binary
    cd cmd/tldrpp
    go build -o ../../bin/tldrpp-go \
        -ldflags "-X main.version=$(git describe --tags --always --dirty) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" .
    cd ../..
    
    # Install to system
    local install_dir="$HOME/.local/bin"
    mkdir -p "$install_dir"
    cp bin/tldrpp-go "$install_dir/tldrpp-go"
    chmod +x "$install_dir/tldrpp-go"
    
    print_status "Go binary installed to $install_dir/tldrpp-go"
    
    # Add to PATH if not already there
    if ! echo "$PATH" | grep -q "$install_dir"; then
        print_warning "Add $install_dir to your PATH:"
        echo "export PATH=\"$install_dir:\$PATH\""
    fi
}

# Function to install Python version
install_python() {
    print_status "Installing Python version..."
    
    if ! command_exists python3; then
        print_error "Python 3 is not installed. Please install Python 3.11+ first."
        return 1
    fi
    
    # Check Python version
    python_version=$(python3 -c "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')")
    required_version="3.11"
    
    if [ "$(printf '%s\n' "$required_version" "$python_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_error "Python version $python_version is too old. Please install Python 3.11+ first."
        return 1
    fi
    
    # Install dependencies
    if command_exists pip; then
        pip install --user -e ".[full]"
    elif command_exists pip3; then
        pip3 install --user -e ".[full]"
    else
        print_error "pip is not installed. Please install pip first."
        return 1
    fi
    
    print_status "Python package installed successfully"
}

# Function to install both versions
install_both() {
    print_status "Installing both Go and Python versions..."
    
    # Install Go version
    if command_exists go; then
        install_go
    else
        print_warning "Go not installed, skipping Go version"
    fi
    
    # Install Python version
    if command_exists python3; then
        install_python
    else
        print_warning "Python 3 not installed, skipping Python version"
    fi
}

# Function to show help
show_help() {
    echo "tldr++ installation script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --go        Install Go version only"
    echo "  --python    Install Python version only"
    echo "  --both      Install both versions (default)"
    echo "  --help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                # Install both versions"
    echo "  $0 --go           # Install Go version only"
    echo "  $0 --python       # Install Python version only"
    echo ""
    echo "Requirements:"
    echo "  Go version: Go 1.22+"
    echo "  Python version: Python 3.11+"
    echo "  Dependencies: pip (for Python version)"
}

# Function to check requirements
check_requirements() {
    print_status "Checking requirements..."
    
    local go_ok=false
    local python_ok=false
    
    # Check Go
    if command_exists go; then
        go_version=$(go version | cut -d' ' -f3 | sed 's/go//')
        required_version="1.22"
        
        if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" = "$required_version" ]; then
            print_status "Go $go_version ✓"
            go_ok=true
        else
            print_warning "Go version $go_version is too old (requires 1.22+)"
        fi
    else
        print_warning "Go not installed"
    fi
    
    # Check Python
    if command_exists python3; then
        python_version=$(python3 -c "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')")
        required_version="3.11"
        
        if [ "$(printf '%s\n' "$required_version" "$python_version" | sort -V | head -n1)" = "$required_version" ]; then
            print_status "Python $python_version ✓"
            python_ok=true
        else
            print_warning "Python version $python_version is too old (requires 3.11+)"
        fi
    else
        print_warning "Python 3 not installed"
    fi
    
    # Check pip
    if command_exists pip || command_exists pip3; then
        print_status "pip ✓"
    else
        print_warning "pip not installed"
    fi
    
    if [ "$go_ok" = false ] && [ "$python_ok" = false ]; then
        print_error "No suitable runtime found. Please install Go 1.22+ or Python 3.11+"
        exit 1
    fi
}

# Main function
main() {
    print_header "tldr++ Installation Script"
    echo ""
    
    local install_go_flag=false
    local install_python_flag=false
    local install_both_flag=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --go)
                install_go_flag=true
                install_both_flag=false
                shift
                ;;
            --python)
                install_python_flag=true
                install_both_flag=false
                shift
                ;;
            --both)
                install_both_flag=true
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
    
    # Check requirements
    check_requirements
    echo ""
    
    # Execute installation
    if [ "$install_both_flag" = true ]; then
        install_both
    else
        if [ "$install_go_flag" = true ]; then
            install_go
        fi
        
        if [ "$install_python_flag" = true ]; then
            install_python
        fi
    fi
    
    echo ""
    print_status "Installation completed successfully!"
    print_status "Run 'tldrpp --help' to get started"
}

# Run main function with all arguments
main "$@"
# Installation Guide

This guide covers how to install tldr++ on different platforms and using different methods.

## Prerequisites

### Go Version (Recommended)
- Go 1.22 or later
- Git (for downloading dependencies)

### Python Version
- Python 3.11 or later
- pip (Python package manager)

## Installation Methods

### 1. Using the Installation Script (Recommended)

The easiest way to install tldr++ is using the provided installation script:

```bash
# Clone the repository
git clone https://github.com/makalin/tldrpp.git
cd tldrpp

# Run the installation script
./scripts/install.sh
```

The script will:
- Check your system requirements
- Build the appropriate version(s) for your system
- Install the binary/package
- Set up the configuration

### 2. Manual Installation

#### Go Version

```bash
# Clone the repository
git clone https://github.com/makalin/tldrpp.git
cd tldrpp

# Build the binary
make build-go

# Install the binary
make install-go
```

#### Python Version

```bash
# Clone the repository
git clone https://github.com/makalin/tldrpp.git
cd tldrpp

# Install the package
pip install -e ".[full]"
```

### 3. Using Package Managers

#### Homebrew (macOS)

```bash
# Add the tap (when available)
brew tap makalin/tldrpp

# Install
brew install tldrpp
```

#### pip (Python)

```bash
# Install from PyPI (when available)
pip install tldrpp[full]
```

## Platform-Specific Instructions

### macOS

1. Install Xcode Command Line Tools:
   ```bash
   xcode-select --install
   ```

2. Install Go (if using Go version):
   ```bash
   brew install go
   ```

3. Install Python (if using Python version):
   ```bash
   brew install python@3.11
   ```

4. Follow the installation script or manual installation steps above.

### Linux (Ubuntu/Debian)

1. Install dependencies:
   ```bash
   sudo apt update
   sudo apt install git build-essential
   ```

2. Install Go (if using Go version):
   ```bash
   # Download and install Go
   wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. Install Python (if using Python version):
   ```bash
   sudo apt install python3.11 python3-pip
   ```

4. Follow the installation script or manual installation steps above.

### Linux (CentOS/RHEL/Fedora)

1. Install dependencies:
   ```bash
   sudo dnf install git gcc
   ```

2. Install Go (if using Go version):
   ```bash
   # Download and install Go
   wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. Install Python (if using Python version):
   ```bash
   sudo dnf install python3.11 python3-pip
   ```

4. Follow the installation script or manual installation steps above.

### Windows

#### Using WSL (Recommended)

1. Install WSL2:
   ```powershell
   wsl --install
   ```

2. Follow the Linux installation instructions above.

#### Using Git Bash

1. Install Git for Windows
2. Install Go or Python following the official guides
3. Use Git Bash to run the installation script

## Post-Installation

### 1. Initialize tldr++

After installation, initialize tldr++ by downloading the page index:

```bash
tldrpp --init
```

### 2. Verify Installation

Test that tldr++ is working correctly:

```bash
# Test basic functionality
tldrpp --help

# Test with a command
tldrpp tar
```

### 3. Configuration

tldr++ will create a configuration file at `~/.config/tldrpp/config.yml`. You can customize:

- Theme (light, dark, solarized)
- Platform filters
- Keybindings
- Cache settings

Example configuration:

```yaml
theme: "dark"
platforms: ["common", "linux"]
confirm_destructive: true
clipboard: true
pager: "less -R"
keymap:
  run: "ctrl+enter"
  copy: "y"
  paste: "p"
cache_ttl_hours: 72
```

## Troubleshooting

### Common Issues

#### "Command not found" after installation

Make sure the installation directory is in your PATH:

```bash
# For Go version
export PATH="$HOME/.local/bin:$PATH"

# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Permission denied errors

Make sure the installation script is executable:

```bash
chmod +x scripts/install.sh
```

#### Go version too old

Update Go to version 1.22 or later:

```bash
# Check current version
go version

# Update Go (example for Linux)
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```

#### Python version too old

Update Python to version 3.11 or later:

```bash
# Check current version
python3 --version

# Update Python (example for Ubuntu)
sudo apt install python3.11 python3-pip
```

### Getting Help

If you encounter issues:

1. Check the [troubleshooting section](#troubleshooting) above
2. Review the [README.md](README.md) for usage instructions
3. Open an issue on [GitHub](https://github.com/makalin/tldrpp/issues)
4. Check the [FAQ](FAQ.md) for common questions

## Uninstallation

### Go Version

```bash
# Remove the binary
rm ~/.local/bin/tldrpp-go

# Remove configuration and cache
rm -rf ~/.config/tldrpp
rm -rf ~/.cache/tldrpp
```

### Python Version

```bash
# Uninstall the package
pip uninstall tldrpp

# Remove configuration and cache
rm -rf ~/.config/tldrpp
rm -rf ~/.cache/tldrpp
```

## Development Installation

For development and contributing:

```bash
# Clone the repository
git clone https://github.com/makalin/tldrpp.git
cd tldrpp

# Set up development environment
make setup-go setup-python

# Run tests
make test

# Build both versions
make build-go build-python
```
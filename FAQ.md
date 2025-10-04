# Frequently Asked Questions (FAQ)

## General Questions

### What is tldr++?

tldr++ is an interactive terminal UI for tldr pages that provides fuzzy search, inline placeholder editing, and command execution capabilities. It's designed to keep you in the terminal while working with cheat-sheets.

### How is tldr++ different from the regular tldr command?

- **Interactive UI**: tldr++ provides a full terminal UI instead of static page output
- **Fuzzy Search**: Search across all pages and platforms with intelligent matching
- **Inline Editing**: Edit placeholders directly in the terminal with live prompts
- **Command Execution**: Execute commands safely with confirmation for destructive operations
- **Multiple Platforms**: Filter by platform (common, linux, osx, sunos, windows, android)

### Which version should I use - Go or Python?

Both versions provide the same functionality:

- **Go version**: Faster startup, single binary, no dependencies
- **Python version**: More features (rich UI), easier to extend, requires Python 3.11+

Choose based on your preferences and system requirements.

## Installation Questions

### Do I need both Go and Python installed?

No, you only need one of them:
- For the Go version: Install Go 1.22+
- For the Python version: Install Python 3.11+

### Can I install both versions?

Yes, you can install both versions. They will coexist without conflicts.

### What if I don't have Go or Python installed?

You can install them using your system's package manager:
- **macOS**: `brew install go` or `brew install python@3.11`
- **Ubuntu/Debian**: `sudo apt install golang-go` or `sudo apt install python3.11`
- **CentOS/RHEL**: `sudo dnf install golang` or `sudo dnf install python3.11`

### How do I update tldr++?

For the Go version:
```bash
go install github.com/makalin/tldrpp/cmd/tldrpp@latest
```

For the Python version:
```bash
pip install --upgrade tldrpp[full]
```

## Usage Questions

### How do I search for commands?

Start typing in the search box. tldr++ will fuzzy-search across:
- Command names
- Descriptions
- Example content

### How do I edit placeholders?

1. Select an example
2. Press Tab to enter edit mode
3. Use arrow keys to navigate between placeholders
4. Type values or use history (↑/↓)

### How do I execute commands?

- **Safe execution**: Ctrl+Enter (shows command first)
- **Copy to clipboard**: y
- **Paste to terminal**: p

### What are destructive commands?

Commands that can cause data loss or system changes:
- `rm`, `dd`, `mkfs`, `iptables`
- `chmod`, `chown`, `kill`
- `shutdown`, `reboot`

These require confirmation before execution.

### How do I change the theme?

Edit your config file at `~/.config/tldrpp/config.yml`:
```yaml
theme: "light"  # or "dark" or "solarized"
```

Or use the command line:
```bash
tldrpp --theme light
```

### How do I filter by platform?

Use the platform filter:
```bash
tldrpp --platform linux
```

Or press 1-6 in the UI to toggle platforms:
- 1: common
- 2: linux
- 3: osx
- 4: sunos
- 5: windows
- 6: android

## Configuration Questions

### Where is the configuration file?

The config file is located at `~/.config/tldrpp/config.yml`.

### What configuration options are available?

```yaml
theme: "dark"                    # UI theme
platforms: ["common", "linux"]    # Default platforms
confirm_destructive: true         # Confirm destructive commands
clipboard: true                   # Enable clipboard integration
pager: "less -R"                 # Pager command
keymap:                          # Custom keybindings
  run: "ctrl+enter"
  copy: "y"
  paste: "p"
cache_ttl_hours: 72              # Cache expiration time
```

### How do I reset the configuration?

Delete the config file and restart tldr++:
```bash
rm ~/.config/tldrpp/config.yml
tldrpp --init
```

## Cache Questions

### Where is the cache stored?

The cache is stored in `~/.cache/tldrpp/pages/`.

### How do I update the cache?

```bash
tldrpp --update
```

Or press 'r' in the UI to refresh.

### How often should I update the cache?

The cache automatically refreshes every 72 hours by default. You can change this in the config file.

### Can I clear the cache?

Yes, delete the cache directory:
```bash
rm -rf ~/.cache/tldrpp
tldrpp --init
```

## Plugin Questions

### What plugins are available?

Currently, the submit plugin is available for submitting examples to tldr-pages.

### How do I use the submit plugin?

```bash
# Initialize a submission
tldrpp plugin submit init

# Validate the example
tldrpp plugin submit validate

# Create a pull request
tldrpp plugin submit create-pr
```

### Do I need GitHub CLI for the submit plugin?

Yes, the submit plugin requires GitHub CLI (`gh`) to create pull requests automatically.

Install it:
- **macOS**: `brew install gh`
- **Linux**: Follow the [GitHub CLI installation guide](https://cli.github.com/manual/installation)

## Troubleshooting Questions

### tldr++ won't start

1. Check if the binary is in your PATH
2. Verify the installation: `tldrpp --version`
3. Check for error messages

### "Command not found" error

Make sure tldr++ is in your PATH:
```bash
# Check if it's installed
which tldrpp

# Add to PATH if needed
export PATH="$HOME/.local/bin:$PATH"
```

### The UI looks broken

1. Check your terminal supports Unicode
2. Try a different theme: `tldrpp --theme light`
3. Ensure your terminal is wide enough (minimum 80 columns)

### Search results are empty

1. Initialize the cache: `tldrpp --init`
2. Update the cache: `tldrpp --update`
3. Check your platform filters

### Commands don't execute

1. Check if the command is destructive (requires confirmation)
2. Verify the command syntax
3. Check your shell permissions

### Performance issues

1. Clear the cache: `rm -rf ~/.cache/tldrpp`
2. Use the Go version for better performance
3. Reduce the number of platforms in your config

## Development Questions

### How do I contribute?

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### How do I run tests?

```bash
# Run all tests
make test

# Run Go tests only
make test-go

# Run Python tests only
make test-python
```

### How do I build from source?

```bash
# Build Go version
make build-go

# Build Python version
make build-python

# Build both
make build-go build-python
```

### How do I set up development environment?

```bash
# Set up Go environment
make setup-go

# Set up Python environment
make setup-python

# Set up both
make setup-go setup-python
```

## License Questions

### What license is tldr++ under?

tldr++ is licensed under the MIT License. See [LICENSE](LICENSE) for details.

### Can I use tldr++ in commercial projects?

Yes, the MIT License allows commercial use.

### Can I modify and distribute tldr++?

Yes, the MIT License allows modification and distribution.

## Support Questions

### Where can I get help?

1. Check this FAQ
2. Read the [README](README.md)
3. Open an issue on [GitHub](https://github.com/makalin/tldrpp/issues)
4. Check the [installation guide](INSTALL.md)

### How do I report bugs?

1. Check if the issue is already reported
2. Provide detailed information:
   - OS and version
   - tldr++ version
   - Steps to reproduce
   - Expected vs actual behavior
3. Open an issue on [GitHub](https://github.com/makalin/tldrpp/issues)

### How do I request features?

1. Check if the feature is already requested
2. Describe the feature and its use case
3. Open an issue on [GitHub](https://github.com/makalin/tldrpp/issues)
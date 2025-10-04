# tldr++ (Interactive Cheat-Sheets)

> **TL;DR:** tldr is static—you still copy-paste. **tldr++** is a terminal UI that lets you fuzzy-search pages, edit placeholders inline, then paste or execute the final command.

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](#license)
![Status](https://img.shields.io/badge/status-alpha-blue)
![Platforms](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20WSL-informational)
![tldr pages](https://img.shields.io/badge/data-tldr%20pages-orange)

---

## Features

* 🔎 **Fuzzy search** across all tldr pages & platforms
* ✏️ **Inline placeholders** (`{{file}}`, `{{port}}`) with live prompts & history
* ⚙️ **One-key actions:** copy to clipboard, paste to shell, or execute safely
* 🗂 **Platforms & aliases:** common/osx/linux/sunos/windows/android
* 🧩 **Plugin hook:** propose new examples back to the official tldr repo
* 💾 **Offline cache** with auto-refresh
* 🎨 **Themes** (light/dark/solarized) & keymap customization

---

## Why?

Copy-pasting from static pages breaks flow. tldr++ keeps you **inside** the terminal: search → fill → run.

---

## Install

Choose one stack. Both behave the same in UI.

### Option A — Go (fast, single binary)

```bash
# requires Go 1.22+
go install github.com/makalin/tldrpp/cmd/tldrpp@latest
# first run downloads page index
tldrpp --init
```

### Option B — Python (rich/textual)

```bash
# Python 3.11+
pipx install "tldrpp[full]"
# or:
python -m pip install --user "tldrpp[full]"
tldrpp --init
```

> Dependencies (Python): `rich`, `textual`, `pyperclip` (optional), `xdg`/`appdirs`.

---

## Quick Start

```bash
tldrpp               # open UI
tldrpp tar           # open UI focused on "tar"
tldrpp --platform linux --theme solarized
```

* Start typing to filter commands/pages.
* Press **Enter** to preview examples.
* Use **Tab** to jump between placeholders and fill values.
* Hit **Ctrl+Enter** to run, **y** to copy, **p** to paste.

---

## UI at a Glance

* **Search** (top): fuzzy across `command`, `desc`, `example`.
* **Pages** (left): grouped by platform; `a` to toggle all/common.
* **Examples** (center): select with arrows; preview updates live.
* **Preview** (bottom): final command with substituted values.
* **Help** (`?`): keymap cheatsheet.

---

## Keybindings

| Action                  | Key                 |
| ----------------------- | ------------------- |
| Accept example          | `Enter`             |
| Edit next placeholder   | `Tab` / `Shift+Tab` |
| Run command (safe)      | `Ctrl+Enter`        |
| Copy to clipboard       | `y`                 |
| Paste to tty*           | `p`                 |
| Toggle platform filters | `1..6` / `a`        |
| Refresh cache           | `r`                 |
| Open in pager           | `o`                 |
| Help                    | `?`                 |
| Quit                    | `q` / `Ctrl+C`      |

* Paste sends keystrokes to the parent TTY (tmux supported).

---

## Safety & Exec Model

* **Dry-run by default:** first run shows the fully rendered command.
* **Confirm before exec:** destructive verbs (rm, dd, mkfs, iptables) trigger a confirm screen.
* **Audit log:** saved under `~/.cache/tldrpp/exec.log`.

---

## Placeholder Editing

tldr examples use `{{…}}`. tldr++ prompts you inline:

* Type value, or press **↑** for recent values
* Use **:file**, **:dir**, **:port**, **:num** suffixes to get validators
* Press **Ctrl+r** for ripgrep-based file search (optional)

---

## Configuration

`~/.config/tldrpp/config.yml`

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

---

## Data & Caching

* Sources: [tldr-pages/tldr](https://github.com/tldr-pages/tldr)
* Cache dir: `~/.cache/tldrpp/pages/`
* Update: background refresh or `tldrpp --update`

---

## Plugin: **Propose Example → tldr**

Opt-in plugin to draft a PR back to `tldr-pages`:

```bash
tldrpp --plugin submit
```

* Opens the currently rendered example as a markdown snippet
* Guides you through page conventions & style checks
* Creates a branch + commit; you confirm before pushing
* Works with GitHub CLI (`gh`) if present

---

## Scripting / Headless Mode

```bash
# render best example for "tar extract", fill placeholders, print
tldrpp render "tar extract" --vars file=archive.tar.gz dest=.
# execute directly (with confirm)
tldrpp exec "ffmpeg convert" --vars in=raw.mov out=out.mp4
```

---

## Development

### Go

```bash
git clone https://github.com/makalin/tldrpp
cd tldrpp
go run ./cmd/tldrpp --dev
```

### Python

```bash
git clone https://github.com/makalin/tldrpp
cd tldrpp
uv venv && uv pip install -e ".[dev]"
python -m tldrpp --dev
```

Tests:

```bash
go test ./...
# or
pytest -q
```

---

## Roadmap

* ⏱ Session snippets & multi-cursor placeholder editing
* 🧠 LLM-assisted example explanations (offline first)
* 🔌 More plugins: export to README.md / gist / Obsidian
* 🧪 Inline “try in tmpdir” sandbox
* 🔒 Per-command allow/deny lists

---

## Contributing

PRs welcome! Please:

1. Open an issue with a brief proposal
2. Follow the existing keymap & UX patterns
3. Add tests and update docs

---

## License

MIT © Mehmet T. AKALIN. See [LICENSE](LICENSE).

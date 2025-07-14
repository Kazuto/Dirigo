# 🛠️ Dirigo (lat. "I lead")

`dirigo` is a simple Go CLI tool that lets you define and quickly navigate to frequently used project directories via short aliases. It supports path lookup, interactive selection (TUI), and adding new entries — all stored in a YAML config file.

---

## 📦 Features

- 🔖 Alias-based path lookup: `dirigo ui` → `/development/dirigo/ui`
- 🧭 Interactive TUI path picker: `dirigo --pick`
- ➕ Add new paths via CLI: `dirigo --add <key> <path>`
- 🗂️ Stores config in `XDG_CONFIG_HOME/dirigo/paths.yml` or `~/.config/dirigo/paths.yml`
- 📝 YAML-based configuration

---

## 🚀 Usage

```bash
# List all paths
dirigo

# Add a new path alias
dirigo --add api ~/projects/backend/api

# Remove a path alias
dirigo --remove api

# Get the full path for an alias
dirigo api

# Use it with cd
cd "$(dirigo api)"
```

or using a shell function

```bash
dirigo() {
  dir=$("$HOME/.local/bin/dirigo" "$@")
  if [ -d "$dir" ]; then
    cd "$dir"
  else
    echo "Dirigo: Could not resolve or enter directory"
  fi
}
```

## 📦 Installation

1. Download the binary
2. Extract it and move to `~/.local/bin`
3. Make it executable `chmod +x ~/.local/bin/dirigo`

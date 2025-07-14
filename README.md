# 🛠️ Dirigo (lat. "I lead")

`dirigo` is a simple Go CLI tool that lets you define and quickly navigate to frequently used project directories via short aliases. It supports path lookup, adding new entries and removing old entries — all stored in a YAML config file.

---

## 📦 Features

- 🔖 Alias-based path lookup: `dirigo ui` → `/development/dirigo/ui`
- ➕ Add new paths via CLI: `dirigo --add <key> <path>`
- 🗂️ Stores config in `XDG_CONFIG_HOME/dirigo/paths.yml` or `~/.config/dirigo/paths.yml`
- 📝 YAML-based configuration

---

## 🚀 Usage

```bash
# List all paths
dirigo [--list]

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
  if [[ "$1" == --* ]]; then
    command dirigo "$@"
  else
    local dir
    dir=$(command dirigo "$1") || return 1
    if [ -d "$dir" ]; then
      cd "$dir" || return
    else
      echo "Directory does not exist: $dir"
      return 1
    fi
  fi
}
```

## 📦 Installation

1. Download the binary
2. Extract it and move to `/usr/local/bin`
3. Make it executable `chmod +x /usr/local/bin/dirigo`

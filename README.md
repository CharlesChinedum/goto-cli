# gotocli

[![CI](https://github.com/CharlesChinedum/goto-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/CharlesChinedum/goto-cli/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A cross-platform CLI tool to save and jump to directories quickly.

## Installation

### Mac & Linux

```bash
sudo mv gotocli /usr/local/bin/
```

### Windows

Move `gotocli.exe` to `C:\Program Files\gotocli\`

---

## Setup Shell Wrapper

### Mac & Linux — add to `~/.zshrc` or `~/.bashrc`

```bash
function goto() {
    local bin=${GOTOCLI_BIN:-gotocli}
    local command=$1
    local name=$2
    local path=$3
    local extra=$4

    if [ "$command" = "jump" ]; then
        TARGET=$("$bin" goto jump "$name" 2>/dev/null)
        if [ -z "$TARGET" ]; then
            echo "No directory found for '$name'"
        else
            cd "$TARGET"
        fi
    elif [ "$command" = "edit" ]; then
        "$bin" goto edit "$name" "$path"
    elif [ "$command" = "rename" ]; then
        "$bin" goto rename "$name" "$path"
    else
        "$bin" goto "$command" "$name" "$path" "$extra"
    fi
}
```

### Windows — add to PowerShell `$PROFILE`

```powershell
function goto {
    param($command, $name, $path, $extra)
    $bin = if ($env:GOTOCLI_BIN) { $env:GOTOCLI_BIN } else { "gotocli" }

    if ($command -eq "jump") {
        $TARGET = & $bin goto jump $name 2>$null
        if (-not $TARGET) {
            Write-Host "No directory found for '$name'"
        } else {
            Set-Location $TARGET
        }
    } elseif ($command -eq "edit") {
        & $bin goto edit $name $path
    } elseif ($command -eq "rename") {
        & $bin goto rename $name $path
    } else {
        & $bin goto $command $name $path $extra
    }
}
```

---

## Usage

### Add a directory

```bash
goto add <name> <path>

# Example
goto add projects /home/user/projects
```

### List all saved directories

```bash
goto list
```

### Jump to a directory

```bash
goto jump <name>

# Example
goto jump projects
```

### Remove a saved directory

```bash
goto remove <name>

# Example
goto remove projects
```

If the name is not saved, the command prints an error to stderr and exits with
status 1. Your saved directories are left unchanged.

### Edit a saved directory path

```bash
goto edit <name> <newpath>

# Example
goto edit projects /home/user/new-projects
```

### Rename a saved directory

```bash
goto rename <oldname> <newname>

# Example
goto rename projects work-projects
```

---

## Commands Summary

| Command       | Usage                             | Action                     |
| ------------- | --------------------------------- | -------------------------- |
| `goto add`    | `goto add <name> <path>`          | Save a directory           |
| `goto list`   | `goto list`                       | List all saved directories |
| `goto jump`   | `goto jump <name>`                | Jump to a directory        |
| `goto remove` | `goto remove <name>`              | Remove a saved directory; fail if not found |
| `goto edit`   | `goto edit <name> <newpath>`      | Update a directory path    |
| `goto rename` | `goto rename <oldname> <newname>` | Rename a saved directory   |

---

## How it works

- The **Go binary** (`gotocli`) handles all data storage and retrieval, saving directories to `~/.goto.json`
- The **shell wrapper** (`goto`) intercepts the `jump` command and runs `cd` in the current shell
- This is the same pattern used by popular tools like `z` and `autojump`

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for how to build from source, the checks CI runs on every pull request, and how to submit changes. This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

---

## Security

Please do not report security issues in public issues. See [SECURITY.md](SECURITY.md) for how to report a vulnerability privately.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

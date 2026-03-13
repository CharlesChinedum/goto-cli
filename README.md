# gotocli

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
    local command=$1
    local name=$2
    local path=$3
    local extra=$4

    if [ "$command" = "jump" ]; then
        TARGET=$(/usr/local/bin/gotocli goto jump "$name" 2>/dev/null)
        if [ -z "$TARGET" ]; then
            echo "No directory found for '$name'"
        else
            cd "$TARGET"
        fi
    elif [ "$command" = "edit" ]; then
        /usr/local/bin/gotocli goto edit "$name" "$path"
    elif [ "$command" = "rename" ]; then
        /usr/local/bin/gotocli goto rename "$name" "$path"
    else
        /usr/local/bin/gotocli goto "$command" "$name" "$path" "$extra"
    fi
}
```

### Windows — add to PowerShell `$PROFILE`

```powershell
function goto {
    param($command, $name, $path, $extra)

    if ($command -eq "jump") {
        $TARGET = & "C:\Program Files\gotocli\gotocli.exe" goto jump $name 2>$null
        if (-not $TARGET) {
            Write-Host "No directory found for '$name'"
        } else {
            Set-Location $TARGET
        }
    } elseif ($command -eq "edit") {
        & "C:\Program Files\gotocli\gotocli.exe" goto edit $name $path
    } elseif ($command -eq "rename") {
        & "C:\Program Files\gotocli\gotocli.exe" goto rename $name $path
    } else {
        & "C:\Program Files\gotocli\gotocli.exe" goto $command $name $path $extra
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
| `goto remove` | `goto remove <name>`              | Remove a saved directory   |
| `goto edit`   | `goto edit <name> <newpath>`      | Update a directory path    |
| `goto rename` | `goto rename <oldname> <newname>` | Rename a saved directory   |

---

## How it works

- The **Go binary** (`gotocli`) handles all data storage and retrieval, saving directories to `~/.goto.json`
- The **shell wrapper** (`goto`) intercepts the `jump` command and runs `cd` in the current shell
- This is the same pattern used by popular tools like `z` and `autojump`

---

## License

MIT

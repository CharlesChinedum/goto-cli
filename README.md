# gotocli

A cross-platform CLI tool to save and jump to directories quickly.

## Installation

### Mac & Linux

```bash
sudo mv gotocli /usr/local/bin/
```

### Windows

Move `gotocli.exe` to `C:\Program Files\gotocli\`

## Setup shell wrapper

### Mac & Linux — add to `~/.zshrc` or `~/.bashrc`

```bash
function goto() {
    local command=$1
    local name=$2
    local path=$3

    if [ "$command" = "jump" ]; then
        TARGET=$(/usr/local/bin/gotocli goto jump "$name" 2>/dev/null)
        if [ -z "$TARGET" ]; then
            echo "No directory found for '$name'"
        else
            cd "$TARGET"
        fi
    else
        /usr/local/bin/gotocli goto "$command" "$name" "$path"
    fi
}
```

### Windows — add to PowerShell `$PROFILE`

```powershell
function goto {
    param($command, $name, $path)

    if ($command -eq "jump") {
        $TARGET = & "C:\Program Files\gotocli\gotocli.exe" goto jump $name 2>$null
        if (-not $TARGET) {
            Write-Host "No directory found for '$name'"
        } else {
            Set-Location $TARGET
        }
    } else {
        & "C:\Program Files\gotocli\gotocli.exe" goto $command $name $path
    }
}
```

## Usage

```bash
goto add projects /home/user/projects   # save a directory
goto list                               # list all saved directories
goto jump projects                      # jump to a directory
goto remove projects                    # remove a saved directory
```

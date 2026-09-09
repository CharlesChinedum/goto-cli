# Contributing to gotocli

Thanks for your interest in contributing! gotocli is a small, single-binary Go CLI, so getting started is quick.

## Prerequisites

- [Go](https://go.dev/dl/) 1.25.1 or newer (the version pinned in `go.mod`)
- Git
- A POSIX shell (bash/zsh) or PowerShell, for exercising the shell wrapper described in the README

## Getting the code

```bash
git clone git@github.com:CharlesChinedum/goto-cli.git
cd goto-cli
```

## Building

The whole tool lives in `app/main.go`.

```bash
go build -o gotocli ./app
```

Then try it out without installing anything:

```bash
./gotocli goto list
./gotocli goto add scratch /tmp
./gotocli goto jump scratch
```

> **Heads up:** the binary reads and writes `~/.goto.json`, the same file your installed copy uses. Back it up (`cp ~/.goto.json ~/.goto.json.bak`) before testing changes that touch storage.

## Project layout

| Path                            | Purpose                                                        |
| ------------------------------- | -------------------------------------------------------------- |
| `app/main.go`                   | Entire CLI: argument parsing, JSON store, all subcommands      |
| `.github/workflows/ci.yml`      | Build, vet, gofmt, test and cross-compile checks on every PR   |
| `.github/workflows/release.yml` | Builds the four platform binaries and attaches them to a release when a `v*` tag is pushed |
| `README.md`                     | User-facing docs, including the shell wrapper                  |

## Before opening a pull request

CI runs exactly these commands (see `.github/workflows/ci.yml`). Run them locally first so your PR is green on the first push:

```bash
go build -o /dev/null ./app   # compiles (plain `go build ./...` clashes with the app/ dir name)
go vet ./...          # static analysis
gofmt -l .            # must print nothing; run `gofmt -w .` to fix
go test ./...         # runs any tests
```

CI also confirms the code still cross-compiles for every release target:

```bash
GOOS=darwin  GOARCH=amd64 go build -o /dev/null app/main.go
GOOS=darwin  GOARCH=arm64 go build -o /dev/null app/main.go
GOOS=linux   GOARCH=amd64 go build -o /dev/null app/main.go
GOOS=windows GOARCH=amd64 go build -o /dev/null app/main.go
```

Also:

- Add a line under `[Unreleased]` in `CHANGELOG.md` describing your change.
- If you add, rename, or change a subcommand, update the **Usage** and **Commands Summary** sections of `README.md` and both shell wrappers (bash/zsh and PowerShell).
- Keep changes cross-platform. Anything touching paths or the home directory must work on macOS, Linux, and Windows.
- Do not commit compiled binaries. The release workflow builds them.

## Submitting changes

1. Fork the repo and create a branch from `main` (the maintainer also uses a `development` branch for staging work).
2. Make your change, run the checks above.
3. Open a pull request. The PR template will ask you to confirm the checklist.
4. A maintainer (see `.github/CODEOWNERS`) will review it.

## Reporting bugs and requesting features

Use the issue templates. Bug reports should include your OS, shell, the exact `goto ...` command you ran, and what happened instead of what you expected.

## Security issues

Please do **not** open a public issue for security problems. See [SECURITY.md](SECURITY.md).

## Code of conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). By participating you agree to abide by it.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).

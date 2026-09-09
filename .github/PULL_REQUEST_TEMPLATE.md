## Summary

<!-- What does this PR change, and why? Link the issue it closes, e.g. "Closes #12". -->

## Type of change

- [ ] Bug fix
- [ ] New feature or subcommand
- [ ] Documentation only
- [ ] CI / tooling

## Checklist

These mirror the gates in `.github/workflows/ci.yml`; run them locally before pushing.

- [ ] `go build -o /dev/null ./app` succeeds
- [ ] `go vet ./...` reports nothing
- [ ] `gofmt -l .` prints nothing (run `gofmt -w .` to fix)
- [ ] `go test ./...` passes
- [ ] Still cross-compiles for all release targets (`GOOS=darwin|linux|windows`, see CONTRIBUTING.md)
- [ ] Added a line under `[Unreleased]` in `CHANGELOG.md`
- [ ] If a subcommand was added or changed: updated `README.md` (Usage + Commands Summary) and **both** shell wrappers (bash/zsh and PowerShell)
- [ ] No compiled binaries or personal `~/.goto.json` data included in the diff

## How I tested it

<!-- OS + shell, and the exact `goto ...` commands you ran. -->

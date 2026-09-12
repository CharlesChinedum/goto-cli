# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Removing an unknown bookmark now reports an error and exits with status 1 without rewriting or creating the store.
- Store writes now go through a temp file, fsync, and rename so `~/.goto.json` is not truncated in place. A failed save prints an error to stderr and exits 1 instead of claiming success.

### Added

- MIT `LICENSE` file (the README previously claimed MIT without shipping the text)
- `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md` (Contributor Covenant 2.1), and `SECURITY.md`
- GitHub issue forms, pull request template, `CODEOWNERS`, and Dependabot configuration
- `.editorconfig`
- CI workflow that builds, vets, gofmt-checks, tests, and cross-compiles on every push and pull request

[Unreleased]: https://github.com/CharlesChinedum/goto-cli/compare/main...HEAD

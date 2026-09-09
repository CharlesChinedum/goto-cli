# Security Policy

## Supported versions

gotocli is distributed as prebuilt binaries attached to GitHub releases. Only the **latest release** receives security fixes. If you are on an older version, please upgrade before reporting.

## What counts as a security issue

gotocli runs locally, has no network access, and stores its data in `~/.goto.json`. Reports we consider in scope include, for example:

- Reading from or writing to files outside `~/.goto.json` in a way the user did not ask for
- Crafted `~/.goto.json` contents or command arguments that cause the shell wrapper to execute unintended commands
- Anything in the release workflow that could let a tampered binary be published

Ordinary bugs (wrong output, crashes on bad input with no wider impact) are welcome as regular [bug reports](https://github.com/CharlesChinedum/goto-cli/issues/new/choose).

## Reporting a vulnerability

**Please do not open a public issue for security problems.**

Report privately through GitHub's private vulnerability reporting:

1. Go to <https://github.com/CharlesChinedum/goto-cli/security/advisories/new>
2. Fill in the form with as much detail as you can: affected version, OS and shell, steps to reproduce, and the impact you observed.

Reports are visible only to the maintainer. No email address is published for this project.

## What to expect

- You will get an acknowledgement within **7 days**.
- The maintainer will work with you to confirm the issue and agree on a fix and disclosure timeline. For a project this size, the target is a fix within **30 days** of confirmation.
- Once a fix is released, a GitHub Security Advisory will be published and you will be credited unless you ask not to be.

Thank you for helping keep gotocli users safe.

# Changelog

All notable changes to this project will be documented in this file.
## [0.2.1] - 2026-06-30

### Bug Fixes

- Replace deprecated brews with homebrew_casks in GoReleaser config
- Pass HOMEBREW_TAP_GITHUB_TOKEN env var to GoReleaser action

## [0.2.0] - 2026-06-30

### Features

- Pin to API v9 with Accept header on every request
- Show API version in `paperless version` and warn on mismatch
- Add Homebrew tap to GoReleaser release pipeline

## [0.1.5] - 2026-06-30

### Features

- Pin to API version 10 and check server compatibility

## [0.1.4] - 2026-06-30

### Features

- Add configure command to write user config file

## [0.1.3] - 2026-06-30

### Bug Fixes

- Translate all user-facing strings from German to English

### Features

- Add CLI binary version to version output

## [0.1.2] - 2026-06-30

### Bug Fixes

- Revert to root main.go, use GitHub Releases as primary install path

## [0.1.1] - 2026-06-30

### Bug Fixes

- Move main to cmd/paperless so go install creates binary named paperless

## [0.1.0] - 2026-06-30

### Bug Fixes

- Correct module path, graceful SSH degradation, trim SKILL.md
- Replace hardcoded install path with standard PREFIX convention
- **ci:** Extract issue body generation to avoid YAML heredoc parse error
- **ci:** Switch schema check to GitHub releases version comparison
- **config:** Use ~/.config/paperless-ngx-cli/ as config dir
- **ci:** Use github changelog provider instead of git-cliff in goreleaser

### Features

- Add generated Paperless-NGX API client
- Add core CLI commands
- Add bulk operations, manage, and version commands
- Fetch schema from Docker image — no running instance required
- **release:** Add git-cliff for changelog generation

### Security

- **security:** Shell-quote args and container name in SSH commands


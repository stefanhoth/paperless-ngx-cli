# Changelog

All notable changes to this project will be documented in this file.
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


# Changelog

All notable changes to this project will be documented in this file.
## [2.0.0] - 2026-08-03

### Bug Fixes

- **release:** Open a PR for the changelog instead of pushing to main (#29)
- **deps:** Update module github.com/oapi-codegen/runtime to v1.6.0 (#30)
- **ci:** Strip v-prefix before pulling ghcr.io paperless-ngx image (#32)
- **ci:** Bypass s6-overlay entrypoint for one-off manage.py commands (#33)
- **ci:** Pretty-print schema before diffing, cap issue body size (#34)

### Features

- Target Paperless-NGX 3.x and API v10 (#36)

## [1.1.2] - 2026-07-16

### Bug Fixes

- **deps:** Bump go directive to 1.26.5 (#13)
- Replace unresolvable :pinDigests preset in renovate.json (#23)
- **deps:** Update module github.com/oapi-codegen/runtime to v1.5.0 (#25)

### Features

- **lint:** Add golangci-lint, gofumpt, and pre-commit hook (#9)
- Add `skill install` command (#27)

## [1.1.1] - 2026-07-15

### Bug Fixes

- Prefix /api on raw api passthrough (#8)

## [1.1.0] - 2026-07-15

### Bug Fixes

- Normalize version display — strip v prefix from latest version

### Features

- Add paperless api raw REST passthrough (#6)

### Ai

- Add mattpocock/skills suite (#5)

## [1.0.0] - 2026-06-30

### Bug Fixes

- Add macOS quarantine removal hook to Homebrew cask

### Features

- Remove SSH/Docker coupling — pure-API CLI (#3)

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


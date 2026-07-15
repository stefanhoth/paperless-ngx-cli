# GitHub repo settings: PR checks & auto-merge

**Date:** 2026-07-15
**Status:** workflows active; one manual step below is still outstanding

## Goal

Every PR against `main` is verified automatically (Lint, Build, Test,
Conventional commit title), and green PRs merge via auto-merge — nothing
lands on `main` unchecked.

## Already in the repo

| File | Purpose |
| --- | --- |
| `.github/workflows/ci.yml` | Lint/Build/Test jobs on every PR and push to main (Test includes a `govulncheck` audit) |
| `.github/workflows/pr-title.yml` | Enforces conventional-commit PR titles (the squash commit on main) |
| `.github/rulesets/main-branch-protection.json` | Importable ruleset requiring those four checks, squash-only, linear history |
| `renovate.json` | Automerge for minor/patch/pin/digest once checks are green |

## Already done

- ✅ **Allow auto-merge** — enabled.
- ✅ **Automatically delete head branches** — enabled.
- ✅ Squash-merge only, merge/rebase commit disabled.

## One-time manual steps still outstanding

1. **Settings → Rules → Rulesets → New ruleset ▾ → Import a ruleset**:
   import `.github/rulesets/main-branch-protection.json`. Until this is
   done, CI results are informational only — a red PR can still merge
   (this happened once while `ci.yml` was first being added).
2. **Secrets** (Settings → Secrets and variables → Actions):
   `HOMEBREW_TAP_GITHUB_TOKEN` — a fine-grained PAT scoped to
   `stefanhoth/homebrew-tap` with `Contents: Read & Write`. Setup steps are
   in [docs/development.md](development.md#homebrew-tap-setup). Required for
   `release.yml` to push the Homebrew cask; confirm it's already set.
3. **Install the Renovate GitHub App** for this repo (if not already
   installed org/account-wide): https://github.com/apps/renovate.

## Maintenance warning

The check names in the ruleset (`Lint`, `Build`, `Test`,
`Conventional commit title`) must match the job `name:` fields in
`ci.yml`/`pr-title.yml` exactly. Renaming or adding a required job means
updating the ruleset under **Settings → Rules → Rulesets** — otherwise PRs
wait forever for a check that never reports.

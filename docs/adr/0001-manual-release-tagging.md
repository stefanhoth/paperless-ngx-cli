---
status: accepted
---

# Releases stay manually tagged, not auto-tagged from conventional commits

Releases are triggered by a human pushing a `vX.Y.Z` git tag, which then runs
GoReleaser + git-cliff via `.github/workflows/release.yml`. This was kept
as-is during the 2026-07 quality-and-automation setup, rather than switching
to full semantic-release automation (auto-tag + release on every merge to
`main`) as the default for a consumable CLI would otherwise suggest.

The CLI pins one major version to one Paperless-NGX REST API version (see
[README's Version Support table](../../README.md#version-support)) and sends
an explicit `Accept: application/json; version=N` header accordingly. A
version bump is a deliberate compatibility statement, not just "however
conventional commits happen to add up." Automated semantic-release derives
the bump from commit types (`feat` → minor, `fix`/`BREAKING CHANGE` → major)
with no human in the loop — a stray `fix!:` or misclassified commit could
trigger an unintended major release with no one noticing until it ships.

## Considered options

- **Manual `git tag` + GoReleaser (chosen)** — keeps the major-version ↔
  API-version mapping under deliberate human control; low overhead for a
  single-maintainer repo cutting releases infrequently.
- **semantic-release, auto-tag on every merge to `main`** — rejected: full
  automation is the default for consumable CLIs, but here it would let commit
  message classification decide a compatibility-sensitive version bump
  unsupervised.

## Consequences

- A release still requires a manual `git tag vX.Y.Z && git push` step;
  nothing auto-ships on merge.
- The CHANGELOG.md commit-back-to-main step in `release.yml`, GoReleaser
  config, and Homebrew tap publishing are all unaffected by this choice and
  keep working exactly as before.
- If release cadence increases enough that manual tagging becomes a
  bottleneck, revisit via a superseding ADR — the fix would likely be a
  human-gated "approve this version bump" step rather than full automation.

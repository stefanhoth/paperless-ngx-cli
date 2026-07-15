# Decisions

A lightweight log of smaller-but-made decisions — the calls that shape the
code or product but don't warrant a full [ADR](adr/).

- **When to use an ADR instead:** genuine architecture decisions (stack,
  hosting, release strategy) go in [docs/adr/](adr/) with `status:`
  frontmatter.
- **Format:** newest first, grouped by date. One short entry per decision —
  the call in bold, and the why. **Add the entry in the same PR that makes
  the decision.**

## 2026-07-15

- **Kept manual `git tag` releases instead of switching to semantic-release
  automation.** See [ADR-0001](adr/0001-manual-release-tagging.md) for the
  full reasoning — recorded here too since it was an explicit fork raised
  during the project-quality-setup retrofit.
- **golangci-lint's curated linter set (errcheck, govet, staticcheck,
  revive, gosec, ineffassign, unused) surfaced 22 pre-existing issues on
  first run.** Fixed the real ones (unchecked errors, unused cobra
  parameters, missing package docs); suppressed the handful of unavoidable
  false positives (gosec G304 on an intentional `--input` file flag and the
  fixed XDG config path, SA1019 on a Paperless API field that's deprecated
  upstream with no documented replacement) with a `//nolint` + justification
  comment at the site, per the zero-warning policy — never disabled the
  rule globally.
- **`go.mod`'s `go` directive was pinned to 1.26.4, one patch behind a
  stdlib security fix (GO-2026-5856, crypto/tls).** `govulncheck` — added as
  part of the new CI `Test` job — caught this on its very first run against
  `main`. Because the branch-protection ruleset didn't exist yet at that
  point, the red `Test` job didn't block the PR that introduced it from
  merging; fixed immediately after in a follow-up PR. This is exactly the
  gap the ruleset (added later in this same retrofit) closes going forward.

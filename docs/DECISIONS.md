# Decisions

A lightweight log of smaller-but-made decisions — the calls that shape the
code or product but don't warrant a full [ADR](adr/).

- **When to use an ADR instead:** genuine architecture decisions (stack,
  hosting, release strategy) go in [docs/adr/](adr/) with `status:`
  frontmatter.
- **Format:** newest first, grouped by date. One short entry per decision —
  the call in bold, and the why. **Add the entry in the same PR that makes
  the decision.**

## 2026-07-28

- **Added `paperless suggest` and `paperless chat` as dedicated commands for
  the v10 AI endpoints, rather than leaving them on the `paperless api`
  passthrough.** Both endpoints return non-JSON error bodies on failure (plain
  text `"AI is required for this feature"` on HTTP 400), so `apiError` from
  `cmd/bulk.go` — already built for that shape — is reused instead of
  `exitOnAPIError`, which assumes JSON. `chat`'s response arrives as
  `text/event-stream` with no `data:` framing (confirmed against upstream's
  `stream_chat_with_documents`, which just yields raw text chunks); since the
  generated client always buffers the full response before returning, the CLI
  prints it as plain text rather than pretending to stream.

## 2026-07-26

- **Bumped `APIVersion` to 10 for Paperless-NGX 3.x — the CLI's first major
  version bump under [ADR-0002](adr/0002-one-major-version-per-api-version.md).**
  Paperless 3.x still accepts `version=9`, but a v9 client gets the legacy
  response shapes (unpaginated tasks, old saved-view fields), so pinning to
  v10 keeps the generated client and the server in agreement. Consequence:
  CLI v2.x cannot talk to a 2.x server at all — those reject `version=10`
  with HTTP 406.
- **Moved `bulk reprocess|delete|merge|rotate` off `/documents/bulk_edit/`
  onto the dedicated endpoints added in API v10.** Upstream still accepts
  them on `bulk_edit` but logs them as deprecated, and they are scheduled to
  go away when v9 is dropped. The metadata operations (`add-tag`,
  `remove-tag`, `set-correspondent`, `set-type`) stay on `bulk_edit`, which
  remains their only endpoint.
- **Replaced the copy-pasted `if err != nil || resp.StatusCode() != 200`
  blocks with `exitOnAPIError`.** They printed `error: <nil>` for every HTTP
  failure — which the version bump would have made the standard experience
  for anyone pointing v2.x at a 2.x server. The helper names the status and
  special-cases 406 with a "server does not support API v10" hint.

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
- **Imported the branch-protection ruleset with `required_review_thread_resolution: true`**,
  not `false` as originally proposed in `.github/rulesets/main-branch-protection.json`.
  Open review threads must now be resolved before a PR can merge. The
  checked-in ruleset file was updated to match the live setting.

---
status: accepted
---

# One major CLI version per Paperless-NGX API version

The CLI targets exactly one Paperless-NGX REST API version per major CLI
version, tracked by the `APIVersion` constant in
[cmd/root.go](../../cmd/root.go) and sent as `Accept: application/json;
version=N` on every request. `paperless version` reports both the CLI's
target and the server's actual API version, and warns on mismatch. Bumping
`APIVersion` is treated as a major-version-bump event for the CLI, recorded
in the README's Version Support table.

This predates the 2026-07 quality-setup work but was never written down as a
decision — it's recorded here retroactively because it's the CLI's most
consequential and least-reversible architectural constraint: the generated
API client (`api/paperless.gen.go`) is regenerated from a schema pinned to a
specific upstream release, and command behavior/JSON shapes are only
guaranteed correct for that one API version.

## Considered options

- **Pin one API version per major CLI version (chosen)** — every command's
  behavior is well-defined against a known schema; upgrading to a new
  Paperless API version is an explicit, tested, versioned CLI release.
- **Negotiate/support multiple API versions at runtime** — rejected: would
  require branching logic per command for schema differences across API
  versions, for a maintenance cost that isn't justified by a single-user CLI
  with one Paperless instance to talk to.

## Consequences

- Older CLI major versions get critical bug fixes for a short transition
  window only, never backported features (per the README).
- The daily `schema-check.yml` drift check exists specifically to catch when
  upstream Paperless ships a new API version, so this pin doesn't silently
  go stale.
- A new stable Paperless API version requires a new major CLI version, a
  schema regeneration (`make generate-docker`), and a README Version Support
  table update — not just a patch release.

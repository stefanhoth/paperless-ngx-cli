# Design: Make the CLI a pure-API tool (remove SSH/Docker coupling)

**Date:** 2026-06-30
**Status:** Approved (approach A)

## Problem

The CLI mixes two unrelated access models in one binary:

- **REST API** (URL + token) — used by almost every command.
- **SSH + `docker exec`** — used by `manage` (fully) and by one line of `version`
  (the installed server version, via `docker inspect`).

This breaks the user's mental model. People expect an API tool, then hit commands
that suddenly require an SSH key, a reachable Docker host, and a specific container
name. The SSH surface is tiny (one command + one enrichment line) but carries a
disproportionate share of the conceptual, deployment, and maintenance cost:

- **UX (primary concern):** two trust/setup models in one tool; "why doesn't this
  command work for me?" failure class.
- **Deployment:** `manage`/`version` assume Docker + SSH host access. Bare-metal,
  k8s, and hosted instances can never use them (see issue #2).
- **Maintenance/security:** shell quoting, exec of external `ssh`, `docker exec`
  string building — a different error/attack class than plain HTTP.

## Decision

Remove the SSH/Docker functionality entirely. The CLI becomes a pure-API tool with
a single trust model: `PAPERLESS_URL` + `PAPERLESS_API_TOKEN`.

Rationale: `manage.py` commands have no API equivalent, but anyone who *can* run
them already has SSH + Docker access to the host — and can run the one-liner
`ssh host docker exec <container> python3 manage.py <cmd>` directly. The CLI wrapper
saves almost nothing while imposing the entire two-model confusion. Removing it
resolves the UX, deployment, and maintenance concerns at once, with the least code.

Alternatives considered:
- **B — keep `manage` under an `admin` namespace, signpost SSH requirement.**
  Reduces confusion but keeps two trust models and the Docker-only deployment limit.
- **C — two separate CLIs (`paperless` + `paperless-admin`).** Conceptually clean
  but overkill: duplicate config loading, two binaries to build/install/version for
  ~70 lines of SSH code.

## `version` command: installed server version stays

`RemoteVersionView` (the `/api/remote_version/` endpoint) returns only `version`
(latest GitHub release) and `update_available` — **not** the installed version.

However, the Paperless API sends the installed app version as a response header:

```
X-Version: 2.20.15
```

So the "Paperless (server)" line is preserved and sourced from
`resp.HTTPResponse.Header.Get("X-Version")` instead of SSH `docker inspect`.
No SSH anywhere in `version`.

## Scope — files to change

### Code (remove SSH)

| File | Change |
|------|--------|
| `cmd/manage.go` | Delete the whole file (command + `shellQuote`). |
| `cmd/manage_test.go` | Delete (only tests `shellQuote`). |
| `cmd/version.go` | Remove `sshInstalledVersion()`; "Paperless (server)" reads `X-Version` header. Drop now-unused imports (`os/exec`, and `strings` if unused). |
| `cmd/config.go` | Remove `sshHost`/`sshUser`/`container` struct fields and their parsing (`PAPERLESS_SSH_HOST`, `PAPERLESS_SSH_USER`, `PAPERLESS_CONTAINER`). Drop now-unused imports (`net/url`, `os/user`). |
| `cmd/configure.go` | Remove the SSH settings block (prompts + writing of SSH/container keys). Update the `Long` description. |
| `cmd/config_test.go` | Remove `TestParseConfig_SSHHostDerivedFromURL`, `TestParseConfig_SSHHostExplicitOverridesURL`, `TestParseConfig_ContainerDefault`, `TestParseConfig_ContainerExplicit`; remove SSH/container env vars from the env-clearing list and the file-values test. |

### Docs

| File | Change |
|------|--------|
| `README.md` | Lines ~5, 10, 22, 42, env table rows 79–81, and the "Management Commands (SSH)" section (110–118). Reword intro to drop "run management commands". |
| `SKILL.md` | Description (4–7), SSH activation note (22–24), "Management Commands (requires SSH)" section (60–70). |
| `docs/development.md` | Remove `manage.go` from the file tree (line ~32). |

### Deliberately untouched

These match the grep but are unrelated to the CLI `manage` command:

- `.github/workflows/schema-check.yml` — uses Docker + `manage.py spectacular` for
  **schema generation**, not the CLI.
- `AGENTS.md` — `make generate-docker` (schema), unrelated.
- `docs/development.md` Docker/Homebrew lines — build tooling, unrelated.
- `CHANGELOG.md` — historical entries; do not rewrite history. A new entry lands
  via the normal release flow (conventional commit).

### Issues

- Close **#2** ("Support bare-metal Paperless installs for `manage` and `version`
  commands") as obsolete — removing the SSH path eliminates the bare-metal problem.

## Verification

1. `go build ./...` succeeds with no unused-import errors.
2. `go test ./...` passes (remaining config tests green).
3. `grep -rinE "ssh|PAPERLESS_SSH|PAPERLESS_CONTAINER|docker exec|manage\.py" cmd/ README.md SKILL.md`
   returns nothing in CLI code/user docs (schema-generation refs excluded).
4. Manual smoke against a real instance: `paperless version` shows the server
   version from the `X-Version` header; no command references SSH.

## Out of scope

- Any replacement/admin tooling for `manage.py` operations.
- Changes to schema generation, release pipeline, or Homebrew packaging.

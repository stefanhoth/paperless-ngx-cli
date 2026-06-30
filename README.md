# paperless-ngx-cli

Your [Paperless-NGX](https://docs.paperless-ngx.com/) document archive, from the terminal.

Search documents, inspect metadata, trigger bulk operations, and run management commands — without opening a browser. Designed for scripting, automation, and AI-assisted workflows.

```bash
paperless search "Steuerbescheid 2024"
paperless bulk reprocess 10,11,12
paperless manage document_retagger
```

The CLI is statically compiled and ships as a single binary for Linux and macOS. No runtime, no dependencies.

## Why this exists

Paperless-NGX has a solid web UI, but the API is the real power interface. Once you expose it on the command line, you can:

- **automate** document processing pipelines via shell scripts or cron jobs
- **integrate** Paperless into AI agents and Claude Code workflows using the bundled [SKILL.md](SKILL.md)
- **bulk-operate** on document sets that would take dozens of clicks in the UI
- **monitor** your instance and trigger reindexing or OCR without SSH-ing into the server manually

The client is generated from Paperless-NGX's own OpenAPI spec, so commands and types stay accurate as the API evolves. A daily CI check detects new Paperless releases and opens a GitHub issue when the schema needs updating.

## How it works

Built with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) — the Go client is generated directly from the Paperless-NGX OpenAPI schema. The generated client gives compile-time safety against API changes: when Paperless updates its schema, `make generate-docker` catches breaking changes at build time rather than at runtime.

## Requirements

- Go 1.21+
- A running Paperless-NGX instance (tested with 2.20.x)
- SSH access to the Docker host (only for `manage` and `version` commands)

## Installation

Download the binary for your platform from [GitHub Releases](https://github.com/stefanhoth/paperless-ngx-cli/releases/latest), extract, and place `paperless` somewhere in your `$PATH`.

**Build from source:**
```bash
git clone https://github.com/stefanhoth/paperless-ngx-cli
cd paperless-ngx-cli
make install   # installs to /usr/local/bin/paperless
```

## Configuration

Set variables in your shell profile, or create a config file at `~/.config/paperless-ngx-cli/config`:

```ini
PAPERLESS_URL=http://paperless.local:8000
PAPERLESS_API_TOKEN=your-token-here
```

Environment variables always take precedence over the config file.

**Variables:**

| Variable | Required | Description |
|---|---|---|
| `PAPERLESS_URL` | Yes | Base URL, e.g. `http://paperless.local:8000` |
| `PAPERLESS_API_TOKEN` | Yes | API token from Paperless Settings → API |
| `PAPERLESS_SSH_HOST` | No | SSH host for `manage`/`version` (defaults to hostname from `PAPERLESS_URL`) |
| `PAPERLESS_SSH_USER` | No | SSH username (defaults to current OS user) |
| `PAPERLESS_CONTAINER` | No | Docker container name (defaults to `paperless-ngx-webserver-1`) |

Get your API token at `http://your-paperless/api/auth/token/` or in the Paperless web UI under Settings → API.

## Usage

```bash
paperless status
paperless docs -n 20
paperless search "Invoice Amazon"
paperless doc 1234
paperless doc 1234 --full-perms
paperless tags
paperless correspondents
paperless types
paperless version
```

### Bulk Operations

```bash
# IDs are comma-separated
paperless bulk reprocess 1,2,3
paperless bulk delete 42
paperless bulk add-tag 10,11,12 7       # add tag ID 7
paperless bulk set-correspondent 5 3    # set correspondent ID 3
paperless bulk rotate 99 90
```

### Management Commands (SSH)

These run `manage.py` inside the Paperless container via SSH:

```bash
paperless manage document_retagger
paperless manage document_sanity_checker
paperless manage document_index reindex
paperless manage document_archiver       # re-run OCR on all documents
```

## Regenerating the API Client

The generated client lives in `api/paperless.gen.go`. Regenerate it when Paperless updates its API:

```bash
# No running instance needed — pulls the official Docker image and exports the schema:
make generate-docker VERSION=v2.20.15
make build

# Or, if you have a running Paperless instance:
make generate   # uses PAPERLESS_URL + PAPERLESS_API_TOKEN
make build
```

The `scripts/fix-schema.py` script patches two known issues in the Paperless OpenAPI schema before generation:
- Renames response types that conflict with oapi-codegen's generated wrapper names
- Fixes `Tag.children` which the API returns as objects but the schema declares as integer IDs

## Project Structure

```
├── api/                    # generated client (do not edit)
│   └── paperless.gen.go
├── cmd/                    # CLI commands
│   ├── root.go             # client setup and config
│   ├── status.go
│   ├── docs.go
│   ├── search.go
│   ├── doc.go
│   ├── list.go             # tags, correspondents, types
│   ├── bulk.go
│   ├── manage.go
│   └── version.go
├── schema/
│   └── paperless.json      # vendored API schema (patched)
├── scripts/
│   └── fix-schema.py       # schema patches for regeneration
├── SKILL.md                # AI assistant skill descriptor
├── Makefile
└── main.go
```

## License

MIT

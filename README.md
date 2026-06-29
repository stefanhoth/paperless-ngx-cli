# paperless-ngx-cli

A typed command-line interface for [Paperless-NGX](https://docs.paperless-ngx.com/), generated from the official OpenAPI spec using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Why

The generated Go client gives you compile-time safety against API changes. When Paperless updates its schema, running `make generate` catches breaking changes at build time rather than at runtime.

## Requirements

- Go 1.21+
- A running Paperless-NGX instance (tested with 2.20.x)
- SSH access to the Docker host (only for `manage` and `version` commands)

## Installation

```bash
git clone https://github.com/yourname/paperless-ngx-cli
cd paperless-ngx-cli
make build
# binary is at ./paperless
```

Or install directly:
```bash
go install github.com/stefanhoth/paperless-ngx-cli@latest
```

## Configuration

All configuration is via environment variables. Add these to your shell profile:

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
make generate   # downloads schema from PAPERLESS_URL, patches it, regenerates client
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

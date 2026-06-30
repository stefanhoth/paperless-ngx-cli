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

## Version Support

| CLI version | Status | Paperless-NGX | API version |
|---|---|---|---|
| **v1.x** | ✅ Active | 2.x stable | v9 |
| v2.x | Planned | 3.x (not yet stable) | v10 |

**One major CLI version per Paperless API version.** The CLI pins to a specific API version and sends `Accept: application/json; version=N` with every request, so responses are always in the expected format even when the server supports multiple API versions.

When Paperless-NGX ships a stable 3.x series with API v10, CLI v2.0.0 will follow. Older major versions receive no backported features, but may receive critical bug fixes for a short transition window.

Run `paperless version` to verify compatibility — it prints the CLI's target API version and warns if your server reports a different API version in its response headers.

## Requirements

- A running Paperless-NGX 2.x instance
- SSH access to the Docker host (only for `manage` and `version` commands)

## Installation

**Homebrew** (macOS / Linux — recommended):
```bash
brew tap stefanhoth/tap
brew install paperless-ngx-cli
```

**Binary** — download from [GitHub Releases](https://github.com/stefanhoth/paperless-ngx-cli/releases/latest), extract, and place `paperless` somewhere in your `$PATH`.

For building from source, see [docs/development.md](docs/development.md).

## Configuration

The easiest way to configure the CLI is the interactive setup command:

```bash
paperless configure
```

This prompts for your Paperless URL and API token and writes them to `~/.config/paperless-ngx-cli/config` with secure permissions (`0600`). Alternatively, set variables in your shell profile or create the config file manually:

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
| `PAPERLESS_CONTAINER` | No | Docker container name (defaults to `paperless-ngx-webserver-1`). Only relevant when Paperless runs as a Docker container — the CLI will use `docker exec` to run management commands inside it. |

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

## Contributing

See [docs/development.md](docs/development.md) for build instructions, project structure, the API client regeneration workflow, and release process.

## License

MIT

# paperless-ngx-cli

Your [Paperless-NGX](https://docs.paperless-ngx.com/) document archive, from the terminal.

Search documents, inspect metadata, and trigger bulk operations — without opening a browser. Designed for scripting, automation, and AI-assisted workflows.

```bash
paperless search "Steuerbescheid 2024"
paperless bulk reprocess 10,11,12
```

The CLI is statically compiled and ships as a single binary for Linux and macOS. No runtime, no dependencies.

## Why this exists

Paperless-NGX has a solid web UI, but the API is the real power interface. Once you expose it on the command line, you can:

- **automate** document processing pipelines via shell scripts or cron jobs
- **integrate** Paperless into AI agents and Claude Code workflows using the bundled [SKILL.md](SKILL.md) (`paperless skill install`)
- **bulk-operate** on document sets that would take dozens of clicks in the UI

The client is generated from Paperless-NGX's own OpenAPI spec, so commands and types stay accurate as the API evolves. A daily CI check detects new Paperless releases and opens a GitHub issue when the schema needs updating.

## Version Support

| CLI version | Status | Paperless-NGX | API version |
|---|---|---|---|
| **v2.x** | ✅ Active | 3.x stable | v10 |
| v1.x | 🛠 Bug fixes only | 2.x | v9 |

**One major CLI version per Paperless API version.** The CLI pins to a specific API version and sends `Accept: application/json; version=N` with every request, so responses are always in the expected format even when the server supports multiple API versions.

CLI v2.x targets Paperless-NGX 3.x (API v10). If your server still runs 2.x, stay on CLI v1.x — v2.x sends `version=10`, which a 2.x server rejects with HTTP 406. Older major versions receive no backported features, but may receive critical bug fixes for a short transition window.

Paperless-NGX 3.x also ships opt-in LLM features (per-document suggestions, document chat). If you run [paperless-gpt](https://github.com/icereed/paperless-gpt) next to your instance, [docs/paperless-3-ai.md](docs/paperless-3-ai.md) covers what changes for you — short version: nothing, until you switch them on.

Run `paperless version` to verify compatibility — it prints the CLI's target API version and warns if your server reports a different API version in its response headers.

## Requirements

- A running Paperless-NGX 3.x instance (for 2.x servers, use CLI v1.x)

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

### AI Suggestions & Chat

Requires Paperless-NGX 3.x with AI enabled server-side — see
[docs/paperless-3-ai.md](docs/paperless-3-ai.md).

```bash
paperless suggest 1234                                    # title/tag/correspondent suggestions
paperless chat "What is the invoice total?" --doc 1234     # ask about one document
paperless chat "Which documents mention BubbleTax?"        # ask across all documents
```

### Raw API Escape Hatch

For anything not covered by the commands above, `paperless api` is a
`gh api`-style passthrough: same auth, raw JSON on stdout.

```bash
paperless api /documents/4028/ --method PATCH --field created=2022-02-08
paperless api /documents/4028/ --method PATCH --input body.json
paperless api "/documents/?created__date=2026-07-08" | jq '.results[].id'
```

### AI Assistant Skill

The CLI ships with a [SKILL.md](SKILL.md) describing itself to AI assistants
like Claude Code. Install it directly from the binary instead of a separate
`skills add` step. It's written to the vendor-neutral
`.agents/skills/paperless-ngx-cli/SKILL.md`, with `.claude/skills/paperless-ngx-cli`
symlinked to it so Claude Code picks it up too:

```bash
paperless skill install          # ./.agents/skills/paperless-ngx-cli/SKILL.md (+ ./.claude/skills symlink)
paperless skill install --user   # ~/.agents/skills/paperless-ngx-cli/SKILL.md (+ ~/.claude/skills symlink)
```

## Contributing

See [docs/development.md](docs/development.md) for build instructions, project structure, the API client regeneration workflow, and release process. Architecture decisions are recorded in [docs/adr/](docs/adr/) and [docs/DECISIONS.md](docs/DECISIONS.md). See [SECURITY.md](SECURITY.md) for the threat model.

## License

MIT

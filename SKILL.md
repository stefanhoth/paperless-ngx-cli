---
name: paperless-ngx-cli
description: >
  CLI for Paperless-NGX document management. Use when the user wants to search,
  list, or manage documents, tags, correspondents, or document types in a
  Paperless-NGX instance. Also handles bulk operations and version checks.
---

# Paperless-NGX CLI

Assumes the `paperless` binary is available in PATH and configured. Configuration
is read from env vars or `~/.config/paperless-ngx-cli/config` (env takes precedence):

```
PAPERLESS_URL           Base URL, e.g. http://paperless.local:8000
PAPERLESS_API_TOKEN     API token from Paperless Settings → API
```

If not yet configured, run `paperless configure` for interactive setup.

---

## Commands

```
paperless status                     System stats (total documents, tags, types, etc.)
paperless docs [-n <count>]          Recent documents (default: 10), newest first
paperless search <query>             Full-text search, up to 20 results
paperless doc <id>                   Document details (title, date, tags, type, pages)
paperless doc <id> --full-perms      Same, plus full permission info
paperless tags                       All tags with document count
paperless correspondents             All correspondents with document count
paperless types                      All document types with document count
paperless configure                  Interactive setup — writes ~/.config/paperless-ngx-cli/config
paperless version                    CLI version, target API version, and Paperless instance version
```

## Bulk Operations

IDs are comma-separated (`1,2,3`). Operations run asynchronously.

```
paperless bulk reprocess <ids>
paperless bulk delete <ids>
paperless bulk merge <ids>
paperless bulk rotate <ids> <90|180|270>
paperless bulk add-tag <ids> <tag_id>
paperless bulk remove-tag <ids> <tag_id>
paperless bulk set-correspondent <ids> <correspondent_id>
paperless bulk set-type <ids> <type_id>
```

Use `tags`, `correspondents`, or `types` to look up numeric IDs first.

## Raw API Escape Hatch

For anything the dedicated commands above don't cover (e.g. setting a single
arbitrary field like `created` on one document), use `paperless api` —
same auth, prints raw JSON for piping into `jq`.

```
paperless api <path> [-X <method>] [-f key=value ...] [--input <file>|-]
```

```
paperless api /documents/4028/ --method PATCH --field created=2022-02-08
paperless api /documents/4028/ --method PATCH --input body.json
paperless api "/documents/?created__date=2026-07-08" | jq '.results[].id'
```

Defaults: `GET` with no body; `POST` when `-f`/`--input` is given and `-X` is
omitted. `-f` and `--input` are mutually exclusive. Non-2xx exits non-zero.

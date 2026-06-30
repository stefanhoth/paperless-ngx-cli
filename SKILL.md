---
name: paperless-ngx-cli
description: >
  CLI for Paperless-NGX document management. Use when the user wants to search,
  list, or manage documents, tags, correspondents, or document types in a
  Paperless-NGX instance. Also handles bulk operations and (when SSH is
  configured) running management commands inside the Paperless container.
---

# Paperless-NGX CLI

Assumes the `paperless` binary is available in PATH and configured. Configuration
is read from env vars or `~/.config/paperless-ngx-cli/config` (env takes precedence):

```
PAPERLESS_URL           Base URL, e.g. http://paperless.local:8000
PAPERLESS_API_TOKEN     API token from Paperless Settings → API
```

If not yet configured, run `paperless configure` for interactive setup.

SSH-based commands (`manage`, and the installed-version part of `version`) are
optional. They activate when `PAPERLESS_SSH_HOST` is set (or derivable from
`PAPERLESS_URL`). Without SSH, those commands print a setup hint and exit cleanly.

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

## Management Commands (requires SSH)

Runs `manage.py` inside the Paperless Docker container via SSH.
Prints a setup hint and exits cleanly if SSH is not configured.

```
paperless manage document_retagger          Re-apply matching rules to all documents
paperless manage document_renamer           Regenerate filenames from storage template
paperless manage document_index reindex     Rebuild the full-text search index
paperless manage document_sanity_checker    Report inconsistencies (missing files, no OCR, etc.)
paperless manage document_archiver          Re-run OCR on all documents
```

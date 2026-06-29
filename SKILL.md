---
name: paperless-ngx-cli
description: >
  CLI for Paperless-NGX document management. Use when the user wants to search,
  list, or manage documents, tags, correspondents, or document types in their
  Paperless-NGX instance. Also handles bulk operations and running management
  commands in the Paperless container via SSH.
---

# Paperless-NGX CLI

A typed Go CLI generated from the Paperless-NGX OpenAPI spec via
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Setup

Build and install:
```bash
make build        # produces ./paperless binary
make install      # copies it to ../bin/paperless (if inside the synology project)
```

Set environment variables (add to `~/.zshenv` or `~/.bashrc`):
```bash
export PAPERLESS_URL=http://your-paperless-host:8000
export PAPERLESS_API_TOKEN=your_api_token_here

# For manage and version commands (SSH into the Docker host):
export PAPERLESS_SSH_HOST=your-docker-host      # defaults to PAPERLESS_URL hostname
export PAPERLESS_SSH_USER=your-ssh-username     # defaults to current OS user
export PAPERLESS_CONTAINER=paperless-ngx-webserver-1  # default
```

## Commands

```
paperless status                          # System stats (documents, tags, etc.)
paperless docs [-n <count>]              # Recent documents (default: 10)
paperless search <query>                 # Full-text search (default: 20 results)
paperless doc <id>                       # Document details
paperless doc <id> --full-perms          # Document details with permissions
paperless tags                           # List all tags with document count
paperless correspondents                 # List all correspondents
paperless types                          # List all document types
paperless version                        # Installed vs. available version (needs SSH)
paperless manage <cmd> [args]            # Run manage.py in container (needs SSH)
```

## Bulk Operations

IDs are comma-separated: `1,2,3`

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

Bulk operations run asynchronously — only a confirmation is returned.

## Management Commands (via SSH)

These require SSH access to the Docker host:

```
paperless manage document_retagger          # Re-apply matching rules
paperless manage document_renamer           # Regenerate filenames
paperless manage document_index reindex     # Rebuild search index
paperless manage document_sanity_checker    # Check for inconsistencies
paperless manage document_archiver          # Re-run OCR on documents
```

## ID Lookup

Use `tags`, `correspondents`, and `types` to get numeric IDs before bulk operations:

```bash
paperless tags                   # find tag ID
paperless bulk add-tag 42 7      # add tag 7 to document 42
```

## Regenerating the API Client

When Paperless-NGX is updated and the API schema changes:

```bash
make generate    # downloads schema, patches it, regenerates api/paperless.gen.go
make install     # rebuild and install
```

Requires `PAPERLESS_URL` and `PAPERLESS_API_TOKEN` to be set for the schema download.

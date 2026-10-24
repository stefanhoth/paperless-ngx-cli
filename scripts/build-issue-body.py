#!/usr/bin/env python3
"""Builds a per-version schema-drift issue body and writes it to /tmp/issue-body.md.

Unlike the old collapsed single-issue format, this is called once per
schema-changing release in the walk from the tracked version to upstream, so
the diff is always against the immediately preceding release — never the
original tracked base — and the fully regenerated schema + client are already
pushed to a branch instead of being re-derived by whoever picks this up.

Reads env vars:
  VERSION        — the release this issue is about (e.g. v3.1.3)
  PREV_VERSION   — the release the diff is against (the one right before it)
  REPO           — owner/repo (e.g. stefanhoth/paperless-ngx-cli)
  BRANCH         — branch holding the regenerated schema/client for VERSION
  DIFF_FILE      — path to the pretty-printed `diff -u` output (optional)
"""
import datetime
import os
import sys
from pathlib import Path

version = os.environ["VERSION"]
prev_version = os.environ["PREV_VERSION"]
repo = os.environ["REPO"]
branch = os.environ["BRANCH"]
diff_file = os.environ.get("DIFF_FILE", "")
today = datetime.date.today().isoformat()

MAX_DIFF_CHARS = 40_000  # keep well under GitHub's 65536-char issue body limit

diff_section = ""
if diff_file and Path(diff_file).exists():
    diff = Path(diff_file).read_text().strip()
    if len(diff) > MAX_DIFF_CHARS:
        diff = diff[:MAX_DIFF_CHARS] + "\n... (truncated — see the full regenerated schema on the branch below)"
    if diff:
        diff_section = f"""
### Schema diff ({prev_version} → {version})

```diff
{diff}
```
"""

raw_base = f"https://raw.githubusercontent.com/{repo}/{branch}"

body = f"""\
## Paperless-NGX {version}: schema change detected

Upstream release **`{version}`** changes the OpenAPI schema relative to the \
release right before it, **`{prev_version}`** — this issue is scoped to just \
that one step, not the whole gap back to whatever this repo currently tracks.
{diff_section}
### Already generated for you

The schema was pulled from the official `{version}` Docker image, patched with \
`scripts/fix-schema.py`, and the client was regenerated with the same pinned \
`oapi-codegen` version the Makefile uses — so implementing this is pulling the \
branch and reviewing it, not redoing the Docker-pull/spectacular/codegen work:

```bash
git fetch origin {branch}
git checkout {branch}
make build   # verify it still compiles
make test
```

Direct links to the regenerated files on that branch:
- [`schema/paperless.json`]({raw_base}/schema/paperless.json)
- [`api/paperless.gen.go`]({raw_base}/api/paperless.gen.go)
- [`.paperless-version`]({raw_base}/.paperless-version)

> **Heads up if other schema-drift issues are open too:** the *diff above* is \
against `{prev_version}` (the immediately preceding release), but this \
*branch* is cut from `main` as it stands today — so if an earlier version's \
PR merges first, this branch will show a conflict on the generated files when \
you go to merge it. Resolve it by merging `main` in and re-running \
`make generate-docker VERSION={version}` (never by hand-editing the generated \
files), the same way the predecessor release was handled.

**Detected:** {today}
"""

out_path = "/tmp/issue-body.md"
with open(out_path, "w") as f:
    f.write(body)

print(f"Written: {out_path}", file=sys.stderr)

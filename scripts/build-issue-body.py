#!/usr/bin/env python3
"""Builds the schema-drift issue body and writes it to /tmp/issue-body.md.

Reads env vars:
  TRACKED_VERSION   — version in .paperless-version (what we last generated against)
  UPSTREAM_VERSION  — latest Paperless-NGX GitHub release tag
  HAS_DIFF          — "true" if a full schema diff is available in /tmp/schema-diff.txt
"""
import datetime
import os
import sys

tracked = os.environ.get("TRACKED_VERSION", "unknown")
upstream = os.environ.get("UPSTREAM_VERSION", "unknown")
has_diff = os.environ.get("HAS_DIFF", "").lower() == "true"
today = datetime.date.today().isoformat()

diff_section = ""
if has_diff:
    try:
        with open("/tmp/schema-diff.txt") as f:
            diff = f.read().strip()
        if diff:
            diff_section = f"""
### Schema diff (first 120 lines)

```diff
{diff}
```
"""
    except FileNotFoundError:
        pass

body = f"""\
## New Paperless-NGX release detected

A new upstream release is available. The vendored schema and generated client \
may need to be updated.

| | Version |
|---|---|
| **Tracked** (last generated) | `{tracked}` |
| **Upstream** (latest release) | `{upstream}` |

**Detected:** {today}
{diff_section}
### To update

```bash
make generate   # downloads schema from your instance, patches it, regenerates client
make build      # verify it compiles
git add schema/paperless.json api/paperless.gen.go .paperless-version
git commit -m "chore: update to Paperless-NGX {upstream}"
```

> **Tip:** Set `PAPERLESS_URL` and `PAPERLESS_API_TOKEN` as repository secrets \
to enable a full schema diff in future checks.
"""

out_path = "/tmp/issue-body.md"
with open(out_path, "w") as f:
    f.write(body)

print(f"Written: {out_path}", file=sys.stderr)

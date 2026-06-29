#!/usr/bin/env python3
"""Builds the schema-drift issue body and writes it to /tmp/issue-body.md.

Reads:
  DEMO_VERSION env var
  /tmp/schema-diff.txt
Writes:
  /tmp/issue-body.md
"""
import datetime
import os
import sys

demo_version = os.environ.get("DEMO_VERSION", "unknown")
today = datetime.date.today().isoformat()

diff_path = "/tmp/schema-diff.txt"
try:
    with open(diff_path) as f:
        diff = f.read()
except FileNotFoundError:
    diff = "(diff file not found)"

body = f"""\
## Paperless-NGX API schema has changed

The daily schema check detected a difference between the vendored `schema/paperless.json`
and the current schema from the [Paperless-NGX demo](https://demo.paperless-ngx.com/).

**Detected:** {today}
**Demo API version:** `{demo_version}`

### Diff (first 120 lines)

```diff
{diff}
```

### To update

```bash
make generate   # downloads new schema, patches it, regenerates api/paperless.gen.go
make build      # verify it compiles
```

Then commit the updated `schema/paperless.json` and `api/paperless.gen.go`.
"""

out_path = "/tmp/issue-body.md"
with open(out_path, "w") as f:
    f.write(body)

print(f"Written: {out_path}", file=sys.stderr)

#!/usr/bin/env python3
"""Patches the Paperless NGX OpenAPI schema before code generation.

Fixes:
  1. Rename schema types that conflict with oapi-codegen's generated response wrappers.
  2. Fix Tag.children which the API returns as Tag objects but the schema declares as int[].
"""
import json, sys
from pathlib import Path

schema_path = Path(__file__).parent.parent / "schema" / "paperless.json"

with open(schema_path) as f:
    content = f.read()

# 1. Rename conflicting response types (schema models that clash with generated wrappers)
renames = {
    "EmailDocumentsResponse": "EmailDocumentsModel",
    "MailAccountProcessResponse": "MailAccountProcessModel",
    "MailAccountTestResponse": "MailAccountTestModel",
}
for old, new in renames.items():
    content = content.replace(old, new)

schema = json.loads(content)

# 2. Fix Tag.children: schema says int[] but API returns full Tag objects
tag = schema.get("components", {}).get("schemas", {}).get("Tag", {})
if "children" in tag.get("properties", {}):
    tag["properties"]["children"] = {"type": "array", "items": {}, "readOnly": True}

with open(schema_path, "w") as f:
    json.dump(schema, f, ensure_ascii=False)

print(f"Schema patched: {schema_path}")

#!/usr/bin/env python3
"""Patches the Paperless NGX OpenAPI schema before code generation.

Fixes:
  1. Rename schema types that conflict with oapi-codegen's generated response wrappers.
  2. Fix Tag.children which the API returns as Tag objects but the schema declares as int[].

Usage:
  fix-schema.py                          # patch schema/paperless.json in-place
  fix-schema.py --input <src> --output <dst>  # patch src, write to dst
"""
import argparse
import json
import sys
from pathlib import Path

ROOT = Path(__file__).parent.parent


def patch(content: str) -> str:
    # Rename types that conflict with oapi-codegen's generated response wrapper names
    renames = {
        "EmailDocumentsResponse": "EmailDocumentsModel",
        "MailAccountProcessResponse": "MailAccountProcessModel",
        "MailAccountTestResponse": "MailAccountTestModel",
    }
    for old, new in renames.items():
        content = content.replace(old, new)

    schema = json.loads(content)

    # Fix Tag.children: schema declares int[] but API returns full Tag objects
    tag = schema.get("components", {}).get("schemas", {}).get("Tag", {})
    if "children" in tag.get("properties", {}):
        tag["properties"]["children"] = {"type": "array", "items": {}, "readOnly": True}

    return json.dumps(schema, ensure_ascii=False)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", default=str(ROOT / "schema" / "paperless.json"))
    parser.add_argument("--output")
    args = parser.parse_args()

    src = Path(args.input)
    dst = Path(args.output) if args.output else src

    with open(src) as f:
        result = patch(f.read())

    with open(dst, "w") as f:
        f.write(result)

    print(f"Patched: {src} → {dst}", file=sys.stderr)


if __name__ == "__main__":
    main()

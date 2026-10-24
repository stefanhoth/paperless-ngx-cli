#!/usr/bin/env python3
"""Version-ordering helpers for the schema-drift workflow.

Subcommands:
  matrix --releases <file> --tracked <vX.Y.Z> --upstream <vX.Y.Z>
      Reads a GitHub releases API response (JSON array, as from
      `gh api repos/paperless-ngx/paperless-ngx/releases --paginate`),
      drops drafts/prereleases and non-`vMAJOR.MINOR.PATCH` tags, sorts
      the rest by semver, and prints a JSON {"include": [...]} matrix
      of every release strictly after `tracked` up to and including
      `upstream`, each paired with the release immediately before it
      in that sorted list (its diff base).

  cmp <vX.Y.Z> <vX.Y.Z>
      Prints -1, 0, or 1 (by semver, independent of any releases
      list) for sorting/comparing two tags.
"""
import argparse
import json
import re
import sys

VERSION_RE = re.compile(r"v(\d+)\.(\d+)\.(\d+)$")


def parse_version(tag: str):
    m = VERSION_RE.fullmatch(tag)
    return tuple(int(g) for g in m.groups()) if m else None


def load_releases(path: str) -> list[str]:
    with open(path) as f:
        releases = json.load(f)
    parsed = []
    for rel in releases:
        if rel.get("draft") or rel.get("prerelease"):
            continue
        tag = rel.get("tag_name", "")
        version = parse_version(tag)
        if version is not None:
            parsed.append((version, tag))
    parsed.sort(key=lambda item: item[0])
    return [tag for _, tag in parsed]


def cmd_matrix(args: argparse.Namespace) -> None:
    versions = load_releases(args.releases)
    for label, tag in (("tracked", args.tracked), ("upstream", args.upstream)):
        if tag not in versions:
            print(f"{label} version {tag!r} not found in fetched releases", file=sys.stderr)
            sys.exit(1)

    start = versions.index(args.tracked)
    end = versions.index(args.upstream)
    include = [
        {"version": versions[i], "prev_version": versions[i - 1]}
        for i in range(start + 1, end + 1)
    ]
    json.dump({"include": include}, sys.stdout)


def cmd_cmp(args: argparse.Namespace) -> None:
    a, b = parse_version(args.v1), parse_version(args.v2)
    if a is None or b is None:
        print(f"not a vMAJOR.MINOR.PATCH tag: {args.v1!r} / {args.v2!r}", file=sys.stderr)
        sys.exit(2)
    print(-1 if a < b else (1 if a > b else 0))


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = parser.add_subparsers(dest="command", required=True)

    matrix_parser = sub.add_parser("matrix")
    matrix_parser.add_argument("--releases", required=True)
    matrix_parser.add_argument("--tracked", required=True)
    matrix_parser.add_argument("--upstream", required=True)
    matrix_parser.set_defaults(func=cmd_matrix)

    cmp_parser = sub.add_parser("cmp")
    cmp_parser.add_argument("v1")
    cmp_parser.add_argument("v2")
    cmp_parser.set_defaults(func=cmd_cmp)

    args = parser.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()

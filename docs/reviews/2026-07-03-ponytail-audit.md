# Ponytail Audit — 2026-07-03

Repo-wide over-engineering audit. Findings ranked biggest cut first.
Scope: complexity only — correctness, security, and performance are out of scope.

## Findings

### 1. `yagni:` 30,280-line generated client for 8 used operations

The vendored OpenAPI schema defines 142 operations; the CLI calls 8:
`bulk_edit`, `correspondents_list`, `document_types_list`, `documents_list`,
`documents_retrieve`, `remote_version_retrieve`, `statistics_retrieve`, `tags_list`.

**Replacement:** add an operation filter to `oapi-codegen.yaml` and regenerate —
pruning drops unused models too, cutting `api/paperless.gen.go` to roughly 2–3k lines:

```yaml
output-options:
  include-operation-ids:
    - bulk_edit
    - correspondents_list
    - document_types_list
    - documents_list
    - documents_retrieve
    - remote_version_retrieve
    - statistics_retrieve
    - tags_list
```

Alternative: hand-roll a ~200-line `net/http` client and drop the
`oapi-codegen/runtime`, `go-jsonmerge`, and `uuid` deps entirely — but that
loses the regenerate-on-schema-drift workflow the repo is built around.
The include-filter keeps it.

**Files:** `oapi-codegen.yaml`, `api/paperless.gen.go`
**Estimated cut:** ~−27,000 lines

### 2. `shrink:` three copy-pasted list commands

`tags`, `correspondents`, and `types` have identical bodies differing only in
the API call. One helper taking a fetch closure that returns
`[]struct{ id int; name string; count string }` collapses them:
`listCmd("tags", fetchTags)`.

**Files:** `cmd/list.go`
**Estimated cut:** ~−40 lines

### 3. `shrink:` 50-line switch of identical cases

Every case in the bulk-operation switch does the same thing: check arg count,
`strconv.Atoi`, set method + param key. Table form plus one shared arg check:

```go
var ops = map[string]struct {
	m     api.MethodEnum
	param string // "" = no extra arg
}{
	"reprocess":         {api.MethodEnumReprocess, ""},
	"delete":            {api.MethodEnumDelete, ""},
	"merge":             {api.MethodEnumMerge, ""},
	"rotate":            {api.MethodEnumRotate, "degrees"},
	"add-tag":           {api.MethodEnumAddTag, "tag"},
	"remove-tag":        {api.MethodEnumRemoveTag, "tag"},
	"set-correspondent": {api.MethodEnumSetCorrespondent, "correspondent"},
	"set-type":          {api.MethodEnumSetDocumentType, "document_type"},
}
```

**Files:** `cmd/bulk.go:43-93`
**Estimated cut:** ~−35 lines

### 4. `native:` Python script that formats a markdown string from env vars

`scripts/build-issue-body.py` (65 lines) builds an issue body from three env
vars. A `cat <<EOF > /tmp/issue-body.md` heredoc step inside
`.github/workflows/schema-check.yml` does the same; the diff-section
conditional is an `if [ "$HAS_DIFF" = true ]` line.

**Files:** `scripts/build-issue-body.py`, `.github/workflows/schema-check.yml`
**Estimated cut:** ~−40 lines net, −1 file

### 5. `delete:` `ctx()` wrapper returning `context.Background()`

Cobra already provides a context: use `cmd.Context()` at call sites.

**Files:** `cmd/root.go:47-49`
**Estimated cut:** −3 lines

### 6. `shrink:` closure that just indexes a map

`get := func(key string) interface{} { return s[key] }` — write
`s["documents_total"]` directly.

**Files:** `cmd/status.go:30`
**Estimated cut:** −1 line

## Out of scope (route to a normal review pass)

- `cmd/doc.go` prints German labels/errors while every other command is English.

## Net

**net: −27,100 lines, −0 deps possible** (−3 deps if the client is hand-rolled
instead of pruned).

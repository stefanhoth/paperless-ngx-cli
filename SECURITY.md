# Security boundaries

paperless-ngx-cli holds one secret — a Paperless-NGX API token — and uses it
to talk to a single, user-configured Paperless-NGX instance over HTTP(S).

```
[user's shell] --config/env--> [paperless CLI] --Authorization: Token--> [Paperless-NGX server]
                                      |
                                      v
                          ~/.config/paperless-ngx-cli/config (0600)
```

## What this protects

| Threat | Mitigation |
|--------|-----------|
| Other local users reading the API token off disk | Config file is written and expected at `0600`; `readConfigFile` (cmd/config.go) warns to stderr if it finds group/other-readable permissions |
| Token or URL leaking into shell history / process list | Config file and env vars (`PAPERLESS_URL`, `PAPERLESS_API_TOKEN`) are the documented ways to supply credentials — no `--token` CLI flag exists |
| A malicious/incorrect `--input` file path or config path escalating beyond a plain file read | `readInputBody`/`readConfigFile` only ever call `os.ReadFile`/`os.Open` — no shell interpolation, no execution of file contents |
| CSRF/host confusion when a raw `api` path or full URL is passed | `normalizeAPIPath` (cmd/api.go) rejects any absolute URL whose scheme+host doesn't match the configured `PAPERLESS_URL` |

## What this does NOT protect against

| Threat | Notes |
|--------|-------|
| A compromised or malicious Paperless-NGX server | The CLI trusts the configured `PAPERLESS_URL` completely — it sends the token and renders whatever JSON comes back. Point it only at a Paperless instance you trust. |
| Plaintext HTTP | If `PAPERLESS_URL` is `http://`, the token and all document data travel unencrypted. The CLI does not enforce HTTPS — that's the operator's/network's responsibility. |
| The `api` command's raw REST passthrough | `paperless api` deliberately allows arbitrary method/path/body against the configured instance (it's a `gh api`-style escape hatch) — treat it as equivalent to holding the API token directly, not a restricted subcommand. |
| Multi-user / shared-machine credential isolation | The config file is a single flat file for a single user; there is no OS keychain integration or per-invocation scoping. |
| Supply-chain integrity of the generated API client | `api/paperless.gen.go` is generated from a schema fetched from either a live Paperless instance or the official `ghcr.io/paperless-ngx/paperless-ngx` image — compromise of either source would flow into the generated client on the next `make generate`. |

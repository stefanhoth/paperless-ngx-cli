# Paperless-NGX 3.x AI features — and what they mean if you run paperless-gpt

Paperless-NGX 3.0 ships its own LLM features. If you already run
[paperless-gpt](https://github.com/icereed/paperless-gpt) alongside your
instance, the obvious question is whether the two collide, and whether the
sidecar is now redundant. Short answer: **no collision, and it doesn't replace
paperless-gpt** — the built-in features are opt-in and suggestion-only.

Everything below was read off the `ghcr.io/paperless-ngx/paperless-ngx:3.0.3`
image and the vendored [schema](../schema/paperless.json) this CLI generates
from.

## What 3.x actually adds

**It is off by default.** `PAPERLESS_AI_ENABLED` defaults to `NO`. It can also
be switched on in the app config (Settings UI) rather than the environment —
both feed the same `AIConfig`. Upgrading a 2.x instance to 3.x therefore
changes nothing about how documents are processed until you turn it on.

**Two new endpoints:**

| Endpoint | Method | What it does |
|---|---|---|
| `/api/documents/{id}/ai_suggestions/` | `GET` | LLM suggestions for one document: `title`, plus `tags`/`correspondents`/`document_types`/`storage_paths` as existing IDs, `suggested_*` as free-text for things that don't exist yet, and `dates` |
| `/api/documents/chat/` | `POST` | RAG chat over your documents (`q`, optional `document_id`). Answers stream back as `text/event-stream` |

Both return HTTP 400 (`"AI is required for this feature"`) while AI is
disabled. `ai_suggestions` additionally requires `change_document` permission
on the document, chat requires `view_document`.

**Backends:** `PAPERLESS_AI_LLM_BACKEND` is `ollama` or `openai-like`;
embeddings (`PAPERLESS_AI_LLM_EMBEDDING_BACKEND`) can be `huggingface`,
`ollama` or `openai-like`. So the same local Ollama that serves paperless-gpt
can serve Paperless. `PAPERLESS_AI_LLM_ALLOW_INTERNAL_ENDPOINTS` defaults to
true, which is what makes a LAN-local Ollama reachable.

**Chat needs an index.** The RAG index (`llmindex_index` task, stored in
`data/llm_index`) is only built when AI is enabled *and* an embedding backend
is configured. Per-document `ai_suggestions` does not depend on it.

**The classic suggestions endpoint is untouched.**
`/api/documents/{id}/suggestions/` still comes from the trained classifier, not
from an LLM. Enabling AI adds a second, separate endpoint — it does not swap
out the old one.

## Why it can't fight with paperless-gpt

`ai_suggestions` is a `GET`. It computes proposals and returns them; applying
them is a click in the web UI. There is no workflow action, no scheduled job
and no bulk operation that writes LLM output to documents on its own.

paperless-gpt, by contrast, is the thing that *writes*: it picks up documents
carrying `MANUAL_TAG` / `AUTO_TAG` (`paperless-gpt` / `paperless-gpt-auto` by
default) and patches title, tags, correspondent, dates and custom fields back
through the API. That division stays intact in 3.x — the two never write the
same field at the same time, because only one of them writes at all.

## What overlaps, and what doesn't

Built-in 3.x covers, on a per-document, manual basis: title, tag,
correspondent, document type, storage path and date suggestions — plus chat,
which paperless-gpt does not do.

No built-in counterpart exists for the parts of paperless-gpt that are usually
the reason people run it:

- **LLM-enhanced OCR** and the alternative OCR providers (Google Document AI,
  Azure Document Intelligence, Docling)
- **Searchable PDF text layers** generated from that OCR
- **Tag-triggered batch and auto-processing** instead of one document at a time
- **Custom field extraction** with append/update/replace write modes
- **Customizable prompts** per suggestion type

Note that 3.x does add a *separate*, non-AI remote OCR path
(`PAPERLESS_REMOTE_OCR_ENGINE`), which is not part of the AI feature set
described here.

## Running both

paperless-gpt's README states it is tested against the 2.20.x series and the
3.0.0 beta, and that it only uses the stable `documents/`, `tags/`,
`correspondents/`, `custom_fields/` and `document_types/` endpoints — none of
which had breaking changes in the v3 migration guide. Check its release notes
for a 3.0 stable confirmation before upgrading a production instance.

Two practical things if you enable both against one Ollama:

- They queue against the same model. `PAPERLESS_AI_LLM_REQUEST_TIMEOUT`
  (default 120s) and `PAPERLESS_AI_LLM_CONTEXT_SIZE` (default 8192) are the
  knobs on the Paperless side; a paperless-gpt auto-processing run and an
  index rebuild at the same time will make both feel slow.
- With the `openai-like` backend, document content leaves your network. That is
  a per-tool decision — enabling it in Paperless says nothing about how
  paperless-gpt is configured, and vice versa.

## Reaching the endpoints from this CLI

There are no dedicated commands for the AI endpoints — they're reachable
through the raw passthrough:

```bash
# Suggestions for one document (plain JSON)
paperless api /documents/1234/ai_suggestions/ | jq

# Ask about one document
paperless api /documents/chat/ --field q="What is the invoice total?" --field document_id=1234
```

Two caveats on the chat call: `--field` sends every value as a JSON string, so
`document_id` arrives as `"1234"` and is coerced server-side; and the
passthrough reads the whole response before printing, so the `text/event-stream`
answer arrives in one go at the end rather than token by token.

While AI is disabled server-side, both return HTTP 400 and the CLI exits
non-zero.

# Changelog

## 1.4.0
- Add SQLite-backed conversation mapping (`/data/gundam.db`) using pure-Go `modernc.org/sqlite`.
- Map external deterministic conversation IDs (from Discord threads or API) to internal Antigravity CLI session IDs.
- Resume existing sessions via `agy --conversation <internal_id>` on follow-up turns to maintain full multi-turn conversational context.
- Update `/api/transcripts` to surface mapped `external_id` alongside transcript logs.

## 1.3.9
- Pass `--print-timeout` to `agy` based on configured `timeout_minutes` to prevent premature 5-minute print-mode aborts.
- Persist brain session transcripts in `/data/brain` across container restarts.
- Add `GET /api/transcripts` endpoint to inspect full JSONL steps, statuses, and error objects.
- Enhance `extractResponseAndError` with broad multi-root directory scanning.

## 1.3.8
- Fix agy conversation flag name to `--conversation` (was `--conversation-id`).

## 1.3.7
- Committed `GEMINI.md` directly into the repository with full system architecture, component topology, invariants, and operational workflows.
- Copied `GEMINI.md` to `/app/GEMINI.md` in Dockerfile so Antigravity CLI natively discovers it on every session.
- Simplified `main.go` to avoid hardcoding markdown rule strings.

## 1.3.6
- Support `conversation_id` in `POST /api/prompt` payload; if not provided, generates a UUIDv4 using `github.com/google/uuid`.
- Pass `--conversation-id` explicitly to `agy` CLI for deterministic conversation scoping and session tracking.
- Directly inspect conversation transcript (`transcript_full.jsonl`) by known `conversation_id` and extract structured error payloads when execution fails.

## 1.3.5
- Removed redundant per-prompt disk file writes in HTTP handler; configuration files are written once on server startup.
- Streamlined `agy` execution logging directly via stdout/stderr.

## 1.3.4
- Added configurable `timeout_minutes` option with a default of 15 minutes for long-running agent tasks and multi-turn MCP operations.

## 1.3.3
- Added `AGENTS.md` repository guidelines defining zero in-image MCP invariant, remote endpoint architecture, and secret isolation.

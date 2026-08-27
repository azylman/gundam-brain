# Changelog

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

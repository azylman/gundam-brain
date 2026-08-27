# Changelog

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

# Changelog

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
- Simplified `main.go` process execution to standard `context.WithTimeout` without `setpgid`.

## 1.3.2
- Removed local MCP package preinstallations from Dockerfile; all MCP servers now operate as remote network endpoints.

## 1.3.1
- Added execution context timeout (2 minutes) and process group cleanup to prevent background MCP child processes from blocking parent CLI termination.

## 1.3.0
- Added built-in tool-first agent system rules (`~/.gemini/rules/agent_core.md` and `/app/GEMINI.md`) ensuring proactive tool invocation over conversational acknowledgement.
- Added `system_prompt` add-on configuration option for custom behavioral guidelines.

## 1.2.2
- Preinstall `@modelcontextprotocol/server-github`, `@pasympa/discord-mcp`, and `@iqai/mcp-discord` in container runtime for instant MCP initialization.

## 1.2.1
- Improved response extractor and log diagnostics.

## 1.2.0
- Add `mcp_config` option to support custom MCP servers.
- Install Node.js, npm, Python, and base tools in runtime image.

## 1.1.0
- Support `api_key` in add-on options to configure Gemini API authentication directly in Home Assistant without hardcoded secrets.

## 1.0.0
- Initial release of Gundam Brain Home Assistant add-on wrapping headless Antigravity CLI.

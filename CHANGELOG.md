# Changelog

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

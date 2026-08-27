# Agent Instructions & Architecture Guidelines: Gundam Brain

## 1. Zero In-Image MCP Installations (CRITICAL INVARIANT)
- **NEVER install MCP servers or packages inside the Dockerfile/image** (no `npm install -g @modelcontextprotocol/...`, `@pasympa/...`, pip packages, etc.).
- Gundam Brain is strictly an agent execution runner. It connects to **remote network MCP endpoints only**.
- All outbound MCP servers must be defined in `mcp_config` as remote URLs:
  - **Discord**: Remote Home Assistant add-on (`http://192.168.1.14:8080/mcp`)
  - **GitHub**: Hosted remote endpoint (`https://api.githubcopilot.com/mcp/`)
  - **Home Assistant**: Remote webhook endpoint (`ha-mcp`)

## 2. Secrets & Credential Isolation
- **NEVER commit tokens, API keys, or secrets into git**.
- Default values in `config.yaml` must always remain empty strings (`""`).
- All credentials (Gemini API keys, Discord tokens, GitHub PATs, HA webhooks) must be configured exclusively in Home Assistant's local add-on options (`/data/options.json`).

## 3. Architecture Separation
- **Inbound Events**: Handled strictly by `discord_funnel` (forwards messages via `POST /api/prompt`).
- **Outbound Actions**: Handled strictly by remote MCP servers over HTTP/SSE.
- **Process Management**: Use standard `context.WithTimeout` for bounded execution. No local stdio child process management is required since all tools run remotely over the network.

## 4. Tool-First Execution
- The agent must always prioritize tool inspection and concrete execution over conversational acknowledgement.

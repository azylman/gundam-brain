# Gundam Brain

Home Assistant add-on running a Go server to execute headless Antigravity CLI (`agy`) agent prompts.

## Features
- **Prompt Execution API**: Exposes `POST /api/prompt` for background goroutine execution.
- **Model Context Protocol (MCP) Support**: Supply standard `mcp_config` JSON to enable custom tools (e.g. Discord bot, Home Assistant, SQLite, etc.).
- **Auto-Authentication**: Provide `api_key` in Home Assistant options to automatically configure Gemini API key auth.
- **Pre-installed Environments**: Node.js, npm, Python 3, pip, git, and Antigravity CLI.

## Configuration Options

| Option | Type | Description |
| --- | --- | --- |
| `port` | int | Container port (default `8080`). |
| `agy_bin` | string | Antigravity CLI binary name or path (default `agy`). |
| `api_key` | string (password) | Gemini API Key for model inference. |
| `model` | string | Model identifier (default `Gemini 3.7 Flash (High)`). |
| `mcp_config` | string | JSON string containing MCP server definitions (`mcpServers`). |

### Example `mcp_config` for Discord

```json
{
  "mcpServers": {
    "discord": {
      "command": "python3",
      "args": ["-c", "..."],
      "env": {
        "DISCORD_BOT_TOKEN": "your-token"
      }
    }
  }
}
```

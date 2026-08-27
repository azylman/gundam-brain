# Gundam Brain Add-on Documentation

## How to Configure MCP Servers

In Home Assistant, navigate to **Settings > Add-ons > Gundam Brain > Configuration**.

In the **mcp_config** field, paste your MCP JSON configuration:

```json
{
  "mcpServers": {
    "my_mcp_server": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-everything"]
    }
  }
}
```

The server automatically creates `/root/.gemini/config/mcp_config.json` inside the container. When `agy` executes prompts, all configured MCP tools are auto-discovered and available to the agent.

# Home Assistant Add-on: Gundam Brain

Gundam Brain is a Home Assistant add-on that runs a Go server exposing a REST API endpoint to execute prompts through a headless Antigravity CLI (`agy`) instance in the background.

## Features

- **Asynchronous Execution:** Accepts requests and spawns a background goroutine immediately (`202 Accepted`).
- **Autonomous Tool Execution:** Invokes `agy` with `--dangerously-skip-permissions` to allow headless non-interactive tool operations without blocking.
- **Robust Process Handling:** Safely captures exit code, standard output, and standard error, preventing panics even if the CLI fails to spawn.
- **Full Output Logging:** Logs the exit code, full `stdout`, and `stderr` to the add-on logs.

## Installation

1. Navigate to **Settings** -> **Add-ons** -> **Add-on Store** in your Home Assistant UI.
2. Click the three dots (top right) -> **Repositories**.
3. Add repository URL: `https://github.com/azylman/gundam-brain`.
4. Find **Gundam Brain** in the list and click **Install**.
5. Start the add-on.

## Configuration

| Option | Type | Default | Description |
|---|---|---|---|
| `port` | `int` | `8080` | HTTP port for the API server |
| `agy_bin` | `string` | `agy` | Path or name of the Antigravity CLI binary |

## API Usage

### Execute Prompt

`POST /api/prompt` (or `POST /prompt` or `POST /`)

**Request Payload:**
```json
{
  "prompt": "Turn off all basement lights and summarize what was changed"
}
```

**Response (`202 Accepted`):**
```json
{
  "status": "accepted",
  "message": "Prompt execution started in background"
}
```

### Health Check

`GET /health` (or `GET /api/health`)

**Response (`200 OK`):**
```json
{
  "status": "healthy"
}
```

## Home Assistant REST Command Example

Add to your Home Assistant `configuration.yaml`:

```yaml
rest_command:
  gundam_brain_prompt:
    url: "http://gundam-brain:8080/api/prompt"
    method: "POST"
    headers:
      content-type: "application/json"
    payload: '{"prompt": "{{ prompt }}"}'
```

Then call it in automations or scripts:

```yaml
action: rest_command.gundam_brain_prompt
data:
  prompt: "Check energy usage and recommend automations"
```

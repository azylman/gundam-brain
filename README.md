# Gundam Brain

Home Assistant add-on and Go server for triggering headless Antigravity CLI (`agy`) tasks.

## Overview

Gundam Brain exposes a lightweight HTTP API designed for Home Assistant integration. When an endpoint receives a prompt, it immediately acknowledges the request and executes the prompt in a background goroutine using the Antigravity CLI (`agy --dangerously-skip-permissions -p "<prompt>"`).

Once execution finishes, it logs:
- Process exit code
- Standard output (`stdout`)
- Standard error (`stderr`)

## Repository Structure

```
.
├── Dockerfile          # Multi-stage Docker build for Home Assistant
├── build.yaml          # Home Assistant build architecture mapping
├── config.yaml         # Home Assistant add-on manifest
├── repository.yaml     # Home Assistant add-on repository manifest
├── DOCS.md             # Detailed add-on documentation and setup
├── CHANGELOG.md        # Add-on version history
├── go.mod              # Go module definition
└── main.go             # Single-file Go HTTP server & background runner
```

## Running Locally

```bash
go run main.go
```

Send a test request:

```bash
curl -X POST http://localhost:8080/api/prompt \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Say hello from Gundam Brain"}'
```

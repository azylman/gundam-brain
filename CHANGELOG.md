# Changelog

## 1.1.0
- Switch base image to Ubuntu.
- Install Antigravity CLI (`agy`), `curl`, `ca-certificates`, `git`, and prerequisites.
- Configure PATH and default port mapping to 8088.

## 1.0.0
- Initial release of Gundam Brain Home Assistant add-on.
- Single-file Go server exposing `/api/prompt`, `/prompt`, and `/`.
- Spawns background goroutine executing `agy --dangerously-skip-permissions -p "<prompt>"`.
- Robust capture and logging of exit code, `stdout`, and `stderr`.
- Health check endpoint `/health`.

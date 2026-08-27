# Changelog

## 1.1.3
- Auto-provision `modelProvider: gemini` in `~/.gemini/antigravity-cli/settings.json` when `api_key` is provided.
- Dynamic runtime injection of `GEMINI_API_KEY`, `ANTIGRAVITY_API_KEY`, and `GOOGLE_API_KEY`.

## 1.1.2
- Support optional `api_key` in Home Assistant add-on options.
- Provide non-blocking empty stdin to background `agy` process to prevent terminal hangs.

## 1.1.1
- Explicit Ubuntu 24.04 runtime base.
- Direct Antigravity CLI (`agy`) binary installation via official bootstrapper script into `/usr/local/bin`.

## 1.1.0
- Switch base image to Ubuntu.
- Install Antigravity CLI (`agy`), `curl`, `ca-certificates`, `git`, and prerequisites.

## 1.0.0
- Initial release of Gundam Brain Home Assistant add-on.

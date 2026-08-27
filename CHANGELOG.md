# Changelog

## 1.1.8
- Added transcript fallback parser to capture model responses even when `agy` non-TTY TUI output fails to flush to stdout.

## 1.1.7
- Use exact model identifier `Gemini 3.7 Flash (High)` for Antigravity CLI.

## 1.1.6
- Add `model` option passed to `--model` flag and `settings.json`.
- Set only `GEMINI_API_KEY` to avoid multi-key collision warnings.

## 1.1.5
- Capture and log `--log-file /tmp/agy.log` for troubleshooting agent executions.

## 1.1.4
- Default model to `gemini-2.5-flash` when configuring Gemini API provider in `settings.json`.

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

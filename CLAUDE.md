# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Smart Suggestion is an AI-powered command line suggestion tool for zsh. It's a fork of zsh-copilot, rewritten in Go, with the key differentiator that it captures terminal output (stdout/stderr) for context-aware intelligent predictions, not just shell history.

**Architecture**: Go CLI binary (`cmd/smart-suggestion/main.go`) + zsh plugin (`smart-suggestion.plugin.zsh`)

## Development Commands

### Build
```bash
./build.sh                    # Builds smart-suggestion binary in project root
go build -o smart-suggestion ./cmd/smart-suggestion/main.go  # Manual build
```

### Test
```bash
go test ./pkg/config/...      # Config package tests
go test ./pkg/privacy/...     # Privacy filter tests
go test ./pkg/...             # All package tests
```

### Config Management
```bash
smart-suggestion config init                    # Print config template to stdout
smart-suggestion config init --file /path       # Write template to file
smart-suggestion config validate                # Validate config at default path
smart-suggestion config validate --file /path   # Validate specific config
```

### Debug
```bash
export SMART_SUGGESTION_DEBUG=true    # Enable logging to /tmp/smart-suggestion.log
tail -f /tmp/smart-suggestion.log     # Monitor debug output
```

### Release
Tag and push to trigger GitHub Actions build workflow (`.github/workflows/release.yml`).

## Architecture

### High-Level Flow

1. **User presses Ctrl+O** (or configured keybinding)
2. **zsh plugin** (`_do_smart_suggestion()`) captures current BUFFER
3. **Calls Go binary**: `smart-suggestion --input "..." --context`
4. **Go binary builds context**:
   - System info (user, shell, directory, OS)
   - Shell aliases (`alias` command output)
   - Shell history (`fc -ln -N`, privacy-filtered)
   - Terminal buffer (from tmux, kitty, or proxy log, privacy-filtered)
5. **Selects provider**: from `--provider` flag or `default_provider` config
6. **Calls AI provider API** with system prompt + context + user input
7. **Parses response**: Extracts command after `</reasoning>` tag, checks prefix (`=` new cmd, `+` completion)
8. **Writes suggestion** to `/tmp/smart_suggestion`
9. **zsh plugin reads** file and displays via zsh-autosuggestions or replaces BUFFER

### Key Components

**cmd/smart-suggestion/main.go** (~2500 lines): Single-file CLI application using cobra
- Root command: Default fetch behavior (AI suggestion)
- `proxy`: Start shell proxy mode for terminal capture
- `rotate-logs`: Log rotation utility
- `version`: Show version
- `update`: Self-update from GitHub releases
- `config init/validate`: Config management

**smart-suggestion.plugin.zsh**: Zsh integration layer
- Key binding (default: Ctrl+O)
- Calls Go binary, reads result from `/tmp/smart_suggestion`
- Proxy mode management (background shell recording)
- Loading animations, auto-update checking

**pkg/config/**: Configuration management
- `config.go`: Load/Save/GetProvider config with 0600 permissions
- `validation.go`: Provider validation, URL format checks, Azure-specific rules
- `security.go`: Security utilities

**pkg/privacy/**: Privacy filtering system
- `filter.go`: Multi-level regex-based filtering (none/basic/moderate/strict)
- Filters API keys, tokens, env vars, passwords, SSH keys, etc.

**pkg/logrotate.go**: Log rotation with size-based rotation, gzip compression, retention limits

### Provider Configuration

Config file location: `~/.config/smart-suggestion/config.json` (or `SMART_SUGGESTION_PROVIDER_FILE`)

Supported providers: `openai`, `openai_compatible`, `azure_openai`, `anthropic`, `gemini`, `deepseek`

**Important**: Azure OpenAI supports either `resource_name` OR `base_url`, but not both (XOR validation in `validation.go`).

**OpenAI-Compatible URL Handling**: If `base_url` contains `/chat/completions`, use as-is; otherwise auto-append `/v1/chat/completions`.

### Proxy Mode

Default: **enabled**. Creates PTY (pseudo-terminal) using `github.com/creack/pty`, spawns nested shell, records all output to session-based logs.

**Security implication**: All terminal activity is logged. Consider privacy implications when modifying proxy mode behavior.

**Process locking**: Prevents duplicate proxy sessions via file-based locking.

### Privacy Filtering

Levels: 0 (none), 1 (basic - default), 2 (moderate), 3 (strict)

Applied to: shell history and terminal buffer before sending to AI providers.

Patterns: API keys, env vars with KEY/TOKEN/SECRET/PASSWORD, JWT, database URLs, bearer tokens, SSH keys, and many service-specific patterns.

## Configuration Variables

Key zsh environment variables (set in `.zshrc` or before plugin load):

- `SMART_SUGGESTION_PROVIDER_FILE`: Config file path (default: `~/.config/smart-suggestion/config.json`)
- `SMART_SUGGESTION_AI_PROVIDER`: Override provider selection
- `SMART_SUGGESTION_KEY`: Keybinding (default: `^o` for Ctrl+O)
- `SMART_SUGGESTION_SEND_CONTEXT`: Include context (default: `true`)
- `SMART_SUGGESTION_PROXY_MODE`: Terminal recording (default: `true`)
- `SMART_SUGGESTION_PRIVACY_FILTER`: Enable filtering (default: `true`)
- `SMART_SUGGESTION_PRIVACY_LEVEL`: `none`/`basic`/`moderate`/`strict`
- `SMART_SUGGESTION_DEBUG`: Debug logging to `/tmp/smart-suggestion.log`

## Adding a New AI Provider

1. Add provider config struct to `pkg/config/config.go`
2. Add validation in `pkg/config/validation.go`:
   - Add to `isValidProvider()` list
   - Add provider-specific validation function
   - Call validator in `Validate()` method
   - Add `ValidateProviderAvailable()` case
3. Implement fetch function in `cmd/smart-suggestion/main.go`:
   - Load config for provider
   - Construct provider-specific HTTP request
   - Make API call with 30s timeout
   - Parse response and extract content
   - Return raw response string
4. Add provider case in provider selection switch
5. Update `defaultSystemPrompt` instructions if needed
6. Add documentation to README.md

## Testing Changes

1. Build: `./build.sh`
2. Reload plugin: `source ~/.zshrc`
3. Test suggestion: Type partial command, press Ctrl+O
4. Check debug log: `tail -f /tmp/smart-suggestion.log` (with `SMART_SUGGESTION_DEBUG=true`)
5. Validate config: `smart-suggestion config validate`

## File Structure Notes

- Main CLI is intentionally a single large file (~2500 lines) - this is by design
- Zsh plugin contains shell integration logic, not just simple sourcing
- Proxy logs use session IDs for isolation
- Config files created with 0600 permissions (owner read/write only)

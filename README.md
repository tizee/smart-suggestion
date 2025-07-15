# Smart Suggestion for Zsh

> [!NOTE]
>
> This project is a fork of [zsh-copilot](https://github.com/Myzel394/zsh-copilot) by [Myzel394](https://github.com/Myzel394).

Get AI-powered command suggestions **directly** in your zsh shell. No complex setup, no external tools - just press `CTRL + O` and get intelligent command suggestions powered by OpenAI, Anthropic Claude, Google Gemini, DeepSeek, or any OpenAI-compatible API.

> [!NOTE]
>
> This project is still in its early stages, and some features may be immature and unstable. I appreciate your understanding.

<https://github.com/user-attachments/assets/90eaf300-c891-4ef9-958d-9890f53f9929>

<https://github.com/user-attachments/assets/0e126456-4e52-4693-a8a8-d3bdb8a133e0>

## Features

- **🚀 Context-aware intelligent prediction**: Predicts the next command you are likely to input based on context (history, aliases, terminal buffer)
- **🤖 Multiple AI Providers**: Support for OpenAI GPT, Anthropic Claude, Google Gemini, DeepSeek, and any OpenAI-compatible API
- **🔒 Privacy Protection**: Built-in privacy filtering to prevent sensitive data (API keys, passwords, tokens) from being sent to AI providers
- **🔧 Highly Configurable**: Customize keybindings, AI provider, context sharing, privacy filtering, and more

## Questions

- Why don't I use [zsh-copilot](https://github.com/Myzel394/zsh-copilot) and instead fork a separate version?

  Because the context of zsh-copilot only includes history commands and does not include the terminal buffer (i.e., the stdout/stderr of history commands), it cannot achieve the context-aware intelligent prediction I want, this is the feature I want the most, and it's also the main reason why I forked. Additionally, since zsh-copilot is written in shell, it's very difficult to concatenate JSON and implement stdio interception. Therefore, I re-implemented almost all logic using Go, which made it too different from the original project to merge back.

## Installation

### Prerequisites

Make sure you have the following installed:

- **zsh** shell
- **[zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions)** plugin
- An API key for one of the supported AI providers

### Method 1: Quick Install (Recommended)

The easiest way to install smart-suggestion is using our installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/yetone/smart-suggestion/main/install.sh | bash
```

This script will:

- Detect your platform (Linux, macOS, Windows)
- Download the appropriate pre-built binary
- Install the plugin to `~/.config/smart-suggestion`
- Configure your `~/.zshrc` automatically with proxy mode enabled by default ⚠️ **See security considerations below**
- Check for zsh-autosuggestions dependency

**Uninstall:**

```bash
curl -fsSL https://raw.githubusercontent.com/yetone/smart-suggestion/main/install.sh | bash -s -- --uninstall
```

### Method 2: Oh My Zsh

1. Clone the repository into your Oh My Zsh custom plugins directory:

```bash
git clone https://github.com/yetone/smart-suggestion ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/smart-suggestion
```

2. Add `smart-suggestion` to your plugins array in `~/.zshrc`:

```bash
omz plugin enable smart-suggestion
```

3. Build the Go binary:

```bash
cd ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/smart-suggestion
./build.sh
```

4. Reload your shell:

```bash
source ~/.zshrc
```

### Method 3: Zinit

1. Add the following to your `~/.zshrc`:

```bash
zinit as"program" atclone'./build.sh' \
    atpull'%atclone' pick"smart-suggestion" src"smart-suggestion.plugin.zsh" for \
        yetone/smart-suggestion
```

2. Update Zinit:

```bash
zi update
```

### Method 4: Manual Installation from Source

1. Clone the repository:

```bash
git clone https://github.com/yetone/smart-suggestion ~/.config/smart-suggestion
```

2. Build the Go binary (requires Go 1.21+):

```bash
cd ~/.config/smart-suggestion
./build.sh
```

3. Add to your `~/.zshrc`:

```bash
source ~/.config/smart-suggestion/smart-suggestion.plugin.zsh
```

4. Reload your shell:

```bash
source ~/.zshrc
```

### Method 5: Manual Installation from Release

1. Download the latest release for your platform from [GitHub Releases](https://github.com/yetone/smart-suggestion/releases)

2. Extract the archive:

```bash
mkdir -p ~/.config/smart-suggestion
tar -xzf smart-suggestion-*.tar.gz -C ~/.config/smart-suggestion --strip-components=1
```

3. Add to your `~/.zshrc`:

```bash
source ~/.config/smart-suggestion/smart-suggestion.plugin.zsh
```

4. Reload your shell:

```bash
source ~/.zshrc
```

## Configuration

### AI Provider Setup

**Important**: Starting with the new version, API keys are now configured via JSON configuration files instead of environment variables for better security. The `SMART_SUGGESTION_PROVIDER_FILE` environment variable should point to your configuration file.

#### Quick Setup

1. **Generate a configuration template**:
```bash
smart-suggestion config init > ~/.config/smart-suggestion/config.json
```

2. **Edit the configuration file** to add your API keys:
```bash
nano ~/.config/smart-suggestion/config.json
```

3. **Set the configuration file path** (if not using the default):
```bash
export SMART_SUGGESTION_PROVIDER_FILE="~/.config/smart-suggestion/config.json"
```

#### Configuration Examples

**OpenAI (default)**:
```json
{
  "openai": {
    "api_key": "your-openai-api-key",
    "base_url": "https://api.openai.com",
    "model": "gpt-4o-mini"
  },
  "default_provider": "openai"
}
```

**OpenAI-Compatible APIs** (Ollama, OpenRouter, etc.):
```json
{
  "openai_compatible": {
    "api_key": "your-api-key",
    "base_url": "https://openrouter.ai/api",
    "model": "anthropic/claude-3.5-sonnet"
  },
  "default_provider": "openai_compatible"
}
```

**Azure OpenAI**:
```json
{
  "azure_openai": {
    "api_key": "your-azure-openai-api-key",
    "resource_name": "your-azure-openai-resource-name",
    "deployment_name": "your-deployment-name",
    "api_version": "2024-10-21",
    "base_url": "https://your-custom-base-url"
  },
  "default_provider": "azure_openai"
}
```

**Azure OpenAI Configuration Notes**:
- `api_key`: Your Azure OpenAI API key (e.g., `c0123456789012345678901234567890`)
- `resource_name`: Your Azure OpenAI resource name (e.g., `awesome-corp` when your endpoint is `https://awesome-corp.openai.azure.com`)
- `deployment_name`: Your deployment name (e.g., `gpt-4o`)
- `api_version`: Optional, defaults to `2024-10-21`
- `base_url`: Optional, defaults to `https://{resource_name}.openai.azure.com`. Use this for custom endpoints.

**Important**: You can use either `resource_name` OR `base_url`, but not both. If you provide `base_url`, `resource_name` will be ignored.

**Anthropic Claude**:
```json
{
  "anthropic": {
    "api_key": "your-anthropic-api-key",
    "base_url": "https://api.anthropic.com",
    "model": "claude-3-5-sonnet-20241022"
  },
  "default_provider": "anthropic"
}
```

**Google Gemini**:
```json
{
  "gemini": {
    "api_key": "your-gemini-api-key",
    "base_url": "https://generativelanguage.googleapis.com",
    "model": "gemini-2.5-flash"
  },
  "default_provider": "gemini"
}
```

**DeepSeek**:
```json
{
  "deepseek": {
    "api_key": "your-deepseek-api-key",
    "base_url": "https://api.deepseek.com",
    "model": "deepseek-chat"
  },
  "default_provider": "deepseek"
}
```

#### Privacy Configuration

**Privacy Filtering** (Enabled by default): Smart Suggestion includes built-in privacy protection to prevent sensitive information from being sent to AI providers.

```json
{
  "privacy_filter": {
    "enabled": true,
    "level": 1,
    "replacement_text": "[REDACTED]"
  }
}
```

**Privacy Filter Levels**:
- `0` (`none`): No filtering (not recommended)
- `1` (`basic`): Filter common API keys, tokens, and secrets (default)
- `2` (`moderate`): Include emails, IPs, and advanced patterns
- `3` (`strict`): Aggressive filtering including potential secrets

**What gets filtered**:
- **API Keys**: OpenAI (sk-*, pk-*), AWS, GitHub, Slack tokens
- **Environment Variables**: `export API_KEY=secret` patterns and **any variable containing KEY/TOKEN/SECRET/PASSWORD**
- **Bearer Tokens**: Authorization headers and tokens
- **Database URLs**: Connection strings with credentials
- **Command Passwords**: `curl -u`, `--password` flags
- **JWT Tokens**: JSON Web Tokens
- **SSH Keys**: Private key markers
- **Service-Specific Keys**: Stripe, Twilio, SendGrid, Mailgun, etc.
- **Cloud Provider Tokens**: DigitalOcean, Vultr, Linode
- **CI/CD Tokens**: GitLab, Jenkins, CI systems
- **Echo Commands**: `echo $API_KEY` and similar sensitive variable reveals
- **Command Output**: Standalone secret values that appear to be API keys or tokens

#### Popular OpenAI-Compatible Services

The `openai_compatible` provider supports many third-party services:

- **OpenRouter**: Access multiple AI models through one API
  ```json
  {
    "openai_compatible": {
      "api_key": "your-openrouter-key",
      "base_url": "https://openrouter.ai/api",
      "model": "anthropic/claude-3.5-sonnet"
    }
  }
  ```

- **Ollama**: Local AI models
  ```json
  {
    "openai_compatible": {
      "base_url": "http://localhost:11434",
      "model": "llama3.2:latest"
    }
  }
  ```

- **Other providers**: LocalAI, vLLM, Text Generation WebUI, etc.

#### Configuration Management

**Validate your configuration**:
```bash
smart-suggestion config validate
```

**View configuration template**:
```bash
smart-suggestion config init
```

**Create configuration file directly**:
```bash
smart-suggestion config init --file ~/.config/smart-suggestion/config.json
```

### Environment Variables

Configure the plugin behavior with these environment variables:

| Variable                           | Description                           | Default       | Options                                                     |
|------------------------------------|---------------------------------------|---------------|-------------------------------------------------------------|
| `SMART_SUGGESTION_PROVIDER_FILE`   | Path to configuration file            | `~/.config/smart-suggestion/config.json` | Any valid JSON file path |
| `SMART_SUGGESTION_AI_PROVIDER`     | AI provider to use                    | `openai` | `openai`, `openai_compatible`, `azure_openai`, `anthropic`, `gemini`, `deepseek` |
| `SMART_SUGGESTION_KEY`             | Keybinding to trigger suggestions     | `^o`          | Any zsh keybinding                                          |
| `SMART_SUGGESTION_SEND_CONTEXT`    | Send shell context to AI ⚠️ **Privacy Risk** | `true`        | `true`, `false`                                             |
| `SMART_SUGGESTION_PRIVACY_FILTER`  | Enable privacy filtering of sensitive data | `true`        | `true`, `false`                                             |
| `SMART_SUGGESTION_PRIVACY_LEVEL`   | Privacy filtering sensitivity level   | `basic`       | `none`, `basic`, `moderate`, `strict`                      |
| `SMART_SUGGESTION_PROXY_MODE`      | Enable proxy mode ⚠️ **Privacy Risk, Shell Nesting** | `true`        | `true`, `false`                                             |
| `SMART_SUGGESTION_DEBUG`           | Enable debug logging                  | `false`       | `true`, `false`                                             |
| `SMART_SUGGESTION_SYSTEM_PROMPT`   | Custom system prompt                  | Built-in      | Any string                                                  |
| `SMART_SUGGESTION_AUTO_UPDATE`     | Enable automatic update checking      | `true`        | `true`, `false`                                             |
| `SMART_SUGGESTION_UPDATE_INTERVAL` | Days between update checks            | 7             | Any positive integer                                        |
| `SMART_SUGGESTION_BINARY`          | Path to the `smart_suggestion` binary | Auto-detected | Any valid filepath to a valid `smart_suggestion` binary     |

If `SMART_SUGGESTION_BINARY` is not specified, we look for one in the following locations:

1. `smart_suggestion` beside the current `smart-suggestion.plugin.zsh`
1. `~/.config/smart-suggestion/smart_suggestion`

**Note**: The configuration file path (`SMART_SUGGESTION_PROVIDER_FILE`) defaults to `~/.config/smart-suggestion/config.json` if not specified.

### Advanced Configuration

#### Multiple Providers in One Config

You can configure multiple providers in a single configuration file:

```json
{
  "openai": {
    "api_key": "your-openai-key",
    "model": "gpt-4o-mini"
  },
  "openai_compatible": {
    "api_key": "your-openrouter-key",
    "base_url": "https://openrouter.ai/api",
    "model": "anthropic/claude-3.5-sonnet"
  },
  "anthropic": {
    "api_key": "your-anthropic-key",
    "model": "claude-3-5-sonnet-20241022"
  },
  "default_provider": "openai_compatible"
}
```

#### Custom API Endpoints

For `openai_compatible` provider, you can use any OpenAI-compatible API:

```json
{
  "openai_compatible": {
    "api_key": "your-api-key",
    "base_url": "https://your-custom-endpoint.com",
    "model": "your-model-name"
  }
}
```

**Important URL Handling Note**: 

The system automatically handles API endpoint paths:
- If your `base_url` contains `/chat/completions`, it will be used as-is
- Otherwise, `/v1/chat/completions` will be automatically appended

Examples:
- `"base_url": "https://openrouter.ai/api"` → `https://openrouter.ai/api/v1/chat/completions`
- `"base_url": "https://ark.cn-beijing.volces.com/api/v3/chat/completions"` → Used as-is
- `"base_url": "http://localhost:11434"` → `http://localhost:11434/v1/chat/completions`

This means you should **NOT** include `/v1` in your base URL for most services, as it will be added automatically.

#### History Lines for Context

Configure how many lines of shell history to include in the context via environment variable:

```bash
export SMART_SUGGESTION_HISTORY_LINES="20"  # Default: 10
```

#### Configuration File Security

- Configuration files are automatically created with `0600` permissions (readable/writable by owner only)
- Store your configuration file in a secure location like `~/.config/smart-suggestion/`
- Never commit configuration files with API keys to version control

### View Current Configuration

To see all available configurations and their current values:

```bash
smart-suggestion
```

To validate your configuration file:

```bash
smart-suggestion config validate
```

## Usage

1. **Start typing a command** or describe what you want to do
2. **Press `CTRL + O`** (or your configured key)
3. **Wait for the AI suggestion** (loading animation will show)
   - _Note: On first use, proxy mode will automatically start in the background to capture terminal context_
4. **The suggestion will appear** as:
   - An autosuggestion you can accept with `→` (for completions)
   - A completely new command that replaces your input (for new commands)

## How It Works

1. **Input Capture**: The plugin captures your current command line input
2. **Proxy Mode (Default)**: Automatically starts a background shell recording session to capture terminal output for better context
3. **Context Collection**: Gathers rich shell context including user info, directory, command history, aliases, and terminal buffer content via proxy mode
4. **AI Processing**: Sends the input and context to your configured AI provider
5. **Smart Response**: AI returns either a completion (`+`) or new command (`=`)
6. **Shell Integration**: The suggestion is displayed using zsh-autosuggestions or replaces your input

### Proxy Mode (New Default)

Smart Suggestion now automatically enables **proxy mode** by default, which provides significantly better context awareness by recording your terminal session. This mode:

- **Automatically starts** when you first use smart suggestions
- **Records terminal output** using the `script` command for maximum compatibility
- **Provides rich context** to the AI including command outputs and error messages
- **Works seamlessly** across different terminal environments

#### ⚠️ Important Security and Privacy Considerations

**Proxy mode creates a nested shell environment** which may cause:
- **Shell nesting issues**: You may need to type `exit` twice to fully close your terminal
- **Privacy risks**: All commands, outputs, and sensitive information are recorded and potentially sent to AI providers
- **Data exposure**: API keys, passwords, personal files, and system information may be logged

**What gets recorded and sent:**
- All shell commands you execute
- Command outputs and error messages
- System context (username, hostname, current directory)
- Environment variables and shell history
- Any sensitive data that appears in your terminal

#### Disabling Proxy Mode

For privacy and to avoid shell nesting issues, you can disable proxy mode:

```bash
export SMART_SUGGESTION_PROXY_MODE=false
```

You can also disable context sending entirely:

```bash
export SMART_SUGGESTION_SEND_CONTEXT=false
```

**Recommended for security-conscious users**: Disable both proxy mode and context sending to minimize data exposure while still getting AI command suggestions based on your current input.

### Privacy Protection

Smart Suggestion includes **built-in privacy filtering** enabled by default to protect your sensitive information:

#### Automatic Privacy Filtering

The tool automatically filters sensitive patterns from shell history and terminal buffer before sending to AI providers:

- **API Keys**: OpenAI, Anthropic, Google, AWS, GitHub, Slack tokens
- **Environment Variables**: `export API_KEY=secret` patterns and **any variable containing KEY/TOKEN/SECRET/PASSWORD**
- **Bearer Tokens**: Authorization headers
- **Database URLs**: Connection strings with credentials (MySQL, PostgreSQL, MongoDB, Redis)
- **Passwords**: Command line password flags
- **JWT Tokens**: JSON Web Tokens
- **SSH Keys**: Private key identifiers
- **Service-Specific Keys**: Stripe, Twilio, SendGrid, Mailgun, Azure, DeepSeek
- **Cloud Provider Tokens**: DigitalOcean, Vultr, Linode
- **CI/CD Tokens**: GitLab, Jenkins, CI systems
- **Secrets**: JWT secrets, encryption keys, session secrets
- **Echo Commands**: `echo $API_KEY` and similar sensitive variable reveals
- **Command Output**: Standalone secret values that appear to be API keys or tokens

#### Privacy Configuration

**Environment Variables**:
```bash
# Disable privacy filtering (not recommended)
export SMART_SUGGESTION_PRIVACY_FILTER=false

# Set privacy level (none, basic, moderate, strict)
export SMART_SUGGESTION_PRIVACY_LEVEL=strict
```

**Configuration File**:
```json
{
  "privacy_filter": {
    "enabled": true,
    "level": 1,
    "replacement_text": "[REDACTED]",
    "custom_patterns": ["my_custom_pattern_\\w+"]
  }
}
```

**Privacy Levels**:
- **`none`**: No filtering (⚠️ not recommended)
- **`basic`**: Filter common secrets (default, recommended)
- **`moderate`**: Include emails, IPs, advanced patterns
- **`strict`**: Aggressive filtering of potential secrets

**Note**: Even with privacy filtering enabled, sensitive information may still be logged locally in debug files and proxy logs. For maximum security, consider disabling context sending entirely.

#### Privacy Level Recommendations

**When both proxy mode and context sending are enabled**, choose your privacy level based on your environment:

**🔒 Strict (Recommended for Production)**:
```bash
export SMART_SUGGESTION_PRIVACY_LEVEL=strict
```
- **Use when**: Production servers, enterprise environments, handling sensitive data
- **Filters**: All API keys, tokens, emails, IPs, long strings (32+ chars), credit card numbers
- **Trade-off**: May filter some legitimate content (file paths, hashes) but maximizes security

**⚖️ Moderate (Balanced)**:
```bash
export SMART_SUGGESTION_PRIVACY_LEVEL=moderate
```
- **Use when**: Personal development, small teams, daily coding work
- **Filters**: Basic patterns plus emails in sensitive contexts, AWS/GitHub tokens, SSH keys
- **Trade-off**: Good privacy protection with minimal false positives

**🚨 Basic (Use with caution)**:
```bash
export SMART_SUGGESTION_PRIVACY_LEVEL=basic  # Default
```
- **Use when**: Personal/learning environments, no real sensitive data
- **Filters**: Common API key formats, bearer tokens, database URLs
- **Trade-off**: Minimal filtering, may miss non-standard sensitive patterns

#### Configuration Examples

**Development Environment**:
```bash
export SMART_SUGGESTION_PROXY_MODE=true
export SMART_SUGGESTION_SEND_CONTEXT=true
export SMART_SUGGESTION_PRIVACY_LEVEL=moderate
```

**Production/Sensitive Environment**:
```bash
export SMART_SUGGESTION_PROXY_MODE=true
export SMART_SUGGESTION_SEND_CONTEXT=true
export SMART_SUGGESTION_PRIVACY_LEVEL=strict
```

**Privacy-First Approach**:
```bash
export SMART_SUGGESTION_PROXY_MODE=false      # No terminal recording
export SMART_SUGGESTION_SEND_CONTEXT=true     # Limited context only
export SMART_SUGGESTION_PRIVACY_LEVEL=basic   # Basic filtering sufficient
```

#### Monitoring Privacy Filtering

Enable debug mode to verify filtering effectiveness:
```bash
export SMART_SUGGESTION_DEBUG=true
tail -f /tmp/smart-suggestion.log
```

If legitimate content is being filtered, adjust the privacy level or add custom patterns to your configuration file.

## Troubleshooting

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
export SMART_SUGGESTION_DEBUG=true
```

Debug logs are written to `/tmp/smart-suggestion.log`.

### Common Issues

1. **"Binary not found" error**: Run `./build.sh` in the plugin directory
2. **No suggestions**: Check your API key and internet connection
3. **Wrong suggestions**: Try adjusting the context settings or system prompt
4. **Key binding conflicts**: Change `SMART_SUGGESTION_KEY` to a different key
5. **Need to type `exit` twice to close terminal**: This is caused by proxy mode creating a nested shell. Disable proxy mode with `export SMART_SUGGESTION_PROXY_MODE=false`
6. **Privacy concerns**: Proxy mode records all terminal activity. Enable privacy filtering with `export SMART_SUGGESTION_PRIVACY_LEVEL=strict` or disable context entirely with `export SMART_SUGGESTION_SEND_CONTEXT=false`
7. **Sensitive data being filtered**: If legitimate content is being filtered, adjust privacy level with `export SMART_SUGGESTION_PRIVACY_LEVEL=basic` or add custom patterns to config

### Build Issues

If the build fails:

```bash
# Check Go installation
go version

# Clean and rebuild
rm -f smart-suggestion
./build.sh
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is open source. Please check the repository for license details.

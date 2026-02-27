# ddg-search

[![CI](https://github.com/Djarvur/ddg-search/actions/workflows/ci.yml/badge.svg)](https://github.com/Djarvur/ddg-search/actions/workflows/ci.yml)
[![Security](https://github.com/Djarvur/ddg-search/actions/workflows/security.yml/badge.svg)](https://github.com/Djarvur/ddg-search/actions/workflows/security.yml)
[![Coveralls](https://coveralls.io/repos/github/Djarvur/ddg-search/badge.svg)](https://coveralls.io/github/Djarvur/ddg-search)
[![Go Report Card](https://goreportcard.com/badge/github.com/Djarvur/ddg-search)](https://goreportcard.com/report/github.com/Djarvur/ddg-search)
[![GoDoc](https://godoc.org/github.com/Djarvur/ddg-search?status.svg)](https://godoc.org/github.com/Djarvur/ddg-search)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Command-line tools for web search and content fetching.

## Tools

### ddg-search - DuckDuckGo Search Client

Search DuckDuckGo from command line (no API key required).

### page-dump - URL to Markdown Converter

Fetch web pages and convert to markdown format.

### ddg-search-mcp - MCP Server

Exposes search and fetch tools via Model Context Protocol (MCP) for use with Claude Code and other MCP-compatible clients.

### perplexity-search - Perplexity API Search Client

Search web using Perplexity API with AI-powered results and citations (requires API key).

## Features

### ddg-search-mcp

- **MCP Protocol Support**: Compatible with Claude Code and other MCP clients
- **Dual Transport**: Supports both stdio and HTTP (Streamable) transports
- **TLS/mTLS**: Secure connections with optional mutual TLS authentication
- **Multiple Tools**:
  - `search`: Web search via DuckDuckGo or Perplexity API
  - `fetch`: URL fetching with HTML-to-markdown conversion
- **Perplexity Integration**: AI-powered search with automatic fallback to DuckDuckGo
- **Flexible Configuration**: Config file, environment variables, and CLI flags
- **Hot Reload**: Configuration reload via SIGHUP signal
- **Health Checks**: Built-in health check endpoint for HTTP transport

### ddg-search

- Markdown output for LLM consumption
- Automatic rate limit detection and retry with exponential backoff
- Configurable retry behavior
- Site-specific, regional, and time-bounded searches

### page-dump

- Fetch any HTTP/HTTPS URL and convert to markdown
- Preserves document structure (headings, links, lists, code blocks)
- Configurable timeout and user-agent
- Clean markdown output to stdout

### perplexity-search

- AI-powered search results with better understanding
- Citations to sources referenced in results
- Configurable model selection (sonar-small-online, sonar-medium-online, sonar-pro-online)
- Automatic retry with exponential backoff for transient errors
- Markdown output for LLM consumption

## Installation

```bash
go install github.com/Djarvur/ddg-search@latest
```

Or build from source:

```bash
git clone https://github.com/Djarvur/ddg-search
cd ddg-search
make build
```

## Usage

### ddg-search-mcp

The MCP server can be configured via config file, environment variables, or CLI flags.

#### Configuration File

Create a config file at `~/.config/ddg-search/config.yaml`:

```yaml
# Server configuration
server:
  # Transport protocol: "stdio" or "http" (default: "stdio")
  protocol: stdio
  # Bind address for HTTP transport (default: "localhost:9100")
  bind_address: localhost:9100
  # TLS configuration
  tls:
    # Enable TLS (default: false)
    enabled: false
    # Path to TLS certificate file
    cert_file: /path/to/cert.pem
    # Path to TLS key file
    key_file: /path/to/key.pem
    # Minimum TLS version (default: "1.2")
    min_version: "1.2"
    # mTLS configuration
    mtls:
      # Enable mTLS (default: false)
      enabled: false
      # Path to CA certificate for client validation
      ca_file: /path/to/ca.pem

# Logging configuration
logging:
  # Log level: "debug", "info", "warn", "error" (default: "info")
  level: info

# Search tool configuration
search:
  # Maximum number of results to return (default: 10)
  max_results: 10
  # Enable safe search (default: false)
  safe_search: false

# Perplexity search configuration
perplexity:
  # Enable Perplexity search (default: false)
  enabled: false
  # Perplexity API access token
  access_token: ""
```

#### Environment Variables

All configuration values can be overridden via environment variables with the `DDG_SEARCH_` prefix:

- `DDG_SEARCH_SERVER_PROTOCOL`
- `DDG_SEARCH_SERVER_BIND_ADDRESS`
- `DDG_SEARCH_SERVER_TLS_ENABLED`
- `DDG_SEARCH_SERVER_TLS_CERT_FILE`
- `DDG_SEARCH_SERVER_TLS_KEY_FILE`
- `DDG_SEARCH_SERVER_TLS_MIN_VERSION`
- `DDG_SEARCH_SERVER_TLS_MTLS_ENABLED`
- `DDG_SEARCH_SERVER_TLS_MTLS_CA_FILE`
- `DDG_SEARCH_LOGGING_LEVEL`
- `DDG_SEARCH_SEARCH_MAX_RESULTS`
- `DDG_SEARCH_SEARCH_SAFE_SEARCH`
- `DDG_SEARCH_PERPLEXITY_ENABLED`
- `DDG_SEARCH_PERPLEXITY_ACCESS_TOKEN`

#### CLI Flags

```bash
# Start with default configuration
ddg-search-mcp

# Start with custom log level
ddg-search-mcp --log-level debug

# Start with custom config file
ddg-search-mcp --config /path/to/config.yaml
```

#### Running as a Service on macOS

To run the MCP server as a background service on macOS:

**Option 1: Using launchd (Recommended)**

Create a launch agent plist file at `~/Library/LaunchAgents/com.ddgsearch.mcp.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.ddgsearch.mcp</string>
    <key>ProgramArguments</key>
    <array>
        <string>/path/to/ddg-search-mcp</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/ddg-search-mcp.stdout.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/ddg-search-mcp.stderr.log</string>
</dict>
</plist>
```

Load the service:
```bash
launchctl load ~/Library/LaunchAgents/com.ddgsearch.mcp.plist
```

Unload the service:
```bash
launchctl unload ~/Library/LaunchAgents/com.ddgsearch.mcp.plist
```

**Option 2: Using nohup (Simple)**

```bash
nohup ddg-search-mcp > /tmp/ddg-search-mcp.log 2>&1 &
```

**Viewing Logs**

For launchd service:
```bash
# View stdout log
tail -f /tmp/ddg-search-mcp.stdout.log

# View stderr log
tail -f /tmp/ddg-search-mcp.stderr.log
```

For nohup:
```bash
tail -f /tmp/ddg-search-mcp.log
```

**Check if Service is Running**

```bash
# Check if the process is running
ps aux | grep ddg-search-mcp

# Check if the port is listening
lsof -i :9100

# Test health check endpoint
curl http://localhost:9100/health
```

#### Using with AI Code Editors

**Claude Code**

Add the MCP server to your Claude Code configuration:

**Project-specific configuration** (`.claude/config.json` in your project):
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": []
    }
  }
}
```

**Global configuration** (`~/.claude/config.json`):
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": []
    }
  }
}
```

**HTTP mode configuration** (if running as a service):
```json
{
  "mcpServers": {
    "ddg-search": {
      "type": "streamable-http",
      "url": "http://localhost:9100/mcp"
    }
  }
}
```

**RooCode**

RooCode uses a different MCP configuration format. Add to your RooCode configuration file at `.roo/mcp.json` in your project:

**Stdio mode**:
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": []
    }
  }
}
```

**HTTP mode**:
```json
{
  "mcpServers": {
    "ddg-search": {
      "type": "streamable-http",
      "url": "http://localhost:9100/mcp"
    }
  }
}
```

**KiloCode**

KiloCode configuration can be managed at two levels:

**Global configuration** (`~/.kilocode/mcp_settings.json`):
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": []
    }
  }
}
```

**Project-level configuration** (`.kilocode/mcp.json` in your project root):
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": []
    }
  }
}
```

**OpenCode**

OpenCode configuration is in `opencode.json` (or `opencode.jsonc`) in your project root:

**Local mode**:
```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "ddg-search": {
      "type": "local",
      "command": ["/path/to/ddg-search-mcp"],
      "enabled": true
    }
  }
}
```

**Remote mode**:
```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "ddg-search": {
      "type": "streamable-http",
      "url": "http://localhost:9100/mcp",
      "enabled": true
    }
  }
}
```

#### Transport Types

MCP supports multiple transport mechanisms for client-server communication:

| Transport | Description | Configuration |
|-----------|-------------|---------------|
| **stdio** | Local, process-based communication via standard input/output. Each editor instance starts its own server process. | Uses `command` and `args` fields (no `type` field needed) |
| **streamable-http** | Enhanced HTTP variant optimized for streaming. Supports partial outputs, long-running operations, and real-time notifications via Server-Sent Events (SSE). | Uses `type: "streamable-http"` with `url` field |

**Stdio Mode** (Recommended for most users):
- Simpler configuration
- Each editor instance starts its own server process
- Better for single-user scenarios
- Lower resource usage when idle
- No `type` field required in configuration

**Streamable HTTP Mode** (Recommended for advanced users):
- Single server instance serves multiple editors
- Better for multi-user scenarios
- Easier to monitor and debug
- Can run as a system service
- Requires managing server lifecycle
- Supports streaming responses for real-time tools
- Requires `type: "streamable-http"` in configuration

#### Configuration Examples

**With Perplexity API enabled**:
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": [],
      "env": {
        "DDG_SEARCH_PERPLEXITY_ACCESS_TOKEN": "your-perplexity-api-token"
      }
    }
  }
}
```

**With custom log level for debugging**:
```json
{
  "mcpServers": {
    "ddg-search": {
      "command": "/path/to/ddg-search-mcp",
      "args": ["--log-level", "debug"]
    }
  }
}
```

**With TLS enabled (HTTP mode)**:
```json
{
  "mcpServers": {
    "ddg-search": {
      "type": "streamable-http",
      "url": "https://localhost:9100/mcp"
    }
  }
}
```

### ddg-search

```bash
ddg-search golang
```

Output (Markdown):

```markdown
1. [The Go Programming Language](https://go.dev/)
   Go is an open source programming language...
2. [Go (programming language) - Wikipedia](https://en.wikipedia.org/wiki/Go_(programming_language))
   ...
```

#### ddg-search Examples

```bash
# Limit results
ddg-search --max-results 5 golang tutorial

# Site-specific search
ddg-search --site github.com docker compose

# Regional search
ddg-search --region uk-en premier league

# Time-bounded search (d=day, w=week, m=month, y=year)
ddg-search --time w news today

# Retry configuration
ddg-search --max-retries 5 --retry-delay 2s --max-retry-delay 60s slow query

# Debug mode
ddg-search --debug golang
```

#### ddg-search Options

| Flag | Description | Default |
|------|-------------|---------|
| `--max-results` | Maximum number of results to return | 10 |
| `--site` | Filter results to a specific domain | - |
| `--region` | Search region (e.g., us-en, uk-en) | us-en |
| `--time` | Time filter: d (day), w (week), m (month), y (year) | - |
| `--safe-search` | Enable safe search | false |
| `--max-retries` | Maximum retry attempts on rate limiting | 3 |
| `--retry-delay` | Initial retry delay | 1s |
| `--max-retry-delay` | Maximum retry delay cap | 30s |
| `--debug` | Enable debug logging to stderr | false |

### page-dump

```bash
page-dump https://example.com
```

Output (Markdown):

```markdown
# Example Domain

This domain is for use in documentation examples...

[Learn more](https://iana.org/domains/example)
```

#### page-dump Options

| Flag | Description | Default |
|------|-------------|---------|
| `--timeout` | Request timeout | 30s |
| `--user-agent` | Custom user agent string | page-dump/1.0 |

### perplexity-search

```bash
# First, set your API key
export PERPLEXITY_API_KEY="your-api-key"

# Then search
perplexity-search "What is Go programming language?"
```

Output (Markdown):

```markdown
Go is an open source programming language that makes it easy to build simple, reliable, and efficient software...
```

#### perplexity-search Examples

```bash
# Use a specific model
perplexity-search --model sonar-pro-online "machine learning fundamentals"
```

#### perplexity-search Options

| Flag | Description | Default |
|------|-------------|---------|
| `--model` | Perplexity model to use | sonar-medium-online |

#### perplexity-search Models

| Model | Description |
|-------|-------------|
| `sonar-small-online` | Faster, lower cost |
| `sonar-medium-online` | Balanced performance and quality |
| `sonar-pro-online` | Higher quality, more expensive |

#### API Key Setup

Perplexity API requires an API key. Get one at https://www.perplexity.ai/settings/api

Set it via environment variable:

```bash
export PERPLEXITY_API_KEY="your-api-key"
```

Or add to `.env` file:

```bash
PERPLEXITY_API_KEY="your-api-key"
```

#### ddg-search vs perplexity-search

| Feature | ddg-search | perplexity-search |
|---------|-------------|------------------|
| API Key Required | No | Yes |
| Result Quality | Standard search results | AI-powered, higher quality |
| Citations | No | Yes, with sources |
| Rate Limits | DuckDuckGo rate limits | Based on API tier |
| Cost | Free | Free tier available, paid for higher limits |
| Speed | Faster (HTML scraping) | Slightly slower (API call) |

Use `ddg-search` for free, fast searches without API key requirements. Use `perplexity-search` for higher-quality, AI-summarized results with citations.

## Rate Limiting

### DuckDuckGo

DuckDuckGo may rate-limit requests. ddg-search automatically:

1. Detects rate-limit responses (HTTP 202, 429, 5xx)
2. Retries with exponential backoff + jitter
3. Fails gracefully with clear error messages after max retries

### Perplexity API

Perplexity API has rate limits based on your API tier. perplexity-search automatically:

1. Detects rate-limit responses (HTTP 429, 5xx)
2. Retries with exponential backoff + jitter
3. Fails gracefully with clear error messages after max retries

If you exceed your rate limit, you'll see an error message indicating the limit. Consider upgrading your Perplexity plan for higher limits at https://www.perplexity.ai/settings/api.

## Use in Skills

These tools are designed for programmatic use in automation skills:

```bash
# Search for information with DuckDuckGo (no API key)
results=$(ddg-search --max-results 3 "$query")

# Search with Perplexity API (requires API key)
results=$(perplexity-search --max-results 3 "$query")

# Fetch and convert a web page
content=$(page-dump "$url")
```

## License

MIT

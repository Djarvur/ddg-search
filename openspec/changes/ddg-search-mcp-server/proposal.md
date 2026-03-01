# MCP Server Implementation Proposal

## Project Overview

**Project Name:** ddg-search-mcp
**Type:** Cloud Code compatible MCP Server
**Goal:** Build an MCP server that exposes web search capabilities (DuckDuckGo + Perplexity) and page fetching to Claude Code and other MCP clients.

---

## 1. Architecture Overview

```
┌───────────────────────────────────────────────────────────────────────┐
│                         ddg-search-mcp                                │
├───────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐           │
│  │   MCP Tool   │     │   MCP Tool   │     │   MCP Tool   │           │
│  │  search      │     │  fetch       │     │  resources   │           │
│  │ (Anthropic)  │     │ (Anthropic)  │     │              │           │
│  └──────┬───────┘     └──────┬───────┘     └──────┬───────┘           │
│         │                    │                    │                   │
│         ▼                    ▼                    ▼                   │
│  ┌─────────────────────────────────────────────────────────────┐      │
│  │                     Tool Handler Layer                      │      │
│  │  - Parameter validation (snake_case for Claude Code)        │      │
│  │  - Output format handling (text/JSON)                       │      │
│  │  - Provider selection (Perplexity vs DuckDuckGo)            │      │
│  │  - Proxy-style request logging                              │      │
│  └──────────────────────────┬──────────────────────────────────┘      │
│                             │                                         │
│         ┌───────────────────┼───────────────────┐                     │
│         ▼                   ▼                   ▼                     │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐              │
│  │  Perplexity │     │ DuckDuckGo  │     │   PageDump  │              │
│  │   Search    │     │   Search    │     │    Fetch    │              │
│  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘              │
│         │                   │                   │                     │
│         └───────────────────┴───────────────────┘                     │
│                             │                                         │
│                             ▼                                         │
│              ┌─────────────────────────┐                              │
│              │    Internal Packages    │                              │
│              │  (search, perplexity,   │                              │
│              │   dump, config)         │                              │
│              └─────────────────────────┘                              │
│                                                                       │
├───────────────────────────────────────────────────────────────────────┤
│                          Transport Layer                              │
│  ┌─────────────────┐              ┌─────────────────┐                 │
│  │     stdio       │              │  StreamableHTTP │                 │
│  │   (default)     │              │   (TCP :9100)   │                 │
│  └─────────────────┘              │  + TLS/mTLS     │                 │
│                                   │  + HTTP/2       │                 │
│                                   └─────────────────┘                 │
│                                                                       │
├───────────────────────────────────────────────────────────────────────┤
│                        Configuration Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐        │
│  │  config.yaml    │  │   env vars      │  │    CLI flags    │        │
│  │ (~/.config/     │  │  DDG_* vars     │  │   --port,       │        │
│  │  ddg-search/    │  │  DDG_* vars     │  │   --transport   │        │
│  │  config.yaml)   │  │                 │  │   --debug       │        │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘        │
└───────────────────────────────────────────────────────────────────────┘
```

---

## 2. MCP Tools Definition

### 2.1 `search` Tool (Anthropic-compatible naming)

**Description:** Search the web using Perplexity (if enabled with valid API key) or DuckDuckGo (fallback). Matches Anthropic's official MCP server naming convention.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "The search query string"
    },
    "max_results": {
      "type": "integer",
      "description": "Maximum number of results (default: 10 for DDG, 5 for Perplexity)",
      "minimum": 1,
      "maximum": 100
    },
    "provider": {
      "type": "string",
      "enum": ["auto", "perplexity", "duckduckgo"],
      "description": "Search provider: auto (default), perplexity, or duckduckgo"
    },
    "model": {
      "type": "string",
      "enum": ["sonar-small-online", "sonar-medium-online", "sonar-pro-online"],
      "description": "Perplexity model (only for perplexity provider)"
    },
    "site": {
      "type": "string",
      "description": "Filter results to specific domain (DDG only)"
    },
    "region": {
      "type": "string",
      "description": "Search region code (DDG only, default: us-en)"
    },
    "time": {
      "type": "string",
      "enum": ["d", "w", "m", "y"],
      "description": "Time filter: d=day, w=week, m=month, y=year (DDG only)"
    },
    "output_format": {
      "type": "string",
      "enum": ["text", "json"],
      "description": "Output format: text (default) or JSON"
    }
  },
  "required": ["query"]
}
```

### 2.2 `fetch` Tool (Anthropic-compatible naming)

**Description:** Fetch a web page and convert it to markdown. Matches Anthropic's official MCP server naming convention.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "url": {
      "type": "string",
      "description": "The URL to fetch (HTTP or HTTPS only)"
    },
    "timeout": {
      "type": "integer",
      "description": "Request timeout in seconds (default: 30)"
    },
    "output_format": {
      "type": "string",
      "enum": ["text", "json"],
      "description": "Output format: text (default) or JSON"
    }
  },
  "required": ["url"]
}
```

> **Note:** MCP resources are implemented for configuration exposure

---

## 3. Configuration System

### 3.1 Configuration File

**Location:** `~/.config/ddg-search/config.yaml`

```yaml
# Server Configuration
server:
  transport: stdio          # stdio | http
  port: 9100              # HTTP server port (default: 9100)
  host: "localhost"       # Bind address

# TLS Configuration
tls:
  enabled: false          # Enable TLS (default: false)
  cert_file: ""           # Path to server certificate
  key_file: ""            # Path to server private key
  combined: ""            # Path to combined cert+key file
  ca_cert: ""             # CA certificate for mTLS
  client_auth: none       # none | request | require

# Perplexity Configuration
perplexity:
  enabled: false          # Enable Perplexity by default
  api_key: ""             # Perplexity API key (can also use DDG_PERPLEXITY_API_KEY env)
  model: "sonar-medium-online"
  max_results: 5
  timeout: 30s

# DuckDuckGo Configuration
duckduckgo:
  max_results: 10
  timeout: 30s
  user_agent: "ddg-search-mcp/1.0"
  region: "us-en"

# Page Dump Configuration
page_dump:
  timeout: 30s
  user_agent: "ddg-search-mcp/1.0"

# Logging Configuration
logging:
  level: info             # debug | info | warn | error
  format: text            # text | json

# Output Configuration
output:
  default_format: text     # text | json
```

### 3.2 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DDG_MCP_TRANSPORT` | Transport mode | `stdio` |
| `DDG_MCP_PORT` | HTTP port | `9100` |
| `DDG_MCP_HOST` | Bind host | `localhost` |
| `DDG_MCP_TLS_ENABLED` | Enable TLS | `false` |
| `DDG_MCP_TLS_CERT_FILE` | TLS certificate path | (empty) |
| `DDG_MCP_TLS_KEY_FILE` | TLS key path | (empty) |
| `DDG_MCP_TLS_COMBINED` | Combined cert+key path | (empty) |
| `DDG_MCP_TLS_CA_CERT` | CA cert for mTLS | (empty) |
| `DDG_MCP_TLS_CLIENT_AUTH` | mTLS mode | `none` |
| `DDG_PERPLEXITY_ENABLED` | Enable Perplexity | `false` |
| `DDG_PERPLEXITY_API_KEY` | Perplexity API key | (empty) |
| `DDG_LOG_LEVEL` | Log level | `info` |
| `DDG_OUTPUT_FORMAT` | Default output format | `text` |

> **Note:** No .env file support - all configuration through config.yaml, environment variables, and CLI flags only.

### 3.3 CLI Flags

```
ddg-search-mcp [flags]

Server Flags:
  --transport string   Transport mode: stdio, http (default "stdio")
  --port int           HTTP server port (default 9100)
  --host string        HTTP server host (default "localhost")
  --config string      Config file path (default "~/.config/ddg-search/config.yaml")

TLS Flags:
  --tls-enabled          Enable TLS (default: false)
  --tls-cert-file string   Server certificate file
  --tls-key-file string    Server private key file
  --tls-combined string     Combined cert+key file
  --tls-ca-cert string      CA certificate for mTLS
  --tls-client-auth string  Client auth mode: none, request, require

Search Flags:
  --max-results int    Default max results
  --region string      Default search region
  --timeout duration   Default timeout

Output Flags:
  --output-format string   Default output format: text, json (default "text")

Logging Flags:
  --log-level string   Log level: debug, info, warn, error (default "info")
  --log-format string Log format: text, json (default "text")

Misc Flags:
  --help, -h           Show help
  --version, -v        Show version
```

### 3.4 Configuration Precedence

CLI flags > Environment variables > Config file > Defaults

---

## 4. Provider Selection Logic

```
┌─────────────────────────────────────────────────────────────┐
│                    Provider Selection                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Input: query, provider=auto (default)                      │
│                                                             │
│  ┌────────────────────────────────────────────────────┐    │
│  │  IF provider != "auto"                              │    │
│  │     THEN use specified provider                    │    │
│  └────────────────────────────────────────────────────┘    │
│                           │                                  │
│                           ▼                                  │
│  ┌────────────────────────────────────────────────────┐    │
│  │  IF provider == "auto"                              │    │
│  │     IF perplexity.enabled AND perplexity.api_key   │    │
│  │        THEN use Perplexity                          │    │
│  │     ELSE use DuckDuckGo                             │    │
│  └────────────────────────────────────────────────────┘    │
│                           │                                  │
│                           ▼                                  │
│  ┌────────────────────────────────────────────────────┐    │
│  │  IF Perplexity returns:                             │    │
│  │     - rate limit error (429)                        │    │
│  │     - quota exceeded (402)                          │    │
│  │     - authentication error (401/403)               │    │
│  │     - network/timeout errors                       │    │
│  │     THEN fallback to DuckDuckGo (ALWAYS, even if  │    │
│  │     provider was explicitly set to perplexity)      │    │
│  │     RETURN success to client (no error)            │    │
│  │     INCLUDE fallback info in response               │    │
│  │     LOG warning (no error to client)               │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.1 Fallback Response Format

When fallback occurs, the response includes metadata indicating the fallback. The client receives a **successful result** with fallback information embedded:

**Text format:**
```
[Results from DuckDuckGo - Perplexity was rate limited]

1. Result Title
   https://example.com
   Description...
```

**JSON format:**
```json
{
  "results": [...],
  "_meta": {
    "provider": "duckduckgo",
    "fallback_reason": "perplexity_rate_limited",
    "original_provider": "perplexity"
  }
}
```

---

## 5. TLS Configuration

The HTTP transport supports optional TLS with mTLS support:

### 5.1 Config Options

```yaml
server:
  transport: http           # stdio | http
  port: 9100              # HTTP server port (default: 9100)
  host: "localhost"       # Bind address

tls:
  enabled: false          # Enable TLS (default: false)
  cert_file: ""           # Path to server certificate
  key_file: ""            # Path to server private key
  combined: ""            # Path to combined cert+key file (overrides cert_file/key_file)
  ca_cert: ""             # CA certificate for mTLS client auth
  client_auth: none       # none | request | require (mTLS mode)
```

### 5.2 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DDG_MCP_TLS_ENABLED` | Enable TLS | `false` |
| `DDG_MCP_TLS_CERT_FILE` | Server certificate path | (empty) |
| `DDG_MCP_TLS_KEY_FILE` | Server key path | (empty) |
| `DDG_MCP_TLS_COMBINED` | Combined cert+key path | (empty) |
| `DDG_MCP_TLS_CA_CERT` | CA cert for mTLS | (empty) |
| `DDG_MCP_TLS_CLIENT_AUTH` | mTLS mode | `none` |

### 5.3 CLI Flags

```
TLS Flags:
  --tls-enabled          Enable TLS (default: false)
  --tls-cert-file string   Server certificate file
  --tls-key-file string   Server private key file
  --tls-combined string   Combined cert+key file
  --tls-ca-cert string   CA certificate for mTLS
  --tls-client-auth string  Client auth mode: none, request, require
```

### 5.4 TLS Modes

| Mode | Description |
|------|-------------|
| `none` | No client certificate required (default) |
| `request` | Request client certificate but don't require |
| `require` | Require valid client certificate (mTLS) |

### 5.5 HTTP/2 Support

HTTP/2 is automatically enabled when TLS is enabled (standard Go net/http behavior). The server will negotiate HTTP/2 with clients that support it.

### 5.6 TLS Certificate Reload on HUP

When SIGHUP is received, the server will reload TLS certificates if the paths have changed. This enables zero-downtime certificate rotation in production environments.

---

## 6. Project Structure

```
cmd/ddg-search-mcp/
└── main.go              # Entry point

internal/
├── config/
    ├── config.go        # Configuration loading (viper)
    └── config_test.go
├── mcp/
    ├── server.go        # MCP server setup
    ├── handlers.go      # Tool handlers
    ├── handlers_test.go
    ├── transport.go     # stdio/HTTP transport setup
    └── resources.go    # MCP resources
├── search/
│   └── (existing)
├── perplexity/
│   └── (existing)
└── dump/
    └── (existing)
```

---

## 6. Dependencies

### Required New Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/mark3labs/mcp-go` | latest | MCP server framework |
| `github.com/spf13/cobra` | ^3.0.0 | CLI framework |
| `github.com/spf13/viper` | ^1.18.0 | Configuration management |

> Note: No .env file support - all configuration through config.yaml, environment variables, and CLI flags.

### Existing Dependencies (Reused)

- `github.com/Djarvur/ddg-search/internal/search` - DuckDuckGo search
- `github.com/Djarvur/ddg-search/internal/perplexity` - Perplexity search
- `github.com/Djarvur/ddg-search/internal/dump` - Page fetching
- `github.com/Djarvur/ddg-search/internal/config` - Shared config

---

## 7. Testing Strategy

### 7.1 Unit Tests

**Location:** `internal/mcp/*_test.go` (public API)
**Location:** `internal/mcp/*_internal_test.go` (private/internal testing)

| Test File | Coverage |
|-----------|----------|
| `config_test.go` | Config loading, env vars, CLI overrides |
| `handlers_test.go` | Tool parameter validation, output formatting |
| `transport_test.go` | Transport selection, server startup |
| `server_test.go` | MCP protocol handling |

### 7.2 Integration Tests

**Location:** `internal/mcp/integration_test.go`

- MCP protocol compliance
- Tool invocation with mocked backends
- Error handling and responses

### 7.3 E2E Tests

**Location:** `e2e/mcp_server_test.go`

| Test | Description |
|------|-------------|
| `TestStdioServer` | Start server in stdio mode, send requests |
| `TestHTTPServer` | Start server on TCP, send HTTP requests |
| `TestSignalHandling` | Verify HUP signal reloads config |
| `TestSignalHandling` | Verify SIGINT/SIGTERM stops server cleanly |
| `TestProviderFallback` | Verify Perplexity→DDG fallback |
| `TestAllParameters` | Test all search parameters |
| `TestOutputFormats` | Test text and JSON output |
| `TestConfigHotReload` | Config file changes picked up on HUP |

### 7.4 TDD Workflow

1. Write failing test first
2. Implement minimal code to pass
3. Refactor
4. Run `mise run lint` before commit
5. Repeat

---

## 8. Implementation Phases

### Phase 1: Foundation (Week 1)
- [ ] Create `cmd/ddg-search-mcp` structure
- [ ] Setup Cobra + Viper configuration
- [ ] Implement config loading with all sources
- [ ] Add logging with slog
- [ ] Basic MCP server setup with mark3labs/mcp-go
- [ ] Unit tests for config

### Phase 2: Transport Layer (Week 1-2)
- [ ] Implement stdio transport
- [ ] Implement HTTP transport (StreamableHTTP)
- [ ] Configurable transport selection
- [ ] Signal handling (HUP for reload, INT/TERM for shutdown)
- [ ] E2E tests for transport

### Phase 3: Tool Handlers (Week 2)
- [ ] Implement `search` tool handler
- [ ] Implement `page_dump` tool handler
- [ ] Provider selection logic
- [ ] Fallback handling
- [ ] Output format handling (text/JSON)

### Phase 4: Testing & Polish (Week 2-3)
- [ ] Complete unit test coverage
- [ ] Integration tests
- [ ] E2E tests (signals, config reload)
- [ ] Run `mise run lint` - fix issues
- [ ] Documentation

---

## 9. Logging Specification (Proxy-style)

Using Go's `log/slog` package with proxy/server-style request logging:

### 9.1 Log Levels

| Level | Usage |
|-------|-------|
| `debug` | Detailed request/response debugging |
| `info` | Normal server operations, successful requests |
| `warn` | Recoverable issues (e.g., Perplexity fallback) |
| `error` | Failed requests, configuration errors |

### 9.2 Request Logging (Proxy-style)

Every MCP request is logged with full details, similar to HTTP proxy servers. **Successful requests are logged at DEBUG level**, while errors and warnings are logged at their respective levels:

```
# Successful request - DEBUG level (not logged by default)
2026/03/01 12:00:00 DEBUG method=tools/call name=search query="golang http client" provider=perplexity status=200 duration=1.234s results=5

# Bad request / error - ERROR level
2026/03/01 12:00:01 ERROR method=tools/call name=search query="" error="query required" status=400 duration=0.001s

# Fallback info - WARN level (but returns success to client)
2026/03/01 12:00:02 WARN method=tools/call name=search query="test" provider=perplexity status=429 error="rate limited" fallback=duckduckgo

# Rate limited but fallback succeeded - INFO level (since result returned)
2026/03/01 12:00:03 INFO method=tools/call name=search query="test" provider=perplexity status=200 fallback=duckduckgo duration=2.5s results=5
```

### 9.3 JSON Log Format

```json
{"time":"2026-03-01T12:00:00.000Z","level":"INFO","method":"tools/call","name":"search","query":"golang http client","provider":"perplexity","status":200,"duration":1.234,"results":5}
{"time":"2026-03-01T12:00:01.000Z","level":"ERROR","method":"tools/call","name":"search","query":"","error":"query required","status":400,"duration":0.001}
```

### 9.4 Log Fields

| Field | Type | Description |
|-------|------|-------------|
| `time` | ISO8601 | Request timestamp |
| `level` | string | Log level (debug/info/warn/error) |
| `method` | string | MCP method (tools/call, tools/list, etc.) |
| `name` | string | Tool name (search, fetch) |
| `query` | string | Search query (sanitized) |
| `url` | string | For fetch, the URL being fetched |
| `provider` | string | Provider used (perplexity/duckduckgo) |
| `status` | int | HTTP-style status code |
| `duration` | float | Request duration in seconds |
| `results` | int | Number of results (for search) |
| `error` | string | Error message if failed |
| `fallback` | string | Fallback provider if applicable |

### 9.5 Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad request (invalid params) |
| 401 | Unauthorized (invalid API key) |
| 402 | Payment required (quota exceeded) |
| 429 | Rate limited |
| 500 | Internal server error |
| 502 | Upstream error (Perplexity/DDG failed) |
| 503 | Service unavailable |

### 9.6 Implementation

```go
// Request logger with proxy-style formatting
func (s *Server) logRequest(logger *slog.Logger, method, toolName string, duration time.Duration, status int, fields ...slog.Attr) {
    attrs := []slog.Attr{
        slog.String("method", method),
        slog.String("name", toolName),
        slog.Int("duration_ms", int(duration.Milliseconds())),
        slog.Int("status", status),
    }
    attrs = append(attrs, fields...)

    switch {
    case status >= 500:
        logger.Error("request failed", attrs...)
    case status >= 400:
        logger.Warn("bad request", attrs...)
    default:
        logger.Info("request completed", attrs...)
    }
}
```

---

## 10. Error Handling

### MCP Error Codes

| Error | MCP Code | Description |
|-------|----------|-------------|
| Invalid query | `-32602` | Invalid params |
| Search failed | `-32603` | Internal error |
| Rate limited | `-32001` | Too many requests |
| Network error | `-32002` | Connection failed |
| Config error | `-32003` | Invalid configuration |

### Error Format

```json
{
  "error": {
    Response "code": "-32001",
    "message": "Rate limited by Perplexity API, falling back to DuckDuckGo",
    "data": {
      "originalProvider": "perplexity",
      "fallbackProvider": "duckduckgo"
    }
  }
}
```

---

## 11. Signal Handling

| Signal | Action |
|--------|--------|
| `SIGINT` (Ctrl-C) | Immediate shutdown |
| `SIGTERM` | Immediate shutdown |
| `SIGHUP` | Reload configuration file |

### HUP Reload Behavior

1. Receive SIGHUP
2. Load config from file (respecting `--config` flag)
3. Update all mutable settings:
   - Log level
   - Output format
   - Perplexity enabled/disabled
4. Keep existing connections alive
5. New requests use new config

---

## 12. Open Questions / Clarifications Needed

- [x] Perplexity API key: Direct `PERPLEXITY_API_KEY` → **Confirmed**
- [x] Tools: Two tools (search, page_dump) with auto-provider → **Confirmed**
- [x] Output: JSON as string in MCP response for Claude Code compatibility → **Confirmed**
- [x] Transport: Both stdio and HTTP, configurable, stdio default → **Confirmed**
- [x] Log levels: debug/info/warn/error configurable → **Confirmed**

---

## 13. Acceptance Criteria

### Must Have
- [ ] MCP server compiles and runs
- [ ] Both stdio and HTTP transports work
- [ ] `search` tool returns results
- [ ] `page_dump` tool fetches pages
- [ ] Perplexity used when API key provided and enabled
- [ ] DuckDuckGo fallback works
- [ ] Configuration via config file, env vars, and CLI
- [ ] Config hot-reload on SIGHUP
- [ ] Immediate shutdown on SIGINT/SIGTERM
- [ ] All tests pass
- [ ] `mise run lint` passes without disabling linters

### Should Have
- [ ] JSON output format
- [ ] Debug logging
- [ ] Proper error messages

### Nice to Have
- [ ] Prometheus metrics

---

## 13. References

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - MCP framework
- [MCP Specification](https://spec.modelcontextprotocol.io/) - Protocol spec
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration
- [slog](https://pkg.go.dev/log/slog) - Structured logging

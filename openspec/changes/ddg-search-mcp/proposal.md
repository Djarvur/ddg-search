# Proposal: ddg-search-mcp - MCP Server for Web Search

## Summary

Create a new MCP (Model Context Protocol) server binary `ddg-search-mcp` that exposes the existing search and page fetch capabilities as MCP tools, enabling Claude Code and other MCP clients to use DuckDuckGo search (with optional Perplexity fallback) and page content fetching.

## Motivation

Currently, the ddg-search project provides CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) that are wrapped as skills for Claude Code. An MCP server would provide:

1. **Native MCP Integration**: Direct integration with Claude Code without skill wrappers
2. **Standard Protocol**: Compatibility with any MCP client (not just Claude Code)
3. **Flexible Transport**: Support for stdio (local) and HTTP/TCP (remote) transports
4. **Intelligent Fallback**: Automatic Perplexity → DuckDuckGo fallback when configured

## Scope

### In Scope

- New binary `ddg-search-mcp` (separate from existing CLI tools)
- Two MCP tools: `search` and `fetch`
- Transport modes: stdio (default), TCP, and HTTP (configurable port)
- TLS support with configurable certificates
- Output formats: JSON and plain text (configurable at server and request level)
- Perplexity integration with automatic fallback to DuckDuckGo
- Configuration via config file + environment variables + CLI flags (cobra + viper)
- Signal handling: graceful stop on SIGINT, config reload on SIGHUP
- Comprehensive test coverage (TDD approach)

### Out of Scope

- Graceful shutdown (just stop on signal)
- .env file support (use config + env + CLI only)
- Modifications to existing CLI tools
- MCP resources or prompts (tools only for v1)

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        ddg-search-mcp Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         cmd/ddg-search-mcp                             │  │
│  │                                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐  ┌────────────────────────┐   │  │
│  │  │ cobra + viper  │──│  MCP Server    │──│ Transport Layer        │   │  │
│  │  │                │  │  (mark3labs)   │  │                        │   │  │
│  │  │ - config file  │  │                │  │ - stdio (default)
  │  │  │ - TCP (port 9100)
  │  │  │ - HTTP (port 8080)
  │  │  │ - TLS (optional)       │   │  │
│  │  └────────────────┘  └───────┬────────┘  └────────────────────────┘   │  │
│  │                              │                                         │  │
│  └──────────────────────────────┼─────────────────────────────────────────┘  │
│                                 │                                            │
│                    ┌────────────┴────────────┐                               │
│                    ▼                         ▼                               │
│           ┌────────────────┐       ┌────────────────┐                        │
│           │    search      │       │     fetch      │                        │
│           │    tool        │       │     tool       │                        │
│           └───────┬────────┘       └───────┬────────┘                        │
│                   │                        │                                 │
│         ┌─────────┴─────────┐              │                                 │
│         ▼                   ▼              ▼                                 │
│  ┌────────────────┐ ┌────────────────┐ ┌────────────────┐                   │
│  │internal/search │ │internal/       │ │  internal/dump │                   │
│  │                │ │  perplexity    │ │                │                   │
│  │ Searcher       │ │ Client         │ │ FetchAndConvert│                   │
│  │ (DuckDuckGo)   │ │ (with fallback)│ │                │                   │
│  └────────────────┘ └────────────────┘ └────────────────┘                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## MCP Tools

### 1. `search` Tool

Web search with intelligent backend selection.

```json
{
  "name": "search",
  "description": "Search the web using DuckDuckGo or Perplexity API",
  "inputSchema": {
    "type": "object",
    "properties": {
      "query": {
        "type": "string",
        "description": "The search query string"
      },
      "backend": {
        "type": "string",
        "enum": ["auto", "ddg", "perplexity"],
        "default": "auto",
        "description": "Search backend: auto (perplexity if configured, else ddg), ddg, or perplexity"
      },
      "max_results": {
        "type": "integer",
        "minimum": 1,
        "maximum": 50,
        "default": 10,
        "description": "Maximum number of results to return"
      },
      "site": {
        "type": "string",
        "description": "Filter results to a specific domain (DuckDuckGo only)"
      },
      "region": {
        "type": "string",
        "default": "us-en",
        "description": "Search region, e.g., us-en, uk-en (DuckDuckGo only)"
      },
      "time_filter": {
        "type": "string",
        "enum": ["d", "w", "m", "y"],
        "description": "Time filter: d (day), w (week), m (month), y (year) (DuckDuckGo only)"
      },
      "safe_search": {
        "type": "boolean",
        "default": false,
        "description": "Enable safe search (DuckDuckGo only)"
      },
      "model": {
        "type": "string",
        "description": "Perplexity model to use (Perplexity only)"
      },
      "format": {
        "type": "string",
        "enum": ["text", "json"],
        "description": "Output format (overrides server default)"
      }
    },
    "required": ["query"]
  }
}
```

**Backend Selection Logic:**

| backend param | Perplexity Configured | Perplexity Available | Result |
|---------------|----------------------|---------------------|--------|
| `auto` | No | - | DuckDuckGo |
| `auto` | Yes | Yes | Perplexity |
| `auto` | Yes | Rate Limited | DuckDuckGo (with note) |
| `ddg` | - | - | DuckDuckGo |
| `perplexity` | No | - | DuckDuckGo (with note) |
| `perplexity` | Yes | Yes | Perplexity |
| `perplexity` | Yes | Rate Limited | DuckDuckGo (with note) |

### 2. `fetch` Tool

Fetch and convert web page content to markdown.

```json
{
  "name": "fetch",
  "description": "Fetch a web page and convert to markdown",
  "inputSchema": {
    "type": "object",
    "properties": {
      "url": {
        "type": "string",
        "format": "uri",
        "description": "The URL to fetch (HTTP or HTTPS only)"
      },
      "timeout": {
        "type": "integer",
        "minimum": 1,
        "maximum": 120,
        "default": 30,
        "description": "Request timeout in seconds"
      },
      "user_agent": {
        "type": "string",
        "description": "Custom user agent string"
      },
      "format": {
        "type": "string",
        "enum": ["text", "json"],
        "description": "Output format (overrides server default)"
      }
    },
    "required": ["url"]
  }
}
```

## Configuration

### Config File Structure

Location: `~/.config/ddg-search/config.yaml`

```yaml
# Server configuration
server:
  transport: stdio          # stdio | tcp | http
  port: 9100                # TCP/HTTP port (default: 9100 for TCP, 8080 for HTTP)

# TLS configuration (optional)
tls:
  enabled: false
  cert_file: ""             # Path to cert file
  key_file: ""              # Path to key file
  cert_key_file: ""         # Path to combined cert+key file (alternative)
  ca_file: ""               # Path to CA cert for mTLS (optional)

# Output configuration
output:
  format: text              # text | json (default format, can be overridden per-request)

# Search configuration
search:
  max_results: 10           # Default max results
  region: us-en             # Default region for DDG
  safe_search: false        # Default safe search

# Perplexity configuration (optional)
perplexity:
  enabled: false            # Enable Perplexity backend
  api_key: ""               # API key (or use PERPLEXITY_API_KEY env)
  model: sonar-medium-online  # Default model

# Retry configuration
retry:
  max_retries: 3
  base_delay: 1s
  max_delay: 30s
  backoff_multiplier: 2.0

# Logging configuration
logging:
  level: info               # debug | info | warn | error
  format: text              # text | json
```

### Environment Variables

All config values can be overridden via environment variables with prefix `DDG_MCP_`:

```bash
DDG_MCP_SERVER_TRANSPORT=tcp
DDG_MCP_SERVER_PORT=9100
DDG_MCP_TLS_ENABLED=true
DDG_MCP_TLS_CERT_FILE=/path/to/cert.pem
DDG_MCP_TLS_KEY_FILE=/path/to/key.pem
DDG_MCP_OUTPUT_FORMAT=json
DDG_MCP_PERPLEXITY_ENABLED=true
DDG_MCP_PERPLEXITY_API_KEY=pplx-xxx
DDG_MCP_LOGGING_LEVEL=debug
```

### CLI Flags

```bash
ddg-search-mcp [flags]

Flags:
      --config string           config file (default is ~/.config/ddg-search/config.yaml)
  -t, --transport string        transport mode: stdio, tcp, http (default: stdio)
  -p, --port int                TCP port (default: 9100)
      --tls                     enable TLS
      --tls-cert string         TLS certificate file
      --tls-key string          TLS key file
      --tls-ca string           CA certificate file for mTLS
  -f, --format string           default output format: text, json (default: text)
      --perplexity              enable Perplexity backend
      --perplexity-key string   Perplexity API key
  -l, --log-level string        log level: debug, info, warn, error (default: info)
  -v, --version                 version info
```

## Response Format

### Text Format (default)

```
# Search Results for "golang best practices"

Backend: DuckDuckGo (fallback from Perplexity: rate limit exceeded)

## Result 1
Title: Go Best Practices - Official Docs
URL: https://go.dev/doc/effective_go
Snippet: Go is a general-purpose language designed with...

## Result 2
Title: Go Code Review Comments
URL: https://github.com/golang/go/wiki/CodeReviewComments
Snippet: This page collects common comments made during...

---
Found 10 results
```

### JSON Format

```json
{
  "query": "golang best practices",
  "backend_used": "ddg",
  "fallback": {
    "from": "perplexity",
    "reason": "rate limit exceeded"
  },
  "results": [
    {
      "title": "Go Best Practices - Official Docs",
      "url": "https://go.dev/doc/effective_go",
      "snippet": "Go is a general-purpose language designed with..."
    },
    {
      "title": "Go Code Review Comments",
      "url": "https://github.com/golang/go/wiki/CodeReviewComments",
      "snippet": "This page collects common comments made during..."
    }
  ],
  "total": 10
}
```

## Signal Handling

| Signal | Behavior |
|--------|----------|
| SIGINT (Ctrl-C) | Stop the server immediately |
| SIGHUP | Reload configuration from file |

## Project Structure

```
cmd/
├── ddg-search/          # Existing CLI (unchanged)
├── page-dump/           # Existing CLI (unchanged)
├── perplexity-search/   # Existing CLI (unchanged)
└── ddg-search-mcp/      # NEW: MCP server binary
    └── main.go

internal/
├── config/              # Existing (unchanged)
├── search/              # Existing (reused)
├── perplexity/          # Existing (reused)
├── dump/                # Existing (reused)
├── mcp/                 # NEW: MCP-specific code
│   ├── server.go        # MCP server setup
│   ├── server_test.go   # Public tests
│   ├── server_internal_test.go  # Private tests
│   ├── tools.go         # Tool definitions
│   ├── tools_test.go
│   ├── tools_internal_test.go
│   ├── search.go        # Search tool handler
│   ├── search_test.go
│   ├── fetch.go         # Fetch tool handler
│   ├── fetch_test.go
│   ├── format.go        # Output formatters
│   ├── format_test.go
│   └── config/          # MCP-specific config
│       ├── config.go
│       └── config_test.go
└── e2e/                 # NEW: End-to-end tests
    ├── e2e_test.go      # Main E2E test suite
    ├── stdio_test.go    # stdio transport tests
    ├── tcp_test.go      # TCP transport tests
    ├── http_test.go     # HTTP transport tests
    ├── tls_test.go      # TLS tests
    └── signal_test.go   # Signal handling tests
```

## Dependencies

Add to go.mod:

```go
require (
    github.com/mark3labs/mcp-go v0.x.x  // MCP implementation
    github.com/spf13/cobra v1.x.x       // CLI framework
    github.com/spf13/viper v1.x.x       // Configuration
)
```

## Implementation Phases

### Phase 1: Core Infrastructure
- [ ] Set up cobra + viper configuration
- [ ] Create MCP server skeleton with stdio transport
- [ ] Implement config file loading and env var binding
- [ ] Add logging with slog

### Phase 2: Search Tool
- [ ] Implement `search` tool definition
- [ ] Integrate internal/search for DuckDuckGo
- [ ] Integrate internal/perplexity with fallback logic
- [ ] Add response formatting (text/json)

### Phase 3: Fetch Tool
- [ ] Implement `fetch` tool definition
- [ ] Integrate internal/dump
- [ ] Add response formatting

### Phase 4: Network Transports
- [ ] Add TCP transport support
- [ ] Add HTTP transport support
- [ ] Implement TLS support
- [ ] Add transport selection logic

### Phase 5: Signal Handling & E2E
- [ ] Implement SIGINT handling
- [ ] Implement SIGHUP config reload
- [ ] Write E2E tests

### Phase 6: Polish
- [ ] Run `mise run lint` and fix issues
- [ ] Update documentation
- [ ] Final testing

## Testing Strategy

### Unit Tests (TDD)
- Tool parameter validation
- Backend selection logic
- Fallback behavior
- Response formatting
- Configuration loading

### Integration Tests
- MCP protocol compliance
- Tool execution with real backends (mocked HTTP)

### E2E Tests
- Server startup/shutdown
- stdio transport communication
- TCP transport communication
- TLS handshake
- Signal handling (SIGINT, SIGHUP)
- Config reload behavior

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| MCP spec changes | Use stable mark3labs/mcp-go library |
| Perplexity API changes | Isolate in internal/perplexity package |
| Rate limiting issues | Implement robust retry with backoff |
| TLS configuration complexity | Provide clear examples and validation |

## Success Criteria

1. ✅ MCP server starts successfully in stdio mode
2. ✅ MCP server starts successfully in TCP mode
3. ✅ TLS works with separate and combined cert files
4. ✅ `search` tool returns results from DuckDuckGo
5. ✅ `search` tool uses Perplexity when configured
6. ✅ `search` tool falls back to DuckDuckGo on Perplexity errors
7. ✅ `fetch` tool returns markdown content
8. ✅ JSON and text output formats work correctly
9. ✅ Config reload on SIGHUP works
10. ✅ Server stops on SIGINT
11. ✅ All linters pass (`mise run lint`)
12. ✅ Test coverage > 80%

## Open Questions

None - all questions have been resolved:
- Tool naming: `search` and `fetch` (snake_case)
- Backend selection: `backend` param with `auto` default, fallback always
- Config location: `~/.config/ddg-search/config.yaml` only
- Fallback notification: In metadata field for JSON, note in text

## References

- [MCP Specification](https://modelcontextprotocol.io/)
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) (official, alternative)
- Existing skills: `skills/ddg-search/`, `skills/perplexity-search/`

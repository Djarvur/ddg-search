# Design: Add MCP Server

## Context

### Background

The project currently provides three CLI tools for web search and page fetching:
- `ddg-search`: DuckDuckGo search (no API key required)
- `perplexity-search`: Perplexity API search (requires API key)
- `page-dump`: Fetch web page and convert to markdown

These tools are well-tested with robust retry logic and error handling in `internal/search/`, `internal/perplexity/`, and `internal/dump/` packages. However, they lack MCP (Model Context Protocol) integration for seamless use with Claude Code and other MCP-compatible AI assistants.

### Current State

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Current Architecture                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                    CLI Tools (No MCP)                             │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │  │
│  │  │ ddg-search   │  │ page-dump    │  │ perplexity-  │       │  │
│  │  │ (CLI)        │  │ (CLI)        │  │ search       │       │  │
│  │  └──────────────┘  └──────────────┘  │ (CLI)        │       │  │
│  │                    └──────────────┘                               │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                   Internal Packages (Reusable)                      │  │
│  │  • search/    - DuckDuckGo client with retry logic            │  │
│  │  • perplexity/ - Perplexity API client with retry logic         │  │
│  │  • dump/      - Page fetch and markdown conversion            │  │
│  │  • config/    - Configuration types                            │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Constraints

1. **Package Selection**: Must use `mark3labs/mcp-go` (selected for fastest implementation and community adoption)
2. **Existing Code**: Must reuse `internal/search/`, `internal/perplexity/`, `internal/dump/` packages
3. **CLI Framework**: Must use cobra+viper (not urfave/cli/v3 used in existing tools)
4. **Logging**: Must use `log/slog` with text handler
5. **Testing**: Strictly TDD approach with public (`*_test.go`) and internal (`*_internal_test.go`) tests
6. **Tooling**: Must use `mise` (not `make`)
7. **Linting**: Must pass `mise run lint` without disabling linters
8. **No Graceful Shutdown**: Not required per requirements

### Stakeholders

- **Claude Code Users**: Need stdio transport for desktop integration
- **Web Client Users**: Need TCP transport with TLS support
- **Developers**: Need clear configuration and testing patterns

## Goals / Non-Goals

**Goals:**

1. Create `ddg-search-mcp` binary implementing MCP server
2. Expose `search` tool with automatic Perplexity → DuckDuckGo fallback
3. Expose `fetch_page` tool for web page fetching
4. Support stdio transport for Claude Desktop compatibility
5. Support TCP transport with configurable host/port (default: 0.0.0.0:9100)
6. Implement TLS support in three modes:
   - Key/Cert files
   - Combined PEM file
   - mTLS with CA certificate
7. Provide flexible configuration via cobra+viper (CLI → ENV → YAML priority)
8. Implement structured logging using `log/slog` text handler
9. Support HUP signal for config reload with validation
10. Support all existing search parameters (max-results, site, region, time, safe-search, model)
11. Support JSON and markdown output formats
12. Follow TDD approach with comprehensive test coverage
13. Ensure all linters pass without disabling

**Non-Goals:**

1. Modify existing CLI tools (`ddg-search`, `perplexity-search`, `page-dump`)
2. Implement graceful shutdown (not required)
3. Support `.env` file configuration (use CLI/ENV/YAML only)
4. Create skill file in `skills/` directory (MCP server is standalone)
5. Implement OAuth authentication (not required for this use case)

## Decisions

### 1. Package Selection: mark3labs/mcp-go

**Decision:** Use `mark3labs/mcp-go` for MCP implementation

**Rationale:**
- Most popular community implementation (8,253 stars)
- Fastest to implement with minimal boilerplate
- Full transport support (stdio, HTTP, SSE)
- Excellent documentation and examples
- Active development and community support

**Alternatives Considered:**
- `modelcontextprotocol/go-sdk`: Official but more verbose, steeper learning curve
- `metoro-io/mcp-golang`: Type-safe but HTTP is stateless
- `ThinkInAIXYZ/go-mcp`: Good for web framework integration but less active

### 2. Tool Architecture: Two Tools with Fallback

**Decision:** Implement two MCP tools:
1. `search`: Web search with automatic provider fallback
2. `fetch_page`: Fetch and convert web page to markdown

**Rationale:**
- Single `search` tool simplifies user experience (no need to choose provider)
- Automatic fallback from Perplexity to DuckDuckGo provides resilience
- Separate `fetch_page` tool keeps concerns clear
- Matches common MCP patterns for search servers

**Tool Schema:**

```
search:
  - query (required): string
  - provider: enum {auto, duckduckgo, perplexity} (default: auto)
  - max-results: integer (default: 10)
  - site: string (filter to domain)
  - region: string (default: us-en)
  - time: string {d, w, m, y}
  - safe-search: boolean (default: false)
  - model: string (perplexity model: sonar-medium-online, etc.)
  - format: enum {json, markdown} (default: markdown)

fetch_page:
  - url (required): string
  - timeout: duration (default: 30s)
  - user-agent: string (default: ddg-search-mcp/1.0)
```

### 3. Configuration: cobra+viper with Priority Chain

**Decision:** Use cobra+viper for configuration with priority: CLI flags → ENV vars → YAML config

**Rationale:**
- Industry-standard Go CLI framework
- Viper provides flexible config binding and automatic ENV variable mapping
- Priority chain allows easy override at any level
- YAML config file for persistent settings

**Config Structure:**

```yaml
# ~/.config/ddg-search/config.yaml
server:
  transport: stdio  # or tcp
  host: 0.0.0.0
  port: 9100
  tls:
    enabled: false
    mode: keycert  # keycert, combined, mtls
    cert_file: /path/to/cert.pem
    key_file: /path/to/key.pem
    ca_file: /path/to/ca.pem  # for mtls

perplexity:
  enabled: false  # disabled by default
  api_key: ""  # can also be set via PERPLEXITY_API_KEY env var
  model: sonar-medium-online
  max_results: 5

search:
  max_results: 10
  region: us-en
  safe_search: false

dump:
  timeout: 30s
  user_agent: ddg-search-mcp/1.0

logging:
  level: info  # debug, info, warn, error
  format: text  # text handler
```

**Environment Variables:**
- `DDG_SEARCH_CONFIG`: Path to config file (default: ~/.config/ddg-search/config.yaml)
- `PERPLEXITY_API_KEY`: Perplexity API key
- `DDG_SEARCH_TRANSPORT`: Transport type (stdio/tcp)
- `DDG_SEARCH_HOST`: Bind host
- `DDG_SEARCH_PORT`: Bind port
- `DDG_SEARCH_TLS_ENABLED`: Enable TLS (true/false)
- `DDG_SEARCH_LOG_LEVEL`: Log level

### 4. Transport Architecture: Dual Mode Support

**Decision:** Support both stdio and TCP transports with runtime selection

**Rationale:**
- stdio required for Claude Desktop integration
- TCP required for web clients and remote access
- Single binary supports both use cases
- Configurable via `--transport` flag

**Transport Selection Logic:**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Transport Selection                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │  stdio Transport                                            │  │
│  │  • For Claude Desktop                                        │  │
│  │  • Full MCP feature support                                     │  │
│  │  • Bidirectional (notifications)                                 │  │
│  │  • No network binding needed                                    │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │  TCP Transport                                              │  │
│  │  • For web clients, remote access                              │  │
│  │  • Full MCP feature support                                     │  │
│  │  • Configurable host:port (default: 0.0.0.0:9100)          │  │
│  │  • Optional TLS (key/cert, combined, mTLS)                     │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  Selection: --transport {stdio|tcp} (default: from config)            │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5. TLS Support: Three Modes

**Decision:** Implement three TLS modes with unified configuration

**Rationale:**
- Different deployment scenarios require different TLS configurations
- Key/Cert: Standard for most deployments
- Combined PEM: Simplified for some infrastructure (e.g., Kubernetes secrets)
- mTLS: Required for mutual authentication scenarios

**TLS Configuration:**

```go
type TLSConfig struct {
    Enabled bool
    Mode    string  // "keycert", "combined", "mtls"
    CertFile string
    KeyFile  string
    CAFile   string  // for mTLS
}
```

**Implementation:**
- Use `tls.LoadX509KeyPair()` for key/cert mode
- Use `tls.X509KeyPair` for combined PEM parsing
- Use `tls.Config` with `ClientCAs` for mTLS

### 6. Logging: log/slog with Text Handler

**Decision:** Use `log/slog` with text handler for structured logging

**Rationale:**
- Go 1.21+ standard library
- Text handler is human-readable (user requirement)
- Structured logging enables better debugging
- Consistent with "proxy servers" logging pattern mentioned

**Log Format:**

```
[INFO] Starting MCP server (transport=tcp, host=0.0.0.0, port=9100)
[INFO] TLS enabled (mode=keycert, cert=/path/to/cert.pem)
[DEBUG] Handling search request (query="golang mcp", provider=auto)
[DEBUG] Trying Perplexity API (model=sonar-medium-online)
[WARN] Perplexity rate limited, falling back to DuckDuckGo
[INFO] Search completed (provider=duckduckgo, results=5, duration=1.2s)
[ERROR] Invalid URL format (url="htp://example.com")
[INFO] Config reload triggered (SIGHUP)
[INFO] Config reloaded and validated
```

**Log Levels:**
- `DEBUG`: Successful requests (user requirement)
- `INFO`: Server lifecycle events, configuration
- `WARN`: Fallback events, retries
- `ERROR`: Failed requests, validation errors

### 7. Search Fallback Logic: Resilient Provider Selection

**Decision:** Implement automatic fallback with explicit override capability

**Rationale:**
- Automatic fallback provides resilience
- Explicit override allows user control
- Fallback events logged for transparency
- No error returned on fallback (user requirement)

**Fallback Logic:**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Search Provider Selection                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │  Provider Selection Flow                                      │  │
│  │                                                             │  │
│  │  1. Check provider parameter:                                  │  │
│  │     • auto → Try Perplexity first, fallback to DDG            │  │
│  │     • perplexity → Try Perplexity only, fallback to DDG          │  │
│  │     • duckduckgo → Use DDG only, no fallback                   │  │
│  │                                                             │  │
│  │  2. For auto/perplexity:                                      │  │
│  │     • Check if Perplexity enabled and API key configured             │  │
│  │     • If yes → Try Perplexity                                  │  │
│  │       ├─ Success → Return results                                 │  │
│  │       └─ Rate limit/Quota exceeded → Fallback to DDG (log)     │  │
│  │     • If no → Use DDG directly                                 │  │
│  │                                                             │  │
│  │  3. Fallback behavior:                                         │  │
│  │     • Log fallback event (WARN level)                              │  │
│  │     • Include fallback reason in response                              │  │
│  │     • No error returned to client                                 │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Fallback Response Example:**

```json
{
  "results": [...],
  "provider_used": "duckduckgo",
  "fallback_reason": "Perplexity API rate limit exceeded (HTTP 429)",
  "original_provider": "perplexity"
}
```

### 8. Config Reload: HUP Signal with Validation

**Decision:** Implement HUP signal handler for config reload with validation

**Rationale:**
- Allows runtime configuration changes without restart
- Validation prevents invalid config from breaking server
- Error logging provides visibility into reload issues

**Reload Flow:**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HUP Signal Handling                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Receive SIGHUP signal                                          │
│  2. Log reload event (INFO)                                         │
│  3. Read config file from disk                                        │
│  4. Validate new config:                                             │
│     • Check file format (YAML parsing)                                 │
│     • Validate TLS files exist (if TLS enabled)                           │
│     • Validate port is available (if TCP)                                │
│  5. If validation fails:                                             │
│     • Log error (ERROR)                                                │
│     • Keep old config active                                           │
│  6. If validation passes:                                            │
│     • Apply new config                                                 │
│     • Log success (INFO)                                              │
│     • Reload TLS listener if needed                                      │
│  7. Continue serving with new config                                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Validation Rules:**
- TLS files must exist and be readable
- Port must be in valid range (1-65535)
- Host must be valid IP or hostname
- Perplexity API key format (if provided)
- Log level must be valid (debug, info, warn, error)

### 9. Testing Strategy: TDD with Comprehensive Coverage

**Decision:** Strict TDD approach with public and internal test separation

**Rationale:**
- TDD ensures code quality from the start
- Public tests document API contracts
- Internal tests verify implementation details
- E2E tests validate complete workflows

**Test Structure:**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Test Organization                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  cmd/ddg-search-mcp/                                               │
│  │  ├── main_test.go          (E2E tests)                         │
│  │  │   • Test server start/stop                                     │
│  │  │   • Test HUP config reload                                     │
│  │  │   • Test tool invocation via stdio                              │
│  │  │   • Test tool invocation via TCP                                 │
│  │  │   • Test TLS modes                                              │
│  │                                                                   │
│  internal/mcp/                                                        │
│  │  ├── server_test.go       (Public API tests)                    │
│  │  │   • Test server initialization                                     │
│  │  │   • Test transport selection                                      │
│  │  │   • Test config loading                                          │
│  │  │   • Test signal handling                                         │
│  │  │                                                                 │
│  │  ├── server_internal_test.go (Internal tests)                    │
│  │  │   • Test fallback logic                                          │
│  │  │   • Test config validation                                      │
│  │  │   • Test TLS configuration                                      │
│  │  │                                                                 │
│  internal/mcp/tools/                                                  │
│  │  ├── search_test.go       (Public API tests)                    │
│  │  │   • Test search tool registration                                │
│  │  │   • Test parameter parsing                                       │
│  │  │   • Test output formatting (json/markdown)                        │
│  │  │                                                                 │
│  │  ├── search_internal_test.go (Internal tests)                    │
│  │  │   • Test Perplexity client integration                            │
│  │  │   • Test DuckDuckGo client integration                           │
│  │  │   • Test fallback logic                                         │
│  │  │   • Test error handling                                        │
│  │  │                                                                 │
│  │  ├── fetch_page_test.go   (Public API tests)                    │
│  │  │   • Test fetch_page tool registration                          │
│  │  │   • Test URL validation                                        │
│  │  │   • Test markdown conversion                                   │
│  │  │                                                                 │
│  │  └── fetch_page_internal_test.go (Internal tests)                    │
│  │       • Test dump package integration                              │
│  │       • Test timeout handling                                     │
│  │       • Test error propagation                                     │
│  │                                                                 │
│  internal/mcp/config/                                                 │
│  │  ├── config_test.go      (Public API tests)                   │
│  │  │   • Test config loading from file/ENV/CLI                    │
│  │  │   • Test config validation                                      │
│  │  │   • Test priority chain                                        │
│  │  │                                                                 │
│  │  └── config_internal_test.go (Internal tests)                    │
│  │       • Test YAML parsing                                        │
│  │       • Test default values                                      │
│  │       • Test type conversion                                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

**E2E Test Scenarios:**
1. Server starts in stdio mode
2. Server starts in TCP mode
3. Server starts with TLS (key/cert mode)
4. Server starts with TLS (combined PEM mode)
5. Server starts with TLS (mTLS mode)
6. Server stops on Ctrl-C
7. Server reloads config on HUP
8. Server rejects invalid config on HUP
9. Search tool returns results from Perplexity
10. Search tool falls back to DuckDuckGo on Perplexity rate limit
11. Search tool respects provider override
12. Fetch page tool returns markdown
13. Logging outputs correct format and levels

## Risks / Trade-offs

### Risk 1: Perplexity Rate Limits Cause Frequent Fallbacks

**Description:** If Perplexity API is frequently rate-limited, the server may fall back to DuckDuckGo often, reducing the value of Perplexity integration.

**Mitigation:**
- Log all fallback events with reason
- Allow explicit provider override to avoid automatic fallback
- Document fallback behavior clearly in tool description
- Consider adding retry delay between Perplexity attempts

### Risk 2: Config Reload on HUP May Disrupt Active Requests

**Description:** Reloading config while requests are in progress could cause inconsistent behavior or connection drops.

**Mitigation:**
- Validate config before applying
- Log reload events for visibility
- Document that reload should be done during low-traffic periods
- Consider implementing graceful connection draining in future if needed

### Risk 3: TLS Configuration Complexity

**Description:** Supporting three TLS modes increases configuration complexity and potential for misconfiguration.

**Mitigation:**
- Clear validation for each mode
- Helpful error messages for misconfiguration
- Documentation examples for each mode
- Default to TLS disabled for simplicity

### Risk 4: Text Logging vs Structured JSON

**Description:** Text logging is human-readable but harder to parse programmatically compared to JSON.

**Trade-off:**
- **Chosen:** Text logging (user requirement)
- **Benefit:** Human-readable, easier to debug manually
- **Drawback:** Harder to parse with log aggregation tools
- **Future:** Could add JSON format option if needed

### Risk 5: No Graceful Shutdown

**Description:** Without graceful shutdown, in-progress requests may be abruptly terminated.

**Mitigation:**
- Document that server should be stopped during low-traffic periods
- Clients should implement retry logic for connection drops
- Consider adding graceful shutdown in future if needed

## Design Decisions (Updated)

### 1. Config File Location: Cross-Platform Support

**Decision:** Support XDG config directories (`~/.config/`) on all platforms (Windows, macOS, Linux).

**Rationale:**
- Follows platform conventions for config file locations
- Windows: `%APPDATA%\ddg-search\config.yaml`
- macOS/Linux: `~/.config/ddg-search/config.yaml`
- Provides better user experience across platforms

### 2. Perplexity Retry Behavior: Immediate Fallback

**Decision:** Fallback immediately on first error without exponential backoff.

**Rationale:**
- Existing Perplexity client already implements retry logic
- Fallback should happen quickly to provide results
- Avoids unnecessary delays for user

### 3. Connection Pooling: Per-Request Connections

**Decision:** Create new connection per request (standard Go HTTP pool behavior).

**Rationale:**
- Go's net/http provides built-in connection pooling
- Simpler implementation
- Sufficient for this use case

### 4. Log Rotation: External Tools

**Decision:** Leave log rotation to external tools (logrotate, journald), use stderr only.

**Rationale:**
- Simpler implementation
- Follows Unix philosophy (do one thing well)
- External tools provide more flexibility
- stderr logging works well with systemd/journald

### 5. Metrics: Prometheus Metrics

**Decision:** Add Prometheus metrics endpoint for monitoring (standard for proxy servers).

**Rationale:**
- Standard practice for proxy servers
- Enables observability
- Out of scope for initial implementation, add in future if needed

### 6. Health Endpoint: Add Health Check

**Decision:** Add `/health` endpoint for load balancer checks.

**Rationale:**
- Required for production deployments
- Simple to implement
- Enables proper load balancer integration

### 7. Concurrent Request Handling: Parallel Execution

**Decision:** Handle concurrent requests in parallel (no rate limiting or queueing).

**Rationale:**
- Simpler implementation
- Go's goroutines handle concurrency well
- Rate limiting adds complexity not required for initial version

### 8. Default Config File: Auto-Create

**Decision:** Create default config file on first run if it doesn't exist.

**Rationale:**
- Better user experience
- Provides working configuration out of the box
- User can still override with custom config

1. **Config File Location:** Should we support XDG config directories (`~/.config/`) on all platforms, or just Unix-like systems?

2. **Perplexity Retry Behavior:** Should we implement exponential backoff for Perplexity rate limits before falling back, or fallback immediately?

3. **Connection Pooling:** For TCP transport, should we implement connection pooling or create new connections per request?

4. **Log Rotation:** Should we implement log rotation for long-running servers, or leave that to external tools?

5. **Metrics:** Should we add Prometheus/OpenTelemetry metrics for monitoring, or is that out of scope?

6. **Health Endpoint:** For TCP transport, should we add a `/health` endpoint for load balancer checks?

7. **Concurrent Request Handling:** How should we handle concurrent requests to the search tools (rate limiting, queueing, parallel execution)?

8. **Default Config File:** Should we create a default config file on first run, or require users to create it manually?

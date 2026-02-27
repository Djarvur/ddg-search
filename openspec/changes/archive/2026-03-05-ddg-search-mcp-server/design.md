## Context

The project currently has three internal libraries for web search and content retrieval:
- `internal/search/` - DuckDuckGo search client with rate limit handling
- `internal/dump/` - Web page fetching and HTML-to-markdown conversion
- `internal/perplexity/` - Perplexity search client

These libraries are used by standalone CLI applications (`ddg-search`, `page-dump`, `perplexity-search`) but are not accessible to Claude Code via MCP. The goal is to build an MCP server that exposes these capabilities as tools compatible with Claude Code's MCP protocol.

**Constraints:**
- Must use Go language
- Must pass `mise lint` (no linter errors)
- Must pass `mise test` (aim for 50% code coverage)
- Must use TDD approach (write tests first)
- Must use `log/slog` for logging
- Must use existing internal libraries without modifications
- Must support stdio and HTTP SSE transports
- Must support TLS/mTLS for HTTP transport

## Goals / Non-Goals

**Goals:**
- Build a Claude Code-compatible MCP server exposing search and fetch tools
- Implement flexible configuration system (config file < env < CLI priority)
- Support both stdio and HTTP SSE transports
- Implement Perplexity search with automatic fallback to DuckDuckGo
- Add TLS and mTLS support for secure HTTP transport
- Ensure comprehensive test coverage (unit, integration, e2e)
- Follow Go best practices and project coding standards

**Non-Goals:**
- Modifying existing internal libraries
- Implementing additional search providers beyond DuckDuckGo and Perplexity
- Building a web UI or dashboard
- Implementing caching or result storage
- Supporting WebSocket transport (only stdio and HTTP SSE)

## Configuration File Structure

The configuration file is located at `~/.config/ddg-search/config.yaml` and uses the following structure:

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

**Environment Variable Mapping:**
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

**Configuration Priority:**
1. CLI parameters (highest priority)
2. Environment variables
3. Configuration file (lowest priority)

## Decisions

### MCP Library Selection

**Decision:** Use `github.com/mark3labs/mcp-go` as the MCP Go library.

**Rationale:**
- Active development and maintenance
- Good documentation and examples
- Supports both stdio and HTTP SSE transports
- Compatible with Claude Code's MCP implementation
- Minimal dependencies

**Alternatives considered:**
- `github.com/modelcontextprotocol/go-sdk` - Less mature, fewer examples
- Custom implementation - Too much complexity, reinventing the wheel

### Configuration Management

**Decision:** Use `github.com/spf13/viper` for configuration management.

**Rationale:**
- Native support for YAML, environment variables, and CLI flags
- Automatic priority handling (config < env < CLI)
- Built-in support for config file watching and reload
- Well-tested and widely used in Go ecosystem
- Integrates well with Cobra for CLI

**Alternatives considered:**
- Custom implementation - More control but more code to maintain
- `github.com/kelseyhightower/envconfig` - Only supports env vars, not config files

### Transport Architecture

**Decision:** Implement transport as an interface with stdio and HTTP SSE implementations.

**Rationale:**
- Clean separation of concerns
- Easy to add new transports in the future
- Testable with mock implementations
- Follows Go interface patterns

**Interface design:**
```go
type Transport interface {
    Start(ctx context.Context) error
    Stop() error
    Send(message []byte) error
    Receive() ([]byte, error)
}
```

### Tool Handler Architecture

**Decision:** Implement tools as handlers with a common interface.

**Rationale:**
- Consistent tool implementation
- Easy to add new tools
- Testable in isolation
- Clear separation between MCP protocol and business logic

**Interface design:**
```go
type ToolHandler interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    Handle(ctx context.Context, params map[string]interface{}) (interface{}, error)
}
```

### Perplexity Fallback Strategy

**Decision:** Implement immediate fallback to DuckDuckGo on any Perplexity error without retries.

**Rationale:**
- Perplexity has rate limits that are not retryable
- Simpler error handling
- Faster response to user
- DuckDuckGo has its own retry logic

**Fallback conditions:**
- No Perplexity access token provided
- Invalid or expired access token
- Rate limit exceeded
- API error
- Network error
- Timeout

### Logging Strategy

**Decision:** Use `log/slog` with structured logging and configurable levels.

**Rationale:**
- Standard library (Go 1.21+)
- Structured logging with key-value pairs
- Configurable log levels
- Good performance
- Easy to integrate with log aggregators

**Log format:**
- Debug: All tool calls with parameters and results
- Info: Server lifecycle events (start, stop, reload)
- Warn: Fallback events, retries
- Error: Failures, invalid configurations

**Important:** All logs must go to stderr. This is required for MCP protocol compatibility and proper integration with AI code editors.

### Testing Strategy

**Decision:** Three-tier testing approach: unit tests, integration tests, e2e tests.

**Rationale:**
- Unit tests for individual components (fast, isolated)
- Integration tests for component interactions (medium speed)
- E2E tests for complete workflows (slower but comprehensive)

**Test organization:**
- Public API tests: Separate test packages (e.g., `config_test` package)
- Internal implementation tests: Separate test packages (e.g., `config_internal_test` package)
- E2E tests: `cmd/ddg-search-mcp/e2e_test.go`

### Signal Handling

**Decision:** Handle SIGINT, SIGTERM for graceful shutdown and SIGHUP for config reload.

**Rationale:**
- Standard Unix signal handling
- Graceful shutdown allows in-flight requests to complete
- Config reload without restart is convenient for production

**Implementation:**
- Use `os/signal` package
- Context cancellation for shutdown
- Config reload on SIGHUP with validation

## Risks / Trade-offs

### Risk: MCP Library Compatibility

**Risk:** The selected MCP library may not be fully compatible with Claude Code's implementation.

**Mitigation:**
- Test with actual Claude Code during development
- Monitor MCP library updates
- Be prepared to switch libraries if needed
- Implement protocol version negotiation

### Risk: Perplexity API Changes

**Risk:** Perplexity API may change, breaking integration.

**Mitigation:**
- Use versioned API endpoints
- Monitor Perplexity API changelog
- Implement graceful degradation (fallback to DuckDuckGo)
- Keep Perplexity client in separate package for easy updates

### Trade-off: No Perplexity Retries

**Trade-off:** Not retrying Perplexity requests may result in suboptimal results for transient errors.

**Mitigation:**
- DuckDuckGo fallback provides reliable results
- Perplexity rate limits are not retryable anyway
- Simpler error handling
- Can add retries later if needed

### Risk: TLS Certificate Management

**Risk:** Managing TLS certificates and rotation adds operational complexity.

**Mitigation:**
- Support HUP signal reload for certificate updates
- Provide clear documentation for certificate setup
- Log certificate expiration warnings
- Consider adding certificate auto-reload in future

### Trade-off: Configuration Complexity

**Trade-off:** Supporting three configuration sources (file, env, CLI) adds complexity.

**Mitigation:**
- Use Viper which handles this complexity
- Clear documentation of priority rules
- Log configuration sources on startup
- Validate configuration early

## Migration Plan

### Deployment Steps

1. **Stage 1-2 (Research):**
   - Research MCP Go libraries and examples
   - Research Claude Code MCP tool calling conventions
   - Document findings

2. **Stage 3 (Internal Library Analysis):**
   - Analyze existing internal libraries
   - Document required parameters and interfaces

3. **Stage 4 (Base Application):**
   - Create `cmd/ddg-search-mcp/` package
   - Implement configuration management
   - Implement signal handling
   - Add logging
   - Write tests

4. **Stage 5 (MCP Server - Mocked):**
   - Implement MCP server with stdio transport
   - Implement search and fetch tools (mocked responses)
   - Add request logging
   - Write tests

5. **Stage 6 (Real Implementation):**
   - Integrate with internal search library
   - Integrate with internal dump library
   - Implement real tool responses
   - Write tests

6. **Stage 7 (Perplexity Integration):**
   - Integrate with internal perplexity library
   - Implement fallback logic
   - Write tests

7. **Stage 8 (HTTP SSE Transport):**
   - Implement HTTP SSE transport
   - Add health check endpoint
   - Write tests

8. **Stage 9 (TLS/mTLS):**
   - Implement TLS support
   - Implement mTLS support
   - Add certificate reload
   - Write tests

### Rollback Strategy

- Keep existing CLI applications unchanged
- MCP server is a new application, no breaking changes
- Can disable MCP server by not deploying it
- Configuration errors prevent startup (fail-fast)

## Open Questions

1. **MCP Library:** Final selection of MCP Go library to be confirmed in Stage 1 research.

2. **Claude Code Compatibility:** Exact MCP protocol version and features supported by Claude Code to be confirmed in Stage 2 research.

3. **Perplexity API Limits:** Exact rate limits and quota details to be documented during implementation.

4. **TLS Configuration:** Specific TLS versions and cipher suites to support (will use Go defaults initially).

5. **Performance Requirements:** Expected request rate and response time targets (not specified, will optimize for typical usage).

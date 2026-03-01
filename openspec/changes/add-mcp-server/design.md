## Context

The ddg-search project currently provides three CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) that expose web search and page fetching capabilities. These tools are not directly accessible through the Model Context Protocol (MCP), which is the standard for integrating tools with AI assistants like Claude Desktop.

The existing codebase has well-structured internal packages:
- `internal/search`: DuckDuckGo search with retry logic
- `internal/perplexity`: Perplexity API client
- `internal/dump`: Page fetching and HTML-to-markdown conversion
- `internal/config`: Configuration types for search options

The project uses `urfave/cli/v3` for the existing CLI tools, but the MCP server requires `cobra+viper` for more flexible configuration management.

## Goals / Non-Goals

**Goals:**
- Create a standalone MCP server binary (`ddg-search-mcp`) that exposes search and fetch capabilities
- Support both stdio and TCP/HTTP transports with runtime configuration
- Implement intelligent Perplexity fallback to DuckDuckGo without errors
- Provide flexible configuration via file, environment variables, and CLI flags
- Support TLS and mTLS for secure connections
- Implement comprehensive logging with slog
- Follow TDD methodology with unit, integration, and E2E tests
- Maintain 50%+ code coverage

**Non-Goals:**
- Modifying existing CLI tools (`ddg-search`, `perplexity-search`, `page-dump`)
- Creating a skill file for the MCP server (no `skills/ddg-search-mcp/`)
- Implementing graceful shutdown (server stops immediately on interrupt)
- Supporting .env files (configuration via config file + env + CLI only)
- Creating resources or prompts (only tools are needed)

## Decisions

### 1. MCP Library: mark3labs/mcp-go

**Decision:** Use `mark3labs/mcp-go` for MCP server implementation.

**Rationale:**
- Simple API with minimal boilerplate
- Supports both stdio and HTTP transports
- Transport-agnostic design allows runtime transport selection
- Good documentation and active development
- Recommended in the project's `mcp-packages.md` research document

**Alternatives considered:**
- `modelcontextprotocol/go-sdk` (official): More verbose API, but official support
- `metoro-io/mcp-golang`: Type-safe with struct-based arguments, but HTTP transport is stateless
- `findleyr/mcp`: No HTTP transport support

### 2. Tool Naming: Simple Names

**Decision:** Use simple tool names matching official MCP style: `search` and `fetch`.

**Rationale:**
- Official MCP examples use simple names (e.g., `fetch` for web fetching)
- More intuitive for AI assistants to use
- Consistent with MCP ecosystem conventions

**Alternatives considered:**
- Descriptive names (`web_search`, `fetch_page`): More explicit but verbose
- Existing skill names (`web-search`, `page-dump`): Hyphenated, less common in Go

### 3. Perplexity Fallback Strategy

**Decision:** The `search` tool tries Perplexity first (if configured) then falls back to DuckDuckGo. No separate `perplexity_search` tool.

**Rationale:**
- Simplifies the API for AI assistants (single search tool)
- Automatic fallback provides better user experience
- Fallback is transparent to the client
- Reduces tool discovery complexity

**Fallback conditions:**
- No Perplexity API key configured → use DDG directly
- Perplexity rate limit exceeded → fallback to DDG
- Perplexity quota exceeded → fallback to DDG
- Perplexity disabled in config → use DDG directly

### 4. Configuration: cobra+viper

**Decision:** Use `cobra+viper` for configuration management despite existing tools using `urfave/cli/v3`.

**Rationale:**
- Explicitly specified in requirements
- Viper provides excellent config file + env + CLI flag integration
- Automatic type conversion and validation
- Widely used in Go projects

**Configuration priority (highest to lowest):**
1. CLI flags
2. Environment variables
3. Configuration file (`~/.config/ddg-search/config.yaml`)
4. Default values

### 5. Output Format: Global Default with Per-Request Override

**Decision:** Support both JSON and plain text output with global default configurable and per-request override.

**Rationale:**
- Flexibility for different use cases
- JSON for programmatic access, text for human-readable output
- Per-request override allows clients to choose format as needed

### 6. TLS Support: Both Separate and Combined Files

**Decision:** Support both separate key/cert files and combined key+cert file, with mTLS via CA cert.

**Rationale:**
- Different deployment scenarios have different file formats
- Combined files are common in some environments (e.g., Kubernetes secrets)
- Separate files are more traditional and easier to manage individually
- mTLS provides additional security for production deployments

**Priority:** Combined file takes precedence if specified, otherwise use separate files.

### 7. Logging: slog with Text Handler

**Decision:** Use standard library `log/slog` with text handler for logging.

**Rationale:**
- Standard library, no additional dependencies
- Structured logging with attributes
- Text handler is human-readable (as specified)
- Consistent with "proxy server" logging style

**Logging behavior:**
- Debug level: Log all requests including successful ones
- Info level: Log info/warn/error, successful requests at debug (not shown)
- Error level: Only log errors
- Bad requests always logged at error level

### 8. Signal Handling: HUP for Config Reload

**Decision:** Handle HUP signal to reload entire configuration file. No graceful shutdown on TERM/INT.

**Rationale:**
- HUP is the standard signal for config reload in Unix systems
- Allows runtime configuration changes without restart
- No graceful shutdown as specified in requirements
- Server exits immediately on interrupt

### 9. Testing Strategy: TDD with Real MCP Client

**Decision:** Follow TDD methodology with unit tests, integration tests, and E2E tests using real MCP client.

**Rationale:**
- TDD ensures code quality from the start
- Real MCP client tests validate actual protocol compliance
- E2E tests cover both transports and signal handling
- Public APIs tested in `_test` package, private in `_internal_test.go`

**Test coverage target:** 50%+

### 10. Binary Location: cmd/ddg-search-mcp/main.go

**Decision:** Place new binary at `cmd/ddg-search-mcp/main.go` parallel to existing commands.

**Rationale:**
- Consistent with existing project structure
- Clear separation of concerns
- Easy to build and install

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    ddg-search-mcp (new binary)                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  cmd/ddg-search-mcp/main.go                              │  │
│  │  - Cobra root command                                    │  │
│  │  - Viper configuration setup                             │  │
│  │  - Signal handling (HUP, INT)                            │  │
│  │  - Transport selection                                   │  │
│  └──────────────────────────────────────────────────────────┘  │
│                              │                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  internal/mcp/ (new package)                             │  │
│  │  - server.go: MCP server setup and tool registration     │  │
│  │  - tools.go: Tool handlers (search, fetch)               │  │
│  │  - config.go: MCP-specific configuration                 │  │
│  │  - logging.go: Logging setup with slog                   │  │
│  │  - tls.go: TLS configuration                             │  │
│  └──────────────────────────────────────────────────────────┘  │
│                              │                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Reused internal packages (no modifications)             │  │
│  │  - internal/search: DuckDuckGo search                    │  │
│  │  - internal/perplexity: Perplexity API                   │  │
│  │  - internal/dump: Page fetching                          │  │
│  │  - internal/config: Configuration types                  │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Package Structure

```
cmd/ddg-search-mcp/
├── main.go                 # Entry point, cobra setup
└── _test/                 # E2E tests

internal/mcp/
├── server.go              # MCP server initialization
├── tools.go               # Tool handlers
├── config.go              # Configuration management
├── logging.go             # Logging setup
├── tls.go                 # TLS configuration
├── server_test.go         # Server tests
├── tools_test.go          # Tool handler tests
├── config_test.go         # Config tests
└── _internal_test.go      # Private function tests
```

### Tool Handler Flow

```
┌─────────────┐
│ MCP Request │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│ Tool Handler (search or fetch)                              │
│ 1. Validate parameters                                      │
│ 2. Apply config defaults                                    │
│ 3. Log request (debug)                                      │
│ 4. Execute tool logic                                       │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│ Search Tool Specific Logic                                  │
│ 1. Check if Perplexity enabled and API key configured       │
│ 2. Try Perplexity search                                    │
│ 3. On error, check if fallback condition                    │
│ 4. If fallback, try DuckDuckGo search                       │
│ 5. Log fallback if occurred                                 │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│ Format Output                                                │
│ 1. Check format parameter (json/text)                       │
│ 2. Apply global default if not specified                    │
│ 3. Format response accordingly                              │
└──────┬──────────────────────────────────────────────────────┘
       │
       ▼
┌─────────────┐
│ MCP Response│
└─────────────┘
```

## Risks / Trade-offs

### Risk 1: Perplexity API Key Exposure
**Risk:** API key stored in config file could be exposed if file permissions are not set correctly.

**Mitigation:**
- Document proper file permissions (0600) for config file
- Support environment variable for API key (more secure)
- Log warning if config file has overly permissive permissions

### Risk 2: Rate Limiting Cascading Failures
**Risk:** If both Perplexity and DuckDuckGo are rate-limited, all searches fail.

**Mitigation:**
- Implement exponential backoff in existing retry logic (already in `internal/search`)
- Log clear error messages indicating both providers failed
- Consider adding a cooldown period after repeated failures

### Risk 3: TLS Configuration Complexity
**Risk:** Supporting both separate and combined key/cert files adds complexity.

**Mitigation:**
- Clear priority: combined file takes precedence
- Comprehensive validation with helpful error messages
- Document all configuration options with examples

### Risk 4: Config Reload Race Conditions
**Risk:** Reloading config on HUP signal could cause race conditions with in-flight requests.

**Mitigation:**
- Use atomic config swap (load new config, then swap pointer)
- Existing requests continue with old config
- New requests use new config
- Log config reload events

### Risk 5: E2E Test Flakiness
**Risk:** E2E tests with real MCP client and signal handling could be flaky.

**Mitigation:**
- Use proper test isolation (unique ports per test)
- Add timeouts to prevent hanging tests
- Mock external services (Perplexity, DuckDuckGo) where possible
- Use test fixtures for predictable responses

### Trade-off 1: No Graceful Shutdown
**Trade-off:** Server exits immediately on interrupt, potentially dropping in-flight requests.

**Rationale:** Specified in requirements. Simpler implementation, acceptable for local tool use.

### Trade-off 2: Single Search Tool
**Trade-off:** No separate Perplexity-only tool, clients cannot explicitly choose provider.

**Rationale:** Simplifies API, automatic fallback is transparent. Clients can disable Perplexity in config if needed.

### Trade-off 3: Config File Location
**Trade-off:** Config file at `~/.config/ddg-search/config.yaml` is shared with potential future tools.

**Rationale:** Standard XDG config location. Could add tool-specific subdirectory if conflicts arise.

## Migration Plan

### Deployment Steps

1. **Add dependencies to go.mod**
   ```bash
   go get github.com/mark3labs/mcp-go
   go get github.com/spf13/cobra
   go get github.com/spf13/viper
   ```

2. **Create new package structure**
   - `cmd/ddg-search-mcp/`
   - `internal/mcp/`

3. **Implement TDD cycle**
   - Write tests first
   - Implement code to pass tests
   - Run `mise lint` after each phase
   - Ensure `mise test` passes

4. **Build and test**
   ```bash
   go build ./cmd/ddg-search-mcp
   mise lint
   mise test
   ```

5. **Integration with Claude Desktop**
   - Add MCP server configuration to Claude Desktop config
   - Test with real Claude Desktop instance

### Rollback Strategy

- Existing CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) are unchanged
- New binary is completely separate
- Rollback is simply not using the new binary
- No database migrations or state changes

## Open Questions

1. **Perplexity Model Selection:** Should we expose all available Perplexity models or limit to a curated list?
   - **Decision:** Expose all models, use `sonar-medium-online` as default

2. **Request ID Generation:** Should we generate request IDs for logging/tracing?
   - **Decision:** Yes, generate UUID for each request for better debugging

3. **Connection Pooling:** Should we implement connection pooling for HTTP client?
   - **Decision:** Use default Go HTTP client behavior, which includes connection pooling

4. **Metrics/Telemetry:** Should we add metrics collection (e.g., request counts, latency)?
   - **Decision:** Out of scope for initial implementation, can be added later if needed

5. **Resource Limits:** Should we add rate limiting or resource limits to prevent abuse?
   - **Decision:** Out of scope for initial implementation, rely on external rate limiting

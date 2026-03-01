# Design: ddg-search-mcp

## Context

The ddg-search project currently provides three CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) that are wrapped as skills for Claude Code. This design introduces a new MCP (Model Context Protocol) server binary that exposes these capabilities as MCP tools, enabling native integration with Claude Code and other MCP clients.

### Current State

- **internal/search**: DuckDuckGo search via HTML scraping with `Searcher` struct providing `Search`, `SearchJSON`, `SearchMarkdown` methods
- **internal/perplexity**: Perplexity API client with `Client.Search` returning `SearchResults` with `Markdown()` formatting
- **internal/dump**: URL fetching and HTML-to-markdown conversion via `FetchAndConvert` function
- **internal/config**: Shared configuration loading (may extend for MCP-specific config)

### Constraints

- Must reuse existing internal packages without modification
- Must pass `go build`, `mise lint`, and `mise test` commands
- Target 50% code coverage minimum
- Follow existing project patterns (cobra/viper from other tools if present)

## Goals / Non-Goals

**Goals:**

- Create standalone `ddg-search-mcp` binary with stdio, TCP, and HTTP transports
- Expose `search` and `fetch` MCP tools using existing internal packages
- Support TLS for TCP/HTTP transport with flexible certificate configuration
- Implement Perplexity → DuckDuckGo fallback with clear metadata
- Support both JSON and text output formats (server default + per-request override)
- Configuration via file + environment variables + CLI flags (cobra + viper)
- Signal handling: SIGINT for stop, SIGHUP for config reload

**Non-Goals:**

- Graceful shutdown with connection draining (just stop on signal)
- .env file support (use config + env + CLI only)
- MCP resources or prompts (tools only for v1)
- Modifications to existing CLI tools

## Decisions

### D1: MCP Library Selection

**Decision:** Use `mark3labs/mcp-go` library.

**Rationale:**
- Mature community library with active maintenance
- Clean Go idiomatic API
- Built-in support for stdio and TCP transports
- Tool definition and handler patterns align with Go conventions

**Alternatives Considered:**
- `modelcontextprotocol/go-sdk`: Official SDK but less mature, more verbose API
- Custom implementation: Too much effort, reinventing the wheel

### D2: Package Structure

**Decision:** Create new `internal/mcp` package with clear separation.

```
internal/mcp/
├── server.go        # MCP server setup and lifecycle
├── tools.go         # Tool definitions (schemas)
├── search.go        # Search tool handler
├── fetch.go         # Fetch tool handler
├── format.go        # Output formatters (text/JSON)
├── config/
│   └── config.go    # MCP-specific configuration
└── *_test.go        # Test files
```

**Rationale:**
- Keeps MCP code isolated from existing packages
- Clear separation of concerns (server, tools, handlers, config)
- Allows independent testing of each component
- Follows existing internal package patterns

**Alternatives Considered:**
- Put everything in `cmd/ddg-search-mcp/main.go`: Too monolithic, hard to test
- Create `pkg/mcp` for exportable code: Not needed, internal is sufficient

### D3: Backend Selection Strategy

**Decision:** Implement backend selection in `search.go` handler with explicit state machine.

```
┌─────────────────┐
│ backend=auto?   │
└────────┬────────┘
         │
    ┌────▼────┐
    │   Yes   │──► Perplexity configured? ──Yes──► Try Perplexity
    └────┬────┘                               │
         │                                    ▼
         │                             ┌──────────────┐
         │                             │   Success?   │
         │                             └──────┬───────┘
         │                                    │
         │                              ┌─────┴─────┐
         │                              │ Yes       │ No
         │                              ▼           ▼
         │                         Return      Fallback to DDG
         │                                     (with note)
         │
         ▼
   backend=perplexity? ──Yes──► Same as above
         │
         │ No
         ▼
    Use DuckDuckGo directly
```

**Rationale:**
- Clear, testable logic
- Fallback always available (DuckDuckGo doesn't require API key)
- Metadata includes fallback information for transparency

### D4: Configuration Hierarchy

**Decision:** Use viper with precedence: CLI flags > env vars > config file.

```go
// Priority order (highest to lowest):
// 1. CLI flags (explicit user choice)
// 2. Environment variables (DDG_MCP_*)
// 3. Config file (~/.config/ddg-search/config.yaml)
// 4. Defaults (hardcoded)
```

**Rationale:**
- Follows 12-factor app principles
- Consistent with cobra/viper best practices
- Allows flexible deployment scenarios

**Alternatives Considered:**
- Config file only: Too rigid for containerized deployments
- Env vars only: No persistent configuration option

### D5: TLS Certificate Configuration

**Decision:** Support three certificate modes:

1. Separate files: `cert_file` + `key_file`
2. Combined file: `cert_key_file` (cert + key in one file)
3. mTLS: above + `ca_file` for client certificate validation

**Rationale:**
- Covers common deployment patterns
- Combined file simplifies single-file configuration
- mTLS support for secure environments

### D6: Output Formatting

**Decision:** Format selection at two levels:

1. Server default: `output.format` config (text or json)
2. Per-request override: `format` parameter in tool calls

**Rationale:**
- Server default for consistent behavior
- Per-request override for flexibility
- Text format is human-readable, JSON is machine-parseable

### D7: Signal Handling

**Decision:**
- **SIGINT**: Immediate stop (no graceful shutdown)
- **SIGHUP**: Reload configuration and apply to new requests

**Rationale:**
- SIGINT behavior matches requirements (no graceful shutdown)
- SIGHUP allows runtime reconfiguration without restart
- Config reload affects new requests; in-flight requests continue with old config

## Risks / Trade-offs

### R1: MCP Specification Changes
- **Risk:** MCP spec evolves, library may lag or break compatibility
- **Mitigation:** Pin library version, monitor releases, abstract library usage behind interfaces

### R2: Perplexity API Rate Limiting
- **Risk:** Perplexity has strict rate limits, may exhaust quickly
- **Mitigation:** Robust fallback to DuckDuckGo with clear metadata; retry with exponential backoff

### R3: HTML Scraping Fragility
- **Risk:** DuckDuckGo HTML structure may change, breaking search
- **Mitigation:** Existing `internal/search` already handles this; MCP layer is insulated

### R4: TLS Configuration Complexity
- **Risk:** Users may misconfigure TLS, leading to security issues
- **Mitigation:** Clear documentation, validation on startup, fail-fast with helpful errors

### R5: Config Reload Race Conditions
- **Risk:** SIGHUP reload may cause race conditions with in-flight requests
- **Mitigation:** Use atomic pointer swap for config; in-flight requests use old config snapshot

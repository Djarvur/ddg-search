# Proposal: Add MCP Server

## Why

Claude Code and other AI assistants need a standardized way to access web search capabilities. Currently, the project provides CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) but lacks an MCP (Model Context Protocol) server that can be easily integrated with Claude Desktop and other MCP-compatible clients. Building an MCP server enables seamless integration with AI assistants while leveraging existing search infrastructure.

## What Changes

- Create new binary `ddg-search-mcp` implementing MCP server using `mark3labs/mcp-go` package
- Implement `search` tool with automatic fallback from Perplexity to DuckDuckGo
- Implement `fetch_page` tool for fetching and converting web pages to markdown
- Support both stdio (for Claude Desktop) and TCP transports
- Add configurable TLS support (key/cert, combined PEM, mTLS with CA cert)
- Implement configuration via cobra+viper (CLI flags → ENV vars → YAML config file)
- Add structured logging using `log/slog` with text handler
- Implement HUP signal handler for config reload with validation
- Support all existing search parameters (max-results, site, region, time, safe-search, model)
- Support JSON and markdown output formats

## Capabilities

### New Capabilities

- `mcp-server`: Model Context Protocol server implementation exposing web search and page fetching capabilities
- `mcp-config`: Flexible configuration system supporting CLI, environment variables, and YAML config files
- `mcp-transport`: Dual transport support (stdio and TCP) with TLS options
- `mcp-logging`: Structured logging for requests, responses, and errors

### Modified Capabilities

None. Existing CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) remain unchanged.

## Impact

- **New Binary**: `cmd/ddg-search-mcp/` - standalone MCP server binary
- **New Packages**:
  - `internal/mcp/` - MCP server implementation
  - `internal/mcp/config/` - MCP-specific configuration
  - `internal/mcp/tools/` - MCP tool handlers
- **New Dependencies**:
  - `github.com/mark3labs/mcp-go` - MCP implementation
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/spf13/viper` - Configuration management
  - `log/slog` - Structured logging (Go 1.21+ standard library)
  - **Existing Code**: Reuse `internal/search/`, `internal/perplexity/`, `internal/dump/` packages
- **Testing**: New test packages for public API (`*_test.go`) and internal tests (`*_internal_test.go`)
- **E2E Tests**: Integration tests for server lifecycle (start, stop, reload, tools)

## Success Criteria

- MCP server starts successfully in both stdio and TCP modes
- `search` tool returns results from Perplexity when API key is configured and enabled
- `search` tool falls back to DuckDuckGo when Perplexity is unavailable or rate-limited
- `search` tool supports all parameters: query, provider, max-results, site, region, time, safe-search, model, format
- `fetch_page` tool successfully fetches and converts web pages to markdown
- Configuration is loaded from CLI flags, environment variables, and YAML config file (in priority order)
- TLS connections work with key/cert, combined PEM, and mTLS modes
- HUP signal triggers config reload with validation
- Logging outputs debug messages for successful requests and error messages for failures
- All tests pass (`mise test`)
- No linter errors remain (`mise run lint`)
- Binary builds successfully (`go build`)

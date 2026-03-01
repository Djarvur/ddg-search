## Why

Claude Code and other AI assistants need a standardized way to access web search and page fetching capabilities through the Model Context Protocol (MCP). Currently, the existing CLI tools (`ddg-search`, `perplexity-search`, `page-dump`) are not directly accessible as MCP tools, requiring manual integration or wrapper scripts. Building a native MCP server enables seamless integration with Claude Desktop and other MCP-compatible clients, providing a consistent interface for web search with intelligent fallback between Perplexity (when available) and DuckDuckGo.

## What Changes

- **New MCP server binary** (`ddg-search-mcp`) that exposes search and page fetching capabilities as MCP tools
- **Two MCP tools**:
  - `search`: Web search that tries Perplexity first (if API key configured) then falls back to DuckDuckGo
  - `fetch`: Fetch web pages and convert to markdown
- **Transport support**: Both stdio (for local CLI integration) and TCP/HTTP (for network access) with configurable selection
- **Configuration system**: Using cobra+viper with config file (`~/.config/ddg-search/config.yaml`), environment variables, and CLI flags
- **TLS support**: Configurable TLS listener with separate or combined key/cert files, and mTLS with CA cert verification
- **Output format**: Configurable JSON or plain text output (global default with per-request override)
- **Signal handling**: Reload entire config on HUP signal
- **Logging**: Using log/slog with text handler, logging all requests (debug for success, error for failures)
- **Comprehensive testing**: TDD approach with unit tests, integration tests, and E2E tests covering both transports and signal handling

## Capabilities

### New Capabilities
- `mcp-server`: Core MCP server implementation with tool registration, transport handling, and configuration management
- `mcp-search-tool`: Web search tool with Perplexity fallback to DuckDuckGo
- `mcp-fetch-tool`: Page fetching and markdown conversion tool
- `mcp-config`: Configuration management with file, environment, and CLI flag support
- `mcp-tls`: TLS and mTLS support for secure connections
- `mcp-logging`: Structured logging with slog for request/response tracking

### Modified Capabilities
None - this is a new standalone binary that reuses existing internal packages without modifying their behavior.

## Impact

- **New binary**: `cmd/ddg-search-mcp/main.go` - new standalone MCP server
- **New dependencies**: `github.com/mark3labs/mcp-go` (MCP server), `github.com/spf13/cobra` (CLI), `github.com/spf13/viper` (config)
- **Reused packages**: `internal/search`, `internal/perplexity`, `internal/dump`, `internal/config` - no modifications required
- **Configuration**: New config file at `~/.config/ddg-search/config.yaml` (separate from existing CLI tools)
- **Testing**: New test packages for MCP server functionality, including E2E tests with real MCP client
- **Build**: New binary target in build system, existing binaries (`ddg-search`, `perplexity-search`, `page-dump`) remain unchanged

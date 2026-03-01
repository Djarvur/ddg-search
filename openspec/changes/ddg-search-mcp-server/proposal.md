## Why

The existing ddg-search CLI tools (ddg-search, perplexity-search, page-dump) are designed for command-line usage but lack integration with AI agent platforms like Claude Code. Building an MCP (Model Context Protocol) server will expose these tools as callable tools for AI agents, enabling seamless web search and page content fetching within agent workflows.

## What Changes

- Create a new MCP server binary (`ddg-search-mcp`) that exposes web search and page fetching as MCP tools
- Support both stdio and TCP transports with configurable selection
- Implement automatic fallback from Perplexity to DuckDuckGo when API limits are exceeded or no API key is configured
- Support configurable TLS for secure TCP connections
- Use existing internal packages for search, Perplexity API client, and page dumping
- Implement comprehensive logging with slog for debugging and monitoring

## Capabilities

### New Capabilities
- `ddg-search-mcp`: MCP server implementation with search and fetch tools, supporting multiple transports (stdio/TCP), TLS, and automatic fallback to DuckDuckGo when Perplexity is unavailable

### Modified Capabilities
None - this is a new capability that leverages existing internal packages without modifying their requirements.

## Impact

**New Code:**
- `cmd/ddg-search-mcp/` - MCP server entry point
- `internal/mcp/` - MCP server implementation with tool handlers
- `internal/mcp/transport/` - Transport layer (stdio, TCP, TLS)
- `internal/mcp/config/` - MCP-specific configuration

**Existing Code (No Changes):**
- `internal/search/` - DuckDuckGo search client (used as-is)
- `internal/perplexity/` - Perplexity API client (used as-is)
- `internal/dump/` - Page dump service (used as-is)
- `internal/config/` - Shared configuration types (used as-is)

**Dependencies:**
- `github.com/mark3labs/mcp-go` - MCP server implementation
- `github.com/spf13/viper` - Config management
- `github.com/spf13/cobra` - CLI interface

**Breaking Changes:** None

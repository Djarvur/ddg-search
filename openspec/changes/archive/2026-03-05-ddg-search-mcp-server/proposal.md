## Why

Claude Code requires MCP (Model Context Protocol) servers to provide tools for web search and content fetching. Currently, the project has internal libraries for DuckDuckGo search, page dumping, and Perplexity search, but no MCP server to expose these capabilities to Claude Code. Building an MCP server will enable seamless integration with Claude Code's tool ecosystem.

## What Changes

- Create a new CLI application `ddg-search-mcp` that implements an MCP server
- Implement configuration management supporting config file, environment variables, and CLI parameters with proper priority (config < env < CLI)
- Implement MCP server with stdio transport initially, then add HTTP SSE transport
- Implement `search` tool supporting DuckDuckGo and Perplexity search with automatic fallback
- Implement `fetch` tool for web page content retrieval
- Add TLS and mTLS support for HTTP transport
- Ensure full compatibility with Claude Code's MCP tool calling conventions

## Capabilities

### New Capabilities
- `mcp-server-core`: Core MCP server implementation with stdio transport, tool registration, and request/response handling
- `mcp-config-management`: Configuration system supporting YAML config file, environment variables, and CLI parameters with HUP signal reload
- `mcp-search-tool`: Search tool implementation supporting DuckDuckGo and Perplexity with automatic fallback behavior
- `mcp-fetch-tool`: Fetch tool for retrieving web page content using existing dump library
- `mcp-perplexity-integration`: Perplexity search integration with token-based authentication and graceful fallback to DuckDuckGo
- `mcp-http-sse`: HTTP Server-Sent Events transport for MCP communication
- `mcp-tls-mtls`: TLS and mutual TLS support for secure HTTP transport

### Modified Capabilities
- None (existing internal libraries will be used as-is without spec changes)

## Impact

- New package: `cmd/ddg-search-mcp/` for the MCP server application
- New internal packages for MCP server logic, configuration, and transport handling
- Dependencies: Add MCP Go library (to be researched in Stage 1)
- Configuration file: `~/.config/ddg-search/config.yaml` for MCP server settings
- No breaking changes to existing code or libraries

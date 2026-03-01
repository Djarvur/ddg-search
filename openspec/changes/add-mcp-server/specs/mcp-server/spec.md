# Spec: MCP Server

Model Context Protocol server implementation exposing web search and page fetching capabilities for Claude Code and other MCP-compatible clients.

## ADDED Requirements

### Requirement: MCP server exposes search tool
The system SHALL expose an MCP tool named `search` that performs web search with automatic provider fallback.

#### Scenario: Search with auto provider and Perplexity available
- **WHEN** client calls `search` tool with `provider="auto"` and Perplexity API key is configured and enabled
- **THEN** system attempts search using Perplexity API first
- **THEN** system returns search results in requested format (json/markdown)
- **THEN** system logs successful request at DEBUG level

#### Scenario: Search with auto provider and Perplexity rate limited
- **WHEN** client calls `search` tool with `provider="auto"` and Perplexity API returns rate limit (HTTP 429) or quota exceeded (HTTP 402)
- **THEN** system automatically falls back to DuckDuckGo search
- **THEN** system logs fallback event at WARN level with reason
- **THEN** system includes `provider_used="duckduckgo"` and `fallback_reason` in response
- **THEN** system returns search results from DuckDuckGo in requested format

#### Scenario: Search with auto provider and no Perplexity configured
- **WHEN** client calls `search` tool with `provider="auto"` and Perplexity is not enabled or API key not configured
- **THEN** system uses DuckDuckGo search directly
- **THEN** system returns search results in requested format
- **THEN** system logs search at DEBUG level

#### Scenario: Search with explicit provider override
- **WHEN** client calls `search` tool with `provider="perplexity"` and Perplexity API key is configured
- **THEN** system uses Perplexity API
- **THEN** system does not fall back to DuckDuckGo even on rate limit
- **THEN** system returns error to client if Perplexity fails

#### Scenario: Search with explicit DuckDuckGo provider
- **WHEN** client calls `search` tool with `provider="duckduckgo"`
- **THEN** system uses DuckDuckGo search
- **THEN** system does not attempt Perplexity API
- **THEN** system returns search results in requested format

#### Scenario: Search with all supported parameters
- **WHEN** client calls `search` tool with query, max-results, site, region, time, safe-search, model, and format parameters
- **THEN** system applies all parameters to the search request
- **THEN** system validates parameters before executing search
- **THEN** system returns results in specified format (json/markdown)

#### Scenario: Search with invalid parameters
- **WHEN** client calls `search` tool with invalid parameters (e.g., negative max-results, invalid format)
- **THEN** system returns error response to client
- **THEN** system logs error at ERROR level with parameter details

### Requirement: MCP server exposes fetch_page tool
The system SHALL expose an MCP tool named `fetch_page` that fetches a web page and converts it to markdown.

#### Scenario: Fetch page with valid URL
- **WHEN** client calls `fetch_page` tool with valid HTTP/HTTPS URL
- **THEN** system fetches the web page content
- **THEN** system converts HTML content to markdown
- **THEN** system returns markdown content to client
- **THEN** system logs successful request at DEBUG level

#### Scenario: Fetch page with invalid URL
- **WHEN** client calls `fetch_page` tool with invalid URL (e.g., missing scheme, invalid format)
- **THEN** system returns error response to client
- **THEN** system logs error at ERROR level with URL details

#### Scenario: Fetch page with custom timeout
- **WHEN** client calls `fetch_page` tool with custom timeout parameter
- **THEN** system applies timeout to the fetch request
- **THEN** system returns error if timeout is exceeded
- **THEN** system logs timeout event at WARN level

#### Scenario: Fetch page with custom user agent
- **WHEN** client calls `fetch_page` tool with custom user-agent parameter
- **THEN** system uses specified user-agent in HTTP request
- **THEN** system logs user-agent at DEBUG level

### Requirement: MCP server supports stdio transport
The system SHALL support stdio transport for Claude Desktop integration.

#### Scenario: Server starts in stdio mode
- **WHEN** server is started with `--transport stdio` or config specifies stdio
- **THEN** server initializes MCP server over stdin/stdout
- **THEN** server does not bind to any network port
- **THEN** server logs transport mode at INFO level
- **THEN** server is ready to receive MCP JSON-RPC messages

#### Scenario: Server receives tool call via stdio
- **WHEN** client sends MCP tool call request via stdin
- **THEN** server processes the request
- **THEN** server executes the appropriate tool handler
- **THEN** server sends response via stdout
- **THEN** server logs request at DEBUG level

### Requirement: MCP server supports TCP transport
The system SHALL support TCP transport for web clients with configurable host and port.

#### Scenario: Server starts in TCP mode with defaults
- **WHEN** server is started with `--transport tcp` without host/port flags
- **THEN** server binds to 0.0.0.0:9100
- **THEN** server logs binding address at INFO level
- **THEN** server is ready to accept TCP connections

#### Scenario: Server starts in TCP mode with custom host/port
- **WHEN** server is started with `--transport tcp --host 127.0.0.1 --port 8080`
- **THEN** server binds to 127.0.0.1:8080
- **THEN** server logs binding address at INFO level
- **THEN** server is ready to accept TCP connections

#### Scenario: Server receives tool call via TCP
- **WHEN** client connects via TCP and sends MCP tool call request
- **THEN** server processes the request
- **THEN** server executes the appropriate tool handler
- **THEN** server sends response via TCP connection
- **THEN** server logs request at DEBUG level

#### Scenario: Server fails to bind TCP port
- **WHEN** server cannot bind to specified host:port (e.g., port in use)
- **THEN** server logs error at ERROR level
- **THEN** server exits with non-zero status code

### Requirement: MCP server supports TLS for TCP transport
The system SHALL support TLS for TCP transport in three modes: key/cert, combined PEM, and mTLS with CA cert.

#### Scenario: Server starts with TLS key/cert mode
- **WHEN** server is started with TLS enabled, mode=keycert, and valid cert/key files
- **THEN** server loads TLS certificate and key from specified files
- **THEN** server accepts TLS connections on TCP port
- **THEN** server logs TLS mode at INFO level

#### Scenario: Server starts with TLS combined PEM mode
- **WHEN** server is started with TLS enabled, mode=combined, and valid combined PEM file
- **THEN** server parses combined PEM file for certificate and key
- **THEN** server accepts TLS connections on TCP port
- **THEN** server logs TLS mode at INFO level

#### Scenario: Server starts with TLS mTLS mode
- **WHEN** server is started with TLS enabled, mode=mtls, and valid cert/key/CA files
- **THEN** server loads client CA certificate for mutual authentication
- **THEN** server requires client certificate for connections
- **THEN** server logs TLS mode at INFO level

#### Scenario: Server starts with invalid TLS configuration
- **WHEN** server is started with TLS enabled but invalid cert/key files (missing, unreadable)
- **THEN** server logs error at ERROR level
- **THEN** server exits with non-zero status code
- **THEN** server does not start

### Requirement: MCP server implements HUP signal handler for config reload
The system SHALL handle SIGHUP signal to reload configuration from file.

#### Scenario: Server receives HUP signal
- **WHEN** running server receives SIGHUP signal
- **THEN** server logs reload event at INFO level
- **THEN** server reads configuration file from disk
- **THEN** server validates new configuration
- **THEN** server applies new configuration if validation passes
- **THEN** server logs successful reload at INFO level

#### Scenario: Config reload with invalid configuration
- **WHEN** server receives HUP signal and new configuration is invalid
- **THEN** server logs validation errors at ERROR level
- **THEN** server keeps old configuration active
- **THEN** server continues serving with old configuration

#### Scenario: Config reload with TLS change
- **WHEN** server receives HUP signal and TLS configuration changes
- **THEN** server validates new TLS configuration
- **THEN** server restarts TLS listener with new configuration if valid
- **THEN** server logs TLS listener restart at INFO level

#### Scenario: Server stops on Ctrl-C
- **WHEN** user sends SIGINT (Ctrl-C) to running server
- **THEN** server logs shutdown event at INFO level
- **THEN** server closes all connections
- **THEN** server exits with zero status code

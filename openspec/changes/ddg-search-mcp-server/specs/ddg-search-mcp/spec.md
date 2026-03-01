## ADDED Requirements

### Requirement: MCP Server Binary
The system SHALL build a binary named `ddg-search-mcp` that implements the Model Context Protocol server.

#### Scenario: Binary builds successfully
- **WHEN** running `go build -o ddg-search-mcp ./cmd/ddg-search-mcp`
- **THEN** the binary is created and executable

### Requirement: Stdio Transport
The system SHALL support stdio transport for local AI agent integration.

#### Scenario: Stdio transport starts successfully
- **WHEN** the server is started with `transport: stdio` configuration
- **THEN** the server reads JSON-RPC requests from stdin and writes responses to stdout

#### Scenario: Stdio transport handles multiple requests
- **WHEN** multiple JSON-RPC requests are sent via stdin
- **THEN** each request is processed and a response is written to stdout

### Requirement: TCP Transport
The system SHALL support TCP transport for remote server deployment.

#### Scenario: TCP transport starts successfully
- **WHEN** the server is started with `transport: tcp` configuration
- **THEN** the server listens on the configured port (default: 9100)

#### Scenario: TCP transport handles multiple connections
- **WHEN** multiple clients connect to the server
- **THEN** each connection is handled concurrently

### Requirement: TLS Support
The system SHALL support TLS encryption for TCP transport.

#### Scenario: TLS with combined cert+key file
- **WHEN** `server.tls.enabled` is true and `server.tls.cert_key_combined` is true
- **THEN** the server uses the combined cert+key file for TLS

#### Scenario: TLS with separate cert and key files
- **WHEN** `server.tls.enabled` is true and `server.tls.cert_key_combined` is false
- **THEN** the server uses separate cert and key files for TLS

#### Scenario: mTLS with CA cert
- **WHEN** `server.tls.ca_file` is configured
- **THEN** the server requires client certificates signed by the CA

### Requirement: Configuration Management
The system SHALL support configuration via config file, environment variables, and CLI flags.

#### Scenario: Config file is loaded
- **WHEN** `~/.config/ddg-search/config.yaml` exists
- **THEN** the server loads configuration from the file

#### Scenario: Environment variables override config
- **WHEN** environment variable `MCP_PORT` is set
- **THEN** the port from env var overrides the config file value

#### Scenario: CLI flags override all
- **WHEN** CLI flag `--port` is provided
- **THEN** the port from CLI flag overrides all other sources

### Requirement: Search Tool
The system SHALL expose a `search` tool that performs web search.

#### Scenario: Search with DDG
- **WHEN** the `search` tool is called with `source: "ddg"` or when Perplexity is unavailable
- **THEN** the server performs a DuckDuckGo search and returns results

#### Scenario: Search with Perplexity
- **WHEN** the `search` tool is called with `source: "perplexity"` and API key is configured
- **THEN** the server performs a Perplexity search and returns results

#### Scenario: Automatic fallback to DDG
- **WHEN** Perplexity API returns rate limit or quota exceeded error
- **THEN** the server falls back to DDG and includes `fallback_used: true` in the response

#### Scenario: Search with output format
- **WHEN** the `search` tool is called with `output_format: "markdown"`
- **THEN** the response includes markdown-formatted results

### Requirement: Fetch Tool
The system SHALL expose a `fetch` tool that fetches web pages and converts to markdown.

#### Scenario: Fetch successful page
- **WHEN** the `fetch` tool is called with a valid URL
- **THEN** the server fetches the page and returns markdown content

#### Scenario: Fetch with custom timeout
- **WHEN** the `fetch` tool is called with `timeout` parameter
- **THEN** the server uses the specified timeout for the request

### Requirement: Logging
The system SHALL log all requests and responses for debugging.

#### Scenario: Debug logging
- **WHEN** `logging.level` is set to "debug"
- **THEN** all requests and responses are logged

#### Scenario: Error logging
- **WHEN** an error occurs during request processing
- **THEN** the error is logged with appropriate context

### Requirement: Config Reload on SIGHUP
The system SHALL reload all configuration when receiving SIGHUP signal.

#### Scenario: Config reload
- **WHEN** the server receives SIGHUP signal
- **THEN** the server reloads all configuration from the config file

### Requirement: Immediate Shutdown on SIGINT
The system SHALL shut down immediately when receiving SIGINT signal.

#### Scenario: Shutdown on ctrl-c
- **WHEN** the server receives SIGINT signal (ctrl-c)
- **THEN** the server closes all connections immediately

### Requirement: Health Check Endpoint
The system SHALL provide a `/health` endpoint for TCP transport.

#### Scenario: Health check
- **WHEN** a client sends a request to `/health` endpoint
- **THEN** the server responds with `{"status": "ok"}`

### Requirement: JSON Output
The system SHALL support JSON output format for all tools.

#### Scenario: Search with JSON output
- **WHEN** the `search` tool is called with `output_format: "json"`
- **THEN** the response includes JSON-formatted results

#### Scenario: Fetch with JSON output
- **WHEN** the `fetch` tool is called with `output_format: "json"`
- **THEN** the response includes JSON-formatted content

### Requirement: Markdown Output
The system SHALL support markdown output format for the fetch tool.

#### Scenario: Fetch with markdown output
- **WHEN** the `fetch` tool is called with `output_format: "markdown"`
- **THEN** the response includes markdown-formatted content

### Requirement: Text Output
The system SHALL support plain text output format.

#### Scenario: Search with text output
- **WHEN** the `search` tool is called with `output_format: "text"`
- **THEN** the response includes plain text-formatted results

### Requirement: Source Parameter
The `search` tool SHALL support a `source` parameter to explicitly select the search provider.

#### Scenario: Search with auto source
- **WHEN** the `search` tool is called with `source: "auto"` (default)
- **THEN** the server automatically selects the best provider (Perplexity if available, DDG otherwise)

#### Scenario: Search with explicit DDG
- **WHEN** the `search` tool is called with `source: "ddg"`
- **THEN** the server performs a DuckDuckGo search

#### Scenario: Search with explicit Perplexity
- **WHEN** the `search` tool is called with `source: "perplexity"`
- **THEN** the server performs a Perplexity search

### Requirement: Default Output Format
The system SHALL support configurable default output format for the search tool.

#### Scenario: Search with default output format
- **WHEN** the `search` tool is called without `output_format` parameter
- **THEN** the server uses the configured default output format

#### Scenario: Search with per-request output format
- **WHEN** the `search` tool is called with `output_format` parameter
- **THEN** the server uses the per-request format instead of default

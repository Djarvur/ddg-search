## ADDED Requirements

### Requirement: MCP server starts with stdio transport
The system SHALL provide an MCP server that communicates via stdio by default.

#### Scenario: Start server in stdio mode
- **WHEN** server is started without transport configuration
- **THEN** server SHALL listen on stdin/stdout for MCP protocol messages

#### Scenario: Server responds to initialize request
- **WHEN** client sends MCP initialize request
- **THEN** server SHALL respond with server capabilities including available tools

### Requirement: MCP server supports TCP transport
The system SHALL support TCP transport for network-based communication.

#### Scenario: Start server in TCP mode
- **WHEN** server is started with `transport: tcp` and `port: 9100`
- **THEN** server SHALL listen on TCP port 9100 for MCP protocol messages

#### Scenario: TCP server accepts connections
- **WHEN** client connects to TCP port and sends MCP request
- **THEN** server SHALL accept the connection and respond to MCP messages

### Requirement: MCP server supports HTTP transport
The system SHALL support HTTP transport for web-based communication.

#### Scenario: Start server in HTTP mode
- **WHEN** server is started with `transport: http` and `port: 8080`
- **THEN** server SHALL listen on HTTP port 8080 for MCP protocol messages

#### Scenario: HTTP server accepts connections
- **WHEN** client sends HTTP request to the server
- **THEN** server SHALL accept the connection and respond to MCP messages

### Requirement: TLS support for TCP transport
The system SHALL support TLS encryption for TCP transport.

#### Scenario: TLS with separate cert and key files
- **WHEN** server is configured with `tls.enabled: true`, `tls.cert_file`, and `tls.key_file`
- **THEN** server SHALL accept TLS connections using the provided certificate

#### Scenario: TLS with combined cert-key file
- **WHEN** server is configured with `tls.enabled: true` and `tls.cert_key_file`
- **THEN** server SHALL accept TLS connections using the combined certificate file

#### Scenario: mTLS with client certificate validation
- **WHEN** server is configured with `tls.enabled: true`, certificate, and `tls.ca_file`
- **THEN** server SHALL require and validate client certificates

### Requirement: Configuration via multiple sources
The system SHALL support configuration from config file, environment variables, and CLI flags.

#### Scenario: Load configuration from file
- **WHEN** config file exists at `~/.config/ddg-search/config.yaml`
- **THEN** server SHALL load configuration values from the file

#### Scenario: Override configuration with environment variables
- **WHEN** environment variable `DDG_MCP_SERVER_PORT` is set to `9200`
- **THEN** server SHALL use port 9200 regardless of config file value

#### Scenario: Override configuration with CLI flags
- **WHEN** CLI flag `--port=9300` is provided
- **THEN** server SHALL use port 9300 regardless of config file or environment variable

### Requirement: Configuration precedence follows standard order
The system SHALL apply configuration with precedence: CLI flags > environment variables > config file > defaults.

#### Scenario: CLI flag takes precedence over environment variable
- **WHEN** `DDG_MCP_SERVER_PORT=9200` and `--port=9300` are both provided
- **THEN** server SHALL use port 9300

### Requirement: Signal handling for server lifecycle
The system SHALL respond to UNIX signals for server control.

#### Scenario: SIGINT stops the server
- **WHEN** SIGINT signal is received
- **THEN** server SHALL stop immediately

#### Scenario: SIGHUP reloads configuration
- **WHEN** SIGHUP signal is received
- **THEN** server SHALL reload configuration from file
- **AND** new requests SHALL use the updated configuration

### Requirement: Server exposes tools list
The system SHALL expose available MCP tools to clients.

#### Scenario: List tools request
- **WHEN** client sends MCP tools/list request
- **THEN** server SHALL return list containing `search` and `fetch` tools with their schemas

### Requirement: Output format configuration
The system SHALL support configurable output formats.

#### Scenario: Default text output format
- **WHEN** `output.format: text` is configured
- **THEN** tool responses SHALL be formatted as human-readable text by default

#### Scenario: Default JSON output format
- **WHEN** `output.format: json` is configured
- **THEN** tool responses SHALL be formatted as JSON by default

#### Scenario: Per-request format override
- **WHEN** tool request includes `format: json` parameter
- **THEN** response SHALL be formatted as JSON regardless of server default

# Spec: MCP Transport

Dual transport support (stdio and TCP) with TLS options for MCP server.

## ADDED Requirements

### Requirement: Transport supports stdio mode for Claude Desktop
The system SHALL support stdio transport for compatibility with Claude Desktop and other MCP clients.

#### Scenario: Server starts in stdio mode
- **WHEN** server is started with `--transport stdio` or config specifies stdio
- **THEN** server initializes MCP server over stdin/stdout
- **THEN** server does not bind to any network port
- **THEN** server is ready to receive MCP JSON-RPC messages
- **THEN** server logs transport mode at INFO level

#### Scenario: Server receives tool call via stdio
- **WHEN** client sends MCP tool call request via stdin
- **THEN** server processes request
- **THEN** server executes appropriate tool handler
- **THEN** server sends response via stdout
- **THEN** server logs request at DEBUG level with tool name and parameters

#### Scenario: Server receives initialize request via stdio
- **WHEN** client sends MCP initialize request via stdin
- **THEN** server responds with server capabilities
- **THEN** server responds with supported tools (search, fetch_page)
- **THEN** server logs initialization at INFO level

#### Scenario: Server receives shutdown request via stdio
- **WHEN** client sends MCP shutdown request or stdin closes
- **THEN** server logs shutdown event at INFO level
- **THEN** server closes gracefully
- **THEN** server exits with zero status code

### Requirement: Transport supports TCP mode for web clients
The system SHALL support TCP transport for web clients and remote access with configurable host and port.

#### Scenario: Server starts in TCP mode with defaults
- **WHEN** server is started with `--transport tcp` without host/port flags
- **THEN** server binds to 0.0.0.0:9100
- **THEN** server logs binding address at INFO level
- **THEN** server is ready to accept TCP connections
- **THEN** server accepts connections from any interface (0.0.0.0)

#### Scenario: Server starts in TCP mode with custom host/port
- **WHEN** server is started with `--transport tcp --host 127.0.0.1 --port 8080`
- **THEN** server binds to 127.0.0.1:8080
- **THEN** server logs binding address at INFO level
- **THEN** server accepts connections only on specified interface
- **THEN** server is ready to accept TCP connections

#### Scenario: Server receives tool call via TCP
- **WHEN** client connects via TCP and sends MCP tool call request
- **THEN** server processes request
- **THEN** server executes appropriate tool handler
- **THEN** server sends response via TCP connection
- **THEN** server logs request at DEBUG level with client address, tool name, and parameters

#### Scenario: Server receives initialize request via TCP
- **WHEN** client connects via TCP and sends MCP initialize request
- **THEN** server responds with server capabilities
- **THEN** server responds with supported tools (search, fetch_page)
- **THEN** server logs client connection at INFO level

#### Scenario: Server handles multiple concurrent TCP connections
- **WHEN** multiple clients connect via TCP simultaneously
- **THEN** server handles each connection independently
- **THEN** server does not block connections waiting for others
- **THEN** server logs concurrent connections at DEBUG level

#### Scenario: Server fails to bind TCP port
- **WHEN** server cannot bind to specified host:port (e.g., port in use, permission denied)
- **THEN** server logs error at ERROR level with binding failure details
- **THEN** server exits with non-zero status code
- **THEN** server does not start

### Requirement: Transport supports TLS for TCP mode
The system SHALL support TLS for TCP transport in three modes: key/cert files, combined PEM file, and mTLS with CA certificate.

#### Scenario: Server starts with TLS key/cert mode
- **WHEN** server is started with TLS enabled, mode=keycert, and valid cert/key files
- **THEN** server loads TLS certificate and key from specified files
- **THEN** server accepts TLS connections on TCP port
- **THEN** server logs TLS mode at INFO level
- **THEN** server requires HTTPS connections

#### Scenario: Server starts with TLS combined PEM mode
- **WHEN** server is started with TLS enabled, mode=combined, and valid combined PEM file
- **THEN** server parses combined PEM file for certificate and key
- **THEN** server accepts TLS connections on TCP port
- **THEN** server logs TLS mode at INFO level
- **THEN** server requires HTTPS connections

#### Scenario: Server starts with TLS mTLS mode
- **WHEN** server is started with TLS enabled, mode=mtls, and valid cert/key/CA files
- **THEN** server loads client CA certificate for mutual authentication
- **THEN** server accepts TLS connections on TCP port
- **THEN** server validates client certificate against CA
- **THEN** server logs TLS mode at INFO level
- **THEN** server requires HTTPS connections with valid client certificate

#### Scenario: Server starts with invalid TLS configuration
- **WHEN** server is started with TLS enabled but invalid cert/key files (missing, unreadable, invalid format)
- **THEN** server logs error at ERROR level with file details
- **THEN** server exits with non-zero status code
- **THEN** server does not start

#### Scenario: Server receives TLS connection with valid certificate
- **WHEN** client connects with valid TLS certificate
- **THEN** server accepts connection
- **THEN** server logs TLS handshake success at DEBUG level
- **THEN** server processes requests normally

#### Scenario: Server receives TLS connection with invalid certificate (mTLS)
- **WHEN** client connects with invalid or missing TLS certificate in mTLS mode
- **THEN** server rejects connection
- **THEN** server logs TLS handshake failure at WARN level
- **THEN** server does not process requests from rejected connection

#### Scenario: Server receives non-TLS connection when TLS enabled
- **WHEN** client connects via HTTP when server requires HTTPS
- **THEN** server rejects connection
- **THEN** server logs protocol mismatch at WARN level
- **THEN** server does not process requests from rejected connection

### Requirement: Transport supports runtime mode selection
The system SHALL support runtime selection between stdio and TCP modes via CLI flag or configuration.

#### Scenario: Transport mode specified via CLI flag
- **WHEN** user starts server with `--transport stdio` or `--transport tcp`
- **THEN** system uses specified transport mode
- **THEN** system ignores transport setting in config file
- **THEN** system logs transport mode at INFO level

#### Scenario: Transport mode specified via environment variable
- **WHEN** `DDG_SEARCH_TRANSPORT` environment variable is set
- **THEN** system uses specified transport mode
- **THEN** system logs transport mode at INFO level

#### Scenario: Transport mode specified via config file
- **WHEN** config file specifies transport mode and no CLI flag or ENV variable
- **THEN** system uses transport mode from config file
- **THEN** system logs transport mode at INFO level

#### Scenario: Invalid transport mode specified
- **WHEN** user specifies invalid transport mode (e.g., "invalid")
- **THEN** system logs error at ERROR level
- **THEN** system exits with non-zero status code
- **THEN** system displays valid options (stdio, tcp)

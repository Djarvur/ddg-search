## ADDED Requirements

### Requirement: Streamable HTTP transport initialization
The system SHALL initialize a Streamable HTTP transport when configured with HTTP protocol.

#### Scenario: HTTP transport enabled
- **WHEN** the protocol is set to "http" in configuration
- **THEN** the system starts an HTTP server
- **AND** the system listens on the configured bind address
- **AND** the system supports Streamable HTTP for MCP communication

#### Scenario: Default bind address
- **WHEN** HTTP transport is enabled
- **AND** no bind address is specified in configuration
- **THEN** the system listens on "localhost:9100"

#### Scenario: Custom bind address
- **WHEN** HTTP transport is enabled
- **AND** a custom bind address is specified in configuration
- **THEN** the system listens on the specified address

### Requirement: MCP endpoint handling
The system SHALL handle MCP connections on a single /mcp endpoint.

#### Scenario: Client connects to MCP endpoint with POST
- **WHEN** a client sends a POST request to the /mcp endpoint
- **THEN** the system processes the JSON-RPC request
- **AND** the system returns the response
- **AND** the system includes MCP-Protocol-Version header in response

#### Scenario: Client connects to MCP endpoint with GET
- **WHEN** a client sends a GET request to the /mcp endpoint with Accept: text/event-stream
- **THEN** the system establishes an SSE connection
- **AND** the system sends initialization messages
- **AND** the system advertises available tools

#### Scenario: Client disconnects
- **WHEN** a client disconnects from the /mcp endpoint
- **THEN** the system cleans up the connection
- **AND** the system logs the disconnection

### Requirement: Protocol configuration
The system SHALL support configurable transport protocol (stdio or http).

#### Scenario: Protocol set to stdio
- **WHEN** the protocol is set to "stdio" in configuration
- **THEN** the system uses stdio transport
- **AND** the system does not start an HTTP server

#### Scenario: Protocol set to http
- **WHEN** the protocol is set to "http" in configuration
- **THEN** the system uses Streamable HTTP transport
- **AND** the system starts an HTTP server
- **AND** the system exposes a single /mcp endpoint

#### Scenario: Default protocol
- **WHEN** no protocol is specified in configuration
- **THEN** the system uses stdio transport as default

### Requirement: HTTP server lifecycle
The system SHALL manage HTTP server lifecycle including startup and shutdown.

#### Scenario: HTTP server starts successfully
- **WHEN** the application starts with HTTP protocol
- **THEN** the HTTP server starts
- **AND** the server logs the bind address and /mcp endpoint
- **AND** the server is ready to accept connections

#### Scenario: HTTP server fails to start
- **WHEN** the application starts with HTTP protocol
- **AND** the bind address is already in use
- **THEN** the server logs an error
- **AND** the application exits with a non-zero status

#### Scenario: HTTP server shuts down gracefully
- **WHEN** the application receives a shutdown signal
- **THEN** the HTTP server stops accepting new connections
- **AND** the server completes in-flight requests
- **AND** the server exits cleanly

### Requirement: Streamable HTTP message format
The system SHALL format messages according to MCP protocol specification.

#### Scenario: Tool call request via POST
- **WHEN** a client sends a tool call request via POST to /mcp
- **THEN** the system receives the JSON-RPC request
- **AND** the system parses the request correctly

#### Scenario: Tool call response via POST
- **WHEN** the system sends a tool call response
- **THEN** the response is formatted according to MCP specification
- **AND** the response includes the tool name and result

#### Scenario: Error response via POST
- **WHEN** the system sends an error response
- **THEN** the error is formatted according to MCP specification
- **AND** the error includes error details

#### Scenario: Server notification via SSE
- **WHEN** the system sends a notification via SSE (GET /mcp)
- **THEN** the notification is formatted according to MCP specification
- **AND** the notification includes the event data

### Requirement: HTTP request logging
The system SHALL log HTTP requests in Combined Log Format at debug level.

**Combined Log Format:**
```
[timestamp] [client] [method] [path] [status] [bytes] [referer] [user-agent]
```

#### Scenario: HTTP request is logged
- **WHEN** an HTTP request is received
- **THEN** the system logs the request at debug level in Combined Log Format
- **AND** the log includes timestamp, client, method, path, status, bytes, referer, and user-agent

#### Scenario: MCP endpoint request is logged
- **WHEN** a request to /mcp is received
- **THEN** the system logs the request at debug level in Combined Log Format
- **AND** the log includes timestamp, client, method, path, status, and bytes

### Requirement: Multiple concurrent connections
The system SHALL support multiple concurrent MCP connections.

#### Scenario: Multiple clients connect
- **WHEN** multiple clients connect to the /mcp endpoint
- **THEN** the system handles all connections concurrently
- **AND** each connection receives its own responses

#### Scenario: One client disconnects
- **WHEN** one client disconnects
- **THEN** the system continues serving other clients
- **AND** the disconnected client's resources are cleaned up

### Requirement: HTTP health check
The system SHALL provide a health check endpoint for HTTP transport.

#### Scenario: Health check request
- **WHEN** a client requests the /health endpoint
- **THEN** the system returns a 200 OK status
- **AND** the response indicates the server is healthy

#### Scenario: Health check during shutdown
- **WHEN** the server is shutting down
- **AND** a health check request is received
- **THEN** the system returns a 503 Service Unavailable status

### Requirement: MCP protocol version header
The system SHALL support the MCP-Protocol-Version HTTP header.

#### Scenario: Client sends protocol version header
- **WHEN** a client sends a request with MCP-Protocol-Version header
- **THEN** the system uses the specified protocol version for processing
- **AND** the system includes the MCP-Protocol-Version header in responses

#### Scenario: Client omits protocol version header
- **WHEN** a client sends a request without MCP-Protocol-Version header
- **THEN** the system assumes protocol version 2025-03-26 for backwards compatibility

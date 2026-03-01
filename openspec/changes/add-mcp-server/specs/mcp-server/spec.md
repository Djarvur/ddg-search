## ADDED Requirements

### Requirement: Server initialization
The system SHALL initialize an MCP server with a name and version using the mark3labs/mcp-go library.

#### Scenario: Server starts with stdio transport
- **WHEN** the server is started without transport configuration
- **THEN** the server SHALL use stdio transport by default
- **AND** the server SHALL listen on stdin/stdout for MCP protocol messages

#### Scenario: Server starts with TCP transport
- **WHEN** the server is started with transport=tcp configuration
- **THEN** the server SHALL use TCP/HTTP transport
- **AND** the server SHALL listen on the configured port (default: 9100)

### Requirement: Tool registration
The system SHALL register all available tools with the MCP server during initialization.

#### Scenario: Tools are registered on startup
- **WHEN** the server starts
- **THEN** the `search` tool SHALL be registered
- **AND** the `fetch` tool SHALL be registered
- **AND** each tool SHALL have a name, description, and input schema

### Requirement: Transport selection
The system SHALL support configurable transport selection via config file, environment variable, or CLI flag.

#### Scenario: Transport selected from config file
- **WHEN** config file specifies transport=tcp
- **THEN** the server SHALL use TCP transport
- **AND** the server SHALL listen on the configured port

#### Scenario: Transport selected from environment variable
- **WHEN** MCP_TRANSPORT environment variable is set to "stdio"
- **THEN** the server SHALL use stdio transport
- **AND** config file transport setting SHALL be overridden

#### Scenario: Transport selected from CLI flag
- **WHEN** --transport flag is set to "tcp"
- **THEN** the server SHALL use TCP transport
- **AND** both config file and environment variable settings SHALL be overridden

### Requirement: Signal handling
The system SHALL handle the HUP signal to reload configuration.

#### Scenario: Config reload on HUP signal
- **WHEN** the server receives a HUP signal
- **THEN** the server SHALL reload the entire configuration file
- **AND** the server SHALL apply new configuration without restarting
- **AND** existing connections SHALL remain active

#### Scenario: Server stops on interrupt
- **WHEN** the server receives SIGINT (Ctrl-C)
- **THEN** the server SHALL stop accepting new connections
- **AND** the server SHALL exit immediately

### Requirement: Server capabilities
The system SHALL advertise its capabilities to MCP clients during initialization.

#### Scenario: Server advertises tools capability
- **WHEN** an MCP client connects
- **THEN** the server SHALL advertise tools capability
- **AND** the server SHALL list all available tools

### Requirement: Error handling
The system SHALL handle errors gracefully and return appropriate MCP error responses.

#### Scenario: Invalid tool call returns error
- **WHEN** a client calls a non-existent tool
- **THEN** the server SHALL return an error response
- **AND** the error SHALL include a descriptive message

#### Scenario: Tool execution error returns error
- **WHEN** a tool execution fails
- **THEN** the server SHALL return an error response
- **AND** the error SHALL include details about the failure

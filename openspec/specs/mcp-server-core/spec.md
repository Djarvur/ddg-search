## ADDED Requirements

### Requirement: MCP server initialization
The system SHALL initialize an MCP server that communicates via stdio transport and registers available tools.

#### Scenario: Server starts successfully
- **WHEN** the ddg-search-mcp application is started
- **THEN** the MCP server initializes with stdio transport
- **AND** the server sends an initialization message to the client
- **AND** the server advertises available tools (search, fetch)

#### Scenario: Server handles initialization request
- **WHEN** the client sends an initialize request
- **THEN** the server responds with server capabilities
- **AND** the server includes supported protocol version
- **AND** the server includes available tools list

### Requirement: Tool registration
The system SHALL register tools with their names, descriptions, and input schemas.

#### Scenario: Search tool is registered
- **WHEN** the server initializes
- **THEN** the search tool is registered with name "search"
- **AND** the tool includes a description
- **AND** the tool includes an input schema defining query, max_results, and other parameters

#### Scenario: Fetch tool is registered
- **WHEN** the server initializes
- **THEN** the fetch tool is registered with name "fetch"
- **AND** the tool includes a description
- **AND** the tool includes an input schema defining url parameter

### Requirement: Tool call handling
The system SHALL handle tool call requests and return appropriate responses.

#### Scenario: Valid tool call is processed
- **WHEN** the client sends a tool call request for a registered tool
- **THEN** the server processes the request
- **AND** the server returns a tool response with the result
- **AND** the response includes the tool name and content

#### Scenario: Unknown tool is called
- **WHEN** the client sends a tool call request for an unknown tool
- **THEN** the server returns an error response
- **AND** the error indicates the tool was not found

#### Scenario: Invalid tool parameters
- **WHEN** the client sends a tool call with invalid parameters
- **THEN** the server returns an error response
- **AND** the error indicates parameter validation failure

### Requirement: Request logging
The system SHALL log all tool calls in Combined Log Format at debug level.

**Important:** All logs MUST go to stderr. This is required for MCP protocol compatibility and proper integration with AI code editors.

**Combined Log Format:**
```
[timestamp] [client] [method] [path] [status] [bytes] [referer] [user-agent]
```

For stdio transport, client will be "stdio" and method/path will be the tool name.

#### Scenario: Successful tool call is logged
- **WHEN** a tool call is successfully processed
- **THEN** the server logs the call at debug level
- **AND** the log includes timestamp, tool name, parameters, and result status

#### Scenario: Failed tool call is logged
- **WHEN** a tool call fails
- **THEN** the server logs the call at debug level
- **AND** the log includes timestamp, tool name, parameters, and error details

### Requirement: Graceful shutdown
The system SHALL handle shutdown signals and terminate gracefully.

#### Scenario: Server receives interrupt signal
- **WHEN** the server receives SIGINT (Ctrl-C)
- **THEN** the server stops accepting new requests
- **AND** the server completes in-flight requests
- **AND** the server exits cleanly

#### Scenario: Server receives termination signal
- **WHEN** the server receives SIGTERM
- **THEN** the server stops accepting new requests
- **AND** the server completes in-flight requests
- **AND** the server exits cleanly

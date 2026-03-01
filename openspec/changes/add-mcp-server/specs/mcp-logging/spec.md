## ADDED Requirements

### Requirement: Logging with slog
The system SHALL use the standard library log/slog package for logging.

#### Scenario: Logger initialization
- **WHEN** the server starts
- **THEN** the system SHALL initialize a slog logger
- **AND** use the text handler as specified

### Requirement: Log level configuration
The system SHALL support configurable log levels via config file, environment variable, or CLI flag.

#### Scenario: Debug log level
- **WHEN** log level is set to "debug"
- **THEN** the system SHALL log all messages including debug
- **AND** log successful tool calls at debug level

#### Scenario: Info log level
- **WHEN** log level is set to "info"
- **THEN** the system SHALL log info, warning, and error messages
- **AND** log successful tool calls at debug level (not shown)

#### Scenario: Warn log level
- **WHEN** log level is set to "warn"
- **THEN** the system SHALL log warning and error messages
- **AND** suppress info and debug messages

#### Scenario: Error log level
- **WHEN** log level is set to "error"
- **THEN** the system SHALL log only error messages
- **AND** suppress info, warning, and debug messages

### Requirement: Request logging
The system SHALL log all incoming tool requests.

#### Scenario: Successful request logged at debug level
- **WHEN** a tool request succeeds
- **AND** log level is "debug"
- **THEN** the system SHALL log the request
- **AND** include tool name, parameters, and duration

#### Scenario: Successful request not logged at info level
- **WHEN** a tool request succeeds
- **AND** log level is "info" or higher
- **THEN** the system SHALL NOT log the successful request

#### Scenario: Failed request logged at error level
- **WHEN** a tool request fails
- **THEN** the system SHALL log the error
- **AND** include tool name, parameters, error message, and duration

### Requirement: Bad request logging
The system SHALL log bad requests (invalid parameters, missing required fields) at error level.

#### Scenario: Missing required parameter
- **WHEN** a request is missing a required parameter
- **THEN** the system SHALL log an error
- **AND** include the tool name and missing parameter

#### Scenario: Invalid parameter value
- **WHEN** a request has an invalid parameter value
- **THEN** the system SHALL log an error
- **AND** include the tool name, parameter name, and invalid value

#### Scenario: Invalid tool name
- **WHEN** a request specifies a non-existent tool
- **THEN** the system SHALL log an error
- **AND** include the invalid tool name

### Requirement: Log format
The system SHALL use the slog text handler for log output.

#### Scenario: Text log format
- **WHEN** logging any message
- **THEN** the output SHALL be in text format
- **AND** include timestamp, level, and message
- **AND** include any structured attributes

### Requirement: Log attributes
The system SHALL include relevant structured attributes in log messages.

#### Scenario: Tool request log attributes
- **WHEN** logging a tool request
- **THEN** the log SHALL include attributes: tool, params, duration_ms
- **AND** include request_id if available

#### Scenario: Error log attributes
- **WHEN** logging an error
- **THEN** the log SHALL include attributes: error, tool (if applicable), params (if applicable)
- **AND** include stack trace if available

#### Scenario: Server startup log attributes
- **WHEN** logging server startup
- **THEN** the log SHALL include attributes: transport, port (if TCP), version

### Requirement: Configuration reload logging
The system SHALL log configuration reload events.

#### Scenario: Successful config reload
- **WHEN** configuration is successfully reloaded on HUP signal
- **THEN** the system SHALL log an info message
- **AND** indicate configuration was reloaded

#### Scenario: Failed config reload
- **WHEN** configuration reload fails on HUP signal
- **THEN** the system SHALL log an error
- **AND** include the reason for failure

### Requirement: Server lifecycle logging
The system SHALL log server lifecycle events.

#### Scenario: Server startup
- **WHEN** the server starts
- **THEN** the system SHALL log an info message
- **AND** include server name, version, and transport

#### Scenario: Server shutdown
- **WHEN** the server shuts down
- **THEN** the system SHALL log an info message
- **AND** indicate the server is stopping

#### Scenario: Signal received
- **WHEN** the server receives a signal (HUP, INT, TERM)
- **THEN** the system SHALL log a debug message
- **AND** include the signal name

### Requirement: Perplexity fallback logging
The system SHALL log when Perplexity search falls back to DuckDuckGo.

#### Scenario: Perplexity rate limit fallback
- **WHEN** Perplexity search fails with rate limit
- **THEN** the system SHALL log a warning
- **AND** indicate "fallback to DuckDuckGo: Perplexity rate limit exceeded"

#### Scenario: Perplexity quota exceeded fallback
- **WHEN** Perplexity search fails with quota exceeded
- **THEN** the system SHALL log a warning
- **AND** indicate "fallback to DuckDuckGo: Perplexity quota exceeded"

#### Scenario: No Perplexity API key
- **WHEN** Perplexity API key is not configured
- **THEN** the system SHALL log a debug message
- **AND** indicate "using DuckDuckGo: no Perplexity API key configured"

### Requirement: TLS logging
The system SHALL log TLS-related events.

#### Scenario: TLS enabled
- **WHEN** TLS is enabled
- **THEN** the system SHALL log an info message
- **AND** indicate TLS is enabled

#### Scenario: mTLS enabled
- **WHEN** mTLS is enabled with CA certificate
- **THEN** the system SHALL log an info message
- **AND** indicate mTLS is enabled

#### Scenario: TLS handshake error
- **WHEN** TLS handshake fails
- **THEN** the system SHALL log an error
- **AND** include details about the failure

### Requirement: Log output destination
The system SHALL write logs to stderr.

#### Scenario: Logs to stderr
- **WHEN** the server is running
- **THEN** all log messages SHALL be written to stderr
- **AND** stdout SHALL be reserved for MCP protocol messages (stdio transport)

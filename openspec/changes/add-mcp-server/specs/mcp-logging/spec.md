# Spec: MCP Logging

Structured logging using log/slog with text handler for requests, responses, and errors.

## ADDED Requirements

### Requirement: Logging uses log/slog with text handler
The system SHALL use Go's log/slog package with text handler for all logging output.

#### Scenario: Server starts with default log level
- **WHEN** server starts without log level configuration
- **THEN** system uses INFO level as default
- **THEN** system logs startup message at INFO level
- **THEN** log output is human-readable with level prefix

#### Scenario: Server starts with debug log level
- **WHEN** server starts with log level configured to debug
- **THEN** system logs all messages at DEBUG level and above
- **THEN** system logs detailed information including request parameters
- **THEN** system logs configured level at INFO level

#### Scenario: Server starts with error log level
- **WHEN** server starts with log level configured to error
- **THEN** system logs only ERROR level messages
- **THEN** system does not log INFO, WARN, or DEBUG messages
- **THEN** system logs configured level at INFO level

### Requirement: Logging logs successful requests at DEBUG level
The system SHALL log successful tool requests at DEBUG level as specified for proxy server behavior.

#### Scenario: Search tool request succeeds
- **WHEN** search tool completes successfully
- **THEN** system logs request at DEBUG level
- **THEN** log includes: query, provider, results count, duration
- **THEN** log format: `[DEBUG] Search completed (provider=duckduckgo, results=5, duration=1.2s)`

#### Scenario: Fetch page tool request succeeds
- **WHEN** fetch_page tool completes successfully
- **THEN** system logs request at DEBUG level
- **THEN** log includes: URL, content length, duration
- **THEN** log format: `[DEBUG] Page fetched (url=https://example.com, length=1234, duration=0.5s)`

#### Scenario: Server initialization succeeds
- **WHEN** server initializes successfully
- **THEN** system logs initialization at INFO level
- **THEN** log includes: transport mode, binding address, TLS status
- **THEN** log format: `[INFO] Starting MCP server (transport=tcp, host=0.0.0.0, port=9100)`

### Requirement: Logging logs bad requests at ERROR level
The system SHALL log failed or invalid requests at ERROR level as specified for proxy server behavior.

#### Scenario: Search tool request fails with invalid parameters
- **WHEN** search tool receives invalid parameters (e.g., negative max_results, invalid format)
- **THEN** system logs error at ERROR level
- **THEN** log includes: parameter name, value, validation error
- **THEN** log format: `[ERROR] Invalid parameter (name=max_results, value=-5, error: must be positive)`

#### Scenario: Search tool request fails with provider error
- **WHEN** search tool fails due to provider error (network, timeout, API error)
- **THEN** system logs error at ERROR level
- **THEN** log includes: provider, error type, error message
- **THEN** log format: `[ERROR] Search failed (provider=perplexity, error=network timeout)`

#### Scenario: Fetch page tool request fails with invalid URL
- **WHEN** fetch_page tool receives invalid URL (missing scheme, invalid format)
- **THEN** system logs error at ERROR level
- **THEN** log includes: URL, validation error
- **THEN** log format: `[ERROR] Invalid URL format (url=htp://example.com, error: unsupported scheme)`

#### Scenario: Fetch page tool request fails with HTTP error
- **WHEN** fetch_page tool fails due to HTTP error (404, 500, timeout)
- **THEN** system logs error at ERROR level
- **THEN** log includes: URL, HTTP status, error message
- **THEN** log format: `[ERROR] HTTP error (url=https://example.com, status=404, error=Not Found)`

### Requirement: Logging logs fallback events at WARN level
The system SHALL log Perplexity to DuckDuckGo fallback events at WARN level.

#### Scenario: Search falls back to DuckDuckGo due to rate limit
- **WHEN** Perplexity API returns rate limit (HTTP 429) and system falls back to DuckDuckGo
- **THEN** system logs fallback at WARN level
- **THEN** log includes: original provider, fallback provider, reason
- **THEN** log format: `[WARN] Falling back to DuckDuckGo (from=perplexity, reason=rate limit exceeded)`

#### Scenario: Search falls back to DuckDuckGo due to quota exceeded
- **WHEN** Perplexity API returns payment required (HTTP 402) and system falls back to DuckDuckGo
- **THEN** system logs fallback at WARN level
- **THEN** log includes: original provider, fallback provider, reason
- **THEN** log format: `[WARN] Falling back to DuckDuckGo (from=perplexity, reason=quota exceeded)`

#### Scenario: Search falls back to DuckDuckGo due to no API key
- **WHEN** Perplexity is not enabled or API key not configured
- **THEN** system logs fallback at INFO level
- **THEN** log includes: using provider directly
- **THEN** log format: `[INFO] Using DuckDuckGo directly (perplexity not configured)`

### Requirement: Logging logs config reload events
The system SHALL log configuration reload events triggered by HUP signal.

#### Scenario: Config reload triggered by HUP signal
- **WHEN** server receives SIGHUP signal
- **THEN** system logs reload event at INFO level
- **THEN** log format: `[INFO] Config reload triggered (SIGHUP)`

#### Scenario: Config reload succeeds
- **WHEN** config file is reloaded and validated successfully
- **THEN** system logs success at INFO level
- **THEN** log format: `[INFO] Config reloaded and validated (file=~/.config/ddg-search/config.yaml)`

#### Scenario: Config reload fails validation
- **WHEN** config file is reloaded but validation fails
- **THEN** system logs error at ERROR level
- **THEN** log includes: validation error, file path
- **THEN** log format: `[ERROR] Config validation failed (error=invalid port, file=~/.config/ddg-search/config.yaml)`

#### Scenario: Config reload with TLS change
- **WHEN** config reload changes TLS configuration
- **THEN** system logs TLS listener restart at INFO level
- **THEN** log format: `[INFO] TLS listener restarted (mode=keycert)`

### Requirement: Logging logs server lifecycle events
The system SHALL log server lifecycle events (start, stop, shutdown) at appropriate levels.

#### Scenario: Server starts successfully
- **WHEN** server starts and is ready to accept connections
- **THEN** system logs ready event at INFO level
- **THEN** log format: `[INFO] Server ready (transport=tcp, address=0.0.0.0:9100)`

#### Scenario: Server stops on Ctrl-C
- **WHEN** user sends SIGINT (Ctrl-C) to running server
- **THEN** system logs shutdown event at INFO level
- **THEN** log format: `[INFO] Server shutting down (signal=SIGINT)`

#### Scenario: Server stops on error
- **WHEN** server encounters fatal error and must stop
- **THEN** system logs shutdown event at ERROR level
- **THEN** log includes: error message, stack trace if available
- **THEN** log format: `[ERROR] Server shutting down (error=binding failed: address already in use)`

### Requirement: Logging format is human-readable text
The system SHALL output logs in human-readable text format with level prefixes.

#### Scenario: Log output format
- **WHEN** system outputs any log message
- **THEN** format is: `[LEVEL] message`
- **THEN** LEVEL is one of: DEBUG, INFO, WARN, ERROR
- **THEN** message is descriptive and includes relevant context

#### Scenario: Log includes contextual information
- **WHEN** system logs events
- **THEN** log includes relevant context (tool name, parameters, duration, error details)
- **THEN** context helps with debugging and troubleshooting
- **THEN** sensitive data (API keys) is not logged

### Requirement: Logging does not expose sensitive information
The system SHALL not log sensitive information such as API keys or passwords.

#### Scenario: Perplexity API key in configuration
- **WHEN** Perplexity API key is configured
- **THEN** system does not log the API key value
- **THEN** system logs that API key is configured (without value)
- **THEN** log format: `[INFO] Perplexity API key configured`

#### Scenario: TLS key/cert file paths in configuration
- **WHEN** TLS key/cert files are configured
- **THEN** system logs file paths but not file contents
- **THEN** log format: `[INFO] TLS key file: /path/to/key.pem`

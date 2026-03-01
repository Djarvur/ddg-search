# Spec: MCP Configuration

Flexible configuration system supporting CLI flags, environment variables, and YAML config files with priority chain.

## ADDED Requirements

### Requirement: Configuration loads from multiple sources with priority
The system SHALL load configuration from three sources in priority order: CLI flags → environment variables → YAML config file.

#### Scenario: CLI flag overrides config file
- **WHEN** user starts server with `--port 8080` flag and config file specifies port 9100
- **THEN** system uses port 8080 (CLI flag has highest priority)
- **THEN** system logs configuration source at DEBUG level

#### Scenario: Environment variable overrides config file
- **WHEN** user sets `DDG_SEARCH_PORT=8080` environment variable and config file specifies port 9100
- **THEN** system uses port 8080 (ENV has higher priority than config file)
- **THEN** system logs configuration source at DEBUG level

#### Scenario: Config file provides defaults
- **WHEN** user starts server without CLI flags or environment variables
- **THEN** system loads configuration from YAML config file
- **THEN** system uses default values for missing config keys
- **THEN** system logs configuration file path at INFO level

#### Scenario: Config file not found
- **WHEN** user starts server and config file does not exist at default location
- **THEN** system uses default configuration values
- **THEN** system logs warning at WARN level about missing config file
- **THEN** system continues with defaults

### Requirement: Configuration supports server settings
The system SHALL support server configuration including transport, host, port, and TLS settings.

#### Scenario: Server transport configuration
- **WHEN** config file specifies `server.transport: stdio`
- **THEN** system starts in stdio mode
- **THEN** system ignores host and port settings
- **THEN** system logs transport mode at INFO level

#### Scenario: Server TCP configuration
- **WHEN** config file specifies `server.transport: tcp` with host and port
- **THEN** system binds to specified host:port
- **THEN** system logs binding address at INFO level

#### Scenario: Server TLS enabled configuration
- **WHEN** config file specifies `server.tls.enabled: true`
- **THEN** system enables TLS for TCP transport
- **THEN** system validates TLS configuration before starting
- **THEN** system logs TLS mode at INFO level

#### Scenario: Server TLS disabled configuration
- **WHEN** config file specifies `server.tls.enabled: false` or TLS section is missing
- **THEN** system starts TCP transport without TLS
- **THEN** system accepts plain HTTP connections

### Requirement: Configuration supports Perplexity settings
The system SHALL support Perplexity configuration including enabled flag, API key, model, and max results.

#### Scenario: Perplexity enabled with API key
- **WHEN** config file specifies `perplexity.enabled: true` and `perplexity.api_key` is set
- **THEN** system uses Perplexity API for search when provider is auto or perplexity
- **THEN** system does not log API key (security)
- **THEN** system logs Perplexity enabled at INFO level

#### Scenario: Perplexity disabled
- **WHEN** config file specifies `perplexity.enabled: false` or Perplexity section is missing
- **THEN** system does not use Perplexity API
- **THEN** system uses DuckDuckGo for all search requests
- **THEN** system logs Perplexity disabled at INFO level

#### Scenario: Perplexity API key from environment variable
- **WHEN** `PERPLEXITY_API_KEY` environment variable is set and config file API key is empty
- **THEN** system uses API key from environment variable
- **THEN** system does not log API key (security)

#### Scenario: Perplexity model configuration
- **WHEN** config file specifies `perplexity.model: sonar-medium-online`
- **THEN** system uses specified model for Perplexity API requests
- **THEN** system validates model name against supported models
- **THEN** system logs model at DEBUG level

### Requirement: Configuration supports search settings
The system SHALL support search configuration including max results, region, safe search, and time filter.

#### Scenario: Search max results configuration
- **WHEN** config file specifies `search.max_results: 10`
- **THEN** system limits search results to 10 unless overridden by tool parameter
- **THEN** system validates max_results is positive integer
- **THEN** system logs max_results at DEBUG level

#### Scenario: Search region configuration
- **WHEN** config file specifies `search.region: us-en`
- **THEN** system uses specified region for DuckDuckGo search
- **THEN** system validates region format
- **THEN** system logs region at DEBUG level

#### Scenario: Search safe search configuration
- **WHEN** config file specifies `search.safe_search: true`
- **THEN** system enables safe search for DuckDuckGo
- **THEN** system logs safe_search at DEBUG level

#### Scenario: Search time filter configuration
- **WHEN** config file specifies `search.time: d` (day)
- **THEN** system applies time filter to DuckDuckGo search
- **THEN** system validates time filter value (d, w, m, y)
- **THEN** system logs time filter at DEBUG level

### Requirement: Configuration supports dump settings
The system SHALL support dump configuration including timeout and user agent.

#### Scenario: Dump timeout configuration
- **WHEN** config file specifies `dump.timeout: 30s`
- **THEN** system uses 30 second timeout for page fetch
- **THEN** system validates timeout is positive duration
- **THEN** system logs timeout at DEBUG level

#### Scenario: Dump user agent configuration
- **WHEN** config file specifies `dump.user_agent: ddg-search-mcp/1.0`
- **THEN** system uses specified user agent for HTTP requests
- **THEN** system logs user agent at DEBUG level

### Requirement: Configuration supports logging settings
The system SHALL support logging configuration including level and format.

#### Scenario: Logging level configuration
- **WHEN** config file specifies `logging.level: debug`
- **THEN** system logs all messages at DEBUG level and above
- **THEN** system validates level is one of: debug, info, warn, error
- **THEN** system logs configured level at INFO level

#### Scenario: Logging format configuration
- **WHEN** config file specifies `logging.format: text`
- **THEN** system uses text handler for slog output
- **THEN** system logs are human-readable with level prefix
- **THEN** system validates format is supported (text)

### Requirement: Configuration supports custom config file path
The system SHALL support custom config file path via `--config` CLI flag or `DDG_SEARCH_CONFIG` environment variable.

#### Scenario: Custom config file specified via CLI flag
- **WHEN** user starts server with `--config /path/to/config.yaml`
- **THEN** system loads configuration from specified file
- **THEN** system logs custom config path at INFO level
- **THEN** system ignores default config location

#### Scenario: Custom config file specified via environment variable
- **WHEN** `DDG_SEARCH_CONFIG` environment variable is set
- **THEN** system loads configuration from specified file
- **THEN** system logs custom config path at INFO level

#### Scenario: Default config file location
- **WHEN** user starts server without custom config path
- **THEN** system loads configuration from `~/.config/ddg-search/config.yaml`
- **THEN** system creates config directory if it does not exist
- **THEN** system logs default config path at DEBUG level

### Requirement: Configuration validates before applying
The system SHALL validate configuration values before applying them to prevent invalid states.

#### Scenario: Valid configuration
- **WHEN** configuration file contains valid values for all settings
- **THEN** system applies configuration without errors
- **THEN** system logs configuration loaded successfully at INFO level

#### Scenario: Invalid TLS configuration
- **WHEN** configuration specifies TLS enabled but cert/key files do not exist
- **THEN** system logs validation error at ERROR level
- **THEN** system exits with non-zero status code
- **THEN** system does not start with invalid configuration

#### Scenario: Invalid port configuration
- **WHEN** configuration specifies port outside valid range (1-65535)
- **THEN** system logs validation error at ERROR level
- **THEN** system exits with non-zero status code
- **THEN** system does not start with invalid configuration

#### Scenario: Invalid log level configuration
- **WHEN** configuration specifies invalid log level (e.g., "invalid")
- **THEN** system logs validation error at ERROR level
- **THEN** system uses default log level (info)
- **THEN** system continues with default level

## ADDED Requirements

### Requirement: Configuration file location
The system SHALL load configuration from `~/.config/ddg-search/config.yaml` by default.

#### Scenario: Config file exists
- **WHEN** the config file exists at the default location
- **THEN** the system SHALL load the configuration from that file

#### Scenario: Config file does not exist
- **WHEN** the config file does not exist at the default location
- **THEN** the system SHALL use default configuration values
- **AND** the system SHALL NOT return an error

#### Scenario: Custom config file location
- **WHEN** a custom config file path is specified via CLI flag
- **THEN** the system SHALL load configuration from the specified path
- **AND** the default location SHALL be ignored

### Requirement: Configuration priority
The system SHALL apply configuration values in the following priority order (highest to lowest):
1. CLI flags
2. Environment variables
3. Configuration file
4. Default values

#### Scenario: CLI flag overrides config file
- **WHEN** a value is set in both config file and CLI flag
- **THEN** the CLI flag value SHALL be used
- **AND** the config file value SHALL be ignored

#### Scenario: Environment variable overrides config file
- **WHEN** a value is set in both config file and environment variable
- **THEN** the environment variable value SHALL be used
- **AND** the config file value SHALL be ignored

#### Scenario: CLI flag overrides environment variable
- **WHEN** a value is set in both environment variable and CLI flag
- **THEN** the CLI flag value SHALL be used
- **AND** the environment variable value SHALL be ignored

#### Scenario: Config file value used when no override
- **WHEN** a value is set only in config file
- **THEN** the config file value SHALL be used

#### Scenario: Default value used when not configured
- **WHEN** a value is not set in config file, environment, or CLI
- **THEN** the default value SHALL be used

### Requirement: Configuration structure
The configuration file SHALL support the following sections:
- `server`: Server configuration (transport, port, TLS)
- `search`: Search tool defaults (max_results, region, time, safe_search)
- `perplexity`: Perplexity configuration (api_key, model, enabled)
- `fetch`: Fetch tool defaults (timeout, user_agent)
- `output`: Output configuration (format)
- `logging`: Logging configuration (level)

#### Scenario: Valid configuration file
- **WHEN** the config file contains all sections
- **THEN** the system SHALL load all configuration values
- **AND** apply them to the appropriate components

#### Scenario: Partial configuration file
- **WHEN** the config file contains only some sections
- **THEN** the system SHALL load the configured sections
- **AND** use default values for missing sections

### Requirement: Server configuration
The system SHALL support the following server configuration options:
- `transport`: Transport type ("stdio" or "tcp", default: "stdio")
- `port`: TCP port number (default: 9100)
- `tls.enabled`: Enable TLS (default: false)
- `tls.cert_file`: Path to TLS certificate file
- `tls.key_file`: Path to TLS key file
- `tls.combined_file`: Path to combined key+cert file
- `tls.ca_file`: Path to CA certificate for mTLS

#### Scenario: Default server configuration
- **WHEN** no server configuration is provided
- **THEN** the system SHALL use stdio transport
- **AND** TLS SHALL be disabled

#### Scenario: TCP transport configuration
- **WHEN** transport is set to "tcp"
- **THEN** the system SHALL use TCP transport
- **AND** listen on the configured port

#### Scenario: TLS enabled with separate files
- **WHEN** TLS is enabled with cert_file and key_file
- **THEN** the system SHALL load the certificate and key from separate files
- **AND** use them for TLS connections

#### Scenario: TLS enabled with combined file
- **WHEN** TLS is enabled with combined_file
- **THEN** the system SHALL load the certificate and key from the combined file
- **AND** use them for TLS connections

#### Scenario: mTLS with CA cert
- **WHEN** TLS is enabled with ca_file
- **THEN** the system SHALL verify client certificates using the CA cert
- **AND** reject connections without valid client certificates

### Requirement: Perplexity configuration
The system SHALL support the following Perplexity configuration options:
- `api_key`: Perplexity API key
- `model`: Default Perplexity model (default: "sonar-medium-online")
- `enabled`: Enable Perplexity search (default: true)

#### Scenario: Perplexity API key from config
- **WHEN** api_key is set in config file
- **THEN** the system SHALL use that API key for Perplexity requests

#### Scenario: Perplexity API key from environment
- **WHEN** PERPLEXITY_API_KEY environment variable is set
- **THEN** the system SHALL use that API key for Perplexity requests
- **AND** override any config file value

#### Scenario: Perplexity disabled
- **WHEN** enabled is set to false
- **THEN** the system SHALL not attempt Perplexity search
- **AND** always use DuckDuckGo

### Requirement: Output configuration
The system SHALL support the following output configuration options:
- `format`: Default output format ("json" or "text", default: "text")

#### Scenario: Default output format
- **WHEN** no output format is specified
- **THEN** the system SHALL use "text" format by default

#### Scenario: JSON output format configured
- **WHEN** format is set to "json" in config
- **THEN** the system SHALL return JSON format by default
- **AND** tools can still override per-request

### Requirement: Logging configuration
The system SHALL support the following logging configuration options:
- `level`: Log level ("debug", "info", "warn", "error", default: "info")

#### Scenario: Debug logging
- **WHEN** level is set to "debug"
- **THEN** the system SHALL log all requests including successful ones
- **AND** include detailed debug information

#### Scenario: Info logging
- **WHEN** level is set to "info"
- **THEN** the system SHALL log informational messages
- **AND** log successful requests at debug level

#### Scenario: Error logging
- **WHEN** level is set to "error"
- **THEN** the system SHALL only log errors
- **AND** suppress info and debug messages

### Requirement: Configuration validation
The system SHALL validate configuration values on load.

#### Scenario: Invalid transport value
- **WHEN** transport is set to an invalid value
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate valid transport options

#### Scenario: Invalid port number
- **WHEN** port is set to an invalid value (e.g., negative or > 65535)
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate valid port range

#### Scenario: Invalid log level
- **WHEN** level is set to an invalid value
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate valid log levels

#### Scenario: Missing TLS files when enabled
- **WHEN** TLS is enabled but cert/key files are not specified
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate required TLS files

### Requirement: Configuration reload
The system SHALL support reloading configuration on HUP signal.

#### Scenario: Successful config reload
- **WHEN** the server receives HUP signal
- **AND** the config file is valid
- **THEN** the system SHALL reload the configuration
- **AND** apply new values to all components

#### Scenario: Config reload with invalid file
- **WHEN** the server receives HUP signal
- **AND** the config file is invalid
- **THEN** the system SHALL log an error
- **AND** continue using the previous configuration

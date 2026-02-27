## ADDED Requirements

### Requirement: Configuration file loading
The system SHALL load configuration from a YAML file at `~/.config/ddg-search/config.yaml`.

#### Scenario: Config file exists and is valid
- **WHEN** the application starts
- **AND** a valid config file exists at `~/.config/ddg-search/config.yaml`
- **THEN** the system loads the configuration
- **AND** the system logs the loaded configuration

#### Scenario: Config file does not exist
- **WHEN** the application starts
- **AND** no config file exists at `~/.config/ddg-search/config.yaml`
- **THEN** the system uses default values
- **AND** the system logs that no config file was found

#### Scenario: Config file is invalid
- **WHEN** the application starts
- **AND** the config file exists but contains invalid YAML
- **THEN** the system logs an error
- **AND** the system exits with a non-zero status

### Requirement: Environment variable configuration
The system SHALL support configuration via environment variables with the prefix `DDG_SEARCH_`.

#### Scenario: Environment variable overrides config file
- **WHEN** a parameter is set in both config file and environment variable
- **THEN** the environment variable value takes precedence
- **AND** the system logs the override

#### Scenario: Environment variable provides default
- **WHEN** a parameter is not in config file
- **AND** the parameter is set via environment variable
- **THEN** the system uses the environment variable value

### Requirement: CLI parameter configuration
The system SHALL support configuration via CLI parameters.

#### Scenario: CLI parameter overrides environment variable
- **WHEN** a parameter is set via both environment variable and CLI parameter
- **THEN** the CLI parameter value takes precedence
- **AND** the system logs the override

#### Scenario: CLI parameter provides default
- **WHEN** a parameter is not in config file or environment
- **AND** the parameter is set via CLI parameter
- **THEN** the system uses the CLI parameter value

### Requirement: Configuration priority
The system SHALL apply configuration in priority order: config file (lowest) < environment variables < CLI parameters (highest).

#### Scenario: All three sources provide values
- **WHEN** a parameter is set in config file, environment variable, and CLI parameter
- **THEN** the CLI parameter value is used
- **AND** the system logs the final value and its source

#### Scenario: Config and env provide values
- **WHEN** a parameter is set in config file and environment variable
- **THEN** the environment variable value is used
- **AND** the system logs the final value and its source

### Requirement: Configuration reload on HUP signal
The system SHALL reload the configuration file when receiving SIGHUP signal.

#### Scenario: HUP signal triggers reload
- **WHEN** the server receives SIGHUP signal
- **THEN** the system reloads the configuration file
- **AND** the system logs the reloaded configuration
- **AND** the server continues running with new configuration

#### Scenario: HUP with invalid config
- **WHEN** the server receives SIGHUP signal
- **AND** the config file is now invalid
- **THEN** the system logs an error
- **AND** the system continues with previous configuration

### Requirement: Configuration logging
The system SHALL log the complete configuration on startup and reload.

#### Scenario: Configuration logged on startup
- **WHEN** the application starts
- **THEN** the system logs all configuration values
- **AND** the system logs the source of each value (default, config, env, CLI)

#### Scenario: Configuration logged on reload
- **WHEN** the configuration is reloaded
- **THEN** the system logs all configuration values
- **AND** the system logs the source of each value

### Requirement: Log level configuration
The system SHALL support configurable log level via config file, environment variable, or CLI parameter.

#### Scenario: Log level set via config file
- **WHEN** log level is specified in config file
- **THEN** the system uses the specified log level
- **AND** the system logs at the configured level

#### Scenario: Log level set via environment variable
- **WHEN** log level is specified via DDG_SEARCH_LOG_LEVEL environment variable
- **THEN** the system uses the specified log level
- **AND** the system logs at the configured level

#### Scenario: Log level set via CLI parameter
- **WHEN** log level is specified via --log-level CLI parameter
- **THEN** the system uses the specified log level
- **AND** the system logs at the configured level

### Requirement: Configuration structure
The system SHALL support a configuration structure with parameters for search, fetch, and server settings.

#### Scenario: Search parameters in config
- **WHEN** config file contains search parameters
- **THEN** the system loads max_results, safe_search, and other search options

#### Scenario: Perplexity parameters in config
- **WHEN** config file contains perplexity parameters
- **THEN** the system loads access_token and enabled flag

#### Scenario: Server parameters in config
- **WHEN** config file contains server parameters
- **THEN** the system loads protocol, bind_address, and TLS settings

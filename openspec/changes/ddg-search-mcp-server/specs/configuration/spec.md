# Configuration Specification

## Overview

This specification defines the configuration system for the MCP server, including file-based configuration, environment variables, and CLI flags.

## ADDED Requirements

### Requirement: Configuration File Support
The server SHALL load configuration from a YAML file.

#### Scenario: Load config from default location
- **WHEN** server starts without --config flag
- **THEN** it SHALL load config from ~/.config/ddg-search/config.yaml

#### Scenario: Load config from custom location
- **WHEN** server starts with --config /path/to/config.yaml
- **THEN** it SHALL load config from the specified path

#### Scenario: Config file does not exist
- **WHEN** the config file does not exist
- **THEN** the server SHALL use default values for all settings

### Requirement: Environment Variable Support
The server SHALL override config file settings with environment variables.

#### Scenario: Transport set via environment variable
- **WHEN** DDG_MCP_TRANSPORT=http is set
- **THEN** the server SHALL use HTTP transport regardless of config file

#### Scenario: Port set via environment variable
- **WHEN** DDG_MCP_PORT=8080 is set
- **THEN** the server SHALL listen on port 8080

#### Scenario: Perplexity API key from environment
- **WHEN** DDG_PERPLEXITY_API_KEY is set
- **THEN** the server SHALL use this API key for Perplexity requests

### Requirement: CLI Flag Support
The server SHALL override all other settings with CLI flags.

#### Scenario: Transport set via CLI flag
- **WHEN** --transport http is specified
- **THEN** the server SHALL use HTTP transport

#### Scenario: Port set via CLI flag
- **WHEN** --port 3000 is specified
- **THEN** the server SHALL listen on port 3000

### Requirement: Configuration Precedence
CLI flags SHALL take highest precedence, followed by environment variables, then config file.

#### Scenario: All sources specify different port
- **WHEN** config file has port 9100, env has DDG_MCP_PORT=8080, CLI has --port 3000
- **THEN** the server SHALL use port 3000

### Requirement: Server Configuration Options
The server SHALL support the following server configuration options.

#### Scenario: Transport option
- **WHEN** server.transport is set to "stdio" or "http"
- **THEN** the server SHALL use the specified transport mode

#### Scenario: Port option
- **WHEN** server.port is set to a valid port number
- **THEN** the server SHALL listen on that port for HTTP transport

#### Scenario: Host option
- **WHEN** server.host is set to a valid hostname or IP
- **THEN** the server SHALL bind to that address

### Requirement: Perplexity Configuration
The server SHALL support Perplexity-specific configuration.

#### Scenario: Perplexity enabled
- **WHEN** perplexity.enabled is true and api_key is set
- **THEN** Perplexity SHALL be available as a search provider

#### Scenario: Perplexity model selection
- **WHEN** perplexity.model is set
- **THEN** the specified model SHALL be used for Perplexity requests

#### Scenario: Perplexity disabled by default
- **WHEN** perplexity.enabled is not set or false
- **THEN** Perplexity SHALL NOT be used unless explicitly requested

### Requirement: DuckDuckGo Configuration
The server SHALL support DuckDuckGo-specific configuration.

#### Scenario: DuckDuckGo max results
- **WHEN** duckduckgo.max_results is set
- **THEN** DuckDuckGo searches SHALL return at most that many results

#### Scenario: DuckDuckGo region
- **WHEN** duckduckgo.region is set
- **THEN** DuckDuckGo searches SHALL be localized to that region

### Requirement: Output Format Configuration
The server SHALL support default output format configuration.

#### Scenario: Default text output
- **WHEN** output.default_format is "text"
- **THEN** tool responses SHALL default to text format

#### Scenario: Default JSON output
- **WHEN** output.default_format is "json"
- **THEN** tool responses SHALL default to JSON format

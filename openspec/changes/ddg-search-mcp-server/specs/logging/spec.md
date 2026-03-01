# Logging Specification

## Overview

This specification defines the logging system for the MCP server, including log levels, formats, and proxy-style request logging.

## ADDED Requirements

### Requirement: Log Level Configuration
The server SHALL support configurable log levels.

#### Scenario: Debug logging enabled
- **WHEN** logging.level is "debug"
- **THEN** all log messages including detailed debug information SHALL be output

#### Scenario: Info logging
- **WHEN** logging.level is "info"
- **THEN** info, warn, and error messages SHALL be output

#### Scenario: Warning logging
- **WHEN** logging.level is "warn"
- **THEN** warn and error messages SHALL be output

#### Scenario: Error logging only
- **WHEN** logging.level is "error"
- **THEN** only error messages SHALL be output

### Requirement: Log Format Configuration
The server SHALL support text and JSON log formats.

#### Scenario: Text log format
- **WHEN** logging.format is "text"
- **THEN** logs SHALL be output in human-readable text format

#### Scenario: JSON log format
- **WHEN** logging.format is "json"
- **THEN** logs SHALL be output as JSON objects

### Requirement: Proxy-Style Request Logging
The server SHALL log MCP requests similar to HTTP proxy servers.

#### Scenario: Log successful request
- **WHEN** a request completes successfully
- **THEN** the server SHALL log the request method, tool name, status, and duration

#### Scenario: Log failed request
- **WHEN** a request fails with error
- **THEN** the server SHALL log the error details with appropriate level

#### Scenario: Log fallback to alternate provider
- **WHEN** Perplexity falls back to DuckDuckGo
- **THEN** the server SHALL log a warning with fallback information

### Requirement: Log Fields
The server SHALL include structured fields in log output.

#### Scenario: Include method field
- **WHEN** logging a request
- **THEN** the log SHALL include the MCP method (e.g., tools/call)

#### Scenario: Include tool name field
- **WHEN** logging a tool request
- **THEN** the log SHALL include the tool name (search, fetch)

#### Scenario: Include status code
- **WHEN** logging a request
- **THEN** the log SHALL include HTTP-style status code

#### Scenario: Include duration
- **WHEN** logging a request
- **THEN** the log SHALL include request duration in milliseconds

#### Scenario: Include provider
- **WHEN** logging a search request
- **THEN** the log SHALL include the provider used (perplexity, duckduckgo)

### Requirement: Log Level for Different Status Codes
The server SHALL use appropriate log levels based on response status.

#### Scenario: 2xx status logged at info
- **WHEN** request returns 2xx success
- **THEN** the log SHALL be at INFO level

#### Scenario: 4xx status logged at warn
- **WHEN** request returns 4xx client error
- **THEN** the log SHALL be at WARN level

#### Scenario: 5xx status logged at error
- **WHEN** request returns 5xx server error
- **THEN** the log SHALL be at ERROR level

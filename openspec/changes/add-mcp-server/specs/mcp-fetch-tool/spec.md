## ADDED Requirements

### Requirement: Fetch tool registration
The system SHALL register a `fetch` tool that fetches web pages and converts them to markdown.

#### Scenario: Fetch tool is available
- **WHEN** an MCP client lists available tools
- **THEN** the `fetch` tool SHALL be present
- **AND** the tool SHALL have a description explaining its functionality

### Requirement: Fetch tool parameters
The `fetch` tool SHALL accept the following parameters:
- `url` (required): The URL to fetch (HTTP or HTTPS only)
- `timeout` (optional): Request timeout duration
- `user_agent` (optional): Custom user agent string
- `format` (optional): Output format: "json" or "text" (default from config)

#### Scenario: Fetch with required parameters
- **WHEN** a client calls `fetch` with only `url` parameter
- **THEN** the system SHALL fetch the page
- **AND** use default values for optional parameters from config

#### Scenario: Fetch with all parameters
- **WHEN** a client calls `fetch` with all parameters
- **THEN** the system SHALL use the provided parameter values
- **AND** override any config file defaults

### Requirement: URL validation
The system SHALL validate URLs before fetching.

#### Scenario: Valid HTTP URL
- **WHEN** `url` parameter is a valid HTTP URL
- **THEN** the system SHALL proceed to fetch the page

#### Scenario: Valid HTTPS URL
- **WHEN** `url` parameter is a valid HTTPS URL
- **THEN** the system SHALL proceed to fetch the page

#### Scenario: Invalid URL scheme
- **WHEN** `url` parameter uses an unsupported scheme (e.g., ftp://)
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "only HTTP and HTTPS URLs are supported"

#### Scenario: Invalid URL format
- **WHEN** `url` parameter is not a valid URL
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "invalid URL format"

### Requirement: Page fetching
The system SHALL fetch web pages using HTTP/HTTPS with configurable timeout and user agent.

#### Scenario: Successful page fetch
- **WHEN** the URL is valid and accessible
- **THEN** the system SHALL fetch the page content
- **AND** use the configured timeout
- **AND** use the configured user agent

#### Scenario: Page fetch timeout
- **WHEN** the page does not respond within the timeout period
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "request timeout"

#### Scenario: Page fetch network error
- **WHEN** a network error occurs during fetch
- **THEN** the system SHALL return an error
- **AND** the error SHALL include details about the network failure

#### Scenario: Page fetch HTTP error
- **WHEN** the server returns an HTTP error status (4xx or 5xx)
- **THEN** the system SHALL return an error
- **AND** the error SHALL include the HTTP status code

### Requirement: HTML to markdown conversion
The system SHALL convert fetched HTML content to markdown format.

#### Scenario: Successful conversion
- **WHEN** HTML content is successfully fetched
- **THEN** the system SHALL convert the HTML to markdown
- **AND** preserve headings, links, lists, and other formatting

#### Scenario: Conversion error
- **WHEN** HTML to markdown conversion fails
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "failed to convert HTML to markdown"

### Requirement: Fetch output format
The system SHALL support both JSON and plain text output formats.

#### Scenario: JSON output format
- **WHEN** `format` parameter is "json"
- **THEN** the system SHALL return results as JSON
- **AND** the JSON SHALL include the markdown content and metadata

#### Scenario: Text output format
- **WHEN** `format` parameter is "text"
- **THEN** the system SHALL return results as plain text
- **AND** the text SHALL be the markdown content

#### Scenario: Default output format from config
- **WHEN** `format` parameter is not provided
- **THEN** the system SHALL use the default format from config
- **AND** the config default SHALL be "text"

### Requirement: Fetch error handling
The system SHALL handle fetch errors gracefully and return appropriate error messages.

#### Scenario: Missing URL parameter
- **WHEN** `url` parameter is not provided
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "url parameter is required"

#### Scenario: Empty URL parameter
- **WHEN** `url` parameter is empty or whitespace only
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "url parameter is required"

#### Scenario: Redirect limit exceeded
- **WHEN** the page redirects too many times
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "maximum redirects exceeded"

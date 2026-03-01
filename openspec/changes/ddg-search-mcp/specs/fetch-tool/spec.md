## ADDED Requirements

### Requirement: Fetch tool accepts URL parameter
The system SHALL provide a fetch tool that accepts a required URL parameter.

#### Scenario: Fetch valid HTTP URL
- **WHEN** fetch tool is called with `url: "https://example.com/page"`
- **THEN** system SHALL fetch the URL and return content as markdown

#### Scenario: Fetch without URL fails
- **WHEN** fetch tool is called without URL parameter
- **THEN** system SHALL return an error indicating URL is required

### Requirement: Fetch tool validates URL scheme
The system SHALL only accept HTTP and HTTPS URLs.

#### Scenario: HTTPS URL accepted
- **WHEN** fetch tool is called with `url: "https://example.com"`
- **THEN** system SHALL fetch the URL successfully

#### Scenario: HTTP URL accepted
- **WHEN** fetch tool is called with `url: "http://example.com"`
- **THEN** system SHALL fetch the URL successfully

#### Scenario: Invalid scheme rejected
- **WHEN** fetch tool is called with `url: "ftp://example.com/file"`
- **THEN** system SHALL return an error indicating unsupported URL scheme

#### Scenario: Invalid URL format rejected
- **WHEN** fetch tool is called with `url: "not-a-url"`
- **THEN** system SHALL return an error indicating invalid URL

### Requirement: Fetch tool supports timeout configuration
The system SHALL allow configurable request timeout.

#### Scenario: Custom timeout
- **WHEN** fetch tool is called with `timeout: 60`
- **THEN** system SHALL use 60 second timeout for the request

#### Scenario: Default timeout
- **WHEN** fetch tool is called without timeout parameter
- **THEN** system SHALL use default timeout of 30 seconds

#### Scenario: Timeout minimum enforced
- **WHEN** fetch tool is called with `timeout: 0`
- **THEN** system SHALL use minimum timeout of 1 second

#### Scenario: Timeout maximum enforced
- **WHEN** fetch tool is called with `timeout: 200`
- **THEN** system SHALL limit timeout to maximum of 120 seconds

### Requirement: Fetch tool supports custom user agent
The system SHALL allow custom user agent configuration.

#### Scenario: Custom user agent
- **WHEN** fetch tool is called with `user_agent: "MyBot/1.0"`
- **THEN** system SHALL use the specified user agent for the request

#### Scenario: Default user agent
- **WHEN** fetch tool is called without user agent parameter
- **THEN** system SHALL use default user agent

### Requirement: Fetch tool supports output format override
The system SHALL allow per-request output format override.

#### Scenario: JSON format override
- **WHEN** fetch tool is called with `format: json`
- **THEN** response SHALL be formatted as JSON regardless of server default

#### Scenario: Text format override
- **WHEN** fetch tool is called with `format: text`
- **THEN** response SHALL be formatted as text regardless of server default

### Requirement: Fetch tool converts HTML to markdown
The system SHALL convert fetched HTML content to markdown format.

#### Scenario: HTML to markdown conversion
- **WHEN** fetch tool retrieves HTML content
- **THEN** system SHALL convert HTML to markdown
- **AND** response SHALL contain markdown-formatted text

#### Scenario: Preserves content structure
- **WHEN** fetch tool converts HTML with headings, lists, and links
- **THEN** markdown output SHALL preserve the document structure

### Requirement: Fetch tool handles errors gracefully
The system SHALL provide clear error messages for fetch failures.

#### Scenario: HTTP error response
- **WHEN** fetch tool receives HTTP 404 response
- **THEN** system SHALL return error indicating resource not found

#### Scenario: Network error
- **WHEN** fetch tool encounters network connection failure
- **THEN** system SHALL return error indicating connection failure

#### Scenario: Timeout error
- **WHEN** fetch tool request exceeds timeout
- **THEN** system SHALL return error indicating timeout

### Requirement: Fetch tool returns structured results
The system SHALL return structured fetch results.

#### Scenario: JSON response structure
- **WHEN** fetch tool returns results in JSON format
- **THEN** response SHALL include `url`, `content` (markdown), and `status_code`

#### Scenario: Text response format
- **WHEN** fetch tool returns results in text format
- **THEN** response SHALL include URL header followed by markdown content

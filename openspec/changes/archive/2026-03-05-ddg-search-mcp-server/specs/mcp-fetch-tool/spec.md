## ADDED Requirements

### Requirement: Fetch tool input schema
The system SHALL define an input schema for the fetch tool that accepts a URL parameter.

#### Scenario: Fetch with required URL parameter
- **WHEN** a fetch tool call includes a URL parameter
- **THEN** the system validates the URL is not empty
- **AND** the system validates the URL is a valid HTTP/HTTPS URL
- **AND** the system processes the fetch request

#### Scenario: Fetch with invalid URL
- **WHEN** a fetch tool call includes an invalid URL
- **THEN** the system returns an error response
- **AND** the error indicates the URL was invalid

#### Scenario: Fetch with missing URL parameter
- **WHEN** a fetch tool call does not include a URL parameter
- **THEN** the system returns an error response
- **AND** the error indicates the URL parameter is required

### Requirement: Web page content retrieval
The system SHALL retrieve web page content using the existing dump library.

#### Scenario: Successful page fetch
- **WHEN** a fetch tool call is received with a valid URL
- **THEN** the system fetches the page content
- **AND** the system converts the content to markdown format
- **AND** the system returns the markdown content in the MCP tool response

#### Scenario: Fetch with redirect
- **WHEN** the URL redirects to another location
- **THEN** the system follows the redirect
- **AND** the system returns content from the final destination

#### Scenario: Fetch with timeout
- **WHEN** the page fetch times out
- **THEN** the system returns an error response
- **AND** the error indicates the fetch timed out

### Requirement: Fetch result formatting
The system SHALL format fetch results according to MCP tool response format.

#### Scenario: Successful fetch returns formatted content
- **WHEN** a fetch completes successfully
- **THEN** the system returns the page content as markdown
- **AND** the system includes the source URL in the response
- **AND** the system formats the content as text in the MCP response

#### Scenario: Fetch with empty content
- **WHEN** a fetch returns an empty page
- **THEN** the system returns a message indicating no content was found
- **AND** the system includes the source URL

### Requirement: Fetch error handling
The system SHALL handle fetch errors gracefully and return appropriate error responses.

#### Scenario: Fetch with network error
- **WHEN** a fetch fails due to network error
- **THEN** the system returns an error response
- **AND** the error indicates a network failure occurred

#### Scenario: Fetch with HTTP error
- **WHEN** the server returns an HTTP error status (4xx, 5xx)
- **THEN** the system returns an error response
- **AND** the error includes the HTTP status code

#### Scenario: Fetch with unsupported content type
- **WHEN** the URL returns an unsupported content type (e.g., binary, PDF)
- **THEN** the system returns an error response
- **AND** the error indicates the content type is not supported

### Requirement: Fetch size limits
The system SHALL enforce reasonable size limits on fetched content.

#### Scenario: Fetch with large page
- **WHEN** the page content exceeds the size limit
- **THEN** the system truncates the content
- **AND** the system includes a note that content was truncated
- **AND** the system returns the truncated content

#### Scenario: Fetch with small page
- **WHEN** the page content is within size limits
- **THEN** the system returns the complete content
- **AND** the system does not truncate

### Requirement: Fetch logging
The system SHALL log all fetch operations at debug level.

#### Scenario: Successful fetch is logged
- **WHEN** a fetch completes successfully
- **THEN** the system logs the fetch at debug level
- **AND** the log includes the URL and content size

#### Scenario: Failed fetch is logged
- **WHEN** a fetch fails
- **THEN** the system logs the fetch at debug level
- **AND** the log includes the URL and error details

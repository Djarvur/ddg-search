## ADDED Requirements

### Requirement: Fetch HTML content from URL

The system SHALL fetch HTML content from a given URL using HTTP GET request.

#### Scenario: Successful fetch
- **WHEN** a valid HTTP URL is provided
- **THEN** system returns the HTML body content

#### Scenario: HTTPS support
- **WHEN** an HTTPS URL is provided
- **THEN** system fetches content over TLS without certificate errors for valid certs

### Requirement: Configurable timeout

The system SHALL support configurable request timeout.

#### Scenario: Default timeout
- **WHEN** no timeout is specified
- **THEN** system uses a default timeout of 30 seconds

#### Scenario: Custom timeout
- **WHEN** `--timeout` flag is specified
- **THEN** system uses the specified timeout value

#### Scenario: Timeout exceeded
- **WHEN** request exceeds the timeout
- **THEN** system returns a timeout error to stderr and exits with non-zero code

### Requirement: Custom user agent

The system SHALL support custom user agent string.

#### Scenario: Default user agent
- **WHEN** no user agent is specified
- **THEN** system uses default user agent `page-dump/1.0`

#### Scenario: Custom user agent
- **WHEN** `--user-agent` flag is specified
- **THEN** system sends the specified user agent header

### Requirement: Follow HTTP redirects

The system SHALL follow HTTP redirects up to a reasonable limit.

#### Scenario: Follow 301 redirect
- **WHEN** server responds with 301 redirect
- **THEN** system follows the redirect and returns content from final URL

#### Scenario: Follow 302 redirect
- **WHEN** server responds with 302 redirect
- **THEN** system follows the redirect and returns content from final URL

#### Scenario: Too many redirects
- **WHEN** redirect chain exceeds 10 hops
- **THEN** system returns a redirect error to stderr and exits with non-zero code

### Requirement: Handle HTTP errors

The system SHALL report HTTP errors clearly.

#### Scenario: 404 Not Found
- **WHEN** server responds with 404
- **THEN** system returns error message to stderr and exits with non-zero code

#### Scenario: 500 Server Error
- **WHEN** server responds with 5xx status
- **THEN** system returns error message to stderr and exits with non-zero code

### Requirement: Validate URL

The system SHALL validate the input URL format.

#### Scenario: Invalid URL format
- **WHEN** an invalid URL is provided
- **THEN** system returns error message to stderr and exits with non-zero code

#### Scenario: Unsupported protocol
- **WHEN** a non-HTTP(S) URL is provided (e.g., ftp://)
- **THEN** system returns error message to stderr and exits with non-zero code

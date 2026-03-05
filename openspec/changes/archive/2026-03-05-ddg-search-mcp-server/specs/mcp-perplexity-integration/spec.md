## ADDED Requirements

### Requirement: Perplexity configuration
The system SHALL support Perplexity configuration via access token and enabled flag.

#### Scenario: Perplexity enabled with token
- **WHEN** Perplexity is enabled in configuration
- **AND** a valid access token is provided
- **THEN** the system uses Perplexity for search requests

#### Scenario: Perplexity disabled
- **WHEN** Perplexity is disabled in configuration
- **THEN** the system uses DuckDuckGo for all search requests
- **AND** the system does not attempt Perplexity search

#### Scenario: Perplexity enabled without token
- **WHEN** Perplexity is enabled in configuration
- **AND** no access token is provided
- **THEN** the system falls back to DuckDuckGo for search requests
- **AND** the system includes a note in the response about the fallback

### Requirement: Perplexity authentication
The system SHALL authenticate Perplexity requests using the provided access token.

#### Scenario: Valid access token
- **WHEN** a valid access token is provided
- **THEN** the system includes the token in Perplexity API requests
- **AND** the system successfully authenticates with Perplexity

#### Scenario: Invalid access token
- **WHEN** an invalid access token is provided
- **THEN** the system receives an authentication error from Perplexity
- **AND** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

#### Scenario: Expired access token
- **WHEN** an expired access token is provided
- **THEN** the system receives an authentication error from Perplexity
- **AND** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

### Requirement: Perplexity rate limit handling
The system SHALL handle Perplexity rate limits by falling back to DuckDuckGo without retry.

#### Scenario: Perplexity rate limit exceeded
- **WHEN** Perplexity returns a rate limit error
- **THEN** the system does not retry Perplexity
- **AND** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

#### Scenario: Perplexity quota exceeded
- **WHEN** Perplexity returns a quota exceeded error
- **THEN** the system does not retry Perplexity
- **AND** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

### Requirement: Perplexity search parameters
The system SHALL pass search parameters to Perplexity API.

#### Scenario: Perplexity search with query
- **WHEN** a search request is made with a query
- **THEN** the system passes the query to Perplexity API
- **AND** the system returns Perplexity search results

#### Scenario: Perplexity search with max_results
- **WHEN** a search request includes max_results parameter
- **THEN** the system passes the max_results to Perplexity API
- **AND** the system limits results to the specified count

#### Scenario: Perplexity search with safe_search
- **WHEN** a search request includes safe_search parameter
- **THEN** the system passes the safe_search setting to Perplexity API
- **AND** the system applies safe search filtering if enabled

### Requirement: Perplexity result formatting
The system SHALL format Perplexity search results according to MCP tool response format.

#### Scenario: Successful Perplexity search
- **WHEN** Perplexity search completes successfully
- **THEN** the system returns results with title, URL, and snippet for each result
- **AND** the system formats results as text content in the MCP response

#### Scenario: Perplexity search with no results
- **WHEN** Perplexity search returns no results
- **THEN** the system returns an empty results message
- **AND** the system indicates no results were found

### Requirement: Perplexity error handling
The system SHALL handle Perplexity errors gracefully and fall back to DuckDuckGo.

#### Scenario: Perplexity API error
- **WHEN** Perplexity returns an API error
- **THEN** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

#### Scenario: Perplexity network error
- **WHEN** Perplexity request fails due to network error
- **THEN** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

#### Scenario: Perplexity timeout
- **WHEN** Perplexity request times out
- **THEN** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

### Requirement: Perplexity fallback behavior
The system SHALL always fall back to DuckDuckGo when Perplexity is unavailable.

#### Scenario: Fallback to DuckDuckGo on Perplexity failure
- **WHEN** Perplexity search fails for any reason
- **THEN** the system executes DuckDuckGo search
- **AND** the system returns DuckDuckGo results
- **AND** the system includes a note in the response about the fallback

#### Scenario: Fallback note in response
- **WHEN** the system falls back to DuckDuckGo
- **THEN** the system includes a clear note in the response
- **AND** the note indicates results are from DuckDuckGo due to Perplexity unavailability

### Requirement: No Perplexity retries
The system SHALL not retry Perplexity requests on failure.

#### Scenario: Perplexity request fails
- **WHEN** a Perplexity request fails
- **THEN** the system does not retry the Perplexity request
- **AND** the system immediately falls back to DuckDuckGo

#### Scenario: Perplexity transient error
- **WHEN** Perplexity returns a transient error (e.g., 5xx)
- **THEN** the system does not retry the Perplexity request
- **AND** the system immediately falls back to DuckDuckGo

## ADDED Requirements

### Requirement: Search tool input schema
The system SHALL define an input schema for the search tool that accepts query, max_results, and other search parameters.

#### Scenario: Search with required query parameter
- **WHEN** a search tool call includes a query parameter
- **THEN** the system validates the query is not empty
- **AND** the system processes the search request

#### Scenario: Search with optional max_results parameter
- **WHEN** a search tool call includes max_results parameter
- **THEN** the system uses the specified max_results value
- **AND** the system limits results to the specified count

#### Scenario: Search without max_results parameter
- **WHEN** a search tool call does not include max_results parameter
- **THEN** the system uses the default max_results value from configuration

#### Scenario: Search with safe_search parameter
- **WHEN** a search tool call includes safe_search parameter
- **THEN** the system uses the specified safe_search value
- **AND** the system applies safe search filtering if enabled

### Requirement: DuckDuckGo search execution
The system SHALL execute DuckDuckGo search when Perplexity is not available or disabled.

#### Scenario: DuckDuckGo search with valid query
- **WHEN** Perplexity is disabled or no token is provided
- **AND** a search tool call is received with a valid query
- **THEN** the system executes DuckDuckGo search
- **AND** the system returns search results in MCP tool response format

#### Scenario: DuckDuckGo search with rate limit
- **WHEN** DuckDuckGo returns a rate limit error
- **THEN** the system retries according to the internal search library's retry logic
- **AND** the system returns results if retry succeeds
- **AND** the system returns an error if all retries fail

### Requirement: Perplexity search execution
The system SHALL execute Perplexity search when enabled and access token is provided.

#### Scenario: Perplexity search with valid token
- **WHEN** Perplexity is enabled in configuration
- **AND** a valid access token is provided
- **AND** a search tool call is received
- **THEN** the system executes Perplexity search
- **AND** the system returns search results in MCP tool response format

#### Scenario: Perplexity search with rate limit
- **WHEN** Perplexity returns a rate limit error
- **THEN** the system does not retry Perplexity
- **AND** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

#### Scenario: Perplexity search with invalid token
- **WHEN** Perplexity returns an authentication error
- **THEN** the system falls back to DuckDuckGo search
- **AND** the system includes a note in the response about the fallback

### Requirement: Search result formatting
The system SHALL format search results according to MCP tool response format.

#### Scenario: Successful search returns formatted results
- **WHEN** a search completes successfully
- **THEN** the system returns results with title, URL, and snippet for each result
- **AND** the system formats results as text content in the MCP response

#### Scenario: Search with no results
- **WHEN** a search returns no results
- **THEN** the system returns an empty results message
- **AND** the system indicates no results were found

### Requirement: Search error handling
The system SHALL handle search errors gracefully and return appropriate error responses.

#### Scenario: Search with invalid query
- **WHEN** a search tool call has an invalid query
- **THEN** the system returns an error response
- **AND** the error indicates the query was invalid

#### Scenario: Search with network error
- **WHEN** a search fails due to network error
- **THEN** the system returns an error response
- **AND** the error indicates a network failure occurred

#### Scenario: Search with timeout
- **WHEN** a search times out
- **THEN** the system returns an error response
- **AND** the error indicates the search timed out

### Requirement: Search parameter defaults
The system SHALL use configuration defaults for optional search parameters.

#### Scenario: Default max_results from config
- **WHEN** max_results is not provided in tool call
- **THEN** the system uses the max_results value from configuration
- **AND** the system uses a sensible default if not configured

#### Scenario: Default safe_search from config
- **WHEN** safe_search is not provided in tool call
- **THEN** the system uses the safe_search value from configuration
- **AND** the system uses a sensible default if not configured

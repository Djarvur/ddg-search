## ADDED Requirements

### Requirement: Search tool registration
The system SHALL register a `search` tool that performs web search with Perplexity fallback to DuckDuckGo.

#### Scenario: Search tool is available
- **WHEN** an MCP client lists available tools
- **THEN** the `search` tool SHALL be present
- **AND** the tool SHALL have a description explaining its functionality

### Requirement: Search tool parameters
The `search` tool SHALL accept the following parameters:
- `query` (required): The search query string
- `max_results` (optional): Maximum number of results to return
- `site` (optional): Filter results to a specific domain
- `region` (optional): Search region (e.g., "us-en", "uk-en")
- `time` (optional): Time filter: d (day), w (week), m (month), y (year)
- `safe_search` (optional): Enable safe search
- `model` (optional): Perplexity model to use (e.g., "sonar-medium-online")
- `format` (optional): Output format: "json" or "text" (default from config)

#### Scenario: Search with required parameters
- **WHEN** a client calls `search` with only `query` parameter
- **THEN** the system SHALL perform the search
- **AND** use default values for optional parameters from config

#### Scenario: Search with all parameters
- **WHEN** a client calls `search` with all parameters
- **THEN** the system SHALL use the provided parameter values
- **AND** override any config file defaults

### Requirement: Perplexity fallback
The system SHALL attempt Perplexity search first when an API key is configured, and fallback to DuckDuckGo on failure.

#### Scenario: Perplexity search succeeds
- **WHEN** Perplexity API key is configured
- **AND** Perplexity search succeeds
- **THEN** the system SHALL return Perplexity results
- **AND** the response SHALL NOT indicate a fallback occurred

#### Scenario: Perplexity search fails with rate limit
- **WHEN** Perplexity API key is configured
- **AND** Perplexity search fails with rate limit error
- **THEN** the system SHALL fallback to DuckDuckGo search
- **AND** the response SHALL indicate "fallback to DuckDuckGo: Perplexity rate limit exceeded"

#### Scenario: Perplexity search fails with quota exceeded
- **WHEN** Perplexity API key is configured
- **AND** Perplexity search fails with quota exceeded error
- **THEN** the system SHALL fallback to DuckDuckGo search
- **AND** the response SHALL indicate "fallback to DuckDuckGo: Perplexity quota exceeded"

#### Scenario: No Perplexity API key configured
- **WHEN** Perplexity API key is not configured
- **THEN** the system SHALL use DuckDuckGo search directly
- **AND** the response SHALL indicate "using DuckDuckGo: no Perplexity API key configured"

#### Scenario: Perplexity disabled in config
- **WHEN** Perplexity is explicitly disabled in config
- **THEN** the system SHALL use DuckDuckGo search directly
- **AND** the system SHALL NOT attempt Perplexity search

### Requirement: Search output format
The system SHALL support both JSON and plain text output formats.

#### Scenario: JSON output format
- **WHEN** `format` parameter is "json"
- **THEN** the system SHALL return results as JSON
- **AND** the JSON SHALL include all result fields (title, url, snippet)

#### Scenario: Text output format
- **WHEN** `format` parameter is "text"
- **THEN** the system SHALL return results as formatted text
- **AND** the text SHALL be markdown-formatted with numbered results

#### Scenario: Default output format from config
- **WHEN** `format` parameter is not provided
- **THEN** the system SHALL use the default format from config
- **AND** the config default SHALL be "text"

### Requirement: Search result structure
The system SHALL return search results in a consistent structure regardless of the search provider.

#### Scenario: DuckDuckGo result structure
- **WHEN** using DuckDuckGo search
- **THEN** results SHALL include title, url, and snippet for each result
- **AND** results SHALL be ordered by relevance

#### Scenario: Perplexity result structure
- **WHEN** using Perplexity search
- **THEN** results SHALL include an AI-generated answer
- **AND** results SHALL include citations with URLs
- **AND** the answer SHALL be formatted as markdown

### Requirement: Search error handling
The system SHALL handle search errors gracefully and return appropriate error messages.

#### Scenario: Empty query
- **WHEN** `query` parameter is empty or whitespace only
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate "query parameter is required"

#### Scenario: Both search providers fail
- **WHEN** Perplexity search fails
- **AND** DuckDuckGo search also fails
- **THEN** the system SHALL return an error
- **AND** the error SHALL include details from both failures

#### Scenario: Invalid parameter value
- **WHEN** a parameter has an invalid value (e.g., negative max_results)
- **THEN** the system SHALL return an error
- **AND** the error SHALL indicate which parameter is invalid

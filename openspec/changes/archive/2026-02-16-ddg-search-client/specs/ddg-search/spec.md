## ADDED Requirements

### Requirement: Search query execution

The system SHALL execute search queries against DuckDuckGo's HTML endpoint and return structured results.

#### Scenario: Basic search returns results

- **WHEN** user provides a search query string
- **THEN** system returns an array of results with title, URL, and snippet

#### Scenario: Empty query returns empty results

- **WHEN** user provides an empty or whitespace-only query
- **THEN** system returns an empty array without making a network request

### Requirement: Result structure

Each search result SHALL contain title, URL, and snippet fields.

#### Scenario: Result contains all fields

- **WHEN** a search result is returned
- **THEN** it includes a non-empty title string
- **AND** it includes a valid URL string
- **AND** it includes a snippet string (may be empty if not available)

### Requirement: JSON output format

The system SHALL output results as a JSON array for programmatic consumption.

#### Scenario: Output is valid JSON array

- **WHEN** search completes successfully
- **THEN** output is a valid JSON array of result objects
- **AND** output contains no metadata wrapper

### Requirement: Site-specific search

The system SHALL support filtering results to a specific domain.

#### Scenario: Site filter limits results

- **WHEN** user specifies a site filter (e.g., `github.com`)
- **THEN** results are restricted to that domain

### Requirement: Regional search

The system SHALL support region-specific search results.

#### Scenario: Region parameter affects results

- **WHEN** user specifies a region (e.g., `us-en`, `uk-en`)
- **THEN** search uses the specified region for localization

### Requirement: Time-bounded search

The system SHALL support limiting results by time period.

#### Scenario: Time filter limits results

- **WHEN** user specifies a time filter (`d`, `w`, `m`, `y`)
- **THEN** results are limited to the specified time period

### Requirement: Safe search control

The system SHALL allow enabling or disabling safe search.

#### Scenario: Safe search can be toggled

- **WHEN** user specifies safe search setting
- **THEN** search respects the safe search preference

### Requirement: Result count limiting

The system SHALL support limiting the number of results returned.

#### Scenario: Max results limits output

- **WHEN** user specifies a maximum result count
- **THEN** system returns at most that many results

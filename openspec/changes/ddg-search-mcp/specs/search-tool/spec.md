## ADDED Requirements

### Requirement: Search tool accepts query parameter
The system SHALL provide a search tool that accepts a required query parameter.

#### Scenario: Search with query only
- **WHEN** search tool is called with `query: "golang best practices"`
- **THEN** system SHALL return search results for the query

#### Scenario: Search without query fails
- **WHEN** search tool is called without query parameter
- **THEN** system SHALL return an error indicating query is required

### Requirement: Search tool supports backend selection
The system SHALL allow selection of search backend via the `backend` parameter.

#### Scenario: Auto backend selection with Perplexity configured
- **WHEN** search tool is called with `backend: auto` and Perplexity is configured
- **THEN** system SHALL use Perplexity as the search backend

#### Scenario: Auto backend selection without Perplexity
- **WHEN** search tool is called with `backend: auto` and Perplexity is not configured
- **THEN** system SHALL use DuckDuckGo as the search backend

#### Scenario: Explicit DuckDuckGo backend
- **WHEN** search tool is called with `backend: ddg`
- **THEN** system SHALL use DuckDuckGo as the search backend

#### Scenario: Explicit Perplexity backend
- **WHEN** search tool is called with `backend: perplexity`
- **THEN** system SHALL use Perplexity as the search backend

### Requirement: Automatic fallback to DuckDuckGo
The system SHALL fall back to DuckDuckGo when Perplexity fails.

#### Scenario: Fallback on Perplexity rate limit
- **WHEN** search tool is called with `backend: auto` or `backend: perplexity`
- **AND** Perplexity returns a rate limit error
- **THEN** system SHALL fall back to DuckDuckGo
- **AND** response SHALL include fallback metadata indicating reason

#### Scenario: Fallback on Perplexity API error
- **WHEN** search tool is called with `backend: auto` or `backend: perplexity`
- **AND** Perplexity returns an API error
- **THEN** system SHALL fall back to DuckDuckGo
- **AND** response SHALL include fallback metadata indicating reason

#### Scenario: Perplexity not configured with explicit backend
- **WHEN** search tool is called with `backend: perplexity` and Perplexity is not configured
- **THEN** system SHALL use DuckDuckGo instead
- **AND** response SHALL include note that Perplexity is not configured

### Requirement: Search tool supports result limit
The system SHALL allow limiting the number of results via `max_results` parameter.

#### Scenario: Limit results to 5
- **WHEN** search tool is called with `max_results: 5`
- **THEN** system SHALL return at most 5 search results

#### Scenario: Default result limit
- **WHEN** search tool is called without `max_results`
- **THEN** system SHALL return default of 10 results

#### Scenario: Maximum result limit enforced
- **WHEN** search tool is called with `max_results: 100`
- **THEN** system SHALL limit results to maximum of 50

### Requirement: Search tool supports DuckDuckGo-specific options
The system SHALL support DuckDuckGo-specific search options.

#### Scenario: Site-specific search
- **WHEN** search tool is called with `site: "github.com"`
- **THEN** system SHALL limit results to the specified domain

#### Scenario: Region-specific search
- **WHEN** search tool is called with `region: "uk-en"`
- **THEN** system SHALL use the specified region for search

#### Scenario: Time-filtered search
- **WHEN** search tool is called with `time_filter: "w"`
- **THEN** system SHALL limit results to past week

#### Scenario: Safe search enabled
- **WHEN** search tool is called with `safe_search: true`
- **THEN** system SHALL enable safe search filtering

### Requirement: Search tool supports Perplexity-specific options
The system SHALL support Perplexity-specific search options.

#### Scenario: Custom Perplexity model
- **WHEN** search tool is called with `model: "sonar-large-online"`
- **THEN** system SHALL use the specified Perplexity model

### Requirement: Search tool supports output format override
The system SHALL allow per-request output format override.

#### Scenario: JSON format override
- **WHEN** search tool is called with `format: json`
- **THEN** response SHALL be formatted as JSON regardless of server default

#### Scenario: Text format override
- **WHEN** search tool is called with `format: text`
- **THEN** response SHALL be formatted as text regardless of server default

### Requirement: Search tool returns structured results
The system SHALL return structured search results.

#### Scenario: JSON response structure
- **WHEN** search tool returns results in JSON format
- **THEN** response SHALL include `query`, `backend_used`, `results` array, and `total` count
- **AND** each result SHALL include `title`, `url`, and `snippet`

#### Scenario: Fallback metadata in JSON response
- **WHEN** search tool falls back from Perplexity to DuckDuckGo
- **AND** response is in JSON format
- **THEN** response SHALL include `fallback` object with `from` and `reason` fields

#### Scenario: Text response format
- **WHEN** search tool returns results in text format
- **THEN** response SHALL include header with query, backend used, and numbered results
- **AND** each result SHALL show title, URL, and snippet

#### Scenario: Fallback note in text response
- **WHEN** search tool falls back from Perplexity to DuckDuckGo
- **AND** response is in text format
- **THEN** response SHALL include note indicating fallback occurred and reason

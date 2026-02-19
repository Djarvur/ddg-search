## ADDED Requirements

### Requirement: Provide perplexity-search CLI command

The system SHALL provide a `perplexity-search` CLI command that enables users to search the web using the Perplexity API.

#### Scenario: Successful search
- **WHEN** a user executes `perplexity-search` with a valid query
- **THEN** the command sends a request to the Perplexity API
- **AND** the command outputs search results in markdown format

#### Scenario: Search with no query
- **WHEN** a user executes `perplexity-search` without a query
- **THEN** the command returns an error indicating a query is required

#### Scenario: Missing API key
- **WHEN** the PERPLEXITY_API_KEY environment variable is not set
- **THEN** the command returns an error directing the user to configure the API key in .env file

#### Scenario: Invalid API key
- **WHEN** the PERPLEXITY_API_KEY is invalid
- **THEN** the command returns an error from the Perplexity API indicating authentication failure

#### Scenario: Rate limit exceeded
- **WHEN** the Perplexity API rate limit is exceeded
- **THEN** the command returns an error with rate limit information from the API

### Requirement: Support configurable search options

The CLI command SHALL support configurable search options via command-line flags.

#### Scenario: Limit results
- **WHEN** a user provides `--max-results N` flag
- **THEN** the command returns at most N search results

#### Scenario: Specify model
- **WHEN** a user provides `--model MODEL` flag
- **THEN** the command uses the specified Perplexity model for the search

#### Scenario: Enable debug mode
- **WHEN** a user provides `--debug` flag
- **THEN** the command outputs debug logging to stderr

### Requirement: Load API key from environment

The system SHALL load the Perplexity API key from the PERPLEXITY_API_KEY environment variable.

#### Scenario: API key from .env
- **WHEN** a .env file exists with PERPLEXITY_API_KEY set
- **THEN** the system loads the API key from the environment

#### Scenario: API key from shell environment
- **WHEN** PERPLEXITY_API_KEY is set in the shell environment
- **THEN** the system uses that API key value

### Requirement: Format output as markdown

The system SHALL format search results as markdown for LLM consumption.

#### Scenario: Result formatting
- **WHEN** search results are received from the API
- **THEN** each result is formatted with title, URL, and snippet in markdown
- **AND** citations are included for sources referenced in the AI-generated content

### Requirement: Handle API errors gracefully

The system SHALL handle Perplexity API errors with clear, actionable error messages.

#### Scenario: Network error
- **WHEN** a network error occurs during API request
- **THEN** the command retries with exponential backoff
- **AND** after max retries, returns a clear error message

#### Scenario: API error response
- **WHEN** the API returns an error response
- **THEN** the command extracts and displays the error message from the API

### Requirement: Provide perplexity-search-skill

The system SHALL provide a Claude Code skill that exposes Perplexity search functionality.

#### Scenario: Skill tool definition
- **WHEN** the skill.yaml is loaded
- **THEN** it contains a tool definition for perplexity-search with appropriate parameters

#### Scenario: Skill invocation
- **WHEN** Claude invokes the perplexity-search tool
- **THEN** the skill executes the perplexity-search binary with the provided query
- **AND** returns the search results to Claude

#### Scenario: Skill location
- **WHEN** the skill is installed
- **THEN** it is located at skills/perplexity-search/

### Requirement: Document installation and usage

The system SHALL provide clear documentation for installation and usage of the perplexity-search command and skill.

#### Scenario: README documentation
- **WHEN** a user reads README.md
- **THEN** it includes a section on perplexity-search with usage examples
- **AND** it documents the API key requirement
- **AND** it compares perplexity-search with ddg-search

#### Scenario: SKILL.md documentation
- **WHEN** a user reads skills/perplexity-search/SKILL.md
- **THEN** it includes installation instructions for the API key
- **AND** it includes usage examples
- **AND** it includes troubleshooting for common issues

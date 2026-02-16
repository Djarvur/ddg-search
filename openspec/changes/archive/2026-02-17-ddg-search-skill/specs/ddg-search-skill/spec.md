## ADDED Requirements

### Requirement: Provide web-search tool

The skill SHALL expose a `web-search` tool that enables Claude to search the web using DuckDuckGo.

#### Scenario: Successful search
- **WHEN** Claude invokes web-search with a query string
- **THEN** the skill executes `ddg-search` binary and returns search results

#### Scenario: Search with no query
- **WHEN** Claude invokes web-search without a query
- **THEN** the skill returns an error indicating query is required

#### Scenario: Binary not found
- **WHEN** the ddg-search binary is not installed
- **THEN** the skill displays install instructions from SKILL.md

### Requirement: Provide page-dump tool

The skill SHALL expose a `page-dump` tool that enables Claude to fetch and convert web pages to markdown.

#### Scenario: Successful page fetch
- **WHEN** Claude invokes page-dump with a valid URL
- **THEN** the skill executes `page-dump` binary and returns markdown content

#### Scenario: Invalid URL
- **WHEN** Claude invokes page-dump with an invalid URL
- **THEN** the skill returns an error from the page-dump binary

#### Scenario: Binary not found
- **WHEN** the page-dump binary is not installed
- **THEN** the skill displays install instructions from SKILL.md

### Requirement: Skill configuration

The skill SHALL be configured via skill.yaml with proper tool definitions.

#### Scenario: Tool definitions present
- **WHEN** the skill.yaml is loaded
- **THEN** it contains tool definitions for web-search and page-dump with appropriate parameters

#### Scenario: Skill location
- **WHEN** the skill is installed
- **THEN** it is located at skills/ddg-search/

### Requirement: Install instructions

The skill SHALL provide install instructions that display when binaries are not found.

#### Scenario: Install instructions available
- **WHEN** a binary is missing
- **THEN** SKILL.md contains build and install commands for ddg-search and page-dump

#### Scenario: Instructions include PATH setup
- **WHEN** user reads install instructions
- **THEN** instructions specify that binaries must be in PATH

#### Scenario: File naming convention
- **WHEN** the skill is created
- **THEN** the documentation file is named SKILL.md (uppercase)

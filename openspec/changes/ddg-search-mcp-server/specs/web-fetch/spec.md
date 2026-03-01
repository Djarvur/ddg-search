# Web Fetch Tool Specification

## Overview

This specification defines the fetch MCP tool that provides web page fetching and conversion to markdown.

## ADDED Requirements

### Requirement: Web Fetch Tool Accepts URL Parameter
The fetch tool SHALL require a URL parameter for the page to fetch.

#### Scenario: Fetch with valid URL
- **WHEN** the tool is called with a valid HTTP or HTTPS URL
- **THEN** the server SHALL fetch the page and return its content

#### Scenario: Fetch with invalid URL
- **WHEN** the tool is called with an invalid or missing URL
- **THEN** the server SHALL return error code -32602 (invalid params)

#### Scenario: Fetch with unsupported protocol
- **WHEN** the tool is called with a URL using a protocol other than HTTP or HTTPS
- **THEN** the server SHALL return error code -32602 (invalid params)

### Requirement: Web Fetch Tool Supports Timeout
The fetch tool SHALL support configurable request timeout.

#### Scenario: Fetch with timeout
- **WHEN** timeout parameter is specified in seconds
- **THEN** the server SHALL abort the request if it exceeds the timeout

#### Scenario: Fetch with default timeout
- **WHEN** timeout is not specified
- **THEN** the server SHALL use the default timeout of 30 seconds

#### Scenario: Fetch times out
- **WHEN** the request exceeds the timeout period
- **THEN** the server SHALL return error with appropriate timeout message

### Requirement: Web Fetch Tool Supports Output Formats
The fetch tool SHALL support text and JSON output formats.

#### Scenario: Fetch with text output
- **WHEN** output_format is "text" or not specified
- **THEN** the page content SHALL be returned as markdown text

#### Scenario: Fetch with JSON output
- **WHEN** output_format is "json"
- **THEN** the result SHALL be returned as structured JSON with content field

### Requirement: Web Fetch Tool Converts HTML to Markdown
The fetch tool SHALL convert fetched HTML pages to markdown format.

#### Scenario: Fetch page converts to markdown
- **WHEN** a valid HTML page is fetched
- **THEN** the content SHALL be converted to readable markdown format

#### Scenario: Fetch preserves links and images
- **WHEN** HTML contains links and images
- **THEN** they SHALL be preserved in the markdown output with proper formatting

### Requirement: Web Fetch Tool Handles Errors
The tool SHALL properly handle and report fetch errors.

#### Scenario: Fetch non-existent page
- **WHEN** the URL returns 404 Not Found
- **THEN** the server SHALL return appropriate error code

#### Scenario: Fetch server error
- **WHEN** the server returns 5xx error
- **THEN** the server SHALL return appropriate error with status code

#### Scenario: Fetch connection error
- **WHEN** connection to the server fails
- **THEN** the server SHALL return appropriate network error message

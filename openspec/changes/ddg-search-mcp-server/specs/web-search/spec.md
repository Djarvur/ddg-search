# Web Search Tool Specification

## Overview

This specification defines the **search** MCP tool that provides web search functionality using Perplexity or DuckDuckGo providers.

## ADDED Requirements

### Requirement: Web Search Tool Accepts Query Parameter
The search tool SHALL require a query parameter for the search string.

#### Scenario: Search with valid query
- **WHEN** the tool is called with a non-empty query string
- **THEN** the server SHALL execute the search and return results

#### Scenario: Search with empty query
- **WHEN** the tool is called with empty or missing query
- **THEN** the server SHALL return error code -32602 (invalid params)

### Requirement: Web Search Tool Supports Max Results
The search tool SHALL support limiting the number of results.

#### Scenario: Search with max_results
- **WHEN** max_results is specified as a positive integer
- **THEN** the server SHALL return at most that many results

#### Scenario: Search with default max_results
- **WHEN** max_results is not specified
- **THEN** the server SHALL use provider-specific defaults (10 for DDG, 5 for Perplexity)

### Requirement: Web Search Tool Supports Provider Selection
The search tool SHALL support selecting the search provider.

#### Scenario: Search with explicit provider
- **WHEN** provider is set to "perplexity" or "duckduckgo"
- **THEN** the server SHALL use the specified provider

#### Scenario: Search with auto provider
- **WHEN** provider is set to "auto" (default)
- **THEN** the server SHALL use Perplexity if API key is configured, otherwise DuckDuckGo

### Requirement: Web Search Tool Supports Perplexity Models
The search tool SHALL support selecting Perplexity models when using perplexity provider.

#### Scenario: Search with sonar-small-online model
- **WHEN** provider is "perplexity" and model is "sonar-small-online"
- **THEN** the server SHALL use the sonar-small-online model

#### Scenario: Search with sonar-medium-online model
- **WHEN** provider is "perplexity" and model is "sonar-medium-online"
- **THEN** the server SHALL use the sonar-medium-online model

#### Scenario: Search with sonar-pro-online model
- **WHEN** provider is "perplexity" and model is "sonar-pro-online"
- **THEN** the server SHALL use the sonar-pro-online model

### Requirement: Web Search Tool Supports DuckDuckGo Filters
The search tool SHALL support DuckDuckGo-specific filters.

#### Scenario: Search with site filter
- **WHEN** site parameter is provided
- **THEN** results SHALL be filtered to the specified domain

#### Scenario: Search with region filter
- **WHEN** region parameter is provided (e.g., "us-en")
- **THEN** results SHALL be localized to the specified region

#### Scenario: Search with time filter
- **WHEN** time parameter is provided (d, w, m, y)
- **THEN** results SHALL be limited to that time period

### Requirement: Web Search Tool Supports Output Formats
The search tool SHALL support text and JSON output formats.

#### Scenario: Search with text output
- **WHEN** output_format is "text" or not specified
- **THEN** results SHALL be returned as human-readable text

#### Scenario: Search with JSON output
- **WHEN** output_format is "json"
- **THEN** results SHALL be returned as structured JSON

### Requirement: Web Search Tool Falls Back on Perplexity Failure
The tool SHALL automatically fall back to DuckDuckGo when Perplexity fails.

#### Scenario: Perplexity rate limited
- **WHEN** Perplexity returns 429 rate limit
- **THEN** the server SHALL fall back to DuckDuckGo and return successful results with fallback metadata

#### Scenario: Perplexity quota exceeded
- **WHEN** Perplexity returns 402 quota exceeded
- **THEN** the server SHALL fall back to DuckDuckGo and return successful results with fallback metadata

#### Scenario: Perplexity authentication error
- **WHEN** Perplexity returns 401/403 authentication error
- **THEN** the server SHALL fall back to DuckDuckGo and return successful results with fallback metadata

#### Scenario: Perplexity network error
- **WHEN** Perplexity request fails due to network error
- **THEN** the server SHALL fall back to DuckDuckGo and return successful results with fallback metadata

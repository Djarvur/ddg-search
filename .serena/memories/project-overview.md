# ddg-search Project Overview

## Purpose

Command-line tools for web search and content fetching. The project provides three main tools:

1. **ddg-search** - DuckDuckGo Search Client (no API key required)
2. **page-dump** - URL to Markdown Converter
3. **perplexity-search** - Perplexity API Search Client (requires API key)

## Tech Stack

- **Language**: Go 1.25.0
- **CLI Framework**: urfave/cli/v3
- **HTTP Client**: go-resty/resty/v2
- **HTML Parsing**: goquery
- **Markdown Conversion**: html-to-markdown/v2

## Project Structure

```
ddg-search/
├── cmd/                    # Entry points for CLI tools
│   ├── ddg-search/        # DuckDuckGo search CLI
│   ├── page-dump/         # URL to markdown converter
│   └── perplexity-search/ # Perplexity API search CLI
├── internal/              # Internal packages
│   ├── config/           # Configuration structures (RetryOptions, SearchOptions, Result)
│   ├── search/           # DuckDuckGo search implementation
│   │   ├── client.go     # HTTP client with retry logic
│   │   ├── parser.go     # HTML result parser
│   │   └── search.go     # Searcher interface
│   ├── perplexity/       # Perplexity API client
│   │   ├── client.go     # API client with retry logic
│   │   └── search.go     # Search implementation
│   └── dump/             # Page fetching and markdown conversion
├── skills/               # Roo skills for automation
│   ├── ddg-search/
│   └── perplexity-search/
├── openspec/             # OpenSpec change management
│   ├── changes/          # Active changes
│   ├── changes/archive/  # Archived changes
│   └── specs/            # Main specifications
└── .github/              # CI/CD workflows
```

## Key Components

### internal/config
- `RetryOptions`: Configures retry behavior (MaxRetries, BaseDelay, MaxDelay, BackoffMultiplier, Debug)
- `SearchOptions`: Search parameters (Query, MaxResults, Site, Region, TimeFilter, SafeSearch)
- `Result`: Search result structure (Title, URL, Snippet)

### internal/search
- `Client`: HTTP client with automatic rate limit detection and retry with exponential backoff
- `Parser`: Parses DuckDuckGo HTML results
- `Searcher`: High-level search interface with JSON and Markdown output

### internal/perplexity
- `Client`: Perplexity API client with retry logic
- `SearchOptions`: Query, MaxResults, Model
- `SearchResults`: Answer, Citations, References

### internal/dump
- `Fetch`: Fetches web pages with configurable timeout and user-agent
- `Convert`: Converts HTML to markdown preserving structure

## Features

### ddg-search
- Markdown output for LLM consumption
- Automatic rate limit detection and retry with exponential backoff
- Site-specific, regional, and time-bounded searches
- Configurable retry behavior

### page-dump
- Fetch any HTTP/HTTPS URL and convert to markdown
- Preserves document structure (headings, links, lists, code blocks)
- Configurable timeout and user-agent

### perplexity-search
- AI-powered search results with citations
- Configurable model selection (sonar-small-online, sonar-medium-online, sonar-pro-online)
- Automatic retry with exponential backoff
- Markdown output for LLM consumption

## Rate Limiting

Both DuckDuckGo and Perplexity clients automatically:
1. Detect rate-limit responses (HTTP 202, 429, 5xx)
2. Retry with exponential backoff + jitter
3. Fail gracefully with clear error messages after max retries

## License

MIT

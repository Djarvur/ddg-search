# ddg-search

[![CI](https://github.com/Djarvur/ddg-search/actions/workflows/ci.yml/badge.svg)](https://github.com/Djarvur/ddg-search/actions/workflows/ci.yml)
[![Security](https://github.com/Djarvur/ddg-search/actions/workflows/security.yml/badge.svg)](https://github.com/Djarvur/ddg-search/actions/workflows/security.yml)
[![Coveralls](https://coveralls.io/repos/github/Djarvur/ddg-search/badge.svg)](https://coveralls.io/github/Djarvur/ddg-search)
[![Go Report Card](https://goreportcard.com/badge/github.com/Djarvur/ddg-search)](https://goreportcard.com/report/github.com/Djarvur/ddg-search)
[![GoDoc](https://godoc.org/github.com/Djarvur/ddg-search?status.svg)](https://godoc.org/github.com/Djarvur/ddg-search)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Command-line tools for web search and content fetching.

## Tools

### ddg-search - DuckDuckGo Search Client

Search DuckDuckGo from the command line (no API key required).

### page-dump - URL to Markdown Converter

Fetch web pages and convert to markdown format.

## Features

### ddg-search

- Markdown output for LLM consumption
- Automatic rate limit detection and retry with exponential backoff
- Configurable retry behavior
- Site-specific, regional, and time-bounded searches

### page-dump

- Fetch any HTTP/HTTPS URL and convert to markdown
- Preserves document structure (headings, links, lists, code blocks)
- Configurable timeout and user-agent
- Clean markdown output to stdout

## Installation

```bash
go install github.com/Djarvur/ddg-search@latest
```

Or build from source:

```bash
git clone https://github.com/Djarvur/ddg-search
cd ddg-search
make build
```

## Usage

### ddg-search

```bash
ddg-search golang
```

Output (Markdown):

```markdown
1. [The Go Programming Language](https://go.dev/)
   Go is an open source programming language...
2. [Go (programming language) - Wikipedia](https://en.wikipedia.org/wiki/Go_(programming_language))
   ...
```

#### ddg-search Examples

```bash
# Limit results
ddg-search --max-results 5 golang tutorial

# Site-specific search
ddg-search --site github.com docker compose

# Regional search
ddg-search --region uk-en premier league

# Time-bounded search (d=day, w=week, m=month, y=year)
ddg-search --time w news today

# Retry configuration
ddg-search --max-retries 5 --retry-delay 2s --max-retry-delay 60s slow query

# Debug mode
ddg-search --debug golang
```

#### ddg-search Options

| Flag | Description | Default |
|------|-------------|---------|
| `--max-results` | Maximum number of results to return | 10 |
| `--site` | Filter results to a specific domain | - |
| `--region` | Search region (e.g., us-en, uk-en) | us-en |
| `--time` | Time filter: d (day), w (week), m (month), y (year) | - |
| `--safe-search` | Enable safe search | false |
| `--max-retries` | Maximum retry attempts on rate limiting | 3 |
| `--retry-delay` | Initial retry delay | 1s |
| `--max-retry-delay` | Maximum retry delay cap | 30s |
| `--debug` | Enable debug logging to stderr | false |

### page-dump

```bash
page-dump https://example.com
```

Output (Markdown):

```markdown
# Example Domain

This domain is for use in documentation examples...

[Learn more](https://iana.org/domains/example)
```

#### page-dump Options

| Flag | Description | Default |
|------|-------------|---------|
| `--timeout` | Request timeout | 30s |
| `--user-agent` | Custom user agent string | page-dump/1.0 |

## Rate Limiting

DuckDuckGo may rate-limit requests. ddg-search automatically:

1. Detects rate-limit responses (HTTP 202, 429, 5xx)
2. Retries with exponential backoff + jitter
3. Fails gracefully with clear error messages after max retries

## Use in Skills

These tools are designed for programmatic use in automation skills:

```bash
# Search for information
results=$(ddg-search --max-results 3 "$query")

# Fetch and convert a web page
content=$(page-dump "$url")
```

## License

MIT

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

### perplexity-search - Perplexity API Search Client

Search the web using Perplexity API with AI-powered results and citations (requires API key).

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

### perplexity-search

- AI-powered search results with better understanding
- Citations to sources referenced in results
- Configurable model selection (sonar, sonar-pro, sonar-reasoning, ...)
- Automatic retry with exponential backoff for transient errors
- Markdown output for LLM consumption

## Installation

`go install` needs the path of each command's `main` package, so install the
tools individually:

```bash
go install github.com/Djarvur/ddg-search/cmd/ddg-search@latest
go install github.com/Djarvur/ddg-search/cmd/page-dump@latest
go install github.com/Djarvur/ddg-search/cmd/perplexity-search@latest
```

Or build from source:

```bash
git clone https://github.com/Djarvur/ddg-search
cd ddg-search
go build -o bin/ddg-search ./cmd/ddg-search
go build -o bin/page-dump ./cmd/page-dump
go build -o bin/perplexity-search ./cmd/perplexity-search
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

### perplexity-search

```bash
# First, set your API key
export PERPLEXITY_API_KEY="your-api-key"

# Then search
perplexity-search "What is Go programming language?"
```

Output (Markdown):

```markdown
Go is an open source programming language that makes it easy to build simple, reliable, and efficient software...

## Sources

1. [The Go Programming Language](https://go.dev/)
2. [Go (programming language) - Wikipedia](https://en.wikipedia.org/wiki/Go_(programming_language))
```

#### perplexity-search Examples

```bash
# Limit the number of sources listed under the answer
perplexity-search --max-results 3 golang tutorial

# Use a specific model
perplexity-search --model sonar-pro "machine learning fundamentals"

# Debug mode
perplexity-search --debug "kubernetes deployment strategies"
```

#### perplexity-search Options

| Flag | Description | Default |
|------|-------------|---------|
| `--max-results` | Maximum number of sources listed with the answer | 5 |
| `--model` | Perplexity model to use | sonar |
| `--debug` | Enable debug logging to stderr | false |

#### perplexity-search Models

| Model | Description |
|-------|-------------|
| `sonar` | Fast and low cost (default) |
| `sonar-pro` | Higher quality, more expensive |
| `sonar-reasoning` | Adds a reasoning pass |
| `sonar-reasoning-pro` | Higher-quality reasoning |
| `sonar-deep-research` | Long-running, exhaustive research |

#### API Key Setup

Perplexity API requires an API key. Get one at https://www.perplexity.ai/settings/api

Set it via environment variable:
```bash
export PERPLEXITY_API_KEY="your-api-key"
```

`perplexity-search` reads the key from the environment only -- it does not load
`.env` files itself. To keep the key in a `.env` file, source it first
(`set -a; . ./.env; set +a`) or use a runner such as `direnv`.

#### ddg-search vs perplexity-search

| Feature | ddg-search | perplexity-search |
|---------|-------------|------------------|
| API Key Required | No | Yes |
| Result Quality | Standard search results | AI-powered, higher quality |
| Citations | No | Yes, with sources |
| Rate Limits | DuckDuckGo rate limits | Based on API tier |
| Cost | Free | Free tier available, paid for higher limits |
| Speed | Faster (HTML scraping) | Slightly slower (API call) |

Use `ddg-search` for free, fast searches without API key requirements. Use `perplexity-search` for higher-quality, AI-summarized results with citations.

## Rate Limiting

### DuckDuckGo

DuckDuckGo may rate-limit requests. ddg-search automatically:

1. Detects rate-limit responses (HTTP 202, 429, 5xx)
2. Retries with exponential backoff + jitter
3. Fails gracefully with clear error messages after max retries

### Perplexity API

Perplexity API has rate limits based on your API tier. perplexity-search automatically:

1. Detects rate-limit responses (HTTP 429, 5xx)
2. Retries with exponential backoff + jitter
3. Fails gracefully with clear error messages after max retries

If you exceed your rate limit, you'll see an error message indicating the limit. Consider upgrading your Perplexity plan for higher limits at https://www.perplexity.ai/settings/api.

## Use in Skills

These tools are designed for programmatic use in automation skills:

```bash
# Search for information with DuckDuckGo (no API key)
results=$(ddg-search --max-results 3 "$query")

# Search with Perplexity API (requires API key)
results=$(perplexity-search --max-results 3 "$query")

# Fetch and convert a web page
content=$(page-dump "$url")
```

## License

MIT

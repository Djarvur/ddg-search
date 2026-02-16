# DuckDuckGo CLI Search Client

## Why

Skills need a reliable way to perform web searches via DuckDuckGo without API keys. DuckDuckGo imposes rate limits that cause intermittent failures, so the client must detect rate-limiting responses and retry automatically with configurable backoff delays. This provides a privacy-focused search capability for programmatic use.

## What Changes

- New Go module `github.com/Djarvur/ddg-search` implementing a DuckDuckGo search client
- CLI binary `ddg-search` that can be invoked from skills and other tools
- Core search functionality using DuckDuckGo HTML endpoint (no API key required)
- Rate limit detection by identifying HTTP 202, 429, and 5xx responses
- Retry mechanism with configurable maximum attempts and exponential backoff with jitter
- Markdown output format as default (JSON also available via API)
- Basic search options: site filter, region, time limit, safe search
- All configuration via CLI arguments only (no config file)
- Verbose CLI parameter names for clarity
- Safe search is a boolean flag, disabled by default
- Debug mode for troubleshooting retry behavior

## Capabilities

### New Capabilities

- `ddg-search`: Core search capability - query DuckDuckGo HTML endpoint, parse results, output structured data (title, URL, snippet). Supports result count limiting, site-specific search, region selection, and time-bounded queries. Outputs markdown-formatted results by default.

- `rate-limit-handler`: Rate limit detection and retry mechanism - detect when DuckDuckGo is rate-limiting requests (HTTP 202, 429, or 5xx responses), wait with configurable exponential backoff and jitter, retry up to configured maximum attempts. Exposes retry configuration (max retries, initial delay, max delay).

- `debug-logging`: Debug output capability - when enabled via `--debug` flag, output verbose logs to stderr including retry attempts, delays, and rate limit detection events. Useful for troubleshooting rate limiting issues.

### Modified Capabilities

(None - this is a new module)

## Impact

- **New Go module**: `github.com/Djarvur/ddg-search` at `/Users/nil/DiskD/W/Djarvur/ddg-search`
- **CLI binary**: `ddg-search` invoked via `go run ./cmd/ddg-search` or built binary
- **Dependencies**:
  - `github.com/urfave/cli/v3` - CLI framework
  - `github.com/go-resty/resty/v2` - HTTP client with retry support
  - `github.com/PuerkitoBio/goquery` - HTML parsing
- **No external API keys required**
- **Integration**: To be used by the `web-search` skill and other tooling
- **Reference implementation**: Inspired by [ddgr](https://github.com/jarun/ddgr) Python CLI

## Implementation Notes

- HTML content-based rate limit detection was attempted but disabled due to false positives (indicator words appearing in legitimate search results)
- Markdown output format chosen as default for better readability in terminal and LLM contexts
- Two-level retry: HTTP client level (status code detection) and search orchestration level (for extensibility)

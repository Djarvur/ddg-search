# DuckDuckGo CLI Search Client - Design

## Context

Building a Go-based CLI tool for DuckDuckGo search that can be invoked by skills programmatically. The tool must handle DuckDuckGo's rate limiting gracefully since skills may make frequent search requests.

**Reference:** [ddgr](https://github.com/jarun/ddgr) is a Python implementation that uses DuckDuckGo's HTML endpoint (`duckduckgo.com/html/`) rather than an API. This approach requires no API key and works over Tor.

**Constraints:**
- No API keys required (uses HTML scraping)
- Must work reliably despite rate limiting
- Output must be machine-readable (JSON) for skills
- Fast implementation using established Go libraries

## Goals / Non-Goals

**Goals:**
- Reliable search via DuckDuckGo HTML endpoint
- Automatic rate limit detection and retry with exponential backoff
- JSON output format for programmatic use
- Configurable retry behavior (max attempts, delays)
- Basic search options (site, region, time limit, safe search)

**Non-Goals:**
- Interactive REPL mode (ddgr has this, we don't need it)
- Browser integration or opening results
- DuckDuckGo Bangs support (can be added later if needed)
- Colored terminal output (skills consume JSON)
- Proxy/Tor support (can be added later if needed)

## Decisions

### 1. CLI Framework: `urfave/cli/v3`

**Choice:** Use `urfave/cli/v3` for CLI parsing.

**Rationale:** Lightweight, well-maintained, simple flag handling. Alternative `cobra` is more powerful but heavier than needed for this tool.

**Alternatives considered:**
- `cobra` - Overkill for single-command tool
- `flag` stdlib - Too basic, lacks help generation

### 2. HTTP Client: `resty/v2`

**Choice:** Use `resty/v2` for HTTP requests with built-in retry support.

**Rationale:** Has retry middleware, easy request customization, and response handling. Integrates well with rate limit handling.

**Alternatives considered:**
- `net/http` stdlib - Requires manual retry implementation
- `req` - Less mature ecosystem

### 3. HTML Parsing: `goquery`

**Choice:** Use `PuerkitoBio/goquery` for HTML parsing.

**Rationale:** jQuery-like API, well-established, handles malformed HTML well (common in search results).

**Alternatives considered:**
- `colly` - Framework overhead, more suited for crawlers
- `x/net/html` - Too low-level

### 4. Package Structure

```
ddg-search/
├── cmd/
│   └── ddg-search/
│       └── main.go      # CLI entry point and command definitions
├── internal/
│   ├── search/
│   │   ├── client.go    # HTTP client with retry
│   │   ├── parser.go    # HTML result parsing
│   │   └── search.go    # Search orchestration
│   └── config/
│       └── config.go    # Configuration types
├── go.mod
└── go.sum
```

### 5. Rate Limit Detection Strategy

**Detection triggers:**
- HTTP status code 429 (Too Many Requests)
- Empty results page when keywords should return results
- HTTP errors (5xx, connection failures)

**Retry strategy:**
- Exponential backoff: `delay = baseDelay * (multiplier ^ attempt)`
- Jitter: Add random 0-500ms to prevent thundering herd
- Configurable: max retries (default: 3), base delay (default: 1s), max delay (default: 30s)

### 6. DuckDuckGo HTML Endpoint

**URL:** `https://html.duckduckgo.com/html/`

**Parameters:**
- `q` - Search query
- `kl` - Region (e.g., `us-en`, `uk-en`)
- `df` - Time filter (`d`=day, `w`=week, `m`=month, `y`=year)
- `p` - Safe search (`1`=on, `-1`=off)
- `sites` - Site-specific search (appended to query)

**Result parsing:**
- Results are in `<div class="result__body">` elements
- Title in `<a class="result__a">`
- URL in `<a class="result__url">`
- Snippet in `<a class="result__snippet">`

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| DuckDuckGo HTML structure changes | Parser isolates CSS selectors; easy to update |
| Aggressive rate limiting blocks all requests | Configurable retry limits; fail gracefully with clear error |
| HTML parsing fails on malformed responses | `goquery` handles malformed HTML; validate parsed results |
| Rate limit detection false positives | Only retry on 429/5xx/empty-results, not on valid empty queries |

## Open Questions

(None - all resolved)

**Resolved:**
- Binary name: `ddg-search`
- JSON output: Just results array (no metadata wrapper)

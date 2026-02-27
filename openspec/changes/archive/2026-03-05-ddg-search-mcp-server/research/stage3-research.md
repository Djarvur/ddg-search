# Stage 3 Research: Internal Library Analysis

## Overview

This document analyzes the three internal libraries that will be used by the MCP server:
1. `internal/search` - DuckDuckGo search client
2. `internal/dump` - Web page fetching and HTML-to-markdown conversion
3. `internal/perplexity` - Perplexity search client

## 1. internal/search Library

### Package: `search`

### Main Type: `Searcher`

```go
type Searcher struct {
    client      *Client
    parser      *Parser
    retryOpts   config.RetryOptions
    debugWriter io.Writer
}
```

### Constructor

```go
func NewSearcher(retryOptions config.RetryOptions) *Searcher
```

**Parameters:**
- `retryOptions` - Configuration for retry behavior

**Returns:**
- `*Searcher` - A new searcher instance

### Main Methods

#### `Search(ctx context.Context, opts config.SearchOptions) ([]config.Result, error)`

Performs a DuckDuckGo search with the given options.

**Parameters:**
- `ctx` - Context for cancellation
- `opts` - Search options (see below)

**Returns:**
- `[]config.Result` - Array of search results
- `error` - Error if search fails

**Behavior:**
- Returns empty array if query is empty
- Implements automatic retry with exponential backoff
- Handles rate limiting by detecting HTML indicators
- Returns `ErrMaxRetries` if all retries exhausted

#### `SearchMarkdown(ctx context.Context, opts config.SearchOptions) (string, error)`

Performs a search and returns Markdown-formatted results.

**Parameters:**
- `ctx` - Context for cancellation
- `opts` - Search options

**Returns:**
- `string` - Markdown formatted results
- `error` - Error if search fails

**Format:**
```
1. [Title](URL)
   Snippet

2. [Title](URL)
   Snippet
...
```

#### `Close()`

Releases resources (closes HTTP client).

### Configuration Types

#### `config.SearchOptions`

```go
type SearchOptions struct {
    Query      string  // Required: The search string
    MaxResults int     // Optional: Max results (0 = unlimited)
    Site       string  // Optional: Filter to specific domain
    Region     string  // Optional: Search region (e.g., "us-en", "uk-en")
    TimeFilter string  // Optional: Time filter ("d", "w", "m", "y")
    SafeSearch bool    // Optional: Enable safe search
}
```

#### `config.RetryOptions`

```go
type RetryOptions struct {
    MaxRetries        int           // Default: 3
    BaseDelay         time.Duration // Default: 1s
    MaxDelay          time.Duration // Default: 30s
    BackoffMultiplier float64       // Default: 2.0
    Debug             bool          // Enable verbose logging
}
```

#### `config.Result`

```go
type Result struct {
    Title   string `json:"title"`
    URL     string `json:"url"`
    Snippet string `json:"snippet"`
}
```

### Errors

- `ErrMaxRetries` - All retry attempts exhausted
- `ErrRateLimited` - Rate limit detected (retryable)

### Usage Example

```go
import (
    "context"
    "github.com/Djarvur/ddg-search/internal/config"
    "github.com/Djarvur/ddg-search/internal/search"
)

func main() {
    searcher := search.NewSearcher(config.DefaultRetryOptions())
    defer searcher.Close()

    opts := config.SearchOptions{
        Query:      "golang mcp server",
        MaxResults: 10,
        SafeSearch: false,
    }

    results, err := searcher.Search(context.Background(), opts)
    if err != nil {
        // handle error
    }

    for _, r := range results {
        fmt.Printf("%s: %s\n", r.Title, r.URL)
    }
}
```

## 2. internal/dump Library

### Package: `dump`

### Main Functions

#### `ValidateURL(rawURL string) (*url.URL, error)`

Parses and validates the URL.

**Parameters:**
- `rawURL` - The URL string to validate

**Returns:**
- `*url.URL` - Parsed URL
- `error` - Error if URL is invalid

**Validation:**
- Must be valid URL format
- Must use HTTP or HTTPS scheme
- Must have a host

**Errors:**
- `ErrInvalidURL` - Invalid URL format
- `ErrUnsupportedScheme` - Scheme not HTTP/HTTPS

#### `Fetch(ctx context.Context, parsedURL *url.URL, cfg Config) (string, error)`

Retrieves the HTML content from the given URL.

**Parameters:**
- `ctx` - Context for cancellation
- `parsedURL` - Validated URL
- `cfg` - Configuration

**Returns:**
- `string` - HTML content
- `error` - Error if fetch fails

**Behavior:**
- Follows redirects (max 10)
- Returns error for HTTP status >= 400
- Uses configurable timeout and user agent

**Errors:**
- Network errors
- `ErrHTTPError` - HTTP error (4xx, 5xx)

#### `Convert(html string) (string, error)`

Transforms HTML content to markdown.

**Parameters:**
- `html` - HTML content

**Returns:**
- `string` - Markdown content
- `error` - Error if conversion fails

#### `FetchAndConvert(ctx context.Context, rawURL string, cfg Config) (string, error)`

Convenience function that fetches a URL and converts HTML to markdown.

**Parameters:**
- `ctx` - Context for cancellation
- `rawURL` - The URL to fetch
- `cfg` - Configuration

**Returns:**
- `string` - Markdown content
- `error` - Error if any step fails

### Configuration

#### `Config`

```go
type Config struct {
    Timeout   time.Duration // Default: 30s
    UserAgent string        // Default: "page-dump/1.0"
}
```

#### `DefaultConfig() Config`

Returns default configuration.

### Constants

```go
const (
    DefaultTimeout     = 30 * time.Second
    DefaultUserAgent   = "page-dump/1.0"
    MaxRedirects       = 10
    HTTPErrorThreshold = 400
)
```

### Usage Example

```go
import (
    "context"
    "github.com/Djarvur/ddg-search/internal/dump"
)

func main() {
    cfg := dump.DefaultConfig()

    markdown, err := dump.FetchAndConvert(context.Background(), "https://example.com", cfg)
    if err != nil {
        // handle error
    }

    fmt.Println(markdown)
}
```

## 3. internal/perplexity Library

### Package: `perplexity`

### Main Type: `Client`

The client is created via `NewClient(accessToken string)` (from client.go).

### Main Method

#### `Search(ctx context.Context, query string, maxResults int, model string) (*SearchResults, error)`

Performs a web search using the Perplexity API.

**Parameters:**
- `ctx` - Context for cancellation
- `query` - The search query (required)
- `maxResults` - Maximum results (currently unused in implementation)
- `model` - Perplexity model to use

**Returns:**
- `*SearchResults` - Search results with AI-generated answer
- `error` - Error if search fails

**Behavior:**
- Returns `ErrQueryEmpty` if query is empty
- Returns `ErrAPI` if API returns an error
- No automatic retry (as per design)

### Types

#### `SearchOptions`

```go
type SearchOptions struct {
    Query      string // Required: The search string
    MaxResults int    // Optional: Max results
    Model      string // Optional: Perplexity model
}
```

#### `SearchResults`

```go
type SearchResults struct {
    Query      string      // Original search query
    Answer     string      // AI-generated answer
    Citations  []string    // URLs cited in the answer
    References []Reference // Formatted citation references
}
```

#### `Reference`

```go
type Reference struct {
    Index int    // Citation number
    URL   string // Citation URL
}
```

### Method

#### `Markdown() string`

Returns the search results formatted as markdown.

**Format:**
```
<AI-generated answer>

## Sources

1. <URL>
2. <URL>
...
```

### Errors

- `ErrQueryEmpty` - Query parameter is empty
- `ErrAPI` - API returned an error

### Usage Example

```go
import (
    "context"
    "github.com/Djarvur/ddg-search/internal/perplexity"
)

func main() {
    client := perplexity.NewClient("your-access-token")

    results, err := client.Search(context.Background(), "golang mcp server", 10, "sonar-small-online")
    if err != nil {
        // handle error
    }

    fmt.Println(results.Markdown())
}
```

## Integration Summary for MCP Server

### Search Tool Parameters

Based on the internal libraries, the search tool should support:

| MCP Parameter | Internal Library | Notes |
|---------------|------------------|-------|
| `query` | `SearchOptions.Query` | Required |
| `max_results` | `SearchOptions.MaxResults` | Optional, default 10 |
| `safe_search` | `SearchOptions.SafeSearch` | Optional, default false |
| `region` | `SearchOptions.Region` | Optional (DuckDuckGo only) |
| `time_filter` | `SearchOptions.TimeFilter` | Optional (DuckDuckGo only) |
| `site` | `SearchOptions.Site` | Optional (DuckDuckGo only) |

**Note:** Perplexity doesn't support all these parameters. The MCP server should:
- Pass supported parameters to Perplexity (query, model)
- Pass all parameters to DuckDuckGo
- Handle Perplexity's different response format (AI answer + citations)

### Fetch Tool Parameters

Based on the internal library, the fetch tool should support:

| MCP Parameter | Internal Library | Notes |
|---------------|------------------|-------|
| `url` | `FetchAndConvert` | Required, must be HTTP/HTTPS |

### Error Handling

All three libraries return errors that should be mapped to MCP error responses:
- `internal/search`: `ErrMaxRetries`, `ErrRateLimited`, network errors
- `internal/dump`: `ErrInvalidURL`, `ErrUnsupportedScheme`, `ErrHTTPError`, network errors
- `internal/perplexity`: `ErrQueryEmpty`, `ErrAPI`, network errors

### Context Usage

All libraries accept `context.Context` for cancellation:
- Use the context from the MCP tool call
- Pass through to library calls
- Handle context cancellation gracefully

### Resource Cleanup

- `internal/search`: Call `Searcher.Close()` when done
- `internal/dump`: No cleanup needed (uses resty client per call)
- `internal/perplexity`: No cleanup needed (uses resty client per call)

## Next Steps

1. ✅ Research complete - internal libraries analyzed
2. Proceed to Stage 4: Base Application
3. Begin implementation

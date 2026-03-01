# ddg-search-mcp Design Context

## Existing Packages

### internal/search
- `Searcher` struct with methods: `Search`, `SearchJSON`, `SearchMarkdown`, `Close`
- Uses `client`, `parser`, `retryOpts`, `debugWriter` fields
- Constructor: `NewSearcher`

### internal/perplexity
- `Client` with `Search` method
- `SearchOptions` struct: `Query`, `MaxResults`, `Model`
- `SearchResults` struct: `Query`, `Answer`, `Citations`, `References`
- `Reference` struct: `Index`, `URL`
- `SearchResults.Markdown()` method for formatting

### internal/dump
- Functions: `Fetch`, `Convert`, `FetchAndConvert`
- `Config` struct: `Timeout`, `UserAgent`
- `DefaultConfig()`, `ValidateURL()` functions
- Constants: `DefaultTimeout`, `DefaultUserAgent`, `MaxRedirects`, `HTTPErrorThreshold`
- Errors: `ErrInvalidURL`, `ErrUnsupportedScheme`, `ErrHTTPError`

## MCP Library Options
1. mark3labs/mcp-go - community library
2. modelcontextprotocol/go-sdk - official SDK

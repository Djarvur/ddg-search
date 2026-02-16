# DDG Search Skill

Web search and page content fetching for Claude Code without API keys.

## Installation

This skill requires two binaries to be installed in your PATH.

### Install via go install

```bash
go install github.com/Djarvur/ddg-search/cmd/ddg-search@latest
go install github.com/Djarvur/ddg-search/cmd/page-dump@latest
```

### Verify installation

```bash
which ddg-search
which page-dump

# Test the binaries
ddg-search --help
page-dump --help
```

## Tools

### web-search

Search the web using DuckDuckGo.

**Parameters:**
- `query` (required): The search query string

**Usage:**
```
ddg-search "golang http client best practices"
```

**Output:** Formatted search results with title, URL, and snippet for each result.

**Options:**
- `--max-results int`: Maximum number of results (default: 10)
- `--site string`: Filter results to a specific domain
- `--region string`: Search region (default: "us-en")
- `--time string`: Time filter: d (day), w (week), m (month), y (year)

### page-dump

Fetch a web page and convert it to markdown.

**Parameters:**
- `url` (required): The URL to fetch (HTTP or HTTPS only)

**Usage:**
```
page-dump "https://example.com/article"
```

**Output:** Markdown-formatted content of the web page.

**Options:**
- `--timeout duration`: Request timeout (default: 30s)
- `--user-agent string`: Custom user agent string (default: page-dump/1.0)

## Usage Examples

### Search for information
```
ddg-search "rust async programming guide"
```

### Fetch and read a web page
```
page-dump "https://doc.rust-lang.org/book/ch10-00-generics.html"
```

### Search and then fetch a result
```
# First search
ddg-search "go context timeout"

# Then fetch an interesting result
page-dump "https://pkg.go.dev/context"
```

## Troubleshooting

### "command not found: ddg-search" or "command not found: page-dump"

The binaries are not in your PATH. Ensure `$(go env GOPATH)/bin` is in your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Rate limiting errors

DuckDuckGo may rate-limit requests if you make too many searches in a short time. Wait a few minutes before retrying.

### Timeout errors

If pages are slow to load, increase the timeout:
```
page-dump --timeout 60s "https://slow-example.com"
```

### Invalid URL errors

Ensure the URL:
- Starts with `http://` or `https://`
- Is properly encoded
- Points to a valid web page (not a file download)

### SSL/TLS errors

For sites with certificate issues, this is expected behavior. The tool requires valid TLS certificates for HTTPS connections.

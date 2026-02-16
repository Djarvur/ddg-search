# Page Dump - URL to Markdown Converter

## Why

Skills and tools need a reliable way to fetch web page content and convert it to markdown for LLM consumption. Unlike plain text extraction (like w3m), markdown preserves document structure (headings, links, lists, code blocks) which is essential for understanding content programmatically.

## What Changes

- New CLI tool `page-dump` in the ddg-search repository
- Accepts a URL, fetches the HTML content, converts to markdown, outputs to stdout
- Static HTTP fetch only (no JavaScript rendering in v1)
- Single binary, pure Go implementation
- Markdown output preserves document structure

## Capabilities

### New Capabilities

- `url-fetch`: HTTP client for fetching web page content - GET request with configurable timeout, user-agent, follows redirects. Returns raw HTML body.

- `html-to-markdown`: HTML to Markdown conversion - transforms HTML documents into clean markdown using JohannesKaufmann/html-to-markdown library. Preserves headings, links, lists, code blocks, tables.

### Modified Capabilities

(None - this is a new tool)

## Impact

- **Location**:
  - `cmd/page-dump/` - CLI entry point only (argument parsing, flag handling)
  - `internal/dump/` - Core logic (URL validation, HTTP fetching, markdown conversion)
- **Dependencies**:
  - `github.com/JohannesKaufmann/html-to-markdown/v2` - HTML to markdown conversion
  - `github.com/go-resty/resty/v2` - HTTP client (already in repo)
  - `github.com/urfave/cli/v3` - CLI framework (already in repo)
- **Binary**: `page-dump` built from `cmd/page-dump/main.go`
- **No external dependencies** - single static binary

## Design Decisions

### Static Only (v1)

JavaScript rendering is deferred to a future version. Rationale:
- Many content sites serve meaningful HTML without JS
- Headless browser (chromedp/rod) adds significant complexity and runtime cost
- Can add `--js` flag in v2 without breaking v1 interface

### Markdown vs Plain Text

Using markdown (not plain text like w3m) because:
- Preserves document structure (headings, lists, links)
- Better for LLM consumption - semantic meaning retained
- Links remain clickable/referable
- Code blocks preserved with language hints

### Same Repository

Adding to ddg-search repository because:
- Shared patterns (CLI, HTTP client)
- Related tools (search + content fetching)
- Single release process

## Future Considerations

- `--js` flag for JavaScript rendering via chromedp or rod
- `--wait` timeout for JS-rendered pages
- `--include-images` to preserve image URLs
- Content extraction (readability mode) to get main content only

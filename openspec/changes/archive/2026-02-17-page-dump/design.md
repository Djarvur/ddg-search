## Context

The ddg-search repository currently provides a DuckDuckGo search CLI (`ddg-search`). We're adding a complementary tool (`page-dump`) that fetches and converts web pages to markdown. Both tools share patterns:
- CLI framework: urfave/cli/v3
- HTTP client: go-resty/resty/v2
- Output to stdout for pipeability

The new tool will live alongside ddg-search in the same repository, sharing infrastructure where possible.

## Goals / Non-Goals

**Goals:**
- Simple CLI that fetches static HTML and converts to markdown
- Single binary with no external runtime dependencies
- Clean markdown output suitable for LLM consumption
- Reuse existing patterns from ddg-search where applicable

**Non-Goals:**
- JavaScript rendering (deferred to v2)
- Content extraction/readability mode (future enhancement)
- Image handling beyond preserving URLs
- Authentication or cookies

## Decisions

### 1. HTML-to-Markdown Library

**Decision:** Use [JohannesKaufmann/html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) v2

**Rationale:**
- Most mature and actively maintained Go library (v2.5.0, Nov 2025)
- Handles complex HTML, uses proper parser (not regex)
- Extensible via custom rules if needed
- Includes CLI reference implementation we can learn from

**Alternatives considered:**
- `tomkosm/html-to-markdown` - less mature, fewer features
- Custom implementation - reinventing the wheel, error-prone

### 2. Package Structure

**Decision:** CLI in `cmd/page-dump/`, core logic in `internal/dump/`

**Rationale:**
- Separation of concerns: CLI handling vs business logic
- Core logic can be tested independently
- Logic can be reused by other packages if needed
- Mirrors `cmd/ddg-search/` structure with internal package

```
cmd/
├── ddg-search/
│   └── main.go
└── page-dump/
    └── main.go        # CLI entry point only
internal/
└── dump/
    └── dump.go        # URL validation, fetching, conversion
```

### 3. HTTP Client Configuration

**Decision:** Use resty with sensible defaults, configurable via flags

**Defaults:**
- Timeout: 30 seconds
- User-Agent: `page-dump/1.0`
- Redirect limit: 10

**Flags:**
- `--timeout duration` - request timeout
- `--user-agent string` - custom user agent

### 4. Error Handling

**Decision:** Exit with non-zero code and error message to stderr

**Rationale:**
- CLI convention for scriptability
- Clear separation between output (stdout) and errors (stderr)

### 5. Output Format

**Decision:** Raw markdown to stdout, no metadata wrapper

**Rationale:**
- Pipe-friendly (can redirect to file or pipe to other tools)
- LLM-friendly (no wrapper to parse)
- Simplicity

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Malformed HTML breaks conversion | html-to-markdown library handles gracefully; add error message if conversion fails |
| Large pages consume memory | Document that tool is for typical web pages, not massive documents |
| Rate limiting by target sites | No built-in retry; users can wrap in retry logic if needed |
| Character encoding issues | resty handles common encodings; library handles HTML entities |

## Migration Plan

Not applicable - this is a new tool with no migration path needed.

## Open Questions

None - scope is well-defined for v1. Future enhancements (JS rendering, readability mode) will require additional design when implemented.

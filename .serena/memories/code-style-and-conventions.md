# Code Style and Conventions for ddg-search

## Go Code Style

### Naming Conventions
- **Package names**: Lowercase, single word (e.g., `search`, `config`, `perplexity`)
- **Exported symbols**: PascalCase (e.g., `NewClient`, `SearchOptions`)
- **Private symbols**: camelCase (e.g., `httpClient`, `retryOptions`)
- **Constants**: PascalCase (e.g., `DefaultMaxRetries`, `apiBaseURL`)
- **Interfaces**: Usually simple names ending with behavior (e.g., `Searcher`)

### Error Handling
- Use sentinel errors for common error cases (e.g., `ErrRateLimited`, `ErrMaxRetries`)
- Wrap errors with context using `fmt.Errorf` or `errors.Wrap`
- Return errors as the last return value
- Check errors immediately after function calls

### Struct Design
- Use pointer receivers for methods that modify the struct
- Use value receivers for methods that don't modify the struct
- Exported fields use PascalCase
- Private fields use camelCase

### Functions and Methods
- Keep functions focused and small
- Use descriptive names that indicate purpose
- Constructor functions use `New` prefix (e.g., `NewClient`, `NewSearcher`)
- Methods that return multiple values follow Go conventions (e.g., `result, err`)

### Comments and Documentation
- Use godoc-style comments for exported symbols
- Package comments at the top of each file
- Function/method comments describe what it does, parameters, and return values
- Inline comments for complex logic

### Testing
- Test files named `*_test.go` in the same package
- Use table-driven tests for multiple test cases
- Test both success and error paths
- Use `t.Run()` for subtests
- Integration tests for external API interactions

### Configuration
- Use structs for configuration options
- Provide default values via `Default*` functions or constants
- Use functional options pattern for complex configuration

## Project-Specific Conventions

### Retry Logic
- Both DuckDuckGo and Perplexity clients use exponential backoff with jitter
- Default retry options: `DefaultRetryOptions()` in `internal/config`
- Retry configuration: `MaxRetries`, `BaseDelay`, `MaxDelay`, `BackoffMultiplier`, `Debug`

### Output Formats
- Search results support both JSON and Markdown output
- Markdown output is optimized for LLM consumption
- JSON output for programmatic use

### Rate Limiting
- Automatic detection of rate limit responses (HTTP 202, 429, 5xx)
- Graceful failure after max retries with clear error messages
- Debug mode logs rate limit information to stderr

### CLI Structure
- Use `urfave/cli/v3` for CLI framework
- Commands defined in `cmd/*/main.go`
- Each tool has a `main()` and a `run*()` function
- Version constants defined in each main.go

## Linting Configuration

The project uses golangci-lint with the following notable settings:
- **Version**: v2.9.0
- **Disabled linters**:
  - `exhaustruct` - too noisy for partial struct initialization
  - `ireturn` - returning interfaces is acceptable
  - `varnamelen` - short variable names are idiomatic in Go
  - `tagliatelle` - JSON tags follow external API conventions
  - `dupl` - duplicate detection too noisy for small project
  - `cyclop` - cyclomatic complexity covered by other linters
  - `depguard` - project-internal imports trigger false positives
  - `wsl` - deprecated, replaced by wsl_v5
- **Settings**:
  - `gocognit.min-complexity`: 15
  - `funlen.lines`: 80
  - `funlen.statements`: 50
- **Exclusions**: errcheck for deferred `Close()` calls

## Testing Requirements

- **Coverage threshold**: 50% minimum (enforced in CI)
- **Race detector**: Always enabled in tests (`-race` flag)
- **Coverage mode**: atomic (`-covermode=atomic`)
- **Test timeout**: 5 minutes (configured in .golangci.yml)

## File Organization

- `cmd/` - Entry points for CLI tools
- `internal/` - Internal packages not meant for external use
- `skills/` - Roo skills for automation
- `openspec/` - OpenSpec change management
- `.github/` - CI/CD workflows

## Dependencies

Key external dependencies:
- `github.com/urfave/cli/v3` - CLI framework
- `github.com/go-resty/resty/v2` - HTTP client
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/JohannesKaufmann/html-to-markdown/v2` - Markdown conversion

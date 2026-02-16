# Implementation Tasks

## 1. Project Setup

- [x] 1.1 Initialize Go module with `go mod init github.com/Djarvur/ddg-search`
- [x] 1.2 Add dependencies: `urfave/cli/v3`, `resty/v2`, `PuerkitoBio/goquery`
- [x] 1.3 Create package structure: `internal/search/`, `internal/config/`, `cmd/`
- [x] 1.4 Add `.golangci.yml` configuration
- [x] 1.5 Add `Makefile` with build, test, lint targets

## 2. Core Search Implementation

- [x] 2.1 Create `internal/config/config.go` with search options struct
- [x] 2.2 Create `internal/search/client.go` with HTTP client setup using resty
- [x] 2.3 Create `internal/search/parser.go` with HTML parsing using goquery
- [x] 2.4 Create `internal/search/search.go` with search orchestration
- [x] 2.5 Implement Result struct with Title, URL, Snippet fields
- [x] 2.6 Implement JSON output formatting

## 3. Rate Limit Handling

- [x] 3.1 Implement rate limit detection (HTTP 429, 5xx, empty results)
- [x] 3.2 Implement exponential backoff with jitter
- [x] 3.3 Add configurable retry parameters (max retries, base delay, max delay)
- [x] 3.4 Implement graceful error messages when retries exhausted

## 4. CLI Implementation

- [x] 4.1 Create `cmd/root.go` with urfave/cli command definition
- [x] 4.2 Create `main.go` entry point
- [x] 4.3 Add verbose CLI flags: `--max-results`, `--site`, `--region`, `--time`, `--safe-search`
- [x] 4.4 Add retry configuration flags: `--max-retries`, `--retry-delay`, `--max-retry-delay`
- [x] 4.5 Implement query argument parsing
- [x] 4.6 Add `--help` and `--version` flags
- [x] 4.7 Add `--debug` flag for verbose logging to stderr
- [x] 4.8 Add debug logs for retry attempts, delays, and rate limit detection

## 5. Testing

- [x] 5.1 Add unit tests for HTML parser
- [x] 5.2 Add unit tests for rate limit detection
- [x] 5.3 Add unit tests for exponential backoff calculation
- [x] 5.4 Add integration tests with mock HTTP server
- [x] 5.5 Add test fixtures for DDG HTML responses

## 6. Documentation

- [x] 6.1 Create README.md with usage examples
- [x] 6.2 Add inline code documentation
- [x] 6.3 Verify CLI help text is clear and complete

## 7. Build & Release

- [x] 7.1 Verify build produces `ddg-search` binary
- [x] 7.2 Test binary with real DuckDuckGo queries
- [x] 7.3 Verify rate limit handling works in practice
- [x] 7.4 Test JSON output format with skills integration

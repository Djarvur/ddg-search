# Implementation Tasks

## 1. Setup

- [x] 1.1 Add `github.com/JohannesKaufmann/html-to-markdown/v2` dependency to go.mod
- [x] 1.2 Create `cmd/page-dump/` directory structure
- [x] 1.3 Create main.go with basic CLI skeleton using urfave/cli

## 2. CLI Implementation

- [x] 2.1 Add `--timeout` flag with default 30s
- [x] 2.2 Add `--user-agent` flag with default `page-dump/1.0`
- [x] 2.3 Add `--help` and `--version` flags
- [x] 2.4 Implement URL argument parsing and validation
- [x] 2.5 Implement HTTP-only URL scheme validation

## 3. URL Fetch Implementation

- [x] 3.1 Create HTTP client with resty and configurable timeout
- [x] 3.2 Implement custom user-agent header support
- [x] 3.3 Implement redirect following with 10-hop limit
- [x] 3.4 Implement HTTP error handling (4xx, 5xx)
- [x] 3.5 Implement timeout error handling
- [x] 3.6 Implement too-many-redirects error handling

## 4. HTML to Markdown Conversion

- [x] 4.1 Create converter using html-to-markdown library
- [x] 4.2 Implement heading conversion (h1-h6)
- [x] 4.3 Implement link conversion
- [x] 4.4 Implement list conversion (ordered and unordered)
- [x] 4.5 Implement code block conversion
- [x] 4.6 Implement paragraph and line break handling
- [x] 4.7 Implement table conversion
- [x] 4.8 Implement stripping of script/style tags

## 5. Output and Error Handling

- [x] 5.1 Output markdown to stdout
- [x] 5.2 Output errors to stderr with non-zero exit code
- [x] 5.3 Handle malformed HTML gracefully

## 6. Testing

- [x] 6.1 Add unit tests for URL validation
- [x] 6.2 Add unit tests for HTTP client configuration
- [x] 6.3 Add unit tests for HTML-to-markdown conversion
- [x] 6.4 Add integration tests with mock HTTP server

## 7. Documentation

- [x] 7.1 Add inline code documentation
- [x] 7.2 Verify CLI help text is clear and complete
- [x] 7.3 Update README.md with page-dump usage

## 8. Build & Release

- [x] 8.1 Verify build produces `page-dump` binary
- [x] 8.2 Test binary with real URLs
- [x] 8.3 Verify error handling in practice

# Implementation Tasks

## 1. CLI Command Setup

- [x] 1.1 Create `cmd/perplexity-search/` directory
- [x] 1.2 Create `cmd/perplexity-search/main.go` with basic CLI structure using urfave/cli
- [x] 1.3 Add flags: `--max-results`, `--model` (optional), `--debug`
- [x] 1.4 Implement command action to parse query and call search function

## 2. Perplexity API Client

- [x] 2.1 Create `internal/perplexity/` package directory
- [x] 2.2 Create `internal/perplexity/client.go` with HTTP client setup
- [x] 2.3 Implement API key loading from environment variable (`PERPLEXITY_API_KEY`)
- [x] 2.4 Implement `Search()` function that calls Perplexity API
- [x] 2.5 Add error handling for: missing API key, invalid key, rate limits, network errors
- [x] 2.6 Implement retry logic with exponential backoff for transient errors

## 3. Search Response Handling

- [x] 3.1 Create `internal/perplexity/search.go` with response structs
- [x] 3.2 Define JSON structures for Perplexity API response
- [x] 3.3 Implement markdown formatter for search results
- [x] 3.4 Add citation formatting for sources referenced in results

## 4. Testing

- [x] 4.1 Create `internal/perplexity/client_test.go` with unit tests
- [x] 4.2 Create `internal/perplexity/search_test.go` with unit tests
- [x] 4.3 Add integration test (mock API or use test API key)
- [x] 4.4 Test error scenarios: missing key, invalid key, rate limit

## 5. Agent Skill

- [x] 5.1 Create `skills/perplexity-search/` directory
- [x] 5.2 Create `skills/perplexity-search/skill.yaml` with tool definition
- [x] 5.3 Create `skills/perplexity-search/SKILL.md` with:
  - [x] 5.3.1 Installation instructions (API key setup)
  - [x] 5.3.2 Usage examples
  - [x] 5.3.3 Tool parameters documentation
  - [x] 5.3.4 Troubleshooting section

## 6. Documentation

- [x] 6.1 Update `README.md` with perplexity-search section
- [x] 6.2 Add usage examples for perplexity-search
- [x] 6.3 Document flags and options
- [x] 6.4 Add comparison between ddg-search and perplexity-search
- [x] 6.5 Document API key requirements and rate limits

## 7. Build Configuration

- [x] 7.1 Update `go.mod` if new dependencies are needed
- [x] 7.2 Verify build process works for new command
- [x] 7.3 Test installation via `go install`

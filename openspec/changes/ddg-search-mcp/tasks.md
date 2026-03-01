## 1. Project Setup

- [ ] 1.1 Add mark3labs/mcp-go, cobra, and viper dependencies to go.mod
- [ ] 1.2 Create cmd/ddg-search-mcp/main.go with cobra root command
- [ ] 1.3 Create internal/mcp package directory structure
- [ ] 1.4 Run `mise run lint` and fix any linter issues
- [ ] 1.5 Run `mise run test` and ensure all tests pass

## 2. Configuration

- [ ] 2.1 Create internal/mcp/config/config.go with MCPConfig struct
- [ ] 2.2 Implement config file loading from ~/.config/ddg-search/config.yaml
- [ ] 2.3 Implement environment variable binding (DDG_MCP_* prefix)
- [ ] 2.4 Implement CLI flag bindings in cobra commands
- [ ] 2.5 Write unit tests for config loading and precedence
- [ ] 2.6 Implement TLS configuration validation
- [ ] 2.7 Run `mise run lint` and fix any linter issues
- [ ] 2.8 Run `mise run test` and ensure all tests pass

## 3. MCP Server Infrastructure

- [ ] 3.1 Create internal/mcp/server.go with MCPServer struct
- [ ] 3.2 Implement stdio transport mode
- [ ] 3.3 Implement TCP transport mode with configurable port
- [ ] 3.4 Implement HTTP transport mode with configurable port
- [ ] 3.5 Implement TLS support for TCP/HTTP transport
- [ ] 3.6 Implement server lifecycle (Start, Stop)
- [ ] 3.7 Write unit tests for server initialization
- [ ] 3.8 Run `mise run lint` and fix any linter issues
- [ ] 3.9 Run `mise run test` and ensure all tests pass

## 4. Tool Definitions

- [ ] 4.1 Create internal/mcp/tools.go with tool registration
- [ ] 4.2 Define search tool schema with input validation
- [ ] 4.3 Define fetch tool schema with input validation
- [ ] 4.4 Write unit tests for tool schema validation
- [ ] 4.5 Run `mise run lint` and fix any linter issues
- [ ] 4.6 Run `mise run test` and ensure all tests pass

## 5. Output Formatting

- [ ] 5.1 Create internal/mcp/format.go with formatters
- [ ] 5.2 Implement text output formatter
- [ ] 5.3 Implement JSON output formatter
- [ ] 5.4 Implement server default format configuration
- [ ] 5.5 Implement per-request format override
- [ ] 5.6 Write unit tests for output formatters
- [ ] 5.7 Run `mise run lint` and fix any linter issues
- [ ] 5.8 Run `mise run test` and ensure all tests pass

## 6. Search Tool Implementation

- [ ] 6.1 Create internal/mcp/search.go with SearchTool handler
- [ ] 6.2 Implement query parameter validation
- [ ] 6.3 Implement backend selection logic (auto, ddg, perplexity)
- [ ] 6.4 Integrate internal/search for DuckDuckGo searches
- [ ] 6.5 Integrate internal/perplexity for Perplexity searches
- [ ] 6.6 Implement Perplexity → DuckDuckGo fallback with metadata
- [ ] 6.7 Implement DuckDuckGo-specific options (site, region, time_filter, safe_search)
- [ ] 6.8 Implement Perplexity-specific options (model)
- [ ] 6.9 Implement result limit enforcement (1-50)
- [ ] 6.10 Write unit tests for search tool handler
- [ ] 6.11 Run `mise run lint` and fix any linter issues
- [ ] 6.12 Run `mise run test` and ensure all tests pass

## 7. Fetch Tool Implementation

- [ ] 7.1 Create internal/mcp/fetch.go with FetchTool handler
- [ ] 7.2 Implement URL parameter validation (HTTP/HTTPS only)
- [ ] 7.3 Integrate internal/dump for URL fetching and conversion
- [ ] 7.4 Implement timeout configuration (1-120 seconds)
- [ ] 7.5 Implement custom user agent support
- [ ] 7.6 Write unit tests for fetch tool handler
- [ ] 7.7 Run `mise run lint` and fix any linter issues
- [ ] 7.8 Run `mise run test` and ensure all tests pass

## 8. Signal Handling

- [ ] 8.1 Implement SIGINT handler for immediate stop
- [ ] 8.2 Implement SIGHUP handler for config reload
- [ ] 8.3 Implement atomic config swap for reload
- [ ] 8.4 Write unit tests for signal handlers
- [ ] 8.5 Run `mise run lint` and fix any linter issues
- [ ] 8.6 Run `mise run test` and ensure all tests pass

## 9. End-to-End Tests

- [ ] 9.1 Create internal/e2e/e2e_test.go with test suite setup
- [ ] 9.2 Write stdio transport E2E tests
- [ ] 9.3 Write TCP transport E2E tests
- [ ] 9.4 Write HTTP transport E2E tests
- [ ] 9.5 Write TLS E2E tests
- [ ] 9.6 Write signal handling E2E tests
- [ ] 9.7 Run `mise run lint` and fix any linter issues
- [ ] 9.8 Run `mise run test` and ensure all tests pass

## 10. Finalization

- [ ] 10.1 Run go build and fix any compilation errors
- [ ] 10.2 Run `mise run lint` and fix all linter issues
- [ ] 10.3 Run `mise run test` and ensure all tests pass
- [ ] 10.4 Verify test coverage meets 50% target
- [ ] 10.5 Update README.md with ddg-search-mcp documentation

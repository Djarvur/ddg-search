# Tasks: Add MCP Server

Implementation tasks for building Claude Code compatible MCP server with mark3labs/mcp-go.

## 1. Project Setup

- [ ] 1.1 Create cmd/ddg-search-mcp directory structure
- [ ] 1.2 Add mark3labs/mcp-go dependency to go.mod
- [ ] 1.3 Add cobra, viper, and slog dependencies to go.mod
- [ ] 1.4 Create .openspec.yaml file for change metadata
- [ ] 1.5 Update .mise.toml with new lint/test tasks if needed
- [ ] 1.6 Run `mise run lint` after setup completion
- [ ] 1.7 Run `mise run test` after setup completion

## 2. Configuration System

- [ ] 2.1 Create internal/mcp/config package
- [ ] 2.2 Define Config struct with all configuration fields (server, perplexity, search, dump, logging)
- [ ] 2.3 Implement config loading with priority chain (CLI → ENV → YAML)
- [ ] 2.4 Implement config validation functions (TLS files, port range, log levels)
- [ ] 2.5 Add default configuration values
- [ ] 2.6 Write public config API tests (config_test.go)
- [ ] 2.7 Write internal config tests (config_internal_test.go)
- [ ] 2.8 Run `mise run lint` after config package completion
- [ ] 2.9 Run `mise run test` after config package completion

## 3. Logging System

- [ ] 3.1 Create internal/mcp/logging package
- [ ] 3.2 Initialize slog logger with text handler
- [ ] 3.3 Implement log level configuration (debug, info, warn, error)
- [ ] 3.4 Add contextual logging helpers (request, response, error, fallback)
- [ ] 3.5 Ensure sensitive data (API keys) is not logged
- [ ] 3.6 Write public logging tests (logging_test.go)
- [ ] 3.7 Write internal logging tests (logging_internal_test.go)
- [ ] 3.8 Run `mise run lint` after logging package completion
- [ ] 3.9 Run `mise run test` after logging package completion

## 4. MCP Server Core

- [ ] 4.1 Create internal/mcp/server package
- [ ] 4.2 Implement Server struct with mark3labs/mcp-go integration
- [ ] 4.3 Implement stdio transport initialization
- [ ] 4.4 Implement TCP transport initialization
- [ ] 4.5 Implement TLS configuration and listener setup
- [ ] 4.6 Implement HUP signal handler for config reload
- [ ] 4.7 Implement server lifecycle (start, stop, reload)
- [ ] 4.8 Write public server tests (server_test.go)
- [ ] 4.9 Write internal server tests (server_internal_test.go)
- [ ] 4.10 Run `mise run lint` after server package completion
- [ ] 4.11 Run `mise run test` after server package completion

## 5. Tool Handlers

- [ ] 5.1 Create internal/mcp/tools package
- [ ] 5.2 Implement search tool handler with provider selection logic
- [ ] 5.3 Implement Perplexity client integration with fallback to DuckDuckGo
- [ ] 5.4 Implement DuckDuckGo client integration
- [ ] 5.5 Implement parameter parsing and validation
- [ ] 5.6 Implement JSON and markdown output formatting
- [ ] 5.7 Implement fetch_page tool handler
- [ ] 5.8 Integrate dump package for page fetching and markdown conversion
- [ ] 5.9 Write public tools tests (search_test.go, fetch_page_test.go)
- [ ] 5.10 Write internal tools tests (search_internal_test.go, fetch_page_internal_test.go)
- [ ] 5.11 Run `mise run lint` after tools package completion
- [ ] 5.12 Run `mise run test` after tools package completion

## 6. CLI Integration

- [ ] 6.1 Create cmd/ddg-search-mcp/main.go
- [ ] 6.2 Implement cobra root command with subcommands/flags
- [ ] 6.3 Add --transport flag (stdio/tcp)
- [ ] 6.4 Add --host and --port flags for TCP transport
- [ ] 6.5 Add --config flag for custom config file path
- [ ] 6.6 Add TLS configuration flags (--tls-enabled, --tls-mode, --tls-cert, --tls-key, --tls-ca)
- [ ] 6.7 Add --debug flag for verbose logging
- [ ] 6.8 Implement config loading and binding to viper
- [ ] 6.9 Implement server start with selected transport
- [ ] 6.10 Implement signal handling (SIGINT for stop, SIGHUP for reload)
- [ ] 6.11 Run `mise run lint` after CLI integration completion
- [ ] 6.12 Run `mise run test` after CLI integration completion

## 7. Testing

- [ ] 7.1 Write E2E tests for server lifecycle (main_test.go)
- [ ] 7.2 Test server start/stop on Ctrl-C
- [ ] 7.3 Test HUP signal config reload
- [ ] 7.4 Test stdio transport with tool invocation
- [ ] 7.5 Test TCP transport with tool invocation
- [ ] 7.6 Test TLS modes (key/cert, combined, mTLS)
- [ ] 7.7 Test search tool with Perplexity (success and fallback scenarios)
- [ ] 7.8 Test search tool with DuckDuckGo
- [ ] 7.9 Test fetch_page tool with valid and invalid URLs
- [ ] 7.10 Test configuration loading and validation
- [ ] 7.11 Test logging at all levels (debug, info, warn, error)
- [ ] 7.12 Achieve 50%+ code coverage
- [ ] 7.13 Run `mise run lint` after testing completion
- [ ] 7.14 Run `mise run test` after testing completion

## 8. Documentation

- [ ] 8.1 Create README.md for ddg-search-mcp
- [ ] 8.2 Document installation and usage
- [ ] 8.3 Document configuration options
- [ ] 8.4 Document tool parameters and behavior
- [ ] 8.5 Document TLS configuration examples
- [ ] 8.6 Document Claude Desktop integration
- [ ] 8.7 Document environment variables
- [ ] 8.8 Run `mise run lint` after documentation completion
- [ ] 8.9 Run `mise run test` after documentation completion

## 9. Build and Lint

- [ ] 9.1 Ensure `go build` succeeds without errors
- [ ] 9.2 Run `mise run lint` and fix all linter errors
- [ ] 9.3 Ensure no linter errors remain
- [ ] 9.4 Ask for user approval before adding any nolint comments
- [ ] 9.5 Run `mise test` and ensure all tests pass
- [ ] 9.6 Verify test coverage meets 50%+ target
- [ ] 9.7 Run `mise run lint` after build completion
- [ ] 9.8 Run `mise run test` after build completion

## 10. Final Verification

- [ ] 10.1 Test stdio mode with Claude Desktop
- [ ] 10.2 Test TCP mode with web client
- [ ] 10.3 Test all TLS modes
- [ ] 10.4 Test HUP config reload
- [ ] 10.5 Verify logging output format
- [ ] 10.6 Verify no sensitive data in logs
- [ ] 10.7 Run final `mise run lint` check
- [ ] 10.8 Run final `mise run test` check
- [ ] 10.9 Confirm all requirements met
- [ ] 10.10 Run `mise run lint` after verification completion
- [ ] 10.11 Run `mise run test` after verification completion

## 1. Project Setup

- [ ] 1.1 Add MCP dependencies to go.mod (mark3labs/mcp-go, cobra, viper)
- [ ] 1.2 Create cmd/ddg-search-mcp directory structure
- [ ] 1.3 Create internal/mcp package directory structure
- [ ] 1.4 Create initial main.go with cobra root command
- [ ] 1.5 Verify build succeeds: `go build ./cmd/ddg-search-mcp`

## 2. Configuration System

- [ ] 2.1 Write tests for config loading (config_test.go)
- [ ] 2.2 Implement config structure and loading (config.go)
- [ ] 2.3 Implement config priority (CLI > env > file > default)
- [ ] 2.4 Write tests for config validation
- [ ] 2.5 Implement config validation
- [ ] 2.6 Write tests for config reload on HUP signal
- [ ] 2.7 Implement config reload on HUP signal
- [ ] 2.8 Run `mise lint` and fix any issues

## 3. Logging System

- [ ] 3.1 Write tests for logging setup (logging_test.go)
- [ ] 3.2 Implement logging setup with slog text handler (logging.go)
- [ ] 3.3 Implement log level configuration
- [ ] 3.4 Write tests for request logging
- [ ] 3.5 Implement request logging (debug for success, error for failure)
- [ ] 3.6 Write tests for bad request logging
- [ ] 3.7 Implement bad request logging
- [ ] 3.8 Run `mise lint` and fix any issues

## 4. TLS Configuration

- [ ] 4.1 Write tests for TLS configuration (tls_test.go)
- [ ] 4.2 Implement TLS configuration structure (tls.go)
- [ ] 4.3 Implement separate key/cert file loading
- [ ] 4.4 Implement combined key/cert file loading
- [ ] 4.5 Implement mTLS with CA certificate
- [ ] 4.6 Write tests for TLS validation
- [ ] 4.7 Implement TLS validation (cert/key match, file existence)
- [ ] 4.8 Run `mise lint` and fix any issues

## 5. MCP Server Core

- [ ] 5.1 Write tests for server initialization (server_test.go)
- [ ] 5.2 Implement MCP server initialization (server.go)
- [ ] 5.3 Write tests for tool registration
- [ ] 5.4 Implement tool registration
- [ ] 5.5 Write tests for transport selection
- [ ] 5.6 Implement transport selection (stdio/TCP)
- [ ] 5.7 Write tests for signal handling (HUP, INT)
- [ ] 5.8 Implement signal handling
- [ ] 5.9 Run `mise lint` and fix any issues

## 6. Search Tool

- [ ] 6.1 Write tests for search tool parameter validation (tools_test.go)
- [ ] 6.2 Implement search tool parameter validation
- [ ] 6.3 Write tests for Perplexity search
- [ ] 6.4 Implement Perplexity search integration
- [ ] 6.5 Write tests for Perplexity fallback to DuckDuckGo
- [ ] 6.6 Implement Perplexity fallback logic (rate limit, quota, no API key)
- [ ] 6.7 Write tests for DuckDuckGo search
- [ ] 6.8 Implement DuckDuckGo search integration
- [ ] 6.9 Write tests for search output formatting (JSON/text)
- [ ] 6.10 Implement search output formatting
- [ ] 6.11 Run `mise lint` and fix any issues

## 7. Fetch Tool

- [ ] 7.1 Write tests for fetch tool parameter validation
- [ ] 7.2 Implement fetch tool parameter validation
- [ ] 7.3 Write tests for URL validation
- [ ] 7.4 Implement URL validation
- [ ] 7.5 Write tests for page fetching
- [ ] 7.6 Implement page fetching integration
- [ ] 7.7 Write tests for HTML to markdown conversion
- [ ] 7.8 Implement HTML to markdown conversion
- [ ] 7.9 Write tests for fetch output formatting (JSON/text)
- [ ] 7.10 Implement fetch output formatting
- [ ] 7.11 Run `mise lint` and fix any issues

## 8. Main Entry Point

- [ ] 8.1 Write tests for main function (main_test.go in _test package)
- [ ] 8.2 Implement main function with cobra setup
- [ ] 8.3 Implement CLI flags for all configuration options
- [ ] 8.4 Implement environment variable binding
- [ ] 8.5 Implement config file loading
- [ ] 8.6 Implement server startup with selected transport
- [ ] 8.7 Run `mise lint` and fix any issues

## 9. Integration Tests

- [ ] 9.1 Write integration tests for stdio transport
- [ ] 9.2 Write integration tests for TCP transport
- [ ] 9.3 Write integration tests for TLS connections
- [ ] 9.4 Write integration tests for mTLS connections
- [ ] 9.5 Write integration tests for config reload
- [ ] 9.6 Run `mise lint` and fix any issues

## 10. E2E Tests

- [ ] 10.1 Write E2E test for stdio transport with real MCP client
- [ ] 10.2 Write E2E test for TCP transport with real MCP client
- [ ] 10.3 Write E2E test for search tool with Perplexity fallback
- [ ] 10.4 Write E2E test for fetch tool
- [ ] 10.5 Write E2E test for server stop on Ctrl-C
- [ ] 10.6 Write E2E test for config reload on HUP signal
- [ ] 10.7 Run `mise lint` and fix any issues

## 11. Documentation

- [ ] 11.1 Create README.md for ddg-search-mcp
- [ ] 11.2 Document configuration options with examples
- [ ] 11.3 Document tool usage with examples
- [ ] 11.4 Document TLS setup with examples
- [ ] 11.5 Document Claude Desktop integration
- [ ] 11.6 Run `mise lint` and fix any issues

## 12. Final Verification

- [ ] 12.1 Run `go build ./cmd/ddg-search-mcp` and verify success
- [ ] 12.2 Run `mise lint` and verify no errors
- [ ] 12.3 Run `mise test` and verify all tests pass
- [ ] 12.4 Verify test coverage is 50%+
- [ ] 12.5 Test binary with stdio transport manually
- [ ] 12.6 Test binary with TCP transport manually
- [ ] 12.7 Test config reload manually
- [ ] 12.8 Test TLS connections manually
- [ ] 12.9 Run `mise lint` one final time after phase completion

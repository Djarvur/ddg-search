## 1. Project Setup

- [ ] 1.1 Create `cmd/ddg-search-mcp/` directory with main.go entry point
- [ ] 1.2 Create `internal/mcp/` directory structure (server, transport, config, tools)
- [ ] 1.3 Add `github.com/mark3labs/mcp-go` dependency to go.mod
- [ ] 1.4 Add `github.com/spf13/viper` dependency to go.mod
- [ ] 1.5 Add `github.com/spf13/cobra` dependency to go.mod
- [ ] 1.6 Create `cmd/ddg-search-mcp/main.go` with basic CLI structure
- [ ] 1.7 Run `mise run test` and fix any failures
- [ ] 1.8 Run `mise run lint` and fix any issues

## 2. Configuration Management

- [ ] 2.1 Create `internal/mcp/config/config.go` with ServerConfig struct
- [ ] 2.2 Implement Viper configuration loading from file, env, and CLI
- [ ] 2.3 Add config validation
- [ ] 2.4 Create `internal/mcp/config/config_test.go` with unit tests
- [ ] 2.5 Run `mise run test` and fix any failures
- [ ] 2.6 Run `mise run lint` and fix any issues

## 3. Stdio Transport

- [ ] 3.1 Create `internal/mcp/transport/stdio.go` with StdioTransport struct
- [ ] 3.2 Implement JSON-RPC request parsing from stdin
- [ ] 3.3 Implement JSON-RPC response writing to stdout
- [ ] 3.4 Create `internal/mcp/transport/stdio_test.go` with unit tests
- [ ] 3.5 Run `mise run test` and fix any failures
- [ ] 3.6 Run `mise run lint` and fix any issues

## 4. TCP Transport

- [ ] 4.1 Create `internal/mcp/transport/tcp.go` with TCPTransport struct
- [ ] 4.2 Implement HTTP handler for MCP protocol
- [ ] 4.3 Add connection pooling for concurrent clients
- [ ] 4.4 Create `internal/mcp/transport/tcp_test.go` with unit tests
- [ ] 4.5 Run `mise run test` and fix any failures
- [ ] 4.6 Run `mise run lint` and fix any issues

## 5. TLS Support

- [ ] 5.1 Create `internal/mcp/transport/tls.go` with TLS configuration
- [ ] 5.2 Implement combined cert+key file support
- [ ] 5.3 Implement separate cert and key file support
- [ ] 5.4 Implement mTLS with CA cert support
- [ ] 5.5 Create `internal/mcp/transport/tls_test.go` with unit tests
- [ ] 5.6 Run `mise run test` and fix any failures
- [ ] 5.7 Run `mise run lint` and fix any issues

## 6. Search Tool Handler

- [ ] 6.1 Create `internal/mcp/tools/search.go` with SearchTool struct
- [ ] 6.2 Implement search tool parameter parsing
- [ ] 6.3 Implement Perplexity API integration
- [ ] 6.4 Implement DDG search fallback
- [ ] 6.5 Implement output format conversion (json, text, markdown)
- [ ] 6.6 Create `internal/mcp/tools/search_test.go` with unit tests
- [ ] 6.7 Run `mise run test` and fix any failures
- [ ] 6.8 Run `mise run lint` and fix any issues

## 7. Fetch Tool Handler

- [ ] 7.1 Create `internal/mcp/tools/fetch.go` with FetchTool struct
- [ ] 7.2 Implement fetch tool parameter parsing
- [ ] 7.3 Implement page dump integration
- [ ] 7.4 Implement output format conversion (markdown, text)
- [ ] 7.5 Create `internal/mcp/tools/fetch_test.go` with unit tests
- [ ] 7.6 Run `mise run test` and fix any failures
- [ ] 7.7 Run `mise run lint` and fix any issues

## 8. Server Implementation

- [ ] 8.1 Create `internal/mcp/server/server.go` with Server struct
- [ ] 8.2 Implement tool registration
- [ ] 8.3 Implement request routing
- [ ] 8.4 Implement logging with slog
- [ ] 8.5 Create `internal/mcp/server/server_test.go` with unit tests
- [ ] 8.6 Run `mise run test` and fix any failures
- [ ] 8.7 Run `mise run lint` and fix any issues

## 9. Signal Handling

- [ ] 9.1 Implement SIGINT handler for graceful shutdown
- [ ] 9.2 Implement SIGHUP handler for config reload
- [ ] 9.3 Create `internal/mcp/signal/signal.go` with signal handling
- [ ] 9.4 Create `internal/mcp/signal/signal_test.go` with unit tests
- [ ] 9.5 Run `mise run test` and fix any failures
- [ ] 9.6 Run `mise run lint` and fix any issues

## 10. E2E Testing

- [ ] 10.1 Write E2E tests for stdio transport
- [ ] 10.2 Write E2E tests for TCP transport
- [ ] 10.3 Write E2E tests for TLS transport
- [ ] 10.4 Write E2E tests for signal handling (ctrl-c, SIGHUP)
- [ ] 10.5 Run `mise run test` and fix any failures
- [ ] 10.6 Run `mise run lint` and fix any issues

## 11. Linting

- [ ] 11.1 Run `mise run lint` and fix any issues
- [ ] 11.2 Ensure no linter errors remain
- [ ] 11.3 Add nolint comments only with user approval

## 12. Build and Verification

- [ ] 12.1 Build `ddg-search-mcp` binary
- [ ] 12.2 Test binary with stdio transport
- [ ] 12.3 Test binary with TCP transport
- [ ] 12.4 Test binary with TLS transport
- [ ] 12.5 Verify config file loading
- [ ] 12.6 Verify environment variable overrides
- [ ] 12.7 Verify CLI flag overrides
- [ ] 12.8 Run `mise run test` and fix any failures
- [ ] 12.9 Run `mise run lint` and fix any issues

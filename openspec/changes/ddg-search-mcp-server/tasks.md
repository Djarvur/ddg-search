# Implementation Tasks

## Phase 1: Foundation

- [ ] **T1.1** Create `cmd/ddg-search-mcp` directory
- [ ] **T1.2** Add dependencies to root go.mod: mark3labs/mcp-go, cobra, viper
- [ ] **T1.3** Implement Cobra root command with all CLI flags
- [ ] **T1.4** Implement Viper configuration loading from config.yaml
- [ ] **T1.5** Add environment variable support for all config options
- [ ] **T1.6** Implement logging with slog (debug, info, warn, error levels)
- [ ] **T1.7** Write unit tests for config loading (config_test.go)
- [ ] **T1.8** Run `mise run lint` and fix issues
- [ ] **T1.9** Write unit tests for private config logic (_internal_test.go)

## Phase 2: MCP Server Setup

- [ ] **T2.1** Create MCP server with mark3labs/mcp-go
- [ ] **T2.2** Define `search` tool with snake_case params
- [ ] **T2.3** Define `fetch` tool with snake_case params
- [ ] **T2.4** Add MCP resources for configuration
- [ ] **T2.5** Write handler tests (handlers_test.go)
- [ ] **T2.6** Run `mise run lint` and fix issues

## Phase 3: Transport Layer

- [ ] **T3.1** Implement stdio transport (default)
- [ ] **T3.2** Implement HTTP transport (StreamableHTTP) on configurable port
- [ ] **T3.3** Make transport configurable (stdio vs http)
- [ ] **T3.4** Implement SIGINT/SIGTERM for immediate stop (no graceful shutdown)
- [ ] **T3.5** Implement SIGHUP for config hot-reload (including TLS cert reload)
- [ ] **T3.6** Implement TLS support with cert/key/combined options
- [ ] **T3.7** Implement mTLS support with CA cert and client auth modes
- [ ] **T3.8** Write transport tests (transport_test.go)
- [ ] **T3.9** Write e2e tests for transport startup/shutdown
- [ ] **T3.10** Write e2e tests for TLS/mTLS
- [ ] **T3.11** Run `mise run lint` and fix issues

## Phase 4: Tool Handlers & Provider Logic

- [ ] **T4.1** Implement search handler with all snake_case parameters
- [ ] **T4.2** Implement fetch handler
- [ ] **T4.3** Implement provider auto-selection logic
- [ ] **T4.4** Implement Perplexity→DuckDuckGo fallback on errors
- [ ] **T4.5** Add output format handling (text/JSON as string)
- [ ] **T4.6** Integrate existing internal/search package
- [ ] **T4.7** Integrate existing internal/perplexity package
- [ ] **T4.8** Integrate existing internal/dump package
- [ ] **T4.9** Write handler unit tests
- [ ] **T4.10** Run `mise run lint` and fix issues

## Phase 5: Testing & Polish

- [ ] **T5.1** Complete unit test coverage for public API (_test.go)
- [ ] **T5.2** Complete unit test coverage for private API (_internal_test.go)
- [ ] **T5.3** Write integration tests (integration_test.go)
- [ ] **T5.4** Write e2e tests for:
  - [ ] Stdio mode requests
  - [ ] HTTP mode requests
  - [ ] Signal handling (INT/TERM - immediate stop)
  - [ ] Config hot-reload (HUP)
  - [ ] Provider fallback (success response with fallback info)
  - [ ] All parameters
  - [ ] Output formats
  - [ ] TLS enabled
  - [ ] mTLS with client auth
- [ ] **T5.5** Run full `mise run lint` - fix all issues
- [ ] **T5.6** Run full `mise run test` - all tests pass

## Phase 6: Documentation & Release

- [ ] **T6.1** Add --version flag with proper output
- [ ] **T6.2** Write README for cmd/ddg-search-mcp
- [ ] **T6.3** Update skills documentation (if needed)
- [ ] **T6.4** Test binary builds correctly
- [ ] **T6.5** Verify manual end-to-end testing

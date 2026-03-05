## Manual Test Plan Format

Each manual test plan should be a markdown document with the following structure:

```markdown
# Stage X Manual Test Plan

## Prerequisites
- List of required setup steps (e.g., config file, certificates, etc.)
- Required tools or dependencies

## Test Cases

### Test Case 1: [Test Name]
**Description:** Brief description of what is being tested

**Steps:**
1. Step 1 with specific command or action
2. Step 2 with specific command or action
3. ...

**Expected Result:**
- Expected output or behavior
- Expected log messages (if applicable)

### Test Case 2: [Test Name]
...

## Cleanup
- Steps to clean up after testing (if applicable)
```

## 1. Stage 1: Research

- [x] 1.1 Research and document best Golang libraries for implementing MCP servers
- [x] 1.2 Research and document good Golang MCP server examples
- [x] 1.3 Document findings in openspec/changes/ddg-search-mcp-server/research/stage1-research.md

## 2. Stage 2: MCP Tool Calling Research

- [x] 2.1 Research Claude Code MCP tool calling conventions
- [x] 2.2 Document web search tool calling format (name, arguments, response)
- [x] 2.3 Document fetch tool calling format (name, arguments, response)
- [x] 2.4 Document findings in openspec/changes/ddg-search-mcp-server/research/stage2-research.md

## 3. Stage 3: Internal Library Analysis

- [x] 3.1 Analyze internal/search library parameters and interfaces
- [x] 3.2 Analyze internal/dump library parameters and interfaces
- [x] 3.3 Analyze internal/perplexity library parameters and interfaces
- [x] 3.4 Document findings in openspec/changes/ddg-search-mcp-server/research/stage3-research.md

## 4. Stage 4: Base Application

- [x] 4.1 Create cmd/ddg-search-mcp directory structure
- [x] 4.2 Design and document config file structure (~/.config/ddg-search/config.yaml)
- [x] 4.3 Implement configuration package with Viper (config < env < CLI priority)
- [x] 4.4 Implement HUP signal handler for config reload
- [x] 4.5 Implement logging with log/slog (configurable log level)
- [x] 4.6 Implement signal handling (SIGINT, SIGTERM for shutdown)
- [x] 4.7 Create main.go with Cobra CLI that stays running until Ctrl-C
- [x] 4.8 Write unit tests for configuration package (config_test package)
- [x] 4.9 Write internal tests for configuration package (config_internal_test package)
- [x] 4.10 Write e2e tests for base application functionality
- [x] 4.11 Create manual test plan for Stage 4
- [x] 4.12 Run mise run lint and fix any issues (no nolint without approval)

## 5. Stage 5: MCP Server with Mocked Tools

- [x] 5.1 Add MCP Go library dependency (based on Stage 1 research)
- [x] 5.2 Implement MCP server core with stdio transport
- [x] 5.3 Implement tool registration interface
- [x] 5.4 Implement search tool with mocked responses
- [x] 5.5 Implement fetch tool with mocked responses
- [x] 5.6 Implement request logging in HTTP proxy access log format (debug level)
- [x] 5.7 Write unit tests for MCP server core (mcp_test package)
- [x] 5.8 Write internal tests for MCP server core (mcp_internal_test package)
- [x] 5.9 Write e2e tests for MCP server with mocked tools
- [x] 5.10 Create manual test plan for Stage 5
- [x] 5.11 Run mise run lint and fix any issues (no nolint without approval)

## 6. Stage 6: Real Tool Implementation

- [x] 6.1 Integrate search tool with internal/search library
- [x] 6.2 Integrate fetch tool with internal/dump library
- [x] 6.3 Implement proper error handling for real tool calls
- [x] 6.4 Implement result formatting for real tool responses
- [x] 6.5 Write unit tests for search tool integration (search_tool_test package)
- [x] 6.6 Write internal tests for search tool integration (search_tool_internal_test package)
- [x] 6.7 Write unit tests for fetch tool integration (fetch_tool_test package)
- [x] 6.8 Write internal tests for fetch tool integration (fetch_tool_internal_test package)
- [x] 6.9 Write e2e tests for real tool implementation
- [x] 6.10 Create manual test plan for Stage 6
- [x] 6.11 Run mise run lint and fix any issues (no nolint without approval)

## 7. Stage 7: Perplexity Integration

- [x] 7.1 Integrate search tool with internal/perplexity library
- [x] 7.2 Implement Perplexity token-based authentication
- [x] 7.3 Implement Perplexity rate limit handling (no retries)
- [x] 7.4 Implement automatic fallback to DuckDuckGo on Perplexity errors
- [x] 7.5 Add fallback note in response when falling back
- [x] 7.6 Write unit tests for Perplexity integration (perplexity_test package)
- [x] 7.7 Write internal tests for Perplexity integration (perplexity_internal_test package)
- [x] 7.8 Write e2e tests for Perplexity integration and fallback
- [x] 7.9 Create manual test plan for Stage 7
- [x] 7.10 Run mise run lint and fix any issues (no nolint without approval)

## 8. Stage 8: HTTP SSE Transport

- [x] 8.1 Implement HTTP SSE transport interface
- [x] 8.2 Implement configurable protocol (stdio/http, default stdio)
- [x] 8.3 Implement default bind address (localhost:9100)
- [x] 8.4 Implement SSE endpoint handling
- [x] 8.5 Implement health check endpoint
- [x] 8.6 Implement multiple concurrent connection support
- [x] 8.7 Write unit tests for HTTP SSE transport (http_sse_test package)
- [x] 8.8 Write internal tests for HTTP SSE transport (http_sse_internal_test package)
- [x] 8.9 Write e2e tests for HTTP SSE transport
- [x] 8.10 Create manual test plan for Stage 8
- [x] 8.11 Run mise run lint and fix any issues (no nolint without approval)

## 9. Stage 9: TLS and mTLS Support

- [x] 9.1 Implement TLS enable/disable configuration (default disabled)
- [x] 9.2 Implement TLS certificate and key file loading
- [x] 9.3 Implement mTLS enable/disable configuration
- [x] 9.4 Implement client certificate validation for mTLS
- [x] 9.5 Implement CA certificate validation for mTLS
- [x] 9.6 Implement TLS configuration parameters (min version, cipher suites)
- [x] 9.7 Implement TLS certificate reload on HUP signal
- [x] 9.8 Write unit tests for TLS/mTLS (tls_test package)
- [x] 9.9 Write internal tests for TLS/mTLS (tls_internal_test package)
- [x] 9.10 Write e2e tests for TLS/mTLS
- [x] 9.11 Create manual test plan for Stage 9
- [x] 9.12 Run mise run lint and fix any issues (no nolint without approval)

## 10. Migrate to Streamable HTTP Transport

- [x] 10.1 Research and understand the new Streamable HTTP transport from mcp-go library
- [x] 10.2 Update internal/mcp/server.go to use server.NewStreamableHTTPServer() instead of SSEServer
- [x] 10.3 Change endpoint from /sse and /message to single /mcp endpoint
- [x] 10.4 Update endpoint handling to support POST and GET methods on /mcp
- [x] 10.5 Add support for MCP-Protocol-Version header
- [x] 10.6 Update health check endpoint to work with new transport
- [x] 10.7 Update logging middleware for new endpoint
- [x] 10.8 Update TLS/mTLS configuration for new transport
- [x] 10.9 Write unit tests for Streamable HTTP transport (streamable_http_test package)
- [x] 10.10 Write internal tests for Streamable HTTP transport (streamable_http_internal_test package)
- [x] 10.11 Write e2e tests for Streamable HTTP transport
- [x] 10.12 Create manual test plan for Stage 10
- [x] 10.13 Run mise run lint and fix any issues (no nolint without approval)
- [x] 10.14 Update documentation to reflect new endpoint (/mcp) and transport

## 11. Final Verification

- [x] 10.1 Run all unit tests (mise test)
- [x] 10.2 Run all integration tests
- [x] 10.3 Run all e2e tests
- [x] 10.4 Verify code coverage meets 50% target (coverage: 68.7%)
- [x] 10.5 Run mise run lint and ensure no errors
- [x] 10.6 Run go build and ensure successful compilation
- [x] 10.7 Complete manual testing following all test plans
- [x] 10.8 Update README.md with MCP server documentation
- [x] 10.9 Create example config file documentation (already documented in README.md)

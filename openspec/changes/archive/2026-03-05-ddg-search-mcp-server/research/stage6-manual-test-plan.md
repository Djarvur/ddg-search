# Stage 6: Real Tool Implementation - Manual Test Plan

## Prerequisites

1. Go 1.23+ installed
2. Network connectivity (for real DuckDuckGo searches and web page fetching)
3. Optional: MCP Inspector or Claude Desktop for testing

## Known Limitations

### 1. E2E Test Status

**Status: FIXED** ✅

The E2E tests in `cmd/ddg-search-mcp/e2e_test.go` have been fixed and are now passing.

**Previous Issue**:

The E2E tests were previously failing due to using shell piping with improper newline handling. The MCP stdio server expects each JSON-RPC request to be terminated with a newline character (`\n`) for proper parsing.

**Fix Applied**:

The E2E tests have been rewritten to use proper stdin/stdout pipes with explicit newline separation:
- Each JSON-RPC request is now sent on its own line with proper newline termination
- Tests use `cmd.StdinPipe()` and `cmd.Stdout = &out` for proper I/O handling
- A build tag `//go:build e2e` has been added to separate e2e tests from unit tests

**Current E2E Test Coverage**:

The following E2E tests are now passing:
1. `TestE2E_BaseApplication_Startup` - Server startup and shutdown
2. `TestE2E_MCPServer_ToolRegistration` - Tool registration and listing
3. `TestE2E_MCPServer_SearchTool` - Search tool execution
4. `TestE2E_MCPServer_FetchTool` - Fetch tool execution
5. `TestE2E_MCPServer_UnknownTool` - Unknown tool error handling
6. `TestE2E_MCPServer_SearchToolMissingQuery` - Missing query parameter validation
7. `TestE2E_MCPServer_FetchToolMissingURL` - Missing URL parameter validation
8. `TestE2E_MCPServer_FetchToolInvalidURL` - Invalid URL format validation
9. `TestE2E_MCPServer_SearchToolWithMaxResults` - Search with max_results parameter
10. `TestE2E_MCPServer_SearchToolWithSafeSearch` - Search with safe_search parameter
11. `TestE2E_Config_WithConfigFile` - Configuration file loading
12. `TestE2E_Config_LogLevelCLI` - Log level via CLI flag
13. `TestE2E_Config_InvalidConfig` - Invalid configuration error handling
14. `TestE2E_Config_LogLevelEnvVar` - Log level via environment variable

**Running E2E Tests**:

```bash
# Run all e2e tests
go test -tags=e2e ./cmd/ddg-search-mcp/...

# Run e2e tests with verbose output
go test -v -tags=e2e ./cmd/ddg-search-mcp/... -run TestE2E
```

**Impact**:

E2E tests now fully verify the MCP server workflow and can be run to ensure proper functionality.

### 2. Network Dependencies
Real tool integration requires network access to:
- DuckDuckGo search service
- Target URLs for the fetch tool

**Impact**: Tests will fail if network is unavailable or if external services are down.

### 3. Rate Limiting
DuckDuckGo may rate limit excessive requests. The internal search library includes retry logic, but aggressive testing may trigger rate limits.

**Impact**: Some search requests may fail with rate limit errors during testing.

---

## Test Cases

### Test Case 1: Real Search Tool Integration

**Objective**: Verify the search tool performs real DuckDuckGo searches and returns markdown-formatted results.

**Shell Command**:
```bash
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n' | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server responds with a successful initialization
- Server response includes protocol version and server capabilities
- Server capabilities indicate tools can be listed

**Note**: The `initialize` response does not include the tools list. To see the registered tools, you must send a separate `tools/list` request.

---

### Test Case 1b: List Available Tools

**Objective**: Verify the server can list registered tools.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/list"}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server responds with a tools list
- Response includes "search" tool with its schema
- Response includes "fetch" tool with its schema

---

### Test Case 2: Real Search Tool Integration

**Objective**: Verify the search tool performs real DuckDuckGo searches and returns markdown-formatted results.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming","max_results":5}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server responds with successful search results
- Results are in markdown format with numbered links
- Each result includes title, URL, and snippet
- No errors in server logs

---

### Test Case 2: Search Tool with Safe Search

**Objective**: Verify the safe_search parameter is properly passed to the search library.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"test query","safe_search":true}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Search completes successfully
- Results are filtered according to safe search settings

---

### Test Case 3: Real Fetch Tool Integration

**Objective**: Verify the fetch tool retrieves and converts web pages to markdown.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server responds with markdown content
- Content is properly converted from HTML
- No errors in server logs

---

### Test Case 4: Search Tool Error Handling - Missing Query

**Objective**: Verify proper error handling for missing query parameter.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server returns error response
- Error message indicates "query parameter is required"
- Server continues running (does not crash)

---

### Test Case 5: Search Tool Error Handling - Empty Query

**Objective**: Verify proper error handling for empty query parameter.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":""}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server returns error response
- Error message indicates "query parameter is required"
- Server continues running

---

### Test Case 6: Fetch Tool Error Handling - Missing URL

**Objective**: Verify proper error handling for missing URL parameter.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fetch","arguments":{}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server returns error response
- Error message indicates "url parameter is required"
- Server continues running

---

### Test Case 7: Fetch Tool Error Handling - Invalid URL Scheme

**Objective**: Verify proper error handling for invalid URL scheme.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"ftp://example.com"}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server returns error response
- Error message indicates "fetch failed"
- Error includes details about invalid URL scheme
- Server continues running

---

### Test Case 8: Multiple Tool Calls

**Objective**: Verify the server handles multiple tool calls correctly.

**Shell Command**:
```bash
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang","max_results":3}}}\n'; printf '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Each tool call is processed independently
- Results are returned for each call in the same order as requests (by `id` field)
- Each response includes the same `id` as the corresponding request for correlation
- No memory leaks or resource exhaustion
- Server remains responsive

---

### Test Case 9: Context Cancellation

**Objective**: Verify context cancellation is handled properly.

**Shell Command**:
```bash
# Start the server in the background and send a request, then close stdin
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang"}}}\n') | timeout 2 go run ./cmd/ddg-search-mcp || true
```

**Expected Results**:
- Server detects stdin closure or timeout
- Server logs context cancellation message
- Server shuts down gracefully
- All resources are cleaned up

---

### Test Case 10: Configuration Integration

**Objective**: Verify configuration is properly loaded and used.

**Create a test config file**:
```bash
cat > /tmp/test-config.yaml <<EOF
logging:
  level: debug
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
EOF
```

**Shell Command**:
```bash
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n' | DDG_SEARCH_CONFIG_FILE=/tmp/test-config.yaml go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server loads configuration from file
- Debug logging is enabled
- Configuration values match the file contents (max_results=10, safe_search=false)

**Note**: ~~The `--config` flag is now supported for specifying a custom config file path. Configuration is loaded via the `DDG_SEARCH_CONFIG_FILE` environment variable.~~ (IRRELEVANT)
---

### Test Case 11: Rate Limit Handling

**Objective**: Verify rate limit errors are handled gracefully.

**Shell Command**:
```bash
# Send multiple rapid search requests
for i in {1..15}; do
  printf '{"jsonrpc":"2.0","id":'"$i"',"method":"tools/call","params":{"name":"search","arguments":{"query":"test query '"$i"'"}}}\n'
done | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Some requests may fail with rate limit errors
- Error messages are clear and informative
- Server continues to function after rate limit
- Retry logic in search library is engaged

---

### Test Case 12: Network Error Handling

**Objective**: Verify network errors are handled gracefully.

**Shell Command**:
```bash
# Use a non-existent domain to simulate a network error
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://this-domain-definitely-does-not-exist-12345.com"}}}\n') | go run ./cmd/ddg-search-mcp
```

**Expected Results**:
- Server returns error response
- Error message indicates network failure
- Server continues running
- No crash or panic

---

## Manual Testing Options

### Option 1: Using Pipes (Simplest for Manual Testing)

```bash
# Single request
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}\n' | go run ./cmd/ddg-search-mcp

# Multiple requests
(printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}\n'; printf '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang"}}}\n') | go run ./cmd/ddg-search-mcp
```

**Important Note**: Always use `printf` with explicit `\n` at the end of each JSON-RPC request. The `echo` command does not add a newline, which causes the MCP server's JSON-RPC parser to fail with a "Parse error".

### Option 2: Using a Request File

```bash
# Create a request file
cat > /tmp/mcp-requests.txt <<EOF
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang"}}}
EOF

# Send requests from file
cat /tmp/mcp-requests.txt | go run ./cmd/ddg-search-mcp
```

### Option 3: Using Claude Desktop (Recommended for Integration Testing)

1. Configure Claude Desktop to use the MCP server
2. Add to Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):
   ```json
   {
     "mcpServers": {
       "ddg-search": {
         "command": "/path/to/ddg-search-mcp",
         "args": ["--config", "/path/to/config.yaml"]
       }
     }
   }
   ```
3. Restart Claude Desktop
4. Use the search and fetch tools through Claude

### Option 4: Using MCP Inspector

```bash
# Install MCP Inspector
npm install -g @modelcontextprotocol/inspector

# Run the inspector with your server
mcp-inspector go run ./cmd/ddg-search-mcp -- go run ./cmd/ddg-search-mcp
```

### Option 5: Using Unit Tests

```bash
# Run short tests (skip integration tests that require network)
go test -v ./internal/mcp/... -short

# Run all tests (including integration tests that require network)
go test -v ./internal/mcp/...
```

---

## Cleanup

After testing:

```bash
# Remove test configuration files
rm -f /tmp/test-config.yaml
rm -f /tmp/mcp-requests.txt

# Verify no processes are left running
ps aux | grep ddg-search-mcp

# Kill any lingering processes
pkill -f ddg-search-mcp
```

---

## Success Criteria

Stage 6 is considered complete when:

1. ✅ Search tool performs real DuckDuckGo searches
2. ✅ Fetch tool retrieves and converts web pages
3. ✅ Error handling works correctly for both tools
4. ✅ Configuration is properly loaded and used
5. ✅ Context cancellation is handled gracefully
6. ✅ Multiple tool calls work correctly
7. ✅ Network errors are handled gracefully
8. ✅ Rate limiting is handled appropriately
9. ✅ Unit tests pass (all tests in `internal/mcp/tools_test.go`)
10. ✅ E2E tests pass (all tests in `cmd/ddg-search-mcp/e2e_test.go` with `-tags=e2e`)
11. ✅ Manual testing confirms expected behavior

---

## Notes

- E2E tests are failing due to a test implementation issue (using `echo` without newlines), not a library limitation
- Manual testing is the primary verification method for Stage 6
- Integration tests that require network access are skipped in short mode
- The real tool integration is complete and working correctly
- All shell commands can be directly copied and executed
- Always use `printf` with explicit `\n` when sending JSON-RPC requests via pipes

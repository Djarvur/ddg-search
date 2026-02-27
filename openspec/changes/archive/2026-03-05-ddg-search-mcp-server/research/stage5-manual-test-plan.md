# Stage 5 Manual Test Plan

## Prerequisites

- Go 1.21+ installed
- Project cloned and dependencies installed (`go mod download`)
- Terminal access for running commands
- Claude Code or another MCP client for testing tool calls

## Known Issues and Limitations

### 1. Stdio Transport Limitations

The MCP server currently only supports stdio transport, which means:
- You cannot connect to the server from another terminal
- The server must be started by the MCP client (e.g., Claude Code)
- Manual testing requires using an MCP client that communicates via stdin/stdout
- The server exits when stdin is closed (client disconnects) - this is expected behavior
- The server exits when interrupted with Ctrl-C - this is expected behavior

### 2. Shutdown Behavior

When the server is interrupted with Ctrl-C:
- The shutdown signal is received and logged
- The "context canceled" error from the MCP library is handled gracefully (not treated as an error)
- The server shuts down cleanly without error messages
- This is expected behavior for stdio-based MCP servers

### 3. Parse Errors on Startup

When the server starts without an MCP client connected:
- You may see "Parse error" messages in the output
- These occur because the MCP server is reading from stdin and receiving invalid JSON
- This is expected behavior and does not indicate a problem

### 4. Debug Logging

Debug logs are only visible when:
- The log level is set to "debug" (via config or CLI flag)
- Tool calls are actually made
- Startup logs are always at INFO level regardless of configuration

## Test Cases

### Test Case 1: MCP Server Initialization

**Description:** Verify the MCP server initializes correctly with stdio transport.

**Steps:**
1. Open a terminal
2. Navigate to the project directory
3. Run: `go run ./cmd/ddg-search-mcp`
4. Observe the output
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts without errors
- Output shows "Starting ddg-search-mcp"
- Output shows "Configuration loaded"
- Output shows "Starting MCP server"
- Output shows "Server ready, waiting for signals..."
- On Ctrl-C, output shows "Received shutdown signal"
- Output shows "Shutting down..."
- Output shows "Shutdown complete"
- Application exits cleanly without errors

**Note:** The "context canceled" error from the MCP library is now handled gracefully and not treated as an error. The server shuts down cleanly when interrupted.

### Test Case 2: Tool Registration

**Description:** Verify the search and fetch tools are registered.

**Steps:**
1. Send an initialize request to the server:
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}' | go run ./cmd/ddg-search-mcp
   ```
2. Send a tools/list request:
   ```bash
   echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | go run ./cmd/ddg-search-mcp
   ```
3. Observe the output

**Expected Result:**
- Server responds to initialize request
- Server responds with tool list
- Tool list includes "search" tool
- Tool list includes "fetch" tool
- Each tool has a description
- Each tool has an input schema
- Server exits cleanly when stdin is closed

**Note:** The MCP server should exit cleanly when stdin is closed (after the echo command completes). This is the expected behavior for stdio-based MCP servers.

### Test Case 3: Search Tool Call

**Description:** Verify the search tool responds with mocked results.

**Steps:**
1. Send a search tool call:
   ```bash
   echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}' | go run ./cmd/ddg-search-mcp
   ```
2. Observe the response

**Expected Result:**
- Tool call succeeds
- Response contains mocked search results
- Results include title, URL, and snippet
- Response is in text format
- Server exits cleanly when stdin is closed

### Test Case 4: Search Tool with Missing Query

**Description:** Verify the search tool validates required parameters.

**Steps:**
1. Send a search tool call without the query parameter:
   ```bash
   echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"search","arguments":{}}}' | go run ./cmd/ddg-search-mcp
   ```
2. Observe the response

**Expected Result:**
- Tool call fails
- Error message indicates query parameter is required
- Server exits cleanly when stdin is closed

### Test Case 5: Fetch Tool Call

**Description:** Verify the fetch tool responds with mocked content.

**Steps:**
1. Send a fetch tool call:
   ```bash
   echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}' | go run ./cmd/ddg-search-mcp
   ```
2. Observe the response

**Expected Result:**
- Tool call succeeds
- Response contains mocked page content
- Response includes source URL
- Response is in markdown format
- Server exits cleanly when stdin is closed

### Test Case 6: Fetch Tool with Missing URL

**Description:** Verify the fetch tool validates required parameters.

**Steps:**
1. Send a fetch tool call without the URL parameter:
   ```bash
   echo '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"fetch","arguments":{}}}' | go run ./cmd/ddg-search-mcp
   ```
2. Observe the response

**Expected Result:**
- Tool call fails
- Error message indicates URL parameter is required
- Server exits cleanly when stdin is closed

### Test Case 7: Request Logging

**Description:** Verify tool calls are logged in Combined Log Format.

**Steps:**
1. Send a search tool call with debug logging:
   ```bash
   echo '{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}' | go run ./cmd/ddg-search-mcp --log-level=debug
   ```
2. Observe the log output

**Expected Result:**
- Log shows tool call in Combined Log Format
- Log includes timestamp, tool name, status, and bytes
- Log format: `[timestamp] [stdio] [CALL] [tool_name] [status] [bytes] [-] [mcp-client]`
- Server exits cleanly when stdin is closed

### Test Case 8: Graceful Shutdown

**Description:** Verify the server shuts down gracefully.

**Steps:**
1. Start the MCP server: `go run ./cmd/ddg-search-mcp`
2. Press Ctrl-C
3. Observe the output

**Expected Result:**
- Server receives shutdown signal
- Output shows "Received shutdown signal"
- Output shows "Shutting down..."
- Output shows "Shutdown complete"
- Application exits cleanly without errors

**Note:** The "context canceled" error from the MCP library is now handled gracefully and not treated as an error. The server shuts down cleanly when interrupted.

### Test Case 9: Multiple Tool Calls

**Description:** Verify the server handles multiple tool calls.

**Steps:**
1. Send multiple tool calls in sequence:
   ```bash
   { echo '{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"search","arguments":{"query":"test1"}}}'; \
     echo '{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"search","arguments":{"query":"test2"}}}'; \
     echo '{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example1.com"}}}'; \
     echo '{"jsonrpc":"2.0","id":12,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example2.com"}}}'; } | go run ./cmd/ddg-search-mcp
   ```
2. Observe all responses

**Expected Result:**
- All tool calls succeed
- Each call returns appropriate mocked response
- Server remains stable
- Server exits cleanly when stdin is closed

### Test Case 10: Configuration Integration

**Description:** Verify the MCP server respects configuration.

**Steps:**
1. Create config file with debug logging (e.g., `~/.config/ddg-search/config.yaml`):
   ```yaml
   server:
     protocol: stdio
     bind_address: localhost:9100
     tls_enabled: false
     mtls_enabled: false
   logging:
     level: debug
   search:
     max_results: 10
     safe_search: false
   perplexity:
     enabled: false
   ```
2. Send a tool call to see debug logs:
   ```bash
   echo '{"jsonrpc":"2.0","id":13,"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}' | go run ./cmd/ddg-search-mcp
   ```
3. Observe the log output

**Expected Result:**
- Server uses configured log level (shown in "Configuration loaded" log)
- Debug logs are visible when tool calls are made
- Startup logs are always at INFO level regardless of configuration
- Server exits cleanly when stdin is closed

## MCP Client Testing

Since the MCP server uses stdio transport, you cannot connect to it from another terminal. To test tool registration and tool calls, you need to use an MCP client that communicates via stdin/stdout.

### Option 1: Using Pipes (Simplest for Manual Testing)

Send JSON-RPC messages directly to the server using pipes:

```bash
# Initialize the server
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}' | go run ./cmd/ddg-search-mcp

# List available tools
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | go run ./cmd/ddg-search-mcp

# Call the search tool
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}' | go run ./cmd/ddg-search-mcp

# Call the fetch tool
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}' | go run ./cmd/ddg-search-mcp
```

### Option 2: Using Claude Code (Recommended for Integration Testing)

Claude Code has built-in MCP client support. To test the server:

1. Configure Claude Code to use the MCP server by adding it to your Claude Code configuration
2. The server will be automatically started when Claude Code needs to use it
3. Tool calls will be made automatically when you ask Claude to search or fetch content

### Option 3: Using the MCP Inspector

The MCP Inspector is a tool for testing MCP servers:

```bash
# Install MCP Inspector
npm install -g @modelcontextprotocol/inspector

# Run the inspector with your server
mcp-inspector go run ./cmd/ddg-search-mcp
```

### Option 4: Using E2E Tests

The project includes E2E tests that can be run:

```bash
go test ./cmd/ddg-search-mcp -v -run TestE2E
```

## Cleanup

After testing:
1. Stop any running MCP server instances
2. Remove test config files if created
3. Verify no test processes are running: `ps aux | grep ddg-search-mcp`

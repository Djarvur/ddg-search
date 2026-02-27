# Stage 8 Manual Test Plan: Streamable HTTP Transport

## Prerequisites

- Go 1.21+ installed
- Project cloned and dependencies installed (`go mod download`)
- Terminal access for running commands
- Text editor for creating config files
- `curl` or similar HTTP client for testing endpoints
- `nc` or similar tool for checking port availability

## Test Cases

### Test Case 1: HTTP Server Startup

**Description:** Verify the application starts successfully with HTTP protocol.

**Steps:**
1. Create config file: `~/.config/ddg-search/config.yaml` with content:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Press Ctrl-C to stop

**Expected Result:**
- Application starts without errors
- Output shows "Starting ddg-search-mcp"
- Output shows "protocol=http"
- Output shows "bind_address=localhost:9100"
- Output shows "HTTP Streamable HTTP server listening"
- Application stops cleanly on Ctrl-C

### Test Case 2: Health Check Endpoint

**Description:** Verify the health check endpoint returns OK when server is running.

**Steps:**
1. Start server with HTTP protocol (from Test Case 1)
2. In another terminal, run: `curl http://localhost:9100/health`
3. Observe the response
4. Stop the server with Ctrl-C

**Expected Result:**
- HTTP response status is 200 OK
- Response body contains "OK"
- Response is quick (< 100ms)

### Test Case 3: Default Bind Address

**Description:** Verify the default bind address is localhost:9100.

**Steps:**
1. Create config file without bind_address:
   ```yaml
   server:
     protocol: http
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Test health check: `curl http://localhost:9100/health`
5. Stop the server

**Expected Result:**
- Application starts successfully
- Output shows "bind_address=localhost:9100"
- Health check returns 200 OK

### Test Case 4: Custom Bind Address

**Description:** Verify custom bind address works correctly.

**Steps:**
1. Create config file with custom address:
   ```yaml
   server:
     protocol: http
     bind_address: 127.0.0.1:9200
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Test health check: `curl http://127.0.0.1:9200/health`
5. Stop the server

**Expected Result:**
- Application starts successfully
- Output shows "bind_address=127.0.0.1:9200"
- Health check returns 200 OK

### Test Case 5: Port Already in Use

**Description:** Verify the application fails gracefully when port is already in use.

**Steps:**
1. Start server on port 9100
2. In another terminal, try to start another server on same port:
   ```bash
   DDG_SEARCH_SERVER_BIND_ADDRESS=localhost:9100 go run ./cmd/ddg-search-mcp
   ```
3. Observe the output
4. Stop both servers

**Expected Result:**
- Second server fails to start
- Error message indicates "address already in use"
- Application exits with non-zero status

### Test Case 6: Protocol Configuration - stdio

**Description:** Verify stdio protocol still works as default.

**Steps:**
1. Remove or rename config file
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Press Ctrl-C to stop

**Expected Result:**
- Application starts successfully
- Output shows "protocol=stdio"
- No HTTP server is started
- Application stops cleanly

### Test Case 7: Protocol Configuration - http

**Description:** Verify HTTP protocol is correctly configured.

**Steps:**
1. Create config file:
   ```yaml
   server:
     protocol: http
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Test health check: `curl http://localhost:9100/health`
5. Stop the server

**Expected Result:**
- Application starts successfully
- Output shows "protocol=http"
- HTTP server is listening
- Health check returns 200 OK

### Test Case 8: Invalid Protocol

**Description:** Verify the application fails with invalid protocol.

**Steps:**
1. Create config file:
   ```yaml
   server:
     protocol: invalid
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output

**Expected Result:**
- Application fails to start
- Error message indicates "invalid protocol"
- Error message shows valid options (stdio, http)

### Test Case 9: MCP Endpoint Access (GET for SSE)

**Description:** Verify the /mcp endpoint is accessible via GET with SSE headers.

**Steps:**
1. Start server with HTTP protocol
2. In another terminal, run: `curl -v -N http://localhost:9100/mcp -H "Accept: text/event-stream" -H "MCP-Protocol-Version: 2025-06-18"`
3. Observe the response
4. Press Ctrl-C to stop curl
5. Stop the server

**Expected Result:**
- MCP endpoint responds with HTTP 200 OK
- Response includes SSE headers:
  - Content-Type: text/event-stream
  - Cache-Control: no-cache
  - Connection: keep-alive
  - Access-Control-Allow-Origin: *
  - MCP-Protocol-Version: 2025-06-18
- Connection stays open (until interrupted)

### Test Case 10: Multiple Concurrent Connections

**Description:** Verify server handles multiple concurrent connections.

**Steps:**
1. Start server with HTTP protocol
2. In multiple terminals (3-5), run: `curl http://localhost:9100/health`
3. Observe all responses
4. Stop the server

**Expected Result:**
- All requests complete successfully
- All responses return 200 OK
- No connection errors or timeouts

### Test Case 12: HTTP Request Logging

**Description:** Verify HTTP requests are logged in Combined Log Format.

**Steps:**
1. Create config file with debug logging:
   ```yaml
   server:
     protocol: http
   logging:
     level: debug
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. In another terminal, run: `curl http://localhost:9100/health`
4. Observe the server output
5. Stop the server

**Expected Result:**
- Server logs the HTTP request at debug level
- Log line includes: timestamp, client IP, method, path, status, bytes, referer, user-agent
- Format matches Combined Log Format: `[timestamp] [client] [method] [path] [status] [bytes] [referer] [user-agent]`

### Test Case 13: Graceful Shutdown

**Description:** Verify HTTP server shuts down gracefully.

**Steps:**
1. Start server with HTTP protocol
2. In another terminal, run: `curl http://localhost:9100/health`
3. Verify it returns 200 OK
4. Send SIGTERM to server: `kill -TERM <pid>`
5. Observe the output
6. Try another health check: `curl http://localhost:9100/health`
7. Observe the response

**Expected Result:**
- Server receives SIGTERM
- Output shows "Shutting down..."
- Output shows "Shutdown complete"
- Second health check fails (connection refused)
- Server exits cleanly

### Test Case 14: Environment Variable Override

**Description:** Verify environment variables override config file for HTTP settings.

**Steps:**
1. Create config file:
   ```yaml
   server:
     protocol: stdio
     bind_address: localhost:9100
   logging:
     level: info
   ```
2. Run: `DDG_SEARCH_SERVER_PROTOCOL=http DDG_SEARCH_SERVER_BIND_ADDRESS=localhost:9300 go run ./cmd/ddg-search-mcp`
3. Observe the output
4. Test health check: `curl http://localhost:9300/health`
5. Stop the server

**Expected Result:**
- Application starts successfully
- Output shows "protocol=http" (env override)
- Output shows "bind_address=localhost:9300" (env override)
- Health check returns 200 OK

### Test Case 15: Invalid Bind Address

**Description:** Verify the application fails with invalid bind address.

**Steps:**
1. Create config file:
   ```yaml
   server:
     protocol: http
     bind_address: invalid-address
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe the output

**Expected Result:**
- Application fails to start
- Error message indicates invalid bind address
- Application exits with non-zero status

### Test Case 16: MCP Endpoint Connection Lifecycle

**Description:** Verify MCP endpoint connection lifecycle (connect, disconnect).

**Steps:**
1. Start server with HTTP protocol and debug logging
2. In another terminal, run: `curl -N http://localhost:9100/mcp -H "Accept: text/event-stream" -H "MCP-Protocol-Version: 2025-06-18"`
3. Observe server logs for connection
4. Press Ctrl-C to stop curl
5. Observe server logs for disconnection
6. Stop the server

**Expected Result:**
- Server logs MCP connection establishment
- Server logs MCP connection disconnection
- Logs include client IP and endpoint path

### Test Case 17: Search Tool via POST

**Description:** Verify the search tool works correctly via POST to the /mcp endpoint.

**Steps:**
1. Start server with HTTP protocol
2. In another terminal, send a search request:
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 1,
       "method": "tools/call",
       "params": {
         "name": "search",
         "arguments": {
           "query": "golang"
         }
       }
     }'
   ```
3. Observe the response
4. Stop the server

**Expected Result:**
- Search request is accepted (HTTP 200)
- Response contains search results with titles, URLs, and snippets
- No errors in server logs

### Test Case 18: Fetch Tool via POST

**Description:** Verify the fetch tool works correctly via POST to the /mcp endpoint.

**Steps:**
1. Start server with HTTP protocol
2. In another terminal, send a fetch request:
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 2,
       "method": "tools/call",
       "params": {
         "name": "fetch",
         "arguments": {
           "url": "https://example.com"
         }
       }
     }'
   ```
3. Observe the response
4. Stop the server

**Expected Result:**
- Fetch request is accepted (HTTP 200)
- Response contains the fetched page content (HTML or markdown)
- No errors in server logs

### Test Case 19: Invalid Tool Call via POST

**Description:** Verify the server handles invalid tool calls gracefully.

**Steps:**
1. Start server with HTTP protocol
2. Send an invalid tool call (non-existent tool):
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 3,
       "method": "tools/call",
       "params": {
         "name": "nonexistent_tool",
         "arguments": {}
       }
     }'
   ```
3. Observe the error response
4. Stop the server

**Expected Result:**
- Request is accepted (HTTP 200)
- Error response indicates tool not found
- Server continues running

### Test Case 20: Search with Missing Query Parameter

**Description:** Verify the server validates required parameters.

**Steps:**
1. Start server with HTTP protocol
2. Send a search request without the required query parameter:
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 4,
       "method": "tools/call",
       "params": {
         "name": "search",
         "arguments": {}
       }
     }'
   ```
3. Observe the error response
4. Stop the server

**Expected Result:**
- Request is accepted (HTTP 200)
- Error response indicates missing required parameter
- Server continues running

### Test Case 21: Multiple Tool Calls

**Description:** Verify multiple tool calls work correctly.

**Steps:**
1. Start server with HTTP protocol
2. Send a search request (ID 5):
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 5,
       "method": "tools/call",
       "params": {
         "name": "search",
         "arguments": {
           "query": "golang"
         }
       }
     }'
   ```
3. Wait for response
4. Send a fetch request (ID 6):
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 6,
       "method": "tools/call",
       "params": {
         "name": "fetch",
         "arguments": {
           "url": "https://example.com"
         }
       }
     }'
   ```
5. Wait for response
6. Send another search request (ID 7)
7. Wait for response
8. Stop the server

**Expected Result:**
- All requests are accepted (HTTP 200)
- All responses are received correctly
- Each response has the correct request ID
- Responses are received in order
- No errors in server logs

### Test Case 22: List Tools via POST

**Description:** Verify the tools/list endpoint works via POST to the /mcp endpoint.

**Steps:**
1. Start server with HTTP protocol
2. Send a tools/list request:
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 8,
       "method": "tools/list"
     }'
   ```
3. Observe the response
4. Stop the server

**Expected Result:**
- Request is accepted (HTTP 200)
- Response includes both "search" and "fetch" tools
- Each tool includes name, description, and input schema

### Test Case 23: Initialize via POST

**Description:** Verify the initialize endpoint works via POST to the /mcp endpoint.

**Steps:**
1. Start server with HTTP protocol
2. Send an initialize request:
   ```bash
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 9,
       "method": "initialize",
       "params": {
         "protocolVersion": "2025-06-18",
         "capabilities": {},
         "clientInfo": {
           "name": "test-client",
           "version": "1.0.0"
         }
       }
     }'
   ```
3. Observe the response
4. Stop the server

**Expected Result:**
- Request is accepted (HTTP 200)
- Response includes server capabilities
- Response includes server info (name, version)

## Cleanup

After testing:
1. Remove test config file: `rm ~/.config/ddg-search/config.yaml`
2. Remove test directory if empty: `rmdir ~/.config/ddg-search`
3. Verify no test processes are running: `ps aux | grep ddg-search-mcp`
4. Verify no ports are in use: `lsof -i :9100` (or other test ports)

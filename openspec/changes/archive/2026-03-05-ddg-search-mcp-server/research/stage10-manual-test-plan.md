# Stage 10 Manual Test Plan - Streamable HTTP Transport Migration

## Prerequisites
- Go 1.25.0 or later installed
- `ddg-search-mcp` binary built and available
- Test certificates for TLS/mTLS testing (located in `internal/mcp/testdata/`)
- curl or similar HTTP client for manual testing

## Test Cases

### Test Case 1: HTTP Server Startup with Streamable HTTP Transport
**Description:** Verify that the HTTP server starts correctly with the new Streamable HTTP transport and exposes the `/mcp` endpoint.

**Steps:**
1. Create a temporary config file with HTTP protocol:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:19130
   logging:
     level: debug
   ```
2. Start the `ddg-search-mcp` server with the config file
3. Wait for server startup (check logs for "HTTP Streamable HTTP server listening")
4. Verify the server is listening on the configured address

**Expected Result:**
- Server logs: "HTTP Streamable HTTP server listening" with address and endpoint
- Server is running and accepting connections
- Endpoint `/mcp` is available

### Test Case 2: POST Request to /mcp Endpoint
**Description:** Verify that POST requests to the `/mcp` endpoint are handled correctly.

**Steps:**
1. Start the server with HTTP protocol
2. Send a POST request to `/mcp` with an InitializeRequest:
   ```bash
   curl -X POST http://localhost:19130/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 1,
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
3. Verify the response is valid JSON-RPC

**Expected Result:**
- Server responds with a valid InitializeResult
- Response includes protocol version and server capabilities
- HTTP status code is 200 OK

### Test Case 3: GET Request to /mcp Endpoint (SSE)
**Description:** Verify that GET requests to the `/mcp` endpoint with `Accept: text/event-stream` establish SSE connections.

**Steps:**
1. Start the server with HTTP protocol
2. Send a GET request to `/mcp` with SSE accept header:
   ```bash
   curl -N http://localhost:19130/mcp \
     -H "Accept: text/event-stream" \
     -H "MCP-Protocol-Version: 2025-06-18"
   ```
3. Verify SSE stream is established

**Expected Result:**
- Server responds with SSE stream
- Initial events include server capabilities and available tools
- Connection stays open for receiving notifications

### Test Case 4: Health Check Endpoint
**Description:** Verify that the health check endpoint works correctly with the new transport.

**Steps:**
1. Start the server with HTTP protocol
2. Request the health check endpoint:
   ```bash
   curl http://localhost:19130/health
   ```
3. Verify the response

**Expected Result:**
- HTTP status code is 200 OK
- Response indicates the server is healthy

### Test Case 5: MCP-Protocol-Version Header Handling
**Description:** Verify that the server handles the MCP-Protocol-Version header correctly.

**Steps:**
1. Start the server with HTTP protocol
2. Send a POST request with `MCP-Protocol-Version: 2025-06-18` header
3. Send another POST request without the header
4. Compare the responses

**Expected Result:**
- Requests with the header use the specified protocol version
- Requests without the header assume protocol version 2025-03-26 for backwards compatibility
- Server includes the header in responses

### Test Case 6: Multiple Concurrent Connections
**Description:** Verify that the server handles multiple concurrent connections correctly.

**Steps:**
1. Start the server with HTTP protocol
2. Open multiple concurrent connections to `/mcp` endpoint (mix of POST and GET)
3. Send tool calls through each connection
4. Verify all connections receive correct responses

**Expected Result:**
- All connections are handled concurrently
- Each connection receives its own responses
- No connection interferes with another

### Test Case 7: TLS Support with Streamable HTTP Transport
**Description:** Verify that TLS works correctly with the new transport.

**Steps:**
1. Create a config file with TLS enabled:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:19131
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
   ```
2. Start the server
3. Send a POST request via HTTPS:
   ```bash
   curl -X POST https://localhost:19131/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     --insecure \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}'
   ```

**Expected Result:**
- Server logs: "HTTPS Streamable HTTP server listening"
- TLS connection is established successfully
- POST request is handled correctly over HTTPS

### Test Case 8: mTLS Support with Streamable HTTP Transport
**Description:** Verify that mTLS works correctly with the new transport.

**Steps:**
1. Create a config file with mTLS enabled:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:19132
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
       mtls:
         enabled: true
         ca_file: internal/mcp/testdata/ca-cert.pem
   ```
2. Start the server
3. Send a POST request via HTTPS with client certificate:
   ```bash
   curl -X POST https://localhost:19132/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     --cert internal/mcp/testdata/client-cert.pem \
     --key internal/mcp/testdata/client-key.pem \
     --insecure \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}'
   ```

**Expected Result:**
- Server logs: "HTTPS Streamable HTTP server listening" with mtls enabled
- mTLS connection is established successfully
- Client certificate is validated
- POST request is handled correctly over mTLS

### Test Case 9: Graceful Shutdown
**Description:** Verify that the server shuts down gracefully with the new transport.

**Steps:**
1. Start the server with HTTP protocol
2. Send SIGTERM to the server process
3. Verify the shutdown process

**Expected Result:**
- Server logs: "Shutting down HTTP server"
- Server stops accepting new connections
- In-flight requests are completed
- Server exits cleanly

### Test Case 10: Tool Calling via Streamable HTTP Transport
**Description:** Verify that tool calling works correctly through the new transport.

**Steps:**
1. Start the server with HTTP protocol
2. Initialize the MCP connection via POST to `/mcp`
3. Send a tool call request:
   ```bash
   curl -X POST http://localhost:19130/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{
       "jsonrpc": "2.0",
       "id": 2,
       "method": "tools/call",
       "params": {
         "name": "search",
         "arguments": {
           "query": "test query"
         }
       }
     }'
   ```
4. Verify the response

**Expected Result:**
- Tool call is processed correctly
- Response includes search results
- HTTP status code is 200 OK

## Cleanup
- Stop the server process if still running
- Remove any temporary config files created
- Clean up any temporary test data

## Notes

- The old `/sse` and `/message` endpoints should no longer be available
- The new `/mcp` endpoint handles both POST and GET methods
- The `MCP-Protocol-Version` header is required for the new protocol
- Backwards compatibility is maintained by assuming version 2025-03-26 if header is missing

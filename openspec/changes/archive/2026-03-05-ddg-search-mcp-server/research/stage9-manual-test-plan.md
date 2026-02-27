# Stage 9 Manual Test Plan: TLS and mTLS Support

## Prerequisites

- Go 1.21+ installed
- Project cloned and dependencies installed (`go mod download`)
- Terminal access for running commands
- Text editor for creating config files
- `openssl` or similar tool for generating test certificates
- `curl` or similar HTTP client for testing endpoints
- Test certificates available in `internal/mcp/testdata/`

## Test Cases

### Test Case 1: TLS Disabled (Default)

**Description:** Verify that application starts successfully with TLS disabled by default.

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
3. Observe output
4. In another terminal, run: `curl http://localhost:9100/health`
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts without errors
- Output shows "TLS disabled" or similar message
- Output shows "HTTP SSE server listening" (not HTTPS)
- Health check returns 200 OK over HTTP
- Application stops cleanly on Ctrl-C

### Test Case 2: TLS Enabled with Valid Certificates

**Description:** Verify that application starts successfully with TLS enabled and valid certificates.

**Steps:**
1. Create config file with TLS enabled:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
       min_version: "1.2"
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe output
4. In another terminal, run: `curl -k https://localhost:9100/health`
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts without errors
- Output shows "TLS enabled"
- Output shows "HTTPS server listening"
- Health check returns 200 OK over HTTPS
- Application stops cleanly on Ctrl-C

### Test Case 3: TLS Enabled with Missing Certificate File

**Description:** Verify that application fails to start when TLS certificate file is missing.

**Steps:**
1. Create config file with TLS enabled but invalid cert path:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: /nonexistent/cert.pem
       key_file: internal/mcp/testdata/server-key.pem
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe output
4. Wait for process to exit

**Expected Result:**
- Application fails to start
- Error message indicates missing certificate file
- Process exits with non-zero status

### Test Case 4: TLS Enabled with Missing Key File

**Description:** Verify that application fails to start when TLS key file is missing.

**Steps:**
1. Create config file with TLS enabled but invalid key path:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: /nonexistent/key.pem
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe output
4. Wait for process to exit

**Expected Result:**
- Application fails to start
- Error message indicates missing key file
- Process exits with non-zero status

### Test Case 5: TLS Enabled with Invalid Certificate

**Description:** Verify that application fails to start when TLS certificate is invalid.

**Steps:**
1. Create config file with TLS enabled and invalid cert:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/ca-cert.pem  # Using CA cert as server cert
       key_file: internal/mcp/testdata/server-key.pem
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe output
4. Wait for process to exit

**Expected Result:**
- Application fails to start
- Error message indicates invalid certificate or key pair
- Process exits with non-zero status

### Test Case 6: mTLS Enabled with Valid Client Certificate

**Description:** Verify that application accepts connections with valid client certificates when mTLS is enabled.

**Steps:**
1. Create config file with mTLS enabled:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
       mtls:
         enabled: true
         ca_file: internal/mcp/testdata/ca-cert.pem
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe output
4. In another terminal, run with client certificate:
   ```bash
   curl -k \
     --cert internal/mcp/testdata/client-cert.pem \
     --key internal/mcp/testdata/client-key.pem \
     https://localhost:9100/health
   ```
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts without errors
- Output shows "mTLS enabled"
- Health check returns 200 OK over HTTPS with client certificate
- Application stops cleanly on Ctrl-C

### Test Case 7: mTLS Enabled Without Client Certificate

**Description:** Verify that application rejects connections without client certificates when mTLS is enabled.

**Steps:**
1. Start server with mTLS enabled (from Test Case 6)
2. In another terminal, run without client certificate:
   ```bash
   curl -k https://localhost:9100/health
   ```
3. Observe response

**Expected Result:**
- Connection fails or returns error
- Server logs rejection of connection without certificate
- Client receives TLS handshake error or HTTP error

### Test Case 8: mTLS Enabled with Invalid Client Certificate

**Description:** Verify that application rejects connections with invalid client certificates when mTLS is enabled.

**Steps:**
1. Start server with mTLS enabled (from Test Case 6)
2. In another terminal, run with server certificate as client certificate (invalid):
   ```bash
   curl -k \
     --cert internal/mcp/testdata/server-cert.pem \
     --key internal/mcp/testdata/server-key.pem \
     https://localhost:9100/health
   ```
3. Observe response

**Expected Result:**
- Connection fails or returns error
- Server logs rejection of connection with invalid certificate
- Client receives TLS handshake error or HTTP error

### Test Case 9: TLS Minimum Version Configuration

**Description:** Verify that TLS minimum version is enforced correctly.

**Steps:**
1. Create config file with TLS min version 1.2:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
       min_version: "1.2"
   logging:
     level: info
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Observe output for min_version configuration
4. Test health check: `curl -k https://localhost:9100/health`
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts without errors
- Output shows configured min_version
- Health check returns 200 OK
- Application stops cleanly on Ctrl-C

### Test Case 10: TLS Certificate Reload on SIGHUP

**Description:** Verify that TLS certificates are reloaded when SIGHUP signal is received.

**Steps:**
1. Start server with TLS enabled (from Test Case 2)
2. Observe output showing server started
3. In another terminal, find the process ID: `ps aux | grep ddg-search-mcp`
4. Send SIGHUP signal: `kill -HUP <pid>`
5. Observe server output for reload message
6. Test health check: `curl -k https://localhost:9100/health`
7. Press Ctrl-C to stop

**Expected Result:**
- Server logs receipt of SIGHUP signal
- Server logs TLS certificate reload
- Health check still works after reload
- Application stops cleanly on Ctrl-C

### Test Case 11: TLS with MCP Endpoint (GET for SSE)

**Description:** Verify that TLS works with the /mcp endpoint via GET with SSE headers.

**Steps:**
1. Start server with TLS enabled (from Test Case 2)
2. In another terminal, run:
   ```bash
   curl -k https://localhost:9100/mcp -H "Accept: text/event-stream" -H "MCP-Protocol-Version: 2025-06-18"
   ```
3. Observe response headers and status
4. Press Ctrl-C to stop server

**Expected Result:**
- MCP endpoint returns 200 OK over HTTPS
- Response includes SSE headers (Content-Type: text/event-stream)
- Response includes MCP-Protocol-Version: 2025-06-18
- Connection is established over TLS

### Test Case 12: mTLS with MCP Endpoint (GET for SSE)

**Description:** Verify that mTLS works with the /mcp endpoint via GET with SSE headers.

**Steps:**
1. Start server with mTLS enabled (from Test Case 6)
2. In another terminal, run with client certificate:
   ```bash
   curl -k \
     --cert internal/mcp/testdata/client-cert.pem \
     --key internal/mcp/testdata/client-key.pem \
     https://localhost:9100/mcp -H "Accept: text/event-stream" -H "MCP-Protocol-Version: 2025-06-18"
   ```
3. Observe response headers and status
4. Press Ctrl-C to stop server

**Expected Result:**
- MCP endpoint returns 200 OK over HTTPS with client certificate
- Response includes SSE headers (Content-Type: text/event-stream)
- Response includes MCP-Protocol-Version: 2025-06-18
- Connection is established over TLS with mTLS

### Test Case 13: TLS Connection Logging at Debug Level

**Description:** Verify that TLS connection details are logged at debug level.

**Steps:**
1. Create config file with debug logging:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
   logging:
     level: debug
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Make a request: `curl -k https://localhost:9100/health`
4. Observe server output for TLS connection details
5. Press Ctrl-C to stop

**Expected Result:**
- Server logs TLS connection at debug level
- Log includes TLS version (e.g., TLS 1.2 or TLS 1.3)
- Log includes cipher suite information
- Application stops cleanly on Ctrl-C

### Test Case 14: mTLS Client Certificate Logging at Debug Level

**Description:** Verify that client certificate details are logged at debug level when mTLS is enabled.

**Steps:**
1. Create config file with mTLS and debug logging:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
     tls:
       enabled: true
       cert_file: internal/mcp/testdata/server-cert.pem
       key_file: internal/mcp/testdata/server-key.pem
       mtls:
         enabled: true
         ca_file: internal/mcp/testdata/ca-cert.pem
   logging:
     level: debug
   ```
2. Run: `go run ./cmd/ddg-search-mcp`
3. Make a request with client certificate (from Test Case 6)
4. Observe server output for client certificate details
5. Press Ctrl-C to stop

**Expected Result:**
- Server logs mTLS connection at debug level
- Log includes client certificate subject or issuer
- Log indicates successful client certificate validation
- Application stops cleanly on Ctrl-C

### Test Case 15: Environment Variable Override for TLS Configuration

**Description:** Verify that TLS configuration can be overridden via environment variables.

**Steps:**
1. Create minimal config file:
   ```yaml
   server:
     protocol: http
     bind_address: localhost:9100
   logging:
     level: info
   ```
2. Run with environment variables:
   ```bash
   DDG_SEARCH_SERVER_TLS_ENABLED=true \
   DDG_SEARCH_SERVER_TLS_CERT_FILE=internal/mcp/testdata/server-cert.pem \
   DDG_SEARCH_SERVER_TLS_KEY_FILE=internal/mcp/testdata/server-key.pem \
   go run ./cmd/ddg-search-mcp
   ```
3. Observe output
4. Test health check: `curl -k https://localhost:9100/health`
5. Press Ctrl-C to stop

**Expected Result:**
- Application starts with TLS enabled via environment variable
- Output shows TLS enabled
- Health check returns 200 OK over HTTPS
- Application stops cleanly on Ctrl-C

## Cleanup

- Stop any running server instances
- Remove temporary config files
- No persistent changes to system configuration

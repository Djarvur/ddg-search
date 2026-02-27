# Streamable HTTP Transport Migration

## Overview

This document outlines the migration from the deprecated HTTP+SSE transport to the new Streamable HTTP transport as per MCP specification version 2025-06-18.

## Current State (Deprecated)

### Protocol Version
- MCP specification version: 2024-11-05 (deprecated)
- Library: `github.com/mark3labs/mcp-go` v0.44.1

### Transport Implementation
- Uses `server.SSEServer` with separate endpoints:
  - `/sse` - SSE endpoint for server-to-client streaming
  - `/message` - POST endpoint for client-to-server requests

### Key Files
- `internal/mcp/server.go` - Main server implementation (lines 310-369)
- `internal/mcp/http_sse_test.go` - HTTP SSE tests
- `internal/mcp/http_sse_internal_test.go` - Internal tests

## Target State (New Standard)

### Protocol Version
- MCP specification version: 2025-06-18 (current)
- Library: `github.com/mark3labs/mcp-go` (same version supports both transports)

### Transport Implementation
- Uses `server.NewStreamableHTTPServer()` with single endpoint:
  - `/mcp` - Single endpoint supporting both POST and GET methods
    - POST: For sending JSON-RPC requests to the server
    - GET (with `Accept: text/event-stream`): For opening SSE stream to receive server notifications

### Key Changes

#### 1. Endpoint Consolidation
**Before:**
```go
mux.Handle("/sse", s.sseServer.SSEHandler())
mux.Handle("/message", s.sseServer.MessageHandler())
```

**After:**
```go
httpServer := server.NewStreamableHTTPServer(s.mcpServer)
// The /mcp endpoint is automatically registered
```

#### 2. Protocol Version Header
The new transport requires handling the `MCP-Protocol-Version` header:
- Clients MUST include `MCP-Protocol-Version: 2025-06-18` in all requests
- Servers MUST include this header in all responses
- If header is missing, assume version `2025-03-26` for backwards compatibility

#### 3. Session Management
The new transport supports session management via `Mcp-Session-Id` header:
- Server MAY assign a session ID during initialization
- Client MUST include session ID in subsequent requests if provided
- Server MAY terminate sessions (respond with 404)

## Implementation Plan

### Phase 1: Code Changes

#### 1. Update `internal/mcp/server.go`
- Replace `server.SSEServer` with `server.NewStreamableHTTPServer()`
- Remove `/sse` and `/message` endpoint handlers
- Update logging to reflect `/mcp endpoint
- Ensure `MCP-Protocol-Version` header handling

#### 2. Update HTTP Server Setup
```go
// Current implementation (lines 310-343)
func (s *Server) serveHTTP(ctx context.Context, bindAddress string) error {
    // Create SSE server
    s.mu.Lock()
    s.sseServer = server.NewSSEServer(s.mcpServer)

    // Create HTTP server with logging middleware
    mux := http.NewServeMux()

    // Health check endpoint
    mux.HandleFunc("/health", s.healthCheckHandler)

    // SSE endpoints - use individual handlers for proper routing
    mux.Handle("/sse", s.sseServer.SSEHandler())
    mux.Handle("/message", s.sseServer.MessageHandler())

    // Wrap with logging middleware
    handler := s.loggingMiddleware(mux)

    s.httpServer = &http.Server{
        Addr:              bindAddress,
        Handler:           handler,
        ReadHeaderTimeout: readHeaderTimeout,
        TLSConfig:         tlsConfig,
    }
    s.mu.Unlock()
    // ...
}

// New implementation
func (s *Server) serveHTTP(ctx context.Context, bindAddress string) error {
    // Get TLS configuration from appConfig
    tlsConfig, err := s.buildTLSConfig()
    if err != nil {
        return fmt.Errorf("failed to build TLS config: %w", err)
    }

    // Create Streamable HTTP server
    s.mu.Lock()
    streamableHTTPServer := server.NewStreamableHTTPServer(s.mcpServer)

    // Create HTTP server with logging middleware
    mux := http.NewServeMux()

    // Health check endpoint
    mux.HandleFunc("/health", s.healthCheckHandler)

    // MCP endpoint (automatically handles POST and GET)
    mux.Handle("/mcp", streamableHTTPServer)

    // Wrap with logging middleware
    handler := s.loggingMiddleware(mux)

    s.httpServer = &http.Server{
        Addr:              bindAddress,
        Handler:           handler,
        ReadHeaderTimeout: readHeaderTimeout,
        TLSConfig:         tlsConfig,
    }
    s.mu.Unlock()
    // ...
}
```

#### 3. Update Server Struct
```go
// Current struct (lines 48-58)
type Server struct {
    mcpServer   *server.MCPServer
    sseServer   *server.SSEServer  // Remove this
    logger      *slog.Logger
    config      *Config
    appConfig   any
    httpServer  *http.Server
    shutdownCtx context.Context
    mu          sync.RWMutex
}

// New struct
type Server struct {
    mcpServer   *server.MCPServer
    streamableHTTPServer *server.StreamableHTTPServer  // Add this
    logger      *slog.Logger
    config      *Config
    appConfig   any
    httpServer  *http.Server
    shutdownCtx context.Context
    mu          sync.RWMutex
}
```

#### 4. Update Shutdown Method
```go
// Current shutdown (lines 231-237)
if s.sseServer != nil {
    err := s.sseServer.Shutdown(ctx)
    if err != nil {
        return fmt.Errorf("SSE server shutdown error: %w", err)
    }
}

// New shutdown
if s.streamableHTTPServer != nil {
    err := s.streamableHTTPServer.Shutdown(ctx)
    if err != nil {
        return fmt.Errorf("Streamable HTTP server shutdown error: %w", err)
    }
}
```

### Phase 2: Test Updates

#### 1. Update `internal/mcp/http_sse_test.go`
- Rename to `streamable_http_test.go`
- Update test cases to use `/mcp` endpoint
- Update test cases to use POST for requests
- Update test cases to use GET with `Accept: text/event-stream` for SSE

#### 2. Update `internal/mcp/http_sse_internal_test.go`
- Rename to `streamable_http_internal_test.go`
- Update internal tests for new transport

#### 3. Update `cmd/ddg-search-mcp/e2e_test.go`
- Update e2e tests to use new endpoint

### Phase 3: Documentation Updates

#### 1. Update `README.md`
- Update endpoint from `/sse` to `/mcp`
- Update transport description to "Streamable HTTP"
- Update protocol version to 2025-06-18

#### 2. Update Configuration Documentation
- Update any references to `/sse` or `/message` endpoints

### Phase 4: Backwards Compatibility (Optional)

If backwards compatibility with old clients is required:
- Keep old endpoints (`/sse` and `/message`) alongside `/mcp`
- Old clients will continue to work
- New clients will use `/mcp`

## Testing Strategy

### Unit Tests
1. Test POST requests to `/mcp` endpoint
2. Test GET requests to `/mcp` endpoint with SSE
3. Test `MCP-Protocol-Version` header handling
4. Test concurrent connections
5. Test TLS/mTLS with new transport

### Integration Tests
1. Test tool calling via POST
2. Test server notifications via SSE
3. Test session management
4. Test error handling

### Manual Testing
1. Test with Claude Desktop (new MCP client)
2. Test with curl:
   ```bash
   # POST request
   curl -X POST http://localhost:9100/mcp \
     -H "Content-Type: application/json" \
     -H "MCP-Protocol-Version: 2025-06-18" \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}'

   # GET request for SSE
   curl -N http://localhost:9100/mcp \
     -H "Accept: text/event-stream" \
     -H "MCP-Protocol-Version: 2025-06-18"
   ```

## Migration Checklist

- [ ] Update `internal/mcp/server.go` to use `NewStreamableHTTPServer()`
- [ ] Update Server struct to replace `sseServer` with `streamableHTTPServer`
- [ ] Update shutdown method
- [ ] Update logging messages
- [ ] Rename `http_sse_test.go` to `streamable_http_test.go`
- [ ] Rename `http_sse_internal_test.go` to `streamable_http_internal_test.go`
- [ ] Update all test cases
- [ ] Update `README.md`
- [ ] Update configuration documentation
- [ ] Run all tests
- [ ] Manual testing with Claude Desktop
- [ ] Update this research document with any findings

## References

- [MCP Specification - Transports (2025-06-18)](https://modelcontextprotocol.io/specification/2025-06-18/basic/transports)
- [MCP-Go Streamable HTTP Documentation](https://mcp-go.dev/transports/http/)
- [MCP-Go GitHub Repository](https://github.com/mark3labs/mcp-go)

# MCP Go Packages Comparison

This document compares Model Context Protocol (MCP) Go packages for building a Cloud Code compatible MCP server based on the existing skills in this repository.

## Overview

The repository contains two skills:
1. **ddg-search** - Web search using DuckDuckGo (no API keys required)
2. **perplexity-search** - Web search using Perplexity API (AI-powered results with citations)

Both skills expose tools via command-line binaries that need to be integrated into an MCP server.

## Comparison Criteria

| Criteria | Description |
|----------|-------------|
| **HTTP Transport** | Support for HTTP/Streamable HTTP transport |
| **Stdio Transport** | Support for stdin/stdout transport |
| **Official Status** | Whether the package is officially maintained |
| **MCP Features** | Support for tools, resources, prompts, notifications |
| **Code Quality** | Benchmark score and code snippet coverage |
| **Ease of Use** | Simplicity of API and documentation |

## Package Analysis

### 1. [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go)

**Status**: Unofficial but well-maintained
**Benchmark Score**: 81.8
**Code Snippets**: 580

#### Features
- **HTTP Transport**: Streamable HTTP transport supported
- **Stdio Transport**: Full stdio support via `server.ServeStdio()`
- **SSE Transport**: Server-Sent Events support
- **MCP Features**: Tools, resources, prompts, notifications
- **Protocol Version**: Supports latest MCP protocol version

#### Pros
- High code snippet coverage (580 snippets)
- Well-documented with comprehensive examples
- Multiple transport options (stdio, SSE, Streamable HTTP)
- Active maintenance by Mark3Labs
- Good balance of features and simplicity

#### Cons
- Unofficial implementation
- Slightly lower benchmark score than some alternatives

#### Example Code
```go
// Stdio Server
s := server.NewMCPServer("My Server", "1.0.0")
server.ServeStdio(s)

// HTTP Server
httpServer := server.NewStreamableHTTPServer(s)
httpServer.Start(":8080")
```

---

### 2. [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk)

**Status**: Official SDK
**Benchmark Score**: 85.7
**Code Snippets**: 366

#### Features
- **HTTP Transport**: Streamable HTTP via `NewStreamableHTTPHandler()`
- **Stdio Transport**: Full stdio support via `StdioTransport`
- **MCP Features**: Tools, resources, prompts, notifications
- **Protocol Version**: Official reference implementation

#### Pros
- **Official SDK** from Model Context Protocol team
- Highest benchmark score (85.7)
- Reference implementation for MCP spec
- Type-safe Go implementation

#### Cons
- Lower code snippet coverage (366 snippets)
- More low-level API (requires more boilerplate)
- Less beginner-friendly examples

#### Example Code
```go
// Stdio Server
server := mcp.NewServer(&mcp.Implementation{Name:"server", Version:"v1.0.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "tool"}, handler)
server.Run(context.Background(), &mcp.StdioTransport{})

// HTTP Server
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
    return server
})
http.ListenAndServe(":8080", handler)
```

---

### 3. [`metoro-io/mcp-golang`](https://github.com/metoro-io/mcp-golang)

**Status**: Unofficial
**Code Snippets**: 79

#### Features
- **HTTP Transport**: Standard HTTP and Gin framework integration
- **Stdio Transport**: Full stdio support
- **MCP Features**: Tools, resources, prompts
- **Type Safety**: Struct-based tool arguments with JSON schema

#### Pros
- Clean, type-safe API
- Struct-based tool definitions with JSON schema tags
- Good for Go developers familiar with struct annotations
- Supports all MCP core features

#### Cons
- Unofficial implementation
- Lower code snippet coverage (79 snippets)
- HTTP transport is stateless (no bidirectional features)
- Less comprehensive documentation

#### Example Code
```go
// Stdio Server
server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
server.RegisterTool("hello", "Say hello", func(args MyArgs) (*mcp_golang.ToolResponse, error) {
    return mcp_golang.NewToolResponse(...), nil
})
server.Serve()
```

---

### 4. [`trpc-group/trpc-mcp-go`](https://github.com/trpc-group/trpc-mcp-go)

**Status**: Unofficial
**Benchmark Score**: 80.3
**Code Snippets**: 73

#### Features
- **HTTP Transport**: Streamable HTTP with graceful shutdown
- **Stdio Transport**: Full stdio support with cross-language compatibility
- **MCP Features**: Tools, resources, prompts
- **Cross-Language**: Examples for TypeScript, Python, Go interoperability

#### Pros
- Comprehensive streaming HTTP support
- Cross-language stdio examples (TS, Python, Go)
- Graceful shutdown handling
- Good for complex server implementations

#### Cons
- Unofficial implementation
- Lower code snippet coverage (73 snippets)
- More complex API than mark3labs/mcp-go
- Less documentation

#### Example Code
```go
// Stdio Server
mcpServer := mcp.NewStdioServer("My-STDIO-Server", "1.0.0")
mcpServer.RegisterTool(greetTool, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    return mcp.NewTextResult("Hello!"), nil
})
mcpServer.Start()

// HTTP Server
mcpServer := mcp.NewServer("Server", "1.0.0",
    mcp.WithServerAddress(":3000"),
    mcp.WithServerPath("/mcp"))
mcpServer.Start()
```

---

## Recommendation for This Project

### **Best Choice: `mark3labs/mcp-go`**

**Reasoning:**

1. **Best Documentation**: 580 code snippets provide comprehensive examples for common use cases
2. **Full Transport Support**: Both HTTP and stdio are first-class citizens
3. **Mature Implementation**: Well-tested with good community adoption
4. **Cloud Code Compatible**: Supports all MCP features needed for Cloud Code integration
5. **Faster Development**: Higher code snippet coverage means faster implementation

### **Alternative: `modelcontextprotocol/go-sdk`**

Choose this if:
- You prefer the official SDK from the spec authors
- You need the most spec-compliant implementation
- You're comfortable with lower-level APIs

---

## Implementation Considerations

### Tool Mapping

Your existing skills map to MCP tools as follows:

| Skill | Tool Name | Parameters | Command |
|-------|-----------|------------|---------|
| ddg-search | `web-search` | `query` (required) | `ddg-search "{{query}}"` |
| ddg-search | `page-dump` | `url` (required) | `page-dump "{{url}}"` |
| perplexity-search | `perplexity-search` | `query` (required), `max-results` (optional), `model` (optional) | `perplexity-search --max-results "{{max-results}}" --model "{{model}}" "{{query}}"` |

### Transport Strategy

For Cloud Code compatibility:
1. **Stdio Transport**: Primary transport for local Cloud Code integration
2. **HTTP Transport**: Secondary transport for remote server deployment

Both `mark3labs/mcp-go` and `modelcontextprotocol/go-sdk` support both transports equally well.

---

## Summary Table

| Package | Official | HTTP | Stdio | Snippets | Score | Recommendation |
|---------|----------|------|-------|----------|-------|----------------|
| mark3labs/mcp-go | ❌ | ✅ | ✅ | 580 | 81.8 | ⭐ **Best Choice** |
| modelcontextprotocol/go-sdk | ✅ | ✅ | ✅ | 366 | 85.7 | ⭐ Alternative |
| metoro-io/mcp-golang | ❌ | ✅ | ✅ | 79 | - | Not recommended |
| trpc-group/trpc-mcp-go | ❌ | ✅ | ✅ | 73 | 80.3 | Not recommended |

---

## Next Steps

1. **Choose `mark3labs/mcp-go`** for faster implementation with better documentation
2. Create MCP server that exposes the three tools from your skills
3. Implement both stdio and HTTP transports
4. Test with Cloud Code client
5. Package as a binary for easy installation

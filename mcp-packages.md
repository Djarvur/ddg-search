# MCP Go Packages Comparison

This document provides a comprehensive comparison of Go packages available for building Model Context Protocol (MCP) servers. The goal is to help choose the right package for implementing an MCP server for web search based on the ddg-search codebase.

## Executive Summary

We analyzed 4 MCP Go packages to build a cloud code compatible MCP server supporting both HTTP and stdio transports:

| Package | Stars | HTTP Support | stdio Support | Type Safety | Ease of Use | Recommendation |
|---------|-------|--------------|---------------|-------------|-------------|----------------|
| **mark3labs/mcp-go** | ~3.5k | ✅ StreamableHTTP, SSE | ✅ Native | High | Easy | **Recommended** |
| **modelcontextprotocol/go-sdk** | ~2.5k | ✅ StreamableHTTP | ✅ Native | High | Medium | Good Alternative |
| **metoro-io/mcp-golang** | ~600 | ✅ HTTP, Gin | ✅ Native | High | Easy | Good Alternative |
| **trpc-group/trpc-mcp-go** | ~200 | ✅ StreamableHTTP, SSE | ✅ Native | High | Medium | For tRPC users |

## Package Details

### 1. mark3labs/mcp-go

**Repository**: https://github.com/mark3labs/mcp-go

**Description**: A Go implementation of the Model Context Protocol for seamless integration between LLM applications and external data sources and tools. Offers a fast, simple, and complete solution.

**Pros**:
- Most popular community implementation (~3.5k stars)
- Excellent documentation with many examples
- Supports multiple transports: stdio, StreamableHTTP, SSE
- In-process communication support via `client.NewInProcessClient(server)`
- Active maintenance and large community
- High benchmark score (81.8)

**Cons**:
- Not an "official" MCP SDK
- Some features may lag behind official spec updates

**Example - stdio server**:
```go
s := server.NewMCPServer("My Server", "1.0.0")
server.ServeStdio(s)
```

**Example - HTTP server**:
```go
s := server.NewMCPServer("My Server", "1.0.0")
httpServer := server.NewStreamableHTTPServer(s)
httpServer.Start(":8080")
```

---

### 2. modelcontextprotocol/go-sdk (Official)

**Repository**: https://github.com/modelcontextprotocol/go-sdk

**Description**: The official Go SDK for the Model Context Protocol, providing APIs for constructing and using MCP clients and servers.

**Pros**:
- Official MCP SDK from the protocol authors
- Supports JSON Schema and JSON-RPC implementations
- StreamableHTTP and stdio transports
- High benchmark score (85.7)
- 6026 code snippets available

**Cons**:
- May have more boilerplate code
- Steeper learning curve compared to community packages
- API design is more low-level

**Example - stdio server**:
```go
server := mcp.NewServer(
    mcp.WithTool(mcp.NewTool("simple-tool", "A simple tool", nil, nil)),
)
server.Run(context.Background(), &mcp.StdioTransport{})
```

**Example - HTTP server**:
```go
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
    return server
}, nil)
http.ListenAndServe(":8080", handler)
```

---

### 3. metoro-io/mcp-golang

**Repository**: https://github.com/metoro-io/mcp-golang

**Description**: An unofficial implementation of the Model Context Protocol in Go, enabling developers to easily write MCP servers and clients with type safety and custom transports.

**Pros**:
- Clean, intuitive API design
- Struct-based argument definitions with JSON schema tags
- Built-in support for Gin framework integration
- Easy to get started

**Cons**:
- Smaller community (~600 stars)
- HTTP transport is stateless (no bidirectional features/notifications)
- Less documentation compared to mark3labs

**Example - stdio server**:
```go
server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
server.RegisterTool("hello", "Say hello", func(args MyArguments) (*ToolResponse, error) {
    return NewToolResponse(NewTextContent("Hello!")), nil
})
server.Serve()
```

**Example - HTTP server**:
```go
transport := http.NewHTTPTransport("/mcp")
transport.WithAddr(":8080")
server := mcp_golang.NewServer(transport)
server.Serve()
```

---

### 4. trpc-group/trpc-mcp-go

**Repository**: https://github.com/trpc-group/trpc-mcp-go

**Description**: A Go implementation of MCP with comprehensive streaming HTTP support, featuring various transport options like Streamable HTTP, STDIO, and SSE.

**Pros**:
- Part of tRPC ecosystem
- Excellent for users already using tRPC
- Good streaming support
- SSE support

**Cons**:
- Less popular (~200 stars)
- Tighter coupling to tRPC framework
- Smaller community

**Example - stdio server**:
```go
server := mcp.NewStdioServer("My-STDIO-Server", "1.0.0")
server.RegisterTool(tool, handler)
server.Start()
```

**Example - HTTP server**:
```go
server := mcp.NewServer("My-Server", "1.0.0",
    mcp.WithServerAddress(":3000"),
    mcp.WithServerPath("/mcp"),
)
server.Start()
```

---

## Comparison Matrix

| Feature | mark3labs/mcp-go | modelcontextprotocol/go-sdk | metoro-io/mcp-golang | trpc-mcp-go |
|---------|------------------|----------------------------|---------------------|-------------|
| **Stars** | ~3.5k | ~2.5k | ~600 | ~200 |
| **Official** | ❌ | ✅ | ❌ | ❌ |
| **stdio** | ✅ | ✅ | ✅ | ✅ |
| **HTTP** | ✅ | ✅ | ✅ | ✅ |
| **SSE** | ✅ | ✅ | ❌ | ✅ |
| **StreamableHTTP** | ✅ | ✅ | ❌ | ✅ |
| **Gin integration** | ❌ | ❌ | ✅ | ❌ |
| **In-process** | ✅ | ❌ | ❌ | ❌ |
| **Type Safety** | High | High | High | High |
| **Documentation** | Excellent | Good | Good | Good |
| **Active Issues** | Low | Low | Medium | Low |

---

## Transport Analysis

### stdio Transport
- All packages support stdio natively
- Best for: Local tools, CLI integrations, development
- Your existing ddg-search CLI already works this way

### HTTP Transport
- **mark3labs/mcp-go**: StreamableHTTP, SSE - Best overall HTTP support
- **modelcontextprotocol/go-sdk**: StreamableHTTP - Official standard
- **metoro-io/mcp-golang**: Basic HTTP, Gin - Stateless (no bidirectional)
- **trpc-mcp-go**: StreamableHTTP, SSE - Good streaming support

---

## Recommendation for ddg-search MCP Server

Based on the analysis, I recommend **mark3labs/mcp-go** for the following reasons:

1. **Best HTTP Support**: StreamableHTTP + SSE provides modern, efficient communication
2. **Most Popular**: Large community and active maintenance
3. **Excellent Documentation**: Easy to get started
4. **In-Process Support**: Can integrate directly with your existing search code
5. **Production Ready**: Used by many production MCP servers
6. **Both Transports**: Native support for both stdio and HTTP

### Runner-up: modelcontextprotocol/go-sdk
- Choose this if you prefer official SDK and don't mind slightly more boilerplate
- Better for strict protocol compliance requirements

---

## Implementation Quick Start (mark3labs/mcp-go)

Based on your existing `internal/search` package:

```go
package main

import (
    "context"
    "log"
    "github.com/mark3labs/mcp-go/server"
    "github.com/your-repo/internal/search"
)

func main() {
    // Initialize your search client
    searcher := search.NewSearcher(search.DefaultClient(), nil)
    defer searcher.Close()

    // Create MCP server
    s := server.NewMCPServer("ddg-search", "1.0.0")

    // Register search tool
    s.AddTool(&mcp.Tool{
        Name:        "ddg-search",
        Description: "Search the web using DuckDuckGo",
        InputSchema: mcp.JSONSchema{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]string{"type": "string"},
                "count": map[string]interface{}{"type": "integer", "default": 10},
            },
            "required": []string{"query"},
        },
    }, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        query := req.Params.Arguments["query"].(string)
        count := 10
        if c, ok := req.Params.Arguments["count"].(float64); ok {
            count = int(c)
        }

        results, err := searcher.Search(context.Background(), query, count)
        if err != nil {
            return nil, err
        }

        // Format results...
    })

    // Start with both transports (separate processes or using SSE)
    server.ServeStdio(s)
}
```

---

## Next Steps

1. **Choose the package** - Based on the analysis above
2. **Create proposal** - Document the implementation plan
3. **Design phase** - Define MCP tools, resources, prompts
4. **Implementation** - Build the MCP server
5. **Testing** - Test both stdio and HTTP transports
6. **Documentation** - Add to your skills system

---

## References

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- [metoro-io/mcp-golang](https://github.com/metoro-io/mcp-golang)
- [trpc-group/trpc-mcp-go](https://github.com/trpc-group/trpc-mcp-go)
- [MCP Specification](https://spec.modelcontextprotocol.io/)

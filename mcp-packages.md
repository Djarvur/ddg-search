# Go MCP Server Packages Comparison

This document compares Go packages for building Model Context Protocol (MCP) servers, specifically for creating a Claude-compatible MCP server based on the existing search skills in this repository.

## Existing Skills to Wrap

The repository contains two search skills that need to be exposed as MCP tools:

| Skill | Tools | Description | API Key Required |
|-------|-------|-------------|------------------|
| `ddg-search` | `web-search`, `page-dump` | DuckDuckGo web search and page fetching | No |
| `perplexity-search` | `perplexity-search` | AI-powered search with citations | Yes |

## Package Comparison

### 1. modelcontextprotocol/go-sdk (Official)

**Repository**: https://github.com/modelcontextprotocol/go-sdk

| Aspect | Details |
|--------|---------|
| **Source Reputation** | High (Official) |
| **Benchmark Score** | 85.7 |
| **Code Snippets** | 366 |
| **Latest Version** | v1.2.0 |
| **Maturity** | Official SDK, actively maintained |

#### Transport Support

| Transport | Server | Client | Notes |
|-----------|--------|--------|-------|
| Stdio | ✅ `StdioTransport` | ✅ `CommandTransport` | Full bidirectional support |
| HTTP | ✅ `StreamableHTTPHandler` | ✅ `StreamableClientTransport` | Streamable HTTP |
| SSE | ❌ | ❌ | Not supported |

#### API Style

```go
// Server creation
server := mcp.NewServer(&mcp.Implementation{Name: "server", Version: "1.0.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "tool"}, handler)

// Stdio transport
server.Run(context.Background(), &mcp.StdioTransport{})

// HTTP transport
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
    return server
}, nil)
http.ListenAndServe(":8080", handler)
```

#### Pros
- **Official SDK** - Maintained by MCP team
- **Comprehensive** - Full protocol support
- **Well-documented** - Extensive examples
- **OAuth support** - Built-in authentication
- **JSON Schema** - Full schema support

#### Cons
- **More verbose** - Requires more boilerplate
- **Explicit API** - Less ergonomic than some alternatives

---

### 2. mark3labs/mcp-go

**Repository**: https://github.com/mark3labs/mcp-go

| Aspect | Details |
|--------|---------|
| **Source Reputation** | High |
| **Benchmark Score** | 81.8 |
| **Code Snippets** | 580 |
| **Maturity** | Active, good documentation |

#### Transport Support

| Transport | Server | Client | Notes |
|-----------|--------|--------|-------|
| Stdio | ✅ `ServeStdio(s)` | ✅ | Simple API |
| HTTP | ✅ `ServeHTTP(s, ":8080")` | ✅ | REST-like |
| SSE | ✅ `ServeSSE(s, ":8080")` | ✅ | Server-Sent Events |

#### API Style

```go
// Server creation
s := server.NewMCPServer("My Server", "1.0.0")

// Choose transport at runtime
switch transport {
case "stdio":
    server.ServeStdio(s)
case "http":
    server.ServeHTTP(s, ":8080")
case "sse":
    server.ServeSSE(s, ":8080")
}
```

#### Pros
- **Very simple API** - Minimal boilerplate
- **Transport-agnostic** - Same server code works with any transport
- **Multiple transports** - Stdio, HTTP, SSE all supported
- **Fast** - Performance-focused
- **Complete solution** - Everything needed in one package

#### Cons
- **Third-party** - Not official
- **Less mature** - Newer than official SDK

---

### 3. metoro-io/mcp-golang

**Repository**: https://github.com/metoro-io/mcp-golang

| Aspect | Details |
|--------|---------|
| **Source Reputation** | High |
| **Code Snippets** | 79 |
| **Maturity** | Active, type-safe |

#### Transport Support

| Transport | Server | Client | Notes |
|-----------|--------|--------|-------|
| Stdio | ✅ `stdio.NewStdioServerTransport()` | ✅ | Full bidirectional support |
| HTTP | ✅ `http.NewHTTPTransport("/mcp")` | ✅ | Stateless, no bidirectional |
| Gin | ✅ `http.NewGinTransport()` | ✅ | Gin framework integration |

#### API Style

```go
// Tool arguments as structs with jsonschema tags
type SearchArgs struct {
    Query string `json:"query" jsonschema:"required,description=The search query"`
}

// Server creation
server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
server.RegisterTool("search", "Search the web", func(args SearchArgs) (*mcp_golang.ToolResponse, error) {
    return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("result")), nil
})

// HTTP transport
transport := http.NewHTTPTransport("/mcp")
transport.WithAddr(":8080")
server := mcp_golang.NewServer(transport)
```

#### Pros
- **Type-safe** - Struct-based arguments with automatic schema generation
- **Low boilerplate** - Minimal code required
- **jsonschema tags** - Automatic schema from struct tags
- **Gin integration** - Easy web framework integration
- **Custom transports** - Extensible transport system

#### Cons
- **HTTP limitations** - HTTP transport is stateless, no bidirectional features
- **Fewer examples** - Less documentation than alternatives
- **Third-party** - Not official

---

### 4. findleyr/mcp

**Repository**: https://github.com/findleyr/mcp

| Aspect | Details |
|--------|---------|
| **Source Reputation** | Medium |
| **Benchmark Score** | 89.5 |
| **Code Snippets** | 120 |
| **Maturity** | Active, well-designed |

#### Transport Support

| Transport | Server | Client | Notes |
|-----------|--------|--------|-------|
| Stdio | ✅ `NewStdioTransport()` | ✅ `NewCommandTransport()` | Full support |
| SSE | ✅ `NewSSEHandler()` | ✅ `NewSSEClientTransport()` | Real-time streaming |
| HTTP | ❌ | ❌ | Not supported |

#### API Style

```go
// Server creation with options
serverOpts := &mcp.ServerOptions{
    Instructions: "A search server",
    PageSize:     100,
}
server := mcp.NewServer("search-server", "v1.0.0", serverOpts)

// Add tools with automatic schema inference
server.AddTools(mcp.NewServerTool(
    "search",
    "Search the web",
    SearchHandler,
))

// Or with custom schema
server.AddTools(mcp.NewServerTool(
    "advanced-search",
    "Advanced search",
    SearchHandler,
    mcp.Input(
        mcp.Property("query", mcp.Description("Search query"), mcp.Required(true)),
        mcp.Property("maxResults", mcp.Default(10)),
    ),
))

// Run server
server.Run(ctx, mcp.NewStdioTransport())
```

#### Pros
- **Type-safe with generics** - Strong typing throughout
- **Automatic schema inference** - Less boilerplate for schemas
- **Bidirectional communication** - Full support for notifications
- **Progress reporting** - Built-in support for long-running operations
- **Tool annotations** - Hints for clients (destructive, idempotent)
- **High benchmark score** - Well-architected

#### Cons
- **No HTTP transport** - Only stdio and SSE
- **Medium reputation** - Less established than official SDK
- **Third-party** - Not official

---

### 5. voocel/mcp-sdk-go

**Repository**: https://github.com/voocel/mcp-sdk-go

| Aspect | Details |
|--------|---------|
| **Source Reputation** | High |
| **Code Snippets** | 7 |
| **Maturity** | Newer, fewer examples |

#### Transport Support

| Transport | Server | Client | Notes |
|-----------|--------|--------|-------|
| Stdio | ❌ | ❌ | Not supported |
| SSE | ✅ `sse.NewServer()` | ✅ | Server-Sent Events |
| Streamable HTTP | ✅ `streamable.NewServer()` | ✅ | Recommended for new projects |

#### API Style

```go
// FastMCP server with chainable API
mcp := server.NewFastMCP("server-name", "1.0.0")

// Register tools with chainable API
mcp.Tool("tool_name", "Tool description").
    WithStringParam("param1", "Param description", true).
    WithIntParam("param2", "Param description", false).
    Handle(func(ctx context.Context, args map[string]interface{}) (*protocol.CallToolResult, error) {
        return protocol.NewToolResultText("result"), nil
    })

// Start SSE server
sseTransport := sse.NewServer(":8080", mcp)
sseTransport.Serve(ctx)
```

#### Pros
- **Chainable API** - Fluent, readable code
- **FastMCP** - Optimized for quick setup
- **Streamable HTTP** - Modern transport option

#### Cons
- **No stdio support** - Critical limitation for local tools
- **Few examples** - Limited documentation
- **Newer** - Less battle-tested
- **Third-party** - Not official

---

## Feature Comparison Matrix

| Feature | Official SDK | mark3labs | metoro-io | findleyr | voocel |
|---------|--------------|-----------|-----------|----------|--------|
| **Official** | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Stdio Transport** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **HTTP Transport** | ✅ (Streamable) | ✅ | ✅ | ❌ | ✅ (Streamable) |
| **SSE Transport** | ❌ | ✅ | ❌ | ✅ | ✅ |
| **Type Safety** | ✅ | ✅ | ✅ | ✅ (Generics) | ✅ |
| **Auto Schema** | ❌ | ❌ | ✅ (Tags) | ✅ (Inference) | ❌ |
| **Bidirectional** | ✅ | ✅ | ⚠️ (HTTP no) | ✅ | ⚠️ |
| **Progress Reporting** | ✅ | ✅ | ❌ | ✅ | ❌ |
| **Tool Annotations** | ✅ | ❌ | ❌ | ✅ | ❌ |
| **OAuth Support** | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Gin Integration** | ❌ | ❌ | ✅ | ❌ | ❌ |
| **API Simplicity** | Medium | High | High | Medium | High |
| **Documentation** | Excellent | Good | Good | Good | Limited |
| **Community** | Large | Growing | Small | Small | Small |

## Recommendations by Use Case

### For This Project (Web Search MCP Server)

**Requirements:**
- Both HTTP and stdio transport support
- Wrap existing CLI tools (`ddg-search`, `perplexity-search`, `page-dump`)
- Fast implementation
- Claude-compatible

**Top Recommendations:**

#### 1. **mark3labs/mcp-go** (Recommended for Speed)

**Why:**
- ✅ Both HTTP and stdio supported
- ✅ Very simple API - minimal boilerplate
- ✅ Transport-agnostic - same code works for both transports
- ✅ Fast to implement
- ✅ Good documentation

**Example Implementation:**
```go
package main

import (
    "context"
    "os/exec"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    s := server.NewMCPServer("ddg-search", "1.0.0")

    // Register tools
    s.AddTool(server.Tool{
        Name:        "web-search",
        Description: "Search the web using DuckDuckGo",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{
                    "type":        "string",
                    "description": "Search query",
                },
            },
            "required": []string{"query"},
        },
    }, func(ctx context.Context, request server.CallToolRequest) (*server.CallToolResult, error) {
        cmd := exec.Command("ddg-search", request.Params.Arguments["query"].(string))
        output, _ := cmd.CombinedOutput()
        return &server.CallToolResult{
            Content: []interface{}{
                map[string]interface{}{
                    "type": "text",
                    "text": string(output),
                },
            },
        }, nil
    })

    // Choose transport at runtime
    if os.Getenv("MCP_TRANSPORT") == "http" {
        server.ServeHTTP(s, ":8080")
    } else {
        server.ServeStdio(s)
    }
}
```

#### 2. **modelcontextprotocol/go-sdk** (Recommended for Stability)

**Why:**
- ✅ Official SDK - long-term support guaranteed
- ✅ Both HTTP and stdio supported
- ✅ Comprehensive feature set
- ✅ Well-tested and documented
- ✅ OAuth support if needed later

**Trade-off:** More verbose API, but more control and stability.

#### 3. **metoro-io/mcp-golang** (Recommended for Type Safety)

**Why:**
- ✅ Both HTTP and stdio supported
- ✅ Type-safe with struct-based arguments
- ✅ Automatic schema generation from tags
- ✅ Low boilerplate

**Trade-off:** HTTP transport is stateless (no bidirectional features).

---

## Decision Framework

### Choose **mark3labs/mcp-go** if:
- ✅ You want the fastest implementation
- ✅ You need both HTTP and stdio
- ✅ You prefer simple, clean APIs
- ✅ You want transport flexibility at runtime

### Choose **modelcontextprotocol/go-sdk** if:
- ✅ You want official support
- ✅ You need maximum stability
- ✅ You might need OAuth later
- ✅ You prefer comprehensive documentation

### Choose **metoro-io/mcp-golang** if:
- ✅ You want type safety with structs
- ✅ You like automatic schema generation
- ✅ You don't need bidirectional HTTP features
- ✅ You might use Gin framework

### Choose **findleyr/mcp** if:
- ✅ You only need stdio and SSE
- ✅ You want automatic schema inference
- ✅ You need progress reporting
- ✅ You want tool annotations

### Avoid **voocel/mcp-sdk-go** if:
- ❌ You need stdio transport (not supported)

---

## Questions for Decision

1. **Priority: Speed vs. Stability?**
   - Speed → mark3labs/mcp-go
   - Stability → modelcontextprotocol/go-sdk

2. **Do you need bidirectional HTTP features?**
   - Yes → mark3labs or official SDK
   - No → metoro-io is also fine

3. **Do you prefer struct-based type safety?**
   - Yes → metoro-io/mcp-golang
   - No → mark3labs or official SDK

4. **Is official support important?**
   - Yes → modelcontextprotocol/go-sdk
   - No → mark3labs is fine

5. **Do you need OAuth authentication?**
   - Yes → modelcontextprotocol/go-sdk only
   - No → any option works

---

## My Recommendation

For this project (web search MCP server wrapping existing CLI tools), I recommend **mark3labs/mcp-go** because:

1. ✅ **Fastest to implement** - Minimal boilerplate
2. ✅ **Both transports supported** - HTTP and stdio
3. ✅ **Transport-agnostic** - Same code works for both
4. ✅ **Simple API** - Easy to wrap CLI tools
5. ✅ **Good documentation** - Plenty of examples
6. ✅ **Active development** - Well-maintained

The official SDK is a close second if you prioritize long-term stability and official support over implementation speed.

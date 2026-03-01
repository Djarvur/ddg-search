# MCP Go Packages Comparison

Comparison of Go packages for building Model Context Protocol (MCP) servers, specifically for creating a Claude Code compatible web search MCP server.

## Requirements

- **Transports**: Both HTTP and stdio must be supported
- **Goal**: Build MCP server for web search based on existing skills
- **Priority**: Get results faster (choose one that's quick to implement)

---

## Packages Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MCP GO PACKAGES LANDSCAPE                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  OFFICIAL                 COMMUNITY                                         │
│  ═════════               ═══════════                                        │
│                                                                             │
│  ┌──────────────────┐    ┌──────────────────┐  ┌──────────────────┐        │
│  │ modelcontextproto│    │ mark3labs/mcp-go │  │ metoro-io        │        │
│  │ col/go-sdk       │    │                  │  │ mcp-golang       │        │
│  │                  │    │ ★ Most Popular   │  │                  │        │
│  │ ★ Official       │    │ ★ High Quality   │  │ ★ Type Safe      │        │
│  │ ★ High Quality   │    │ ★ Fast           │  │ ★ Low Boilerplate│        │
│  └──────────────────┘    └──────────────────┘  └──────────────────┘        │
│         │                        │                      │                   │
│         │                        │                      │                   │
│  ┌──────┴──────┐          ┌──────┴──────┐        ┌──────┴──────┐           │
│  │ Benchmark:  │          │ Benchmark:  │        │ Benchmark:  │           │
│  │   85.7      │          │   81.8      │        │   N/A       │           │
│  │ Snippets:   │          │ Snippets:   │        │ Snippets:   │           │
│  │   366       │          │   580       │        │   79        │           │
│  └─────────────┘          └─────────────┘        └─────────────┘           │
│                                                                             │
│  OTHER OPTIONS                                                              │
│  ══════════════                                                             │
│                                                                             │
│  ┌──────────────────┐    ┌──────────────────┐                              │
│  │ trpc-group       │    │ findleyr/mcp     │                              │
│  │ trpc-mcp-go      │    │                  │                              │
│  │                  │    │ ★ Highest Score  │                              │
│  │ ★ Streaming HTTP │    │ ★ Simple API     │                              │
│  │ ★ SSE Support    │    │                  │                              │
│  └──────────────────┘    └──────────────────┘                              │
│         │                        │                                          │
│  ┌──────┴──────┐          ┌──────┴──────┐                                  │
│  │ Benchmark:  │          │ Benchmark:  │                                  │
│  │   80.3      │          │   89.5      │                                  │
│  │ Snippets:   │          │ Snippets:   │                                  │
│  │   73        │          │   120       │                                  │
│  └─────────────┘          └─────────────┘                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Detailed Comparison

### 1. modelcontextprotocol/go-sdk (Official)

**Repository**: https://github.com/modelcontextprotocol/go-sdk

| Aspect | Details |
|--------|---------|
| **Status** | ✅ Official SDK from Anthropic |
| **Benchmark Score** | 85.7 |
| **Code Snippets** | 366 |
| **Source Reputation** | High |
| **Versions** | v0.2.0, v0.4.0, v1.0.0, v1.2.0 |

**Transport Support**:
```
┌─────────────────┐     ┌─────────────────┐
│    stdio        │     │     HTTP        │
│   ✅ Full       │     │   ✅ Full       │
│   StdioTransport│     │ StreamableHTTP  │
└─────────────────┘     └─────────────────┘
```

**Pros**:
- ✅ Official SDK - best long-term support
- ✅ High code quality and documentation
- ✅ OAuth authentication support
- ✅ JSON Schema and JSON RPC implementations
- ✅ Semantic versioning with stable releases

**Cons**:
- ⚠️ Slightly more verbose API
- ⚠️ Less community examples than mcp-go

**Example - Basic Server**:
```go
package main

import (
    "context"
    "log"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
    Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
    Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (
    *mcp.CallToolResult, Output, error,
) {
    return nil, Output{Greeting: "Hi " + input.Name}, nil
}

func main() {
    server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
    mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        log.Fatal(err)
    }
}
```

**Example - HTTP Server**:
```go
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
    return server
}, nil)

go http.ListenAndServe(":8080", handler)
```

---

### 2. mark3labs/mcp-go

**Repository**: https://github.com/mark3labs/mcp-go

| Aspect | Details |
|--------|---------|
| **Status** | ✅ Most popular community SDK |
| **Benchmark Score** | 81.8 |
| **Code Snippets** | 580 |
| **Source Reputation** | High |

**Transport Support**:
```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    stdio        │     │   SSE HTTP      │     │ Streamable HTTP │
│   ✅ Full       │     │   ✅ Full       │     │   ✅ Full       │
│   ServeStdio    │     │   SSEServer     │     │ StreamableHTTP  │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

**Pros**:
- ✅ Most code snippets/examples available
- ✅ Clean, intuitive API
- ✅ Multiple HTTP transport options (SSE, Streamable HTTP)
- ✅ Good documentation with dedicated website
- ✅ Active community

**Cons**:
- ⚠️ Not official - may lag behind spec changes
- ⚠️ More dependencies

**Example - Basic Server**:
```go
package main

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    s := server.NewMCPServer(
        "Demo 🚀",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

    tool := mcp.NewTool("hello_world",
        mcp.WithDescription("Say hello to someone"),
        mcp.WithString("name",
            mcp.Required(),
            mcp.Description("Name of the person to greet"),
        ),
    )

    s.AddTool(tool, helloHandler)

    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, err := request.RequireString("name")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
```

**Example - HTTP Server**:
```go
// SSE Transport
sseServer := server.NewSSEServer(s, "/mcp")
http.Handle("/mcp", sseServer)
http.ListenAndServe(":8080", nil)

// Streamable HTTP Transport (Modern)
httpServer := server.NewStreamableHTTPServer(s,
    server.WithEndpointPath("/mcp"),
    server.WithHeartbeatInterval(30*time.Second),
)
httpServer.Start(":8080")
```

---

### 3. metoro-io/mcp-golang

**Repository**: https://github.com/metoro-io/mcp-golang

| Aspect | Details |
|--------|---------|
| **Status** | ✅ Active community SDK |
| **Code Snippets** | 79 |
| **Source Reputation** | High |

**Transport Support**:
```
┌─────────────────┐     ┌─────────────────┐
│    stdio        │     │     HTTP        │
│   ✅ Full       │     │   ⚠️ Stateless  │
│ StdioServerTrans│     │  HTTPTransport  │
└─────────────────┘     └─────────────────┘
```

**Pros**:
- ✅ Type-safe tool definitions using Go structs
- ✅ Low boilerplate code
- ✅ Automatic schema generation from struct tags
- ✅ Gin framework integration available

**Cons**:
- ⚠️ HTTP transport is stateless (no bidirectional features)
- ⚠️ Fewer examples and documentation
- ⚠️ Smaller community

**Example - Basic Server**:
```go
package main

import (
    "fmt"
    "github.com/metoro-io/mcp-golang"
    "github.com/metoro-io/mcp-golang/transport/stdio"
)

type MyArguments struct {
    Submitter string `json:"submitter" jsonschema:"required,description=The name"`
}

func main() {
    done := make(chan struct{})

    server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
    err := server.RegisterTool("hello", "Say hello", func(arguments MyArguments) (*mcp_golang.ToolResponse, error) {
        return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Submitter))), nil
    })
    if err != nil {
        panic(err)
    }

    err = server.Serve()
    if err != nil {
        panic(err)
    }

    <-done
}
```

**Example - HTTP Server**:
```go
// Standard HTTP
transport := http.NewHTTPTransport("/mcp")
transport.WithAddr(":8080")
server := mcp_golang.NewServer(transport)

// Or with Gin framework
transport := http.NewGinTransport()
router := gin.Default()
router.POST("/mcp", transport.Handler())
server := mcp_golang.NewServer(transport)
```

---

### 4. trpc-group/trpc-mcp-go

**Repository**: https://github.com/trpc-group/trpc-mcp-go

| Aspect | Details |
|--------|---------|
| **Status** | ✅ Active |
| **Benchmark Score** | 80.3 |
| **Code Snippets** | 73 |
| **Source Reputation** | Medium |

**Transport Support**:
```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    stdio        │     │   SSE HTTP      │     │ Streamable HTTP │
│   ✅ Full       │     │   ✅ Full       │     │   ✅ Full       │
│ NewStdioServer  │     │ WithPostSSE     │     │ NewServer       │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

**Pros**:
- ✅ Comprehensive streaming HTTP support
- ✅ Type-safe handlers with struct input/output
- ✅ Automatic schema generation
- ✅ Clean separation of stdio and HTTP servers

**Cons**:
- ⚠️ Less popular/smaller community
- ⚠️ Fewer examples

**Example - Server with Type Safety**:
```go
type WeatherInput struct {
    Location string `json:"location" jsonschema:"required,description=City name"`
    Units    string `json:"units,omitempty" jsonschema:"description=Temperature units"`
}

type WeatherOutput struct {
    Temperature float64 `json:"temperature" jsonschema:"description=Current temperature"`
}

weatherTool := mcp.NewTool(
    "get_weather",
    mcp.WithDescription("Get weather information"),
    mcp.WithInputStruct[WeatherInput](),
    mcp.WithOutputStruct[WeatherOutput](),
)

weatherHandler := mcp.NewTypedToolHandler(func(ctx context.Context, req *mcp.CallToolRequest, input WeatherInput) (WeatherOutput, error) {
    return WeatherOutput{Temperature: 22.5}, nil
})
```

---

### 5. findleyr/mcp

**Repository**: https://github.com/findleyr/mcp

| Aspect | Details |
|--------|---------|
| **Status** | ✅ Active |
| **Benchmark Score** | 89.5 (Highest!) |
| **Code Snippets** | 120 |
| **Source Reputation** | Medium |

**Transport Support**:
```
┌─────────────────┐
│    stdio        │
│   ✅ Full       │
│ StdioTransport  │
└─────────────────┘
```

**Pros**:
- ✅ Highest benchmark score (89.5)
- ✅ Clean, simple API
- ✅ Automatic schema inference
- ✅ Tool annotations support

**Cons**:
- ⚠️ HTTP transport not clearly documented
- ⚠️ Smaller community
- ⚠️ Less examples

**Example - Server with Annotations**:
```go
server := mcp.NewServer("search-server", "v1.0.0", &mcp.ServerOptions{
    Instructions: "A search server",
    PageSize:     100,
})

// Tool with automatic schema inference
server.AddTools(mcp.NewServerTool(
    "search",
    "Search the web",
    SearchHandler,
))

// Tool with annotations
server.AddTools(mcp.NewServerTool(
    "risky-op",
    "Potentially risky operation",
    RiskyHandler,
    mcp.Annotations(&mcp.ToolAnnotations{
        DestructiveHint: true,
        IdempotentHint:  false,
    }),
))

server.Run(ctx, mcp.NewStdioTransport())
```

---

## Feature Comparison Matrix

| Feature | go-sdk (official) | mcp-go | mcp-golang | trpc-mcp-go | findleyr/mcp |
|---------|:-----------------:|:------:|:----------:|:-----------:|:------------:|
| **stdio transport** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **HTTP transport** | ✅ | ✅ | ⚠️ Stateless | ✅ | ❓ |
| **SSE support** | ❓ | ✅ | ❌ | ✅ | ❓ |
| **Streamable HTTP** | ✅ | ✅ | ❌ | ✅ | ❓ |
| **Type safety** | ✅ | ⚠️ | ✅ | ✅ | ✅ |
| **Schema generation** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Resources support** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Prompts support** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **OAuth support** | ✅ | ❓ | ❌ | ❓ | ❓ |
| **Documentation** | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ | ⭐ |
| **Community size** | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐ | ⭐ |
| **Code examples** | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ | ⭐⭐ |
| **Benchmark** | 85.7 | 81.8 | N/A | 80.3 | 89.5 |

---

## Recommendation for Web Search MCP Server

### For Getting Results Faster: **mark3labs/mcp-go** 🏆

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        RECOMMENDATION: mcp-go                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  WHY IT'S THE BEST CHOICE FOR THIS PROJECT:                                │
│                                                                             │
│  1. MOST EXAMPLES ⭐                                                         │
│     ├── 580 code snippets available                                         │
│     ├── Dedicated documentation website (mcp-go.dev)                        │
│     └── Active community sharing patterns                                   │
│                                                                             │
│  2. TRANSPORT FLEXIBILITY ⭐                                                 │
│     ├── stdio: ServeStdio(s) - one line                                    │
│     ├── SSE: NewSSEServer(s, "/mcp") - simple                              │
│     └── HTTP: NewStreamableHTTPServer(s) - modern                          │
│                                                                             │
│  3. SIMPLE API ⭐                                                            │
│     ├── Tool creation: mcp.NewTool()                                        │
│     ├── Add handlers: s.AddTool(tool, handler)                              │
│     └── Start server: server.ServeStdio(s)                                  │
│                                                                             │
│  4. MAPPING TO EXISTING SKILLS                                              │
│     ┌────────────────────┐     ┌────────────────────┐                       │
│     │  ddg-search skill  │ ──▶ │  web-search tool   │                       │
│     │  page-dump skill   │ ──▶ │  page-dump tool    │                       │
│     │  perplexity skill  │ ──▶ │  perplexity tool   │                       │
│     └────────────────────┘     └────────────────────┘                       │
│                                                                             │
│  5. QUICK START PATH                                                        │
│     Day 1: Implement stdio server with ddg-search tool                      │
│     Day 2: Add page-dump and perplexity tools                               │
│     Day 3: Add HTTP transport support                                       │
│     Day 4: Testing and documentation                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Alternative: Official SDK for Long-term Stability

If long-term maintenance and official support are priorities over speed:

**modelcontextprotocol/go-sdk** - Official SDK with:
- Guaranteed spec compliance
- OAuth support for future auth needs
- Semantic versioning stability

---

## Implementation Sketch for ddg-search MCP Server

```go
package main

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"

    "github.com/Djarvur/ddg-search/internal/search"
    "github.com/Djarvur/ddg-search/internal/dump"
)

func main() {
    // Create MCP server
    s := server.NewMCPServer(
        "ddg-search",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

    // Add web-search tool (maps to ddg-search skill)
    webSearchTool := mcp.NewTool("web-search",
        mcp.WithDescription("Search the web using DuckDuckGo"),
        mcp.WithString("query",
            mcp.Required(),
            mcp.Description("The search query string"),
        ),
        mcp.WithNumber("max-results",
            mcp.Description("Maximum number of results (default: 10)"),
        ),
        mcp.WithString("site",
            mcp.Description("Filter results to a specific domain"),
        ),
    )
    s.AddTool(webSearchTool, handleWebSearch)

    // Add page-dump tool (maps to page-dump skill)
    pageDumpTool := mcp.NewTool("page-dump",
        mcp.WithDescription("Fetch a web page and convert to markdown"),
        mcp.WithString("url",
            mcp.Required(),
            mcp.Description("The URL to fetch"),
        ),
    )
    s.AddTool(pageDumpTool, handlePageDump)

    // Choose transport based on environment
    transport := os.Getenv("MCP_TRANSPORT")
    switch transport {
    case "http":
        httpServer := server.NewStreamableHTTPServer(s)
        httpServer.Start(":8080")
    default:
        server.ServeStdio(s)
    }
}

func handleWebSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    query, err := req.RequireString("query")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    // Use existing search.Client from internal/search
    client := search.NewClient()
    results, err := client.Search(ctx, query)
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    return mcp.NewToolResultText(formatResults(results)), nil
}

func handlePageDump(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    url, err := req.RequireString("url")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    // Use existing dump functionality from internal/dump
    content, err := dump.FetchAndConvert(ctx, url)
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    return mcp.NewToolResultText(content), nil
}
```

---

## Questions for Decision

Before proceeding with implementation, please confirm:

1. **Primary transport**: Should we prioritize stdio (for Claude Code) or HTTP (for web services)?
2. **Package choice**: Do you agree with **mark3labs/mcp-go** for faster implementation, or prefer the official SDK?
3. **Tool scope**: Should we include all three skills (ddg-search, page-dump, perplexity) or start with a subset?
4. **Configuration**: Environment variables or config file for settings like API keys?

---

## References

- [MCP Specification](https://modelcontextprotocol.io/)
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - Documentation: https://mcp-go.dev
- [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official SDK
- [metoro-io/mcp-golang](https://github.com/metoro-io/mcp-golang)
- [trpc-group/trpc-mcp-go](https://github.com/trpc-group/trpc-mcp-go)
- [findleyr/mcp](https://github.com/findleyr/mcp)

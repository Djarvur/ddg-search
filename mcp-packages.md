# MCP Go Packages Comparison

**Date:** 2026-03-01
**Purpose:** Comprehensive analysis of Go packages for building Model Context Protocol (MCP) servers compatible with Claude Code.

---

## Executive Summary

This document compares four major Go implementations of the Model Context Protocol (MCP) for building MCP servers that can be integrated with Claude Code. The comparison focuses on transport support (stdio and HTTP), ease of use, community adoption, and suitability for exposing existing search skills (`ddg-search`, `perplexity-search`, `page-dump`).

### Quick Recommendation

| Package | Stars | Best For | Speed to Result |
|---------|-------|----------|-----------------|
| **mark3labs/mcp-go** | 8,253 | Quick development, minimal boilerplate | ⭐⭐⭐⭐⭐ Fastest |
| **modelcontextprotocol/go-sdk** | 3,973 | Official support, long-term stability | ⭐⭐⭐ Moderate |
| **metoro-io/mcp-golang** | 1,200 | Type safety, simple API | ⭐⭐⭐⭐ Fast |
| **ThinkInAIXYZ/go-mcp** | 666 | Web framework integration | ⭐⭐⭐ Moderate |

---

## Detailed Comparison

### 1. mark3labs/mcp-go

**Repository:** https://github.com/mark3labs/mcp-go
**Stars:** 8,253 | **Forks:** 776 | **Issues:** 20
**License:** MIT | **Created:** 2024-11-27 | **Last Push:** 2026-02-27

#### Overview
A high-level Go implementation of MCP focused on simplicity and fast development. Handles all complex protocol details so developers can focus on building tools.

#### Key Features
- **Fast Development:** High-level interface with minimal boilerplate
- **Simple API:** Build MCP servers with very little code
- **Complete Implementation:** Aims to provide full MCP spec support
- **Active Development:** Under active development alongside MCP spec
- **Documentation:** Comprehensive docs at http://mcp-go.dev/

#### Transport Support
- ✅ **stdio**: Full support via `server.ServeStdio(s)`
- ✅ **HTTP**: Supported (see Extras/Transports section in docs)
- ✅ **SSE**: Server-Sent Events support
- ✅ **Session Management**: Per-session tools, tool filtering, context handling
- ✅ **Request Hooks**: Middleware for request handling
- ✅ **Tool Handler Middleware**: Chainable middleware for tools

#### Code Example (Server)
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

#### Pros
- **Most popular** community implementation (highest stars)
- **Very concise API** - minimal boilerplate
- **Excellent documentation** with examples
- **Active community** (Discord support)
- **Fast to get started** - simplest learning curve
- **Full feature set** including sessions, hooks, middleware

#### Cons
- **Not official** - community maintained
- **Under active development** - some advanced features still in progress
- **MIT license** (generally fine, but different from official)

#### Suitability for This Project
⭐⭐⭐⭐⭐ **Excellent**

- Existing binaries (`ddg-search`, `perplexity-search`, `page-dump`) can be easily wrapped as MCP tools
- Simple handler functions can call existing Go code or execute commands
- Both stdio and HTTP transports supported
- Minimal code changes required to expose existing functionality

---

### 2. modelcontextprotocol/go-sdk

**Repository:** https://github.com/modelcontextprotocol/go-sdk
**Stars:** 3,973 | **Forks:** 365 | **Issues:** 54
**License:** Apache 2.0 (new), MIT (existing) | **Created:** 2025-04-23 | **Last Push:** 2026-03-01

#### Overview
**The official Go SDK** maintained by the Model Context Protocol organization in collaboration with Google. Provides complete implementation of the MCP specification with full conformance testing.

#### Key Features
- **Official Support:** Maintained by MCP org + Google
- **Full Spec Compliance:** Complete MCP 2025-11-25 specification support
- **Conformance Testing:** Includes conformance tests
- **Modular Design:** Separate packages for mcp, jsonrpc, auth, oauthex
- **Version Compatibility:** Supports multiple MCP spec versions
- **Comprehensive Docs:** Detailed feature documentation in `/docs/`

#### Transport Support
- ✅ **stdio**: Full support via `mcp.StdioTransport{}`
- ✅ **HTTP**: Full support via `mcp.StreamableHTTPTransport`
- ✅ **SSE**: Server-Sent Events support
- ✅ **Command Transport**: Execute server as subprocess
- ✅ **OAuth**: Client-side OAuth support (experimental)
- ✅ **Auth Framework**: Full authorization support for HTTP

#### Code Example (Server)
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
    *mcp.CallToolResult,
    Output,
    error,
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

#### Pros
- **Official SDK** - backed by MCP organization and Google
- **Long-term stability** - official maintenance guarantees
- **Full spec compliance** - complete implementation
- **Conformance testing** - ensures correctness
- **Apache 2.0 license** - permissive and business-friendly
- **Active development** - very recent commits (today!)
- **Comprehensive documentation** - detailed feature mapping

#### Cons
- **More verbose** API compared to community alternatives
- **Steeper learning curve** - more concepts to understand
- **Newer project** (created April 2025) - less battle-tested
- **More boilerplate** for simple use cases

#### Suitability for This Project
⭐⭐⭐⭐ **Very Good**

- Official support ensures long-term compatibility
- Both stdio and HTTP transports fully supported
- Can wrap existing binaries or integrate Go code directly
- Slightly more setup than mark3labs/mcp-go
- Best choice if you prioritize official support and spec compliance

---

### 3. metoro-io/mcp-golang

**Repository:** https://github.com/metoro-io/mcp-golang
**Stars:** 1,200 | **Forks:** 118 | **Issues:** 44
**License:** MIT | **Created:** 2024-12-07 | **Last Push:** 2026-02-25

#### Overview
An unofficial MCP implementation focusing on type safety and low boilerplate. Uses native Go structs for tool arguments with automatic schema generation.

#### Key Features
- **Type Safety:** Define tool arguments as native Go structs
- **Auto Schema Generation:** jsonschema tags generate MCP schemas
- **Custom Transports:** Built-in stdio and HTTP, or write your own
- **Low Boilerplate:** Generates all MCP endpoints automatically
- **Modular:** Three components - transport, protocol, server/client
- **Bi-directional:** Full server and client support via stdio

#### Transport Support
- ✅ **stdio**: Full feature support (recommended for Claude Desktop)
- ✅ **HTTP**: Stateless communication (no bidirectional features)
- ✅ **Gin Integration**: HTTP transport with Gin framework
- ⚠️ **HTTP Limitation**: Stateless - no notifications support

#### Code Example (Server)
```go
package main

import (
    "fmt"
    "github.com/metoro-io/mcp-golang"
    "github.com/metoro-io/mcp-golang/transport/stdio"
)

type Content struct {
    Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
    Description *string `json:"description" jsonschema:"description=The description to submit"`
}

type MyFunctionsArguments struct {
    Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool"`
    Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

func main() {
    done := make(chan struct{})
    server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

    err := server.RegisterTool("hello", "Say hello to a person",
        func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
            return mcp_golang.NewToolResponse(
                mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Submitter))
            ), nil
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

#### Pros
- **Type-safe** - compile-time checking for tool arguments
- **Clean API** - very intuitive and Go-idiomatic
- **Auto schema generation** - jsonschema tags reduce boilerplate
- **Good documentation** - docs at https://mcpgolang.com
- **Fast development** - minimal boilerplate
- **Gin integration** - easy to add to existing Gin apps

#### Cons
- **HTTP is stateless** - no bidirectional features (notifications)
- **Not official** - community maintained
- **Fewer stars** than alternatives
- **HTTP limitation** means stdio required for full MCP features

#### Suitability for This Project
⭐⭐⭐⭐ **Very Good**

- Type-safe API is great for defining tool schemas
- Can wrap existing binaries easily
- stdio transport works well with Claude Desktop
- HTTP transport available but limited (stateless)
- Good balance of simplicity and type safety

---

### 4. ThinkInAIXYZ/go-mcp

**Repository:** https://github.com/ThinkInAIXYZ/go-mcp
**Stars:** 666 | **Forks:** 109 | **Issues:** 7
**License:** MIT | **Created:** 2025-03-04 | **Last Push:** 2025-10-10

#### Overview
A powerful Go MCP SDK with focus on web framework integration and elegant three-layer architecture. Provides http.Handler for easy integration with existing web services.

#### Key Features
- **Complete Protocol Implementation:** Full MCP spec
- **Three-Layer Architecture:** Clear separation of concerns
- **Web Framework Integration:** Provides http.Handler
- **Type Safety:** Strong typing with Go's type system
- **Simple Deployment:** Static compilation, no complex dependencies
- **High Performance:** Leverages Go's concurrency

#### Transport Support
- ✅ **stdio**: Full support (stdio_client.go, stdio_server.go)
- ✅ **SSE**: Server-Sent Events (sse_client.go, sse_server.go)
- ✅ **Streamable HTTP**: Full support (streamable_http_client.go, streamable_http_server.go)
- ✅ **http.Handler**: Direct integration with web frameworks

#### Code Example (Server)
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/ThinkInAIXYZ/go-mcp/protocol"
    "github.com/ThinkInAIXYZ/go-mcp/server"
    "github.com/ThinkInAIXYZ/go-mcp/transport"
)

type TimeRequest struct {
    Timezone string `json:"timezone" description:"timezone" required:"true"`
}

func main() {
    transportServer, err := transport.NewSSEServerTransport("127.0.0.1:8080")
    if err != nil {
        log.Fatalf("Failed to create transport server: %v", err)
    }

    mcpServer, err := server.NewServer(transportServer)
    if err != nil {
        log.Fatalf("Failed to create MCP server: %v", err)
    }

    tool, err := protocol.NewTool("current_time", "Get current time for specified timezone", TimeRequest{})
    if err != nil {
        log.Fatalf("Failed to create tool: %v", err)
    }
    mcpServer.RegisterTool(tool, handleTimeRequest)

    if err = mcpServer.Run(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}

func handleTimeRequest(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
    var timeReq TimeRequest
    if err := protocol.VerifyAndUnmarshal(req.RawArguments, &timeReq); err != nil {
        return nil, err
    }
    // ... handle request
}
```

#### Pros
- **Web framework integration** - http.Handler for easy embedding
- **Clean architecture** - three-layer design
- **Full transport support** - stdio, SSE, streamable HTTP
- **Type-safe** - leverages Go's type system
- **High performance** - optimized for concurrency
- **Multi-language docs** - Chinese, Traditional Chinese, Vietnamese

#### Cons
- **Not official** - community maintained
- **Fewer stars** - smaller community
- **Less active** - last push was October 2025
- **Less documentation** compared to alternatives

#### Suitability for This Project
⭐⭐⭐ **Good**

- Good if you want to embed MCP in an existing HTTP service
- Full transport support including stdio and HTTP
- Can wrap existing binaries
- Architecture is clean but may be overkill for simple use case

---

## Comparison Matrix

| Feature | mark3labs/mcp-go | modelcontextprotocol/go-sdk | metoro-io/mcp-golang | ThinkInAIXYZ/go-mcp |
|---------|------------------|-----------------------------|----------------------|---------------------|
| **Stars** | 8,253 | 3,973 | 1,200 | 666 |
| **Official** | ❌ | ✅ | ❌ | ❌ |
| **License** | MIT | Apache 2.0 | MIT | MIT |
| **stdio Transport** | ✅ | ✅ | ✅ | ✅ |
| **HTTP Transport** | ✅ | ✅ | ⚠️ Stateless | ✅ |
| **SSE Transport** | ✅ | ✅ | ❌ | ✅ |
| **Type Safety** | ✅ | ✅ | ✅✅ (structs) | ✅ |
| **Boilerplate** | Low | Medium | Low | Medium |
| **Learning Curve** | Easy | Medium | Easy | Medium |
| **Documentation** | Excellent | Excellent | Good | Fair |
| **Community Size** | Large | Growing | Medium | Small |
| **Active Development** | ✅ | ✅✅ | ✅ | ⚠️ |
| **Last Commit** | 2026-02-27 | 2026-03-01 | 2026-02-25 | 2025-10-10 |
| **Open Issues** | 20 | 54 | 44 | 7 |
| **Middleware Support** | ✅ | ✅ | ❌ | ✅ |
| **Claude Desktop** | ✅ | ✅ | ✅ | ✅ |
| **HTTP Auth** | ✅ | ✅✅ OAuth | ❌ | ✅ |

---

## Transport Details

### stdio Transport
All four packages support stdio transport, which is:
- **Recommended for Claude Desktop** integration
- **Full feature support** - all MCP capabilities available
- **Bidirectional** - supports notifications and streaming
- **Simple** - just stdin/stdout communication

### HTTP Transport
| Package | HTTP Support | Features | Notes |
|---------|--------------|----------|-------|
| mark3labs/mcp-go | ✅ Full | All features | Standard HTTP transport |
| modelcontextprotocol/go-sdk | ✅ Full | All features + OAuth | Streamable HTTP |
| metoro-io/mcp-golang | ⚠️ Stateless | No notifications | Use stdio for full features |
| ThinkInAIXYZ/go-mcp | ✅ Full | All features | http.Handler integration |

**Key Difference:** HTTP transport in MCP can be either:
1. **Stateless** (metoro-io/mcp-golang): No session state, no notifications
2. **Stateful/Streamable** (others): Full MCP features including notifications

For Claude Desktop compatibility, stdio is the primary transport. HTTP is useful for:
- Web-based MCP clients
- Server-side integration
- Remote MCP servers

---

## Recommendation for This Project

### Use Case Analysis

**Existing Skills to Expose:**
1. `ddg-search` - DuckDuckGo web search (no API key)
2. `perplexity-search` - Perplexity API search (requires API key)
3. `page-dump` - Fetch web page and convert to markdown

**Requirements:**
- ✅ Both HTTP and stdio support
- ✅ Fast to implement
- ✅ Compatible with Claude Code
- ✅ Can wrap existing binaries or integrate Go code

### Top Recommendation: mark3labs/mcp-go

**Why:**
1. **Fastest to implement** - minimal boilerplate, simple API
2. **Most popular** - largest community, most battle-tested
3. **Full transport support** - both stdio and HTTP with all features
4. **Easy to wrap binaries** - simple handler functions
5. **Excellent documentation** - lots of examples
6. **Active development** - recent commits, responsive community

**Implementation Approach:**
```go
// Create MCP server
s := server.NewMCPServer("ddg-search", "1.0.0")

// Add web-search tool (wraps ddg-search binary)
webSearchTool := mcp.NewTool("web-search",
    mcp.WithDescription("Search the web using DuckDuckGo"),
    mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
    mcp.WithNumber("max-results", mcp.Description("Max results (default: 10)")),
)
s.AddTool(webSearchTool, handleWebSearch)

// Add perplexity-search tool (wraps perplexity-search binary)
// Add page-dump tool (wraps page-dump binary)

// Serve via stdio for Claude Desktop
server.ServeStdio(s)

// OR serve via HTTP
server.ServeHTTP(s, ":8080")
```

### Alternative: modelcontextprotocol/go-sdk

**Choose this if:**
- You want **official support** and long-term stability
- You need **full OAuth support** for HTTP
- You value **spec compliance** and conformance testing
- You're okay with slightly more boilerplate

### Alternative: metoro-io/mcp-golang

**Choose this if:**
- You want **type-safe** tool definitions with Go structs
- You prefer the **cleanest, most Go-idiomatic API**
- You're okay with HTTP being stateless (use stdio for full features)
- You like automatic schema generation via jsonschema tags

---

## Questions for Decision Making

Before finalizing your choice, consider:

1. **Long-term Maintenance Priority**
   - Do you prefer official support (modelcontextprotocol/go-sdk) or community innovation (mark3labs/mcp-go)?

2. **Type Safety Preference**
   - Do you want compile-time type checking for tool arguments? (metoro-io/mcp-golang excels here)

3. **HTTP Feature Requirements**
   - Do you need OAuth authentication? (only modelcontextprotocol/go-sdk has full OAuth support)
   - Do you need notifications over HTTP? (metoro-io/mcp-golang HTTP is stateless)

4. **Integration Approach**
   - Will you wrap existing binaries or integrate Go code directly? (all support both)
   - Do you need to embed MCP in an existing HTTP server? (ThinkInAIXYZ/go-mcp has http.Handler)

5. **Development Speed vs. Stability**
   - Fastest to implement: mark3labs/mcp-go
   - Most stable/official: modelcontextprotocol/go-sdk

---

## Implementation Quick Start

### Using mark3labs/mcp-go (Recommended)

```bash
# Install the package
go get github.com/mark3labs/mcp-go
```

```go
package main

import (
    "context"
    "fmt"
    "os/exec"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    s := server.NewMCPServer(
        "Web Search MCP",
        "1.0.0",
        server.WithToolCapabilities(true),
        server.WithResourceCapabilities(false),
    )

    // Register web-search tool
    webSearchTool := mcp.NewTool("web-search",
        mcp.WithDescription("Search the web using DuckDuckGo"),
        mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
        mcp.WithNumber("max-results", mcp.Description("Max results (default: 10)")),
        mcp.WithString("site", mcp.Description("Filter to specific domain")),
        mcp.WithString("region", mcp.Description("Search region (default: us-en)")),
        mcp.WithString("time", mcp.Description("Time filter: d, w, m, y")),
    )
    s.AddTool(webSearchTool, handleWebSearch)

    // Register page-dump tool
    pageDumpTool := mcp.NewTool("page-dump",
        mcp.WithDescription("Fetch a web page and convert to markdown"),
        mcp.WithString("url", mcp.Required(), mcp.Description("URL to fetch")),
    )
    s.AddTool(pageDumpTool, handlePageDump)

    // Register perplexity-search tool
    perplexityTool := mcp.NewTool("perplexity-search",
        mcp.WithDescription("Search using Perplexity API (requires PERPLEXITY_API_KEY)"),
        mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
        mcp.WithNumber("max-results", mcp.Description("Max results (default: 5)")),
        mcp.WithString("model", mcp.Description("Model: sonar-medium-online, sonar-small-online, sonar-pro-online")),
    )
    s.AddTool(perplexityTool, handlePerplexitySearch)

    // Serve via stdio (for Claude Desktop)
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}

func handleWebSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    query, _ := request.Params.Arguments.String("query")

    cmd := exec.CommandContext(ctx, "ddg-search", query)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
    }

    return mcp.NewToolResultText(string(output)), nil
}

func handlePageDump(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    url, _ := request.Params.Arguments.String("url")

    cmd := exec.CommandContext(ctx, "page-dump", url)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Page dump failed: %v", err)), nil
    }

    return mcp.NewToolResultText(string(output)), nil
}

func handlePerplexitySearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    query, _ := request.Params.Arguments.String("query")

    cmd := exec.CommandContext(ctx, "perplexity-search", query)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Perplexity search failed: %v", err)), nil
    }

    return mcp.NewToolResultText(string(output)), nil
}
```

### Claude Desktop Configuration

```json
{
  "mcpServers": {
    "web-search": {
      "command": "/path/to/your/mcp-server-binary",
      "args": [],
      "env": {
        "PERPLEXITY_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

---

## Additional Resources

### Official Documentation
- MCP Specification: https://modelcontextprotocol.io
- MCP Go (mark3labs): http://mcp-go.dev/
- MCP Golang (metoro-io): https://mcpgolang.com
- Official Go SDK: https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk

### Community
- mark3labs/mcp-go Discord: https://discord.gg/RqSS2NQVsY
- metoro-io/mcp-golang Discord: https://discord.gg/33saRwE3pT

### Related Projects
- ktr0731/go-mcp (deprecated - use official SDK)
- google/adk-go (Agent Development Kit, not MCP-specific)

---

## Conclusion

For building a Claude Code compatible MCP server based on your existing search skills, **mark3labs/mcp-go** is the recommended choice due to:

1. ✅ Fastest implementation time
2. ✅ Largest community and most battle-tested
3. ✅ Full stdio and HTTP support
4. ✅ Simple, clean API
5. ✅ Excellent documentation
6. ✅ Active development

**Alternative choice:** `modelcontextprotocol/go-sdk` if official support and long-term stability are higher priorities than development speed.

---

*Document generated: 2026-03-01*

# Stage 1 Research: MCP Go Libraries and Examples

## MCP Go Libraries Research

### 1. github.com/mark3labs/mcp-go (Selected)

**Status:** ✅ Selected for implementation

**Description:** A Go implementation of the Model Context Protocol (MCP) that provides both server and client functionality.

**Key Features:**
- Active development and maintenance
- Supports both stdio and HTTP SSE transports
- Compatible with Claude Code's MCP implementation
- Good documentation and examples
- Minimal dependencies
- Clean API design

**Pros:**
- Well-structured codebase
- Comprehensive examples
- Active community
- Supports both server and client modes
- Good error handling

**Cons:**
- Relatively new project
- Limited third-party integrations

**Repository:** https://github.com/mark3labs/mcp-go

**Documentation:** https://pkg.go.dev/github.com/mark3labs/mcp-go

### 2. github.com/modelcontextprotocol/go-sdk

**Status:** ❌ Not selected

**Description:** Official Go SDK for the Model Context Protocol.

**Key Features:**
- Official implementation
- Supports MCP protocol specification
- Server and client support

**Pros:**
- Official implementation
- Follows specification closely

**Cons:**
- Less mature than alternatives
- Fewer examples
- Less active development
- Limited documentation

**Repository:** https://github.com/modelcontextprotocol/go-sdk

### 3. Custom Implementation

**Status:** ❌ Not selected

**Description:** Implement MCP protocol from scratch.

**Pros:**
- Full control over implementation
- No external dependencies

**Cons:**
- High complexity
- Reinventing the wheel
- Maintenance burden
- Risk of protocol incompatibility

## MCP Server Examples Research

### Example 1: mark3labs/mcp-go Examples

**Location:** https://github.com/mark3labs/mcp-go/tree/main/examples

**Available Examples:**
- `stdio-server` - Basic stdio transport server
- `http-server` - HTTP SSE transport server
- `tools-example` - Server with custom tools
- `resources-example` - Server with resources

**Key Takeaways:**
- Clean separation between transport and business logic
- Tool registration is straightforward
- Error handling is well-structured
- Configuration via environment variables

### Example 2: Simple MCP Server Pattern

Based on the examples, a typical MCP server follows this pattern:

```go
package main

import (
    "context"
    "log/slog"
    "os"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    // Create server
    s := server.NewMCPServer(
        "ddg-search-mcp",
        "1.0.0",
        server.WithToolCapabilities(true),
    )

    // Register tools
    s.AddTool(mcp.Tool{
        Name:        "search",
        Description: "Search the web",
        InputSchema: searchInputSchema(),
    }, handleSearch)

    // Start server
    if err := s.Serve(); err != nil {
        slog.Error("Server error", "error", err)
        os.Exit(1)
    }
}
```

## Recommendations

### Selected Library: github.com/mark3labs/mcp-go

**Rationale:**
1. **Active Development:** Regular updates and bug fixes
2. **Transport Support:** Both stdio and HTTP SSE out of the box
3. **Claude Code Compatibility:** Tested and verified to work with Claude Code
4. **Documentation:** Good examples and API documentation
5. **Community:** Active community for support
6. **Minimal Dependencies:** Keeps the project lightweight

### Implementation Approach

1. **Use stdio transport initially** (Stage 5)
2. **Add HTTP SSE transport later** (Stage 8)
3. **Follow the example patterns** from the repository
4. **Implement tool handlers** as separate functions
5. **Use structured logging** with log/slog

### Dependencies to Add

```go
require (
    github.com/mark3labs/mcp-go v0.0.0-latest
    github.com/spf13/viper v1.18.2
    github.com/spf13/cobra v1.8.0
)
```

## Next Steps

1. ✅ Research complete - library selected
2. Proceed to Stage 2: MCP Tool Calling Research
3. Analyze internal libraries in Stage 3
4. Begin implementation in Stage 4

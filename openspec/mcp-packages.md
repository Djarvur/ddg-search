# Go MCP Server Packages — Comparison & Analysis

> **Purpose:** Evaluate Go packages for building a Claude Code compatible MCP server  
> that wraps our existing web search tools (`ddg-search`, `perplexity-search`, `page-dump`).  
> **Requirements:** Must support both **stdio** and **HTTP** transports.  
> **Date:** 2026-03-01

---

## Packages Evaluated

| # | Package | Stars | Releases | Dependents | Latest | License |
|---|---------|-------|----------|------------|--------|---------|
| 1 | [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) | 8.3k | 70 (v0.44.1) | 2,800 | Active | MIT |
| 2 | [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) | 4.0k | 19 (v1.4.0) | 912 | Active | Apache 2.0 + MIT |
| 3 | [metoro-io/mcp-golang](https://github.com/metoro-io/mcp-golang) | 1.2k | 18 | 198 | Active | MIT |
| 4 | [ThinkInAIXYZ/go-mcp](https://github.com/ThinkInAIXYZ/go-mcp) | 666 | 41 | — | Active | MIT |

---

## Feature Matrix

| Feature | mcp-go | go-sdk (official) | mcp-golang | go-mcp |
|---------|--------|--------------------|------------|--------|
| **MCP Spec Version** | 2025-11-25 (backward compat to 2024-11-05) | 2025-11-25 (partial: no client OAuth / sampling w/ tools) | Not specified | Not specified |
| **Stdio Transport** | ✅ | ✅ | ✅ | ✅ |
| **SSE Transport** | ✅ | ❌ (superseded by streamable-HTTP) | ❌ (planned) | ✅ |
| **Streamable HTTP** | ✅ | ✅ | ❌ | ✅ |
| **HTTP (stateless)** | — | — | ✅ (plain HTTP + Gin) | — |
| **Tools** | ✅ | ✅ | ✅ | ✅ |
| **Resources** | ✅ | ✅ | ✅ | ✅ |
| **Prompts** | ✅ | ✅ | ✅ | ✅ |
| **Client SDK** | ✅ | ✅ | ✅ | ✅ |
| **Session Management** | ✅ (per-session tools, filtering) | ✅ (via transport) | ❌ | ✅ |
| **Task (async tools)** | ✅ | ✅ | ❌ | ❌ |
| **Middleware / Hooks** | ✅ (request hooks + tool middleware) | ❌ | ❌ | ✅ (global middleware) |
| **OAuth Support** | ❌ | ✅ (experimental) | ❌ | ❌ |
| **Notifications** | ✅ | ✅ | ❌ (HTTP only → no notifications) | ✅ |
| **Auto Schema from Structs** | ❌ (builder pattern) | ✅ (jsonschema tags) | ✅ (jsonschema tags) | ✅ (tags) |
| **Semver / Stable API** | ❌ (v0.x) | ✅ (v1.x) | ❌ (v0.x) | ❌ (v0.x) |
| **Maintained by** | Community (Ed Zynda) | Anthropic + Google | Metoro.io | ThinkInAI |

---

## Transport Requirements Check

Our MCP server needs **both stdio and HTTP** support:

| Transport | mcp-go | go-sdk | mcp-golang | go-mcp |
|-----------|--------|--------|------------|--------|
| Stdio ✅ | ✅ | ✅ | ✅ | ✅ |
| HTTP ✅ | ✅ streamable-HTTP | ✅ streamable-HTTP | ✅ stateless HTTP/Gin | ✅ SSE + streamable-HTTP |
| **Both** | **✅** | **✅** | **✅** (limited — no SSE/streaming) | **✅** |

---

## API Ergonomics — Side-by-Side Tool Registration

### 1. mark3labs/mcp-go

```go
tool := mcp.NewTool("web_search",
    mcp.WithDescription("Search the web"),
    mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
    mcp.WithNumber("max_results", mcp.Description("Max results")),
)

s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    query, _ := req.RequireString("query")
    // ... call searcher.SearchMarkdown(ctx, opts) ...
    return mcp.NewToolResultText(output), nil
})

server.ServeStdio(s)          // stdio
// OR
server.NewStreamableHTTPServer(s) // HTTP
```

**Pros:** Builder pattern is explicit; rich middleware/hooks system.  
**Cons:** More boilerplate per tool; no auto schema generation from structs.

### 2. modelcontextprotocol/go-sdk (Official)

```go
type SearchInput struct {
    Query      string `json:"query" jsonschema:"the search query,required"`
    MaxResults int    `json:"max_results" jsonschema:"max results to return"`
}

func WebSearch(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, error) {
    // ... call searcher.SearchMarkdown(ctx, opts) ...
    return nil, nil // result returned via second value
}

server := mcp.NewServer(&mcp.Implementation{Name: "ddg-search-mcp", Version: "v1.0.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "web_search", Description: "Search the web"}, WebSearch)

server.Run(ctx, &mcp.StdioTransport{})          // stdio
// OR
server.Run(ctx, &mcp.StreamableHTTPTransport{})  // HTTP
```

**Pros:** Typed struct args → auto JSON schema; clean API; official support; v1.x stable.  
**Cons:** Fewer community examples; no SSE (only streamable-HTTP); younger ecosystem.

### 3. metoro-io/mcp-golang

```go
type SearchArgs struct {
    Query string `json:"query" jsonschema:"required,description=Search query"`
}

server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
server.RegisterTool("web_search", "Search the web", func(args SearchArgs) (*mcp_golang.ToolResponse, error) {
    // ... call searcher ...
    return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(output)), nil
})
server.Serve()
```

**Pros:** Minimal boilerplate; struct-based schemas.  
**Cons:** HTTP transport is stateless only (no SSE/streaming); no notifications; smaller community.

### 4. ThinkInAIXYZ/go-mcp

```go
type SearchReq struct {
    Query string `json:"query" description:"search query" required:"true"`
}

tool, _ := protocol.NewTool("web_search", "Search the web", SearchReq{})
mcpServer.RegisterTool(tool, func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
    var args SearchReq
    protocol.VerifyAndUnmarshal(req.RawArguments, &args)
    // ...
    return &protocol.CallToolResult{Content: []protocol.Content{&protocol.TextContent{Text: output}}}, nil
})
```

**Pros:** Full transport support (SSE, streamable-HTTP, stdio); Gin integration; middleware.  
**Cons:** Smaller community; Chinese documentation dominant; more verbose result construction.

---

## Analysis for Our Use Case

### What we need to expose as MCP tools:

| Tool | Source | Arguments |
|------|--------|-----------|
| `web_search` | `internal/search` | query, max_results, site, region, time, safe_search |
| `perplexity_search` | `internal/perplexity` | query, max_results, model |
| `page_dump` | `internal/dump` | url, timeout, user_agent |

### Key Decision Factors:

1. **Speed to implement** — We have 3 simple tools, no complex state management needed
2. **Both transports** — stdio (for Claude Code CLI) + HTTP (for remote use)
3. **Future-proofing** — Which will be the standard in 6–12 months?
4. **API stability** — We don't want to chase breaking changes

---

## Recommendation

### 🏆 Top Pick: `modelcontextprotocol/go-sdk` (Official SDK)

| Factor | Score | Reasoning |
|--------|-------|-----------|
| Speed to implement | ⭐⭐⭐⭐⭐ | Struct-based tool args → minimal boilerplate for our 3 tools |
| Transport support | ⭐⭐⭐⭐ | Stdio ✅ + Streamable HTTP ✅ (no SSE but SSE is deprecated) |
| API stability | ⭐⭐⭐⭐⭐ | v1.x with semver; official so unlikely to be abandoned |
| Future-proofing | ⭐⭐⭐⭐⭐ | The official SDK = the reference implementation going forward |
| Community | ⭐⭐⭐ | Growing fast (4k ★, 912 dependents), backed by Anthropic + Google |
| Claude Code compat | ⭐⭐⭐⭐⭐ | Official SDK → guaranteed protocol compliance |

### 🥈 Strong Alternative: `mark3labs/mcp-go`

| Factor | Score | Reasoning |
|--------|-------|-----------|
| Speed to implement | ⭐⭐⭐⭐ | More boilerplate (builder pattern) but excellent docs/examples |
| Transport support | ⭐⭐⭐⭐⭐ | Stdio + SSE + Streamable HTTP |
| API stability | ⭐⭐⭐ | v0.x — breaking changes still happen |
| Future-proofing | ⭐⭐⭐⭐ | Huge community; acknowledged by official SDK |
| Community | ⭐⭐⭐⭐⭐ | 8.3k ★, 2.8k dependents, 170 contributors |
| Claude Code compat | ⭐⭐⭐⭐ | Battle-tested with Claude Desktop |

### Not Recommended for This Project

- **mcp-golang** — HTTP transport is stateless only, no SSE/streaming, smallest community of the four
- **go-mcp** — Good tech but very small community, Chinese-focused docs, higher implementation overhead

---

## Decision Matrix Summary

```
                    Speed   Transport   Stability   Future   Community   TOTAL
go-sdk (official)    5         4           5          5         3         22/25  ← RECOMMENDED
mcp-go (mark3labs)   4         5           3          4         5         21/25
go-mcp (ThinkInAI)   3         5           2          3         2         15/25
mcp-golang (metoro)  4         2           2          2         2         12/25
```

---

## Questions for Decision

1. **Which package do you want to use?**
   - `modelcontextprotocol/go-sdk` (official, fastest implementation, v1.x stable) — **recommended**
   - `mark3labs/mcp-go` (largest community, most features, v0.x)

2. **HTTP transport flavor:**
   - Streamable HTTP (modern, spec ≥2025-03-26) — default for both top picks
   - SSE (legacy, only available with `mcp-go`)
   - Do you need backward-compatible SSE support?

3. **Scope of tools to expose:**
   - All 3 (`web_search`, `perplexity_search`, `page_dump`)?
   - Or start with `web_search` only?

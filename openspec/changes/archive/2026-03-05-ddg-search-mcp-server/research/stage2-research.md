# Stage 2 Research: Claude Code MCP Tool Calling Conventions

## MCP Tool Calling Overview

The Model Context Protocol (MCP) defines a standard way for AI assistants to call tools provided by servers. Claude Code uses MCP to communicate with tool servers.

## MCP Tool Call Flow

### 1. Initialization
```
Client → Server: initialize request
Server → Client: initialize response (with tools list)
```

### 2. Tool Call
```
Client → Server: tools/call request
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search",
    "arguments": {
      "query": "golang mcp server"
    }
  }
}
```

### 3. Tool Response
```
Server → Client: tools/call response
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Search results..."
      }
    ]
  }
}
```

## Web Search Tool Specification

### Tool Name
`search`

### Description
Search the web for information using DuckDuckGo or Perplexity.

### Input Schema (JSON Schema)
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "The search query string"
    },
    "max_results": {
      "type": "integer",
      "description": "Maximum number of results to return",
      "default": 10,
      "minimum": 1,
      "maximum": 50
    },
    "safe_search": {
      "type": "boolean",
      "description": "Enable safe search filtering",
      "default": false
    }
  },
  "required": ["query"]
}
```

### Request Format
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search",
    "arguments": {
      "query": "golang mcp server",
      "max_results": 5,
      "safe_search": true
    }
  }
}
```

### Response Format (Success)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "## Search Results\n\n### 1. Title\n**URL:** https://example.com/article\n**Snippet:** This is a brief description of the search result.\n\n### 2. Title\n**URL:** https://example.com/another\n**Snippet:** Another search result description.\n\n..."
      }
    ]
  }
}
```

### Response Format (Error)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "Query parameter is required"
    }
  }
}
```

### Response Format (Fallback)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "## Search Results\n\n**Note:** Results from DuckDuckGo (Perplexity unavailable due to rate limit).\n\n### 1. Title\n**URL:** https://example.com/article\n**Snippet:** This is a brief description of the search result.\n\n..."
      }
    ]
  }
}
```

## Fetch Tool Specification

### Tool Name
`fetch`

### Description
Fetch and convert a web page to markdown format.

### Input Schema (JSON Schema)
```json
{
  "type": "object",
  "properties": {
    "url": {
      "type": "string",
      "description": "The URL to fetch",
      "format": "uri"
    }
  },
  "required": ["url"]
}
```

### Request Format
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "fetch",
    "arguments": {
      "url": "https://example.com/article"
    }
  }
}
```

### Response Format (Success)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "# Article Title\n\nThis is the markdown content of the fetched page.\n\n## Section\n\nContent here...\n\n**Source:** https://example.com/article"
      }
    ]
  }
}
```

### Response Format (Error)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "details": "URL parameter is required and must be a valid HTTP/HTTPS URL"
    }
  }
}
```

### Response Format (Network Error)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "error": {
    "code": -32603,
    "message": "Internal error",
    "data": {
      "details": "Failed to fetch URL: connection timeout"
    }
  }
}
```

## MCP Error Codes

| Code | Name | Description |
|------|------|-------------|
| -32700 | Parse error | Invalid JSON was received |
| -32600 | Invalid Request | The JSON sent is not a valid Request object |
| -32601 | Method not found | The method does not exist |
| -32602 | Invalid params | Invalid method parameters |
| -32603 | Internal error | Internal JSON-RPC error |

## Tool Registration Format

When the server initializes, it advertises available tools:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": {}
    },
    "serverInfo": {
      "name": "ddg-search-mcp",
      "version": "1.0.0"
    }
  }
}
```

Tools are listed via the `tools/list` method:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}
```

Response:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "name": "search",
        "description": "Search the web for information",
        "inputSchema": {
          "type": "object",
          "properties": {
            "query": {
              "type": "string",
              "description": "The search query string"
            },
            "max_results": {
              "type": "integer",
              "description": "Maximum number of results to return",
              "default": 10
            },
            "safe_search": {
              "type": "boolean",
              "description": "Enable safe search filtering",
              "default": false
            }
          },
          "required": ["query"]
        }
      },
      {
        "name": "fetch",
        "description": "Fetch and convert a web page to markdown",
        "inputSchema": {
          "type": "object",
          "properties": {
            "url": {
              "type": "string",
              "description": "The URL to fetch",
              "format": "uri"
            }
          },
          "required": ["url"]
        }
      }
    ]
  }
}
```

## Content Types

MCP supports multiple content types in tool responses:

### Text Content
```json
{
  "type": "text",
  "text": "Plain text content"
}
```

### Image Content (not used in this project)
```json
{
  "type": "image",
  "data": "base64-encoded-image-data",
  "mimeType": "image/png"
}
```

### Resource Content (not used in this project)
```json
{
  "type": "resource",
  "uri": "file:///path/to/resource"
}
```

## Implementation Notes

1. **JSON-RPC 2.0**: MCP uses JSON-RPC 2.0 for all communication
2. **Content Array**: Tool responses must have a `content` array with at least one item
3. **Error Handling**: Always return proper error codes and messages
4. **Validation**: Validate all input parameters before processing
5. **Fallback Notes**: When falling back from Perplexity to DuckDuckGo, include a note in the response
6. **Source Attribution**: For fetch tool, include the source URL in the response

## Next Steps

1. ✅ Research complete - tool calling conventions documented
2. Proceed to Stage 3: Internal Library Analysis
3. Begin implementation in Stage 4

package mcp_test

import (
	"context"
	"testing"

	"github.com/Djarvur/ddg-search/internal/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// TestSearchToolDefinition verifies the search tool definition is correct.
func TestSearchToolDefinition(t *testing.T) {
	t.Parallel()

	tool := mcp.SearchTool()
	require.Equal(t, "search", tool.Name)
	require.Equal(t, "Search the web for information using DuckDuckGo or Perplexity", tool.Description)
	require.Equal(t, "object", tool.InputSchema.Type)

	props := tool.InputSchema.Properties// Check query property
	queryProp, ok := props["query"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "string", queryProp["type"])
	require.Equal(t, "The search query string", queryProp["description"])

	// Check max_results property
	maxResultsProp, ok := props["max_results"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "number", maxResultsProp["type"])
	require.Equal(t, "Maximum number of results to return (default: 10)", maxResultsProp["description"])

	// Check safe_search property
	safeSearchProp, ok := props["safe_search"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "boolean", safeSearchProp["type"])
	require.Equal(t, "Enable safe search filtering (default: false)", safeSearchProp["description"])

	// Check required fields
	require.Contains(t, tool.InputSchema.Required, "query")
}

// TestHandleSearch_MissingQuery verifies that missing query returns an error.
func TestHandleSearch_MissingQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "query parameter is required")
}

// TestHandleSearch_EmptyQuery verifies that empty query returns an error.
func TestHandleSearch_EmptyQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"query": "",
			},
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "query parameter is required")
}

// TestHandleSearch_InvalidQueryType verifies that non-string query returns an error.
func TestHandleSearch_InvalidQueryType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"query": 123,
			},
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "query parameter is required")
}

// TestHandleSearch_ValidQuery verifies that a valid query is processed.
// Note: This is an integration test that requires network access.
func TestHandleSearch_ValidQuery(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"query":       "golang",
				"max_results": 5.0,
			},
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.NotEmpty(t, textContent.Text)
}

// TestHandleSearch_WithMaxResults verifies max_results parameter is handled.
func TestHandleSearch_WithMaxResults(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"query":       "golang",
				"max_results": 3.0,
			},
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.NotEmpty(t, textContent.Text)
}

// TestHandleSearch_WithSafeSearch verifies safe_search parameter is handled.
func TestHandleSearch_WithSafeSearch(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"query":       "golang",
				"safe_search": true,
			},
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.NotEmpty(t, textContent.Text)
}

// TestHandleSearch_NilArguments verifies nil arguments are handled.
func TestHandleSearch_NilArguments(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: nil,
		},
	}

	result, err := mcp.HandleSearch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "query parameter is required")
}

// TestFetchToolDefinition verifies the fetch tool definition is correct.
func TestFetchToolDefinition(t *testing.T) {
	t.Parallel()

	tool := mcp.FetchTool()
	require.Equal(t, "fetch", tool.Name)
	require.Equal(t, "Fetch and convert a web page to markdown format", tool.Description)
	require.Equal(t, "object", tool.InputSchema.Type)

	props := tool.InputSchema.Properties// Check url property
	urlProp, ok := props["url"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "string", urlProp["type"])
	require.Equal(t, "The URL of the web page to fetch", urlProp["description"])

	// Check required fields
	require.Contains(t, tool.InputSchema.Required, "url")
}

// TestHandleFetch_MissingURL verifies that missing URL returns an error.
func TestHandleFetch_MissingURL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{},
		},
	}

	result, err := mcp.HandleFetch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "url parameter is required")
}

// TestHandleFetch_EmptyURL verifies that empty URL returns an error.
func TestHandleFetch_EmptyURL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"url": "",
			},
		},
	}

	result, err := mcp.HandleFetch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "url parameter is required")
}

// TestHandleFetch_InvalidURLType verifies that non-string URL returns an error.
func TestHandleFetch_InvalidURLType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"url": 123,
			},
		},
	}

	result, err := mcp.HandleFetch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "url parameter is required")
}

// TestHandleFetch_ValidURL verifies that a valid URL is processed.
// Note: This is an integration test that requires network access.
func TestHandleFetch_ValidURL(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"url": "https://example.com",
			},
		},
	}

	result, err := mcp.HandleFetch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.NotEmpty(t, textContent.Text)
}

// TestHandleFetch_InvalidURLScheme verifies that invalid URL scheme returns an error.
func TestHandleFetch_InvalidURLScheme(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: map[string]any{
				"url": "ftp://example.com",
			},
		},
	}

	result, err := mcp.HandleFetch(ctx, nil, request)
	// HandleFetch returns both result and error
	require.Error(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "fetch failed")
}

// TestHandleFetch_NilArguments verifies nil arguments are handled.
func TestHandleFetch_NilArguments(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request := mcpgo.CallToolRequest{
		Params: mcpgo.CallToolParams{
			Arguments: nil,
		},
	}

	result, err := mcp.HandleFetch(ctx, nil, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.IsError)
	textContent, ok := mcpgo.AsTextContent(result.Content[0])
	require.True(t, ok)
	require.Contains(t, textContent.Text, "url parameter is required")
}

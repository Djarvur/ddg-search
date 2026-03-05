package mcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Djarvur/ddg-search/internal/config"
	"github.com/Djarvur/ddg-search/internal/dump"
	"github.com/Djarvur/ddg-search/internal/mcpconfig"
	"github.com/Djarvur/ddg-search/internal/perplexity"
	"github.com/Djarvur/ddg-search/internal/search"
	"github.com/mark3labs/mcp-go/mcp"
)

// SearchTool returns the search tool definition.
func SearchTool() mcp.Tool {
	return mcp.Tool{
		Name:        "search",
		Description: "Search the web for information using DuckDuckGo or Perplexity",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "The search query string",
				},
				"max_results": map[string]any{
					"type":        "number",
					"description": "Maximum number of results to return (default: 10)",
				},
				"safe_search": map[string]any{
					"type":        "boolean",
					"description": "Enable safe search filtering (default: false)",
				},
			},
			Required: []string{"query"},
		},
	}
}

// HandleSearch handles search tool calls with real DuckDuckGo or Perplexity search.
//nolint:gocognit
func HandleSearch(ctx context.Context, appConfig any, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments as map
	var args map[string]any

	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]any); ok {
			args = argsMap
		}
	}

	// Extract query parameter
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query parameter is required and must be a non-empty string"), nil
	}

	// Extract optional parameters
	maxResults := 10
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	safeSearch := false
	if ss, ok := args["safe_search"].(bool); ok {
		safeSearch = ss
	}

	// Try Perplexity search if enabled
	var (
		results      string
		fallbackNote string
	)

	// Type assert appConfig to get Perplexity configuration
	cfg, ok := appConfig.(*mcpconfig.Config)
	if ok && cfg.Perplexity.Enabled && cfg.Perplexity.AccessToken != "" {
		// Use Perplexity search
		perplexityClient := perplexity.NewClient(cfg.Perplexity.AccessToken)

		perplexityResults, err := perplexityClient.Search(ctx, query, "sonar")
		if err == nil {
			// Perplexity search succeeded, format results as markdown
			results = formatPerplexityResults(perplexityResults)
			results = "*Search powered by Perplexity AI*\n\n" + results

			return mcp.NewToolResultText(results), nil
		}

		// Perplexity failed, check if we should fall back
		if isPerplexityError(err) {
			fallbackNote = "*Note: Perplexity search failed, falling back to DuckDuckGo.*"
		} else {
			return mcp.NewToolResultError("Perplexity search failed: " + err.Error()), nil
		}
	}

	// Fall back to DuckDuckGo search
	retryOpts := config.DefaultRetryOptions()

	searcher := search.NewSearcher(retryOpts)
	defer searcher.Close()

	// Perform search
	searchOpts := config.SearchOptions{
		Query:      query,
		MaxResults: maxResults,
		SafeSearch: safeSearch,
	}

	ddgResults, err := searcher.SearchMarkdown(ctx, searchOpts)
	if err != nil {
		return mcp.NewToolResultError("search failed: " + err.Error()), nil
	}

	results = ddgResults
	if fallbackNote != "" {
		results = "*Search powered by DuckDuckGo*\n\n" + fallbackNote + "\n\n" + ddgResults
	} else {
		results = "*Search powered by DuckDuckGo*\n\n" + results
	}

	return mcp.NewToolResultText(results), nil
}

// formatPerplexityResults formats Perplexity search results as markdown.
func formatPerplexityResults(results *perplexity.SearchResults) string {
	var md string

	md += "# " + results.Query + "\n\n"
	md += results.Answer + "\n\n"

	if len(results.References) > 0 {
		md += "## References\n\n"

		var mdSb128 strings.Builder
		for _, ref := range results.References {
			mdSb128.WriteString(fmt.Sprintf("%d. %s\n", ref.Index, ref.URL))
		}

		md += mdSb128.String()
	}

	return md
}

// isPerplexityError checks if an error indicates we should fall back to DuckDuckGo.
func isPerplexityError(err error) bool {
	if err == nil {
		return false
	}

	// Check for specific Perplexity errors that should trigger fallback
	return errors.Is(err, perplexity.ErrUnauthorized) ||
		errors.Is(err, perplexity.ErrPaymentRequired) ||
		errors.Is(err, perplexity.ErrBadRequest) ||
		errors.Is(err, perplexity.ErrServer) ||
		errors.Is(err, perplexity.ErrNetwork)
}

// FetchTool returns the fetch tool definition.
func FetchTool() mcp.Tool {
	return mcp.Tool{
		Name:        "fetch",
		Description: "Fetch and convert a web page to markdown format",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"url": map[string]any{
					"type":        "string",
					"description": "The URL of the web page to fetch",
				},
			},
			Required: []string{"url"},
		},
	}
}

// HandleFetch handles fetch tool calls with real web page fetching.
func HandleFetch(ctx context.Context, _ any, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments as map
	var args map[string]any

	if request.Params.Arguments != nil {
		if argsMap, ok := request.Params.Arguments.(map[string]any); ok {
			args = argsMap
		}
	}

	// Extract URL parameter
	url, ok := args["url"].(string)
	if !ok || url == "" {
		return mcp.NewToolResultError("url parameter is required and must be a non-empty string"), nil
	}

	// Fetch and convert the URL
	cfg := dump.DefaultConfig()

	markdown, err := dump.FetchAndConvert(ctx, url, cfg)
	if err != nil {
		return mcp.NewToolResultError("fetch failed: " + err.Error()), err
	}

	return mcp.NewToolResultText(markdown), nil
}

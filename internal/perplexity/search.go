package perplexity

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// SearchOptions configures a Perplexity search query.
type SearchOptions struct {
	// Query is the search string.
	Query string
	// Model specifies the Perplexity model to use.
	Model string
}

// Search performs a web search using the Perplexity API.
func (c *Client) Search(ctx context.Context, query string, model string) (*SearchResults, error) {
	if query == "" {
		return nil, ErrQueryEmpty
	}

	// Build request body for Perplexity Chat Completions API
	reqBody := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": query,
			},
		},
	}

	// Make the request
	resp, err := c.Do(ctx, c.httpClient.R().
		SetBody(reqBody).
		SetContext(ctx).
		SetHeader("Content-Type", "application/json"), "POST")
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	// Parse response
	var apiResponse APIResponse

	err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check for API errors in response
	if apiResponse.Error != nil {
		return nil, fmt.Errorf("%w: %s", ErrAPI, apiResponse.Error.Message)
	}

	// Check if we have choices
	if len(apiResponse.Choices) == 0 {
		return nil, ErrNoChoices
	}

	// Extract answer and citations from the response
	answer := apiResponse.Choices[0].Message.Content
	citations := apiResponse.Citations

	// Convert to SearchResults
	results := &SearchResults{
		Query:      query,
		Answer:     answer,
		Citations:  citations,
		References: make([]Reference, len(citations)),
	}

	// Create references from citations
	for i, citation := range citations {
		results.References[i] = Reference{
			Index: i + 1,
			URL:   citation,
		}
	}

	return results, nil
}

// APIResponse represents the raw Perplexity Chat Completions API response.
type APIResponse struct {
	ID        string   `json:"id"`
	Citations []string `json:"citations"`
	Choices   []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// SearchResults represents the results of a Perplexity search.
type SearchResults struct {
	// Query is the original search query.
	Query string
	// Answer is the AI-generated answer from Perplexity.
	Answer string
	// Citations are the URLs cited in the answer.
	Citations []string
	// References are formatted citation references.
	References []Reference
}

// Reference represents a cited source.
type Reference struct {
	// Index is the citation number.
	Index int
	// URL is the citation URL.
	URL string
}

// Markdown returns the search results formatted as markdown.
func (r *SearchResults) Markdown() string {
	var sb strings.Builder

	// Write answer
	sb.WriteString(r.Answer)
	sb.WriteString("\n\n")

	// Write citations
	if len(r.Citations) > 0 {
		sb.WriteString("## Sources\n\n")

		for _, ref := range r.References {
			sb.WriteString(fmt.Sprintf("%d. %s\n", ref.Index, ref.URL))
		}
	}

	return sb.String()
}

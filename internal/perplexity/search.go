package perplexity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Search performs a web search using the Perplexity API.
//
// maxResults caps the number of sources reported alongside the answer; a
// non-positive value keeps every source the API returned.
func (c *Client) Search(ctx context.Context, query string, maxResults int, model string) (*SearchResults, error) {
	if query == "" {
		return nil, ErrQueryEmpty
	}

	// Build request body
	reqBody := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": query,
			},
		},
		"max_tokens":  maxTokens,
		"temperature": temperature,
		"stream":      false,
	}

	// Make the request
	resp, err := c.Do(ctx, http.MethodPost, chatCompletionsPath, c.httpClient.R().
		SetBody(reqBody).
		SetHeader("Content-Type", "application/json"))
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

	if len(apiResponse.Choices) == 0 {
		return nil, ErrNoChoices
	}

	return newSearchResults(query, &apiResponse, maxResults), nil
}

// APIResponse represents the raw Perplexity API response.
type APIResponse struct {
	// Choices carries the generated answer; the text lives at
	// Choices[0].Message.Content.
	Choices []Choice `json:"choices"`
	// SearchResults describes the sources behind the answer. It supersedes
	// Citations, which Perplexity has deprecated.
	SearchResults []SearchResult `json:"search_results"`
	// Citations is the deprecated, URL-only form of SearchResults.
	Citations []string `json:"citations"`
	Error     *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// Choice is a single completion returned by the API.
type Choice struct {
	// Message holds the assistant's reply.
	Message struct {
		// Role is the author of the message, normally "assistant".
		Role string `json:"role"`
		// Content is the generated answer text.
		Content string `json:"content"`
	} `json:"message"`
	// FinishReason explains why generation stopped.
	FinishReason string `json:"finish_reason"`
}

// SearchResult is a single source the answer was grounded in.
type SearchResult struct {
	// Title is the source's page title.
	Title string `json:"title"`
	// URL is the source's address.
	URL string `json:"url"`
	// Date is the source's publication date, if the API reported one.
	Date string `json:"date"`
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
	// Title is the source's page title, empty when the API did not report one.
	Title string
}

// newSearchResults converts an API response into SearchResults, keeping at most
// maxResults sources. A non-positive maxResults keeps every source.
func newSearchResults(query string, resp *APIResponse, maxResults int) *SearchResults {
	refs := resp.references()
	if maxResults > 0 && len(refs) > maxResults {
		refs = refs[:maxResults]
	}

	citations := make([]string, len(refs))

	for i := range refs {
		refs[i].Index = i + 1
		citations[i] = refs[i].URL
	}

	return &SearchResults{
		Query:      query,
		Answer:     resp.Choices[0].Message.Content,
		Citations:  citations,
		References: refs,
	}
}

// references builds the source list, preferring search_results over the
// deprecated citations field.
func (r *APIResponse) references() []Reference {
	if len(r.SearchResults) > 0 {
		refs := make([]Reference, len(r.SearchResults))
		for i, sr := range r.SearchResults {
			refs[i] = Reference{URL: sr.URL, Title: sr.Title}
		}

		return refs
	}

	refs := make([]Reference, len(r.Citations))
	for i, citation := range r.Citations {
		refs[i] = Reference{URL: citation}
	}

	return refs
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
			if ref.Title == "" {
				fmt.Fprintf(&sb, "%d. %s\n", ref.Index, ref.URL)

				continue
			}

			fmt.Fprintf(&sb, "%d. [%s](%s)\n", ref.Index, ref.Title, ref.URL)
		}
	}

	return sb.String()
}

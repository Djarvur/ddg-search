package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
)

// Searcher orchestrates the search process.
type Searcher struct {
	client      *Client
	parser      *Parser
	retryOpts   config.RetryOptions
	debugWriter io.Writer
}

// NewSearcher creates a new searcher with the given retry options.
func NewSearcher(retryOptions config.RetryOptions) *Searcher {
	debugWriter := io.Discard
	if retryOptions.Debug {
		debugWriter = os.Stderr
	}

	return &Searcher{
		client:      NewClient(retryOptions),
		parser:      NewParser(),
		retryOpts:   retryOptions,
		debugWriter: debugWriter,
	}
}

// Search performs a DuckDuckGo search with the given options.
func (s *Searcher) Search(ctx context.Context, opts config.SearchOptions) ([]config.Result, error) {
	if strings.TrimSpace(opts.Query) == "" {
		return []config.Result{}, nil
	}

	params := s.buildSearchParams(opts)
	s.debugf("search: query=%q, max_results=%d", params.Get("q"), opts.MaxResults)

	var lastErr error

	for attempt := 0; attempt <= s.retryOpts.MaxRetries; attempt++ {
		results, retryable, err := s.searchAttempt(ctx, params, opts.MaxResults, attempt)
		if !retryable {
			return results, err
		}

		lastErr = err

		if attempt < s.retryOpts.MaxRetries {
			delay := s.calculateDelay(attempt)
			s.debugf("search: attempt %d: waiting %v before retry", attempt+1, delay)

			select {
			case <-ctx.Done():
				s.debugf("search: context cancelled during retry delay")

				return nil, fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
			}
		}
	}

	s.debugf("search: all %d attempts exhausted", s.retryOpts.MaxRetries+1)

	if lastErr != nil {
		return nil, fmt.Errorf("search failed after max retries: %w", errors.Join(ErrMaxRetries, lastErr))
	}

	return nil, ErrMaxRetries
}

// SearchJSON performs a search and returns JSON-formatted results.
func (s *Searcher) SearchJSON(ctx context.Context, opts config.SearchOptions) (string, error) {
	results, err := s.Search(ctx, opts)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(data), nil
}

// SearchMarkdown performs a search and returns Markdown-formatted results.
func (s *Searcher) SearchMarkdown(ctx context.Context, opts config.SearchOptions) (string, error) {
	results, err := s.Search(ctx, opts)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "No results found.", nil
	}

	var sb strings.Builder
	for i, r := range results {
		fmt.Fprintf(&sb, "%d. [%s](%s)\n", i+1, r.Title, r.URL)

		if r.Snippet != "" {
			fmt.Fprintf(&sb, "   %s\n", r.Snippet)
		}
	}

	return sb.String(), nil
}

// Close releases resources.
func (s *Searcher) Close() {
	s.client.Close()
}

// buildSearchParams constructs URL parameters from search options.
func (s *Searcher) buildSearchParams(opts config.SearchOptions) url.Values {
	query := opts.Query
	if opts.Site != "" {
		query = fmt.Sprintf("site:%s %s", opts.Site, query)
	}

	params := url.Values{}
	params.Set("q", query)

	if opts.Region != "" {
		params.Set("kl", opts.Region)
	}

	if opts.TimeFilter != "" {
		params.Set("df", opts.TimeFilter)
	}

	if opts.SafeSearch {
		params.Set("p", "1")
	} else {
		params.Set("p", "-1")
	}

	return params
}

// searchAttempt executes a single search attempt.
// Returns (results, retryable, error). If retryable is true, the caller should retry.
func (s *Searcher) searchAttempt(
	ctx context.Context, params url.Values, maxResults, attempt int,
) ([]config.Result, bool, error) {
	s.debugf("search: attempt %d: executing search request", attempt+1)

	req := s.client.httpClient.R().
		SetContext(ctx).
		SetQueryParamsFromValues(params)

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		if errors.Is(err, ErrMaxRetries) {
			s.debugf("search: attempt %d: HTTP client exhausted retries", attempt+1)

			return nil, true, err
		}

		return nil, false, fmt.Errorf("search request failed: %w", err)
	}

	body := resp.String()

	// Check for rate limit page in HTML content
	if indicator := s.parser.FindRateLimitIndicator(body); indicator != "" {
		s.debugf("search: attempt %d: HTML rate limit detected - found indicator: %q", attempt+1, indicator)

		return nil, true, ErrRateLimited
	}

	// Parse results
	results, err := s.parser.Parse(body, maxResults)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse results: %w", err)
	}

	s.debugf("search: attempt %d: success, got %d results", attempt+1, len(results))

	return results, false, nil
}

// calculateDelay computes the retry delay with exponential backoff.
func (s *Searcher) calculateDelay(attempt int) time.Duration {
	delay := float64(s.retryOpts.BaseDelay)
	for range attempt {
		delay *= s.retryOpts.BackoffMultiplier
	}

	if delay > float64(s.retryOpts.MaxDelay) {
		delay = float64(s.retryOpts.MaxDelay)
	}

	return time.Duration(delay)
}

// debugf logs a debug message if debug mode is enabled.
func (s *Searcher) debugf(format string, args ...any) {
	if s.debugWriter != io.Discard {
		_, _ = fmt.Fprintf(s.debugWriter, "[DEBUG] "+format+"\n", args...)
	}
}

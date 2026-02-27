package search

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
)

func TestSearchIntegrationBasic(t *testing.T) {
	t.Parallel()

	html, err := os.ReadFile("testdata/results_basic.html")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(html)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	parser := NewParser()
	searcher := &Searcher{
		client:      client,
		parser:      parser,
		retryOpts:   opts,
		debugWriter: io.Discard,
	}

	searchOpts := config.SearchOptions{
		Query:      "test query",
		MaxResults: 10,
	}

	results, err := searcher.Search(context.Background(), searchOpts)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Search() returned %d results, want 3", len(results))
	}

	// Verify first result
	if results[0].Title != "Example Result 1" {
		t.Errorf("results[0].Title = %q, want %q", results[0].Title, "Example Result 1")
	}

	if results[0].URL != "https://example.com/page1" {
		t.Errorf("results[0].URL = %q, want %q", results[0].URL, "https://example.com/page1")
	}
}

func TestSearchIntegrationEmptyResults(t *testing.T) {
	t.Parallel()

	html, err := os.ReadFile("testdata/results_empty.html")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(html)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	searcher := &Searcher{
		client:      client,
		parser:      NewParser(),
		retryOpts:   opts,
		debugWriter: io.Discard,
	}

	searchOpts := config.SearchOptions{
		Query:      "nonexistent query",
		MaxResults: 10,
	}

	results, err := searcher.Search(context.Background(), searchOpts)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Search() returned %d results, want 0", len(results))
	}
}

func TestSearchIntegrationMaxResults(t *testing.T) {
	t.Parallel()

	html, err := os.ReadFile("testdata/results_basic.html")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(html)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	searcher := &Searcher{
		client:      client,
		parser:      NewParser(),
		retryOpts:   opts,
		debugWriter: io.Discard,
	}

	searchOpts := config.SearchOptions{
		Query:      "test query",
		MaxResults: 2, // Request only 2 results
	}

	results, err := searcher.Search(context.Background(), searchOpts)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Search() returned %d results, want 2", len(results))
	}
}

func TestSearchIntegrationRetryOn500(t *testing.T) {
	t.Parallel()

	html, err := os.ReadFile("testdata/results_basic.html")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 2 {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(html)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	searcher := &Searcher{
		client:      client,
		parser:      NewParser(),
		retryOpts:   opts,
		debugWriter: io.Discard,
	}

	searchOpts := config.SearchOptions{
		Query:      "test query",
		MaxResults: 10,
	}

	results, err := searcher.Search(context.Background(), searchOpts)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Search() returned %d results, want 3", len(results))
	}

	if atomic.LoadInt32(&requestCount) != 2 {
		t.Errorf("expected 2 requests, got %d", requestCount)
	}
}

func TestSearchMarkdown(t *testing.T) {
	t.Parallel()

	html, err := os.ReadFile("testdata/results_basic.html")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(html)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	searcher := &Searcher{
		client:      client,
		parser:      NewParser(),
		retryOpts:   opts,
		debugWriter: io.Discard,
	}

	searchOpts := config.SearchOptions{
		Query:      "test query",
		MaxResults: 10,
	}

	output, err := searcher.SearchMarkdown(context.Background(), searchOpts)
	if err != nil {
		t.Fatalf("SearchMarkdown() error = %v", err)
	}

	// Verify markdown format
	if len(output) == 0 {
		t.Error("SearchMarkdown() returned empty output")
	}

	// Should contain numbered list
	if output[0] != '1' {
		t.Errorf("SearchMarkdown() output should start with '1', got %q", output[0])
	}
}

func TestSearchJSON(t *testing.T) {
	t.Parallel()

	html, err := os.ReadFile("testdata/results_basic.html")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(html)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	searcher := &Searcher{
		client:      client,
		parser:      NewParser(),
		retryOpts:   opts,
		debugWriter: io.Discard,
	}

	searchOpts := config.SearchOptions{
		Query:      "test query",
		MaxResults: 10,
	}

	output, err := searcher.SearchJSON(context.Background(), searchOpts)
	if err != nil {
		t.Fatalf("SearchJSON() error = %v", err)
	}

	// Verify JSON format
	if len(output) == 0 {
		t.Error("SearchJSON() returned empty output")
	}

	if output[0] != '[' {
		t.Errorf("SearchJSON() output should start with '[', got %q", output[0])
	}
}

package config_test

import (
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
)

func TestDefaultRetryOptions(t *testing.T) {
	t.Parallel()

	opts := config.DefaultRetryOptions()

	if opts.MaxRetries != config.DefaultMaxRetries {
		t.Errorf("Expected MaxRetries %d, got %d", config.DefaultMaxRetries, opts.MaxRetries)
	}

	if opts.BaseDelay != 1*time.Second {
		t.Errorf("Expected BaseDelay 1s, got %v", opts.BaseDelay)
	}

	if opts.MaxDelay != config.DefaultMaxDelaySecs*time.Second {
		t.Errorf("Expected MaxDelay %ds, got %v", config.DefaultMaxDelaySecs, opts.MaxDelay)
	}

	if opts.BackoffMultiplier != config.DefaultBackoffMult {
		t.Errorf("Expected BackoffMultiplier %f, got %f", config.DefaultBackoffMult, opts.BackoffMultiplier)
	}

	if opts.Debug != false {
		t.Errorf("Expected Debug false, got %v", opts.Debug)
	}
}

func TestSearchOptions(t *testing.T) {
	t.Parallel()

	opts := config.SearchOptions{
		Query:      "test query",
		MaxResults: 10,
		Site:       "example.com",
		Region:     "us-en",
		TimeFilter: "d",
		SafeSearch: true,
	}

	if opts.Query != "test query" {
		t.Errorf("Expected Query 'test query', got %q", opts.Query)
	}

	if opts.MaxResults != 10 {
		t.Errorf("Expected MaxResults 10, got %d", opts.MaxResults)
	}

	if opts.Site != "example.com" {
		t.Errorf("Expected Site 'example.com', got %q", opts.Site)
	}

	if opts.Region != "us-en" {
		t.Errorf("Expected Region 'us-en', got %q", opts.Region)
	}

	if opts.TimeFilter != "d" {
		t.Errorf("Expected TimeFilter 'd', got %q", opts.TimeFilter)
	}

	if opts.SafeSearch != true {
		t.Errorf("Expected SafeSearch true, got %v", opts.SafeSearch)
	}
}

func TestSearchOptions_ZeroValues(t *testing.T) {
	t.Parallel()

	opts := config.SearchOptions{}

	if opts.Query != "" {
		t.Errorf("Expected Query empty, got %q", opts.Query)
	}

	if opts.MaxResults != 0 {
		t.Errorf("Expected MaxResults 0, got %d", opts.MaxResults)
	}

	if opts.Site != "" {
		t.Errorf("Expected Site empty, got %q", opts.Site)
	}

	if opts.Region != "" {
		t.Errorf("Expected Region empty, got %q", opts.Region)
	}

	if opts.TimeFilter != "" {
		t.Errorf("Expected TimeFilter empty, got %q", opts.TimeFilter)
	}

	if opts.SafeSearch != false {
		t.Errorf("Expected SafeSearch false, got %v", opts.SafeSearch)
	}
}

func TestRetryOptions(t *testing.T) {
	t.Parallel()

	opts := config.RetryOptions{
		MaxRetries:        5,
		BaseDelay:         2 * time.Second,
		MaxDelay:          60 * time.Second,
		BackoffMultiplier: 3.0,
		Debug:             true,
	}

	if opts.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries 5, got %d", opts.MaxRetries)
	}

	if opts.BaseDelay != 2*time.Second {
		t.Errorf("Expected BaseDelay 2s, got %v", opts.BaseDelay)
	}

	if opts.MaxDelay != 60*time.Second {
		t.Errorf("Expected MaxDelay 60s, got %v", opts.MaxDelay)
	}

	if opts.BackoffMultiplier != 3.0 {
		t.Errorf("Expected BackoffMultiplier 3.0, got %f", opts.BackoffMultiplier)
	}

	if opts.Debug != true {
		t.Errorf("Expected Debug true, got %v", opts.Debug)
	}
}

func TestResult(t *testing.T) {
	t.Parallel()

	result := config.Result{
		Title:   "Test Title",
		URL:     "https://example.com",
		Snippet: "Test snippet",
	}

	if result.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got %q", result.Title)
	}

	if result.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got %q", result.URL)
	}

	if result.Snippet != "Test snippet" {
		t.Errorf("Expected Snippet 'Test snippet', got %q", result.Snippet)
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	if config.DefaultMaxRetries != 3 {
		t.Errorf("Expected DefaultMaxRetries 3, got %d", config.DefaultMaxRetries)
	}

	if config.DefaultMaxDelaySecs != 30 {
		t.Errorf("Expected DefaultMaxDelaySecs 30, got %d", config.DefaultMaxDelaySecs)
	}

	if config.DefaultBackoffMult != 2.0 {
		t.Errorf("Expected DefaultBackoffMult 2.0, got %f", config.DefaultBackoffMult)
	}
}

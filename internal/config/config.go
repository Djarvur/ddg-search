// Package config defines configuration types for the ddg-search search client.
package config

import "time"

// Default configuration values.
const (
	DefaultMaxRetries   = 3
	DefaultMaxDelaySecs = 30
	DefaultBackoffMult  = 2.0
)

// SearchOptions configures a DuckDuckGo search query.
type SearchOptions struct {
	// Query is the search string.
	Query string
	// MaxResults limits the number of results returned (0 = unlimited).
	MaxResults int
	// Site filters results to a specific domain.
	Site string
	// Region specifies the search region (e.g., "us-en", "uk-en").
	Region string
	// TimeFilter limits results by time: "d" (day), "w" (week), "m" (month), "y" (year).
	TimeFilter string
	// SafeSearch enables safe search when true.
	SafeSearch bool
}

// RetryOptions configures the retry behavior for rate-limited requests.
type RetryOptions struct {
	// MaxRetries is the maximum number of retry attempts (default: 3).
	MaxRetries int
	// BaseDelay is the initial retry delay (default: 1s).
	BaseDelay time.Duration
	// MaxDelay is the maximum retry delay cap (default: 30s).
	MaxDelay time.Duration
	// BackoffMultiplier is the exponential backoff multiplier (default: 2.0).
	BackoffMultiplier float64
	// Debug enables verbose logging to stderr for troubleshooting.
	Debug bool
}

// DefaultRetryOptions returns the default retry configuration.
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxRetries:        DefaultMaxRetries,
		BaseDelay:         1 * time.Second,
		MaxDelay:          DefaultMaxDelaySecs * time.Second,
		BackoffMultiplier: DefaultBackoffMult,
	}
}

// Result represents a single search result.
type Result struct {
	// Title is the page title.
	Title string `json:"title"`
	// URL is the result link.
	URL string `json:"url"`
	// Snippet is a brief description of the result.
	Snippet string `json:"snippet"`
}

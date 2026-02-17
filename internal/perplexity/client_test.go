// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	if ErrRateLimited == nil {
		t.Error("Expected ErrRateLimited to be defined")
	}

	if ErrMaxRetries == nil {
		t.Error("Expected ErrMaxRetries to be defined")
	}

	if ErrUnauthorized == nil {
		t.Error("Expected ErrUnauthorized to be defined")
	}

	if ErrBadRequest == nil {
		t.Error("Expected ErrBadRequest to be defined")
	}

	if ErrPaymentRequired == nil {
		t.Error("Expected ErrPaymentRequired to be defined")
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	if jitterMs != 500 {
		t.Errorf("Expected jitterMs 500, got %d", jitterMs)
	}

	if apiBaseURL != "https://api.perplexity.ai" {
		t.Errorf("Expected apiBaseURL 'https://api.perplexity.ai', got %q", apiBaseURL)
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	apiKey := "test-api-key"
	retryOpts := DefaultRetryOptions()

	client := NewClient(apiKey, retryOpts)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey %s, got %s", apiKey, client.apiKey)
	}

	if client.retryOptions.MaxRetries != retryOpts.MaxRetries {
		t.Errorf("Expected MaxRetries %d, got %d", retryOpts.MaxRetries, client.retryOptions.MaxRetries)
	}
}

func TestDefaultRetryOptions(t *testing.T) {
	t.Parallel()

	opts := DefaultRetryOptions()

	if opts.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3, got %d", opts.MaxRetries)
	}

	if opts.BaseDelay != 1*time.Second {
		t.Errorf("Expected BaseDelay 1s, got %v", opts.BaseDelay)
	}

	if opts.MaxDelay != 30*time.Second {
		t.Errorf("Expected MaxDelay 30s, got %v", opts.MaxDelay)
	}

	if opts.BackoffMultiplier != 2.0 {
		t.Errorf("Expected BackoffMultiplier 2.0, got %f", opts.BackoffMultiplier)
	}
}

func TestCalculateDelay(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key", RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Second,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
	})

	tests := []struct {
		name     string
		attempt  int
		minDelay time.Duration
		maxDelay time.Duration
	}{
		{"first attempt", 0, 1 * time.Second, 1*time.Second + 500*time.Millisecond},
		{"second attempt", 1, 2 * time.Second, 2*time.Second + 500*time.Millisecond},
		{"third attempt", 2, 4 * time.Second, 4*time.Second + 500*time.Millisecond},
		{"capped delay", 10, 30 * time.Second, 30*time.Second + 500*time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := client.calculateDelay(tt.attempt)
			if delay < tt.minDelay || delay > tt.maxDelay {
				t.Errorf("Delay %v out of expected range [%v, %v]", delay, tt.minDelay, tt.maxDelay)
			}
		})
	}
}

func TestCheckAPIError(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key", DefaultRetryOptions())

	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{"unauthorized", 401, ErrUnauthorized},
		{"bad request", 400, ErrBadRequest},
		{"payment required", 402, ErrPaymentRequired},
		{"rate limited", 429, ErrRateLimited},
		{"server error", 500, errors.New("server error: 500")}, // Returns a generic error
		{"success", 200, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock response
			resp := &resty.Response{}
			resp.RawResponse = &http.Response{StatusCode: tt.statusCode}

			err := client.checkAPIError(resp)
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("checkAPIError() = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("checkAPIError() = %v, want nil", err)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   int
		err      error
		expected bool
	}{
		{"rate limit status", 429, nil, true},
		{"internal server error", 500, nil, true},
		{"bad gateway", 502, nil, true},
		{"service unavailable", 503, nil, true},
		{"gateway timeout", 504, nil, true},
		{"network error", 0, errors.New("network error"), true},
		{"success", 200, nil, false},
		{"not found", 404, nil, false},
		{"bad request", 400, nil, false},
		{"unauthorized", 401, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resp *resty.Response
			if tt.status > 0 {
				resp = &resty.Response{}
				resp.RawResponse = &http.Response{StatusCode: tt.status}
			}

			got := isRateLimited(resp, tt.err)
			if got != tt.expected {
				t.Errorf("isRateLimited() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewClient_WithDebug(t *testing.T) {
	t.Parallel()

	retryOpts := DefaultRetryOptions()
	retryOpts.Debug = true

	client := NewClient("test-key", retryOpts)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.debugWriter == nil {
		t.Error("Expected debugWriter to be set when Debug is true")
	}
}

func TestNewClient_WithoutDebug(t *testing.T) {
	t.Parallel()

	retryOpts := DefaultRetryOptions()
	retryOpts.Debug = false

	client := NewClient("test-key", retryOpts)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.debugWriter == nil {
		t.Error("Expected debugWriter to be set (to io.Discard) when Debug is false")
	}
}

func TestRetryOptions(t *testing.T) {
	t.Parallel()

	opts := RetryOptions{
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

func TestSearchWithEmptyQuery(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key", DefaultRetryOptions())

	ctx := context.Background()

	_, err := client.Search(ctx, "", 5, "sonar-medium-online")
	if err == nil {
		t.Error("Expected error for empty query")
	}

	if err.Error() != "query cannot be empty" {
		t.Errorf("Expected 'query cannot be empty' error, got %v", err)
	}
}

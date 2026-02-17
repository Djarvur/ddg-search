// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
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
	// Note: checkAPIError is tested indirectly through integration tests
	// This is a placeholder for future unit tests with proper mocking
	t.Skip("checkAPIError requires proper resty.Response mocking")
}

func TestSearchWithEmptyQuery(t *testing.T) {
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

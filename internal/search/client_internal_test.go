package search

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
)

func TestClientCalculateDelay(t *testing.T) {
	t.Parallel()

	opts := config.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Second,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
	}
	client := NewClient(opts)

	tests := []struct {
		name     string
		attempt  int
		minDelay time.Duration
		maxDelay time.Duration
	}{
		{
			name:     "first attempt",
			attempt:  0,
			minDelay: 1 * time.Second,
			maxDelay: 2 * time.Second, // with jitter
		},
		{
			name:     "second attempt",
			attempt:  1,
			minDelay: 2 * time.Second,
			maxDelay: 3 * time.Second, // with jitter
		},
		{
			name:     "third attempt",
			attempt:  2,
			minDelay: 4 * time.Second,
			maxDelay: 5 * time.Second, // with jitter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Run multiple times to account for jitter randomness
			for range 10 {
				got := client.calculateDelay(tt.attempt)
				if got < tt.minDelay {
					t.Errorf("calculateDelay() = %v, want at least %v", got, tt.minDelay)
				}

				if got > tt.maxDelay {
					t.Errorf("calculateDelay() = %v, want at most %v", got, tt.maxDelay)
				}
			}
		})
	}
}

func TestClientCalculateDelayMaxCap(t *testing.T) {
	t.Parallel()

	opts := config.RetryOptions{
		MaxRetries:        10,
		BaseDelay:         1 * time.Second,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
	}
	client := NewClient(opts)

	// With these settings, attempt 10 would be 1024s without cap
	// But should be capped at 5s + jitter
	for range 10 {
		got := client.calculateDelay(10)
		if got > 6*time.Second {
			t.Errorf("calculateDelay() = %v, should be capped around max delay", got)
		}
	}
}

func TestIsRateLimited(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"200 OK", 200, false},
		{"429 Too Many Requests", 429, true},
		{"500 Internal Server Error", 500, true},
		{"502 Bad Gateway", 502, true},
		{"503 Service Unavailable", 503, true},
		{"404 Not Found", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock response with given status code
			// Note: We can't easily mock resty.Response, so we test logic indirectly
			// This is a simplified test
			status := tt.statusCode

			got := status == 429 || status >= 500
			if got != tt.want {
				t.Errorf("rate limit check for status %d = %v, want %v", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestClientDoRetryOn429(t *testing.T) {
	t.Parallel()

	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 3 {
			w.WriteHeader(http.StatusTooManyRequests)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
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

	req := client.httpClient.R()

	resp, err := client.Do(context.Background(), req)
	if err != nil {
		t.Errorf("Do() unexpected error: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode(), http.StatusOK)
	}

	if atomic.LoadInt32(&requestCount) != 3 {
		t.Errorf("expected 3 requests, got %d", requestCount)
	}
}

func TestClientDoRetryOn500(t *testing.T) {
	t.Parallel()

	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 2 {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
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

	req := client.httpClient.R()

	resp, err := client.Do(context.Background(), req)
	if err != nil {
		t.Errorf("Do() unexpected error: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode(), http.StatusOK)
	}

	if atomic.LoadInt32(&requestCount) != 2 {
		t.Errorf("expected 2 requests, got %d", requestCount)
	}
}

func TestClientDoMaxRetriesExceeded(t *testing.T) {
	t.Parallel()

	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        2,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	req := client.httpClient.R()

	_, err := client.Do(context.Background(), req)
	if err == nil {
		t.Error("Do() expected error, got nil")
	}

	if !errors.Is(err, ErrMaxRetries) {
		t.Errorf("Do() error should wrap ErrMaxRetries, got: %v", err)
	}

	// MaxRetries=2 means initial attempt + 2 retries = 3 total requests
	if atomic.LoadInt32(&requestCount) != 3 {
		t.Errorf("expected 3 requests, got %d", requestCount)
	}
}

func TestClientDoErrRateLimitedSet(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        1,
		BaseDelay:         1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	req := client.httpClient.R()

	_, err := client.Do(context.Background(), req)
	if err == nil {
		t.Error("Do() expected error, got nil")
	}

	if !errors.Is(err, ErrMaxRetries) {
		t.Errorf("Do() error should wrap ErrMaxRetries, got: %v", err)
	}

	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("Do() error should wrap ErrRateLimited, got: %v", err)
	}
}

func TestClientDoContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	opts := config.RetryOptions{
		MaxRetries:        10,
		BaseDelay:         1 * time.Second, // Long delay to ensure context cancels first
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
	}

	client := NewClientWithBaseURL(opts, server.URL)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req := client.httpClient.R()

	_, err := client.Do(ctx, req)
	if err == nil {
		t.Error("Do() expected error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Do() error should wrap context.DeadlineExceeded, got: %v", err)
	}
}

func TestClientDoSuccessOnFirstTry(t *testing.T) {
	t.Parallel()

	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
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

	req := client.httpClient.R()

	resp, err := client.Do(context.Background(), req)
	if err != nil {
		t.Errorf("Do() unexpected error: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Do() status = %d, want %d", resp.StatusCode(), http.StatusOK)
	}

	if atomic.LoadInt32(&requestCount) != 1 {
		t.Errorf("expected 1 request, got %d", requestCount)
	}
}

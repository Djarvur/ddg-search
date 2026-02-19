package perplexity

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	if opts.MaxRetries != defaultMaxRetries {
		t.Errorf("Expected MaxRetries %d, got %d", defaultMaxRetries, opts.MaxRetries)
	}

	if opts.BaseDelay != defaultBaseDelay {
		t.Errorf("Expected BaseDelay %v, got %v", defaultBaseDelay, opts.BaseDelay)
	}

	if opts.MaxDelay != defaultMaxDelay {
		t.Errorf("Expected MaxDelay %v, got %v", defaultMaxDelay, opts.MaxDelay)
	}

	if opts.BackoffMultiplier != defaultBackoffMult {
		t.Errorf("Expected BackoffMultiplier %f, got %f", defaultBackoffMult, opts.BackoffMultiplier)
	}
}

func TestCalculateDelay(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key", RetryOptions{
		MaxRetries:        defaultMaxRetries,
		BaseDelay:         defaultBaseDelay,
		MaxDelay:          defaultMaxDelay,
		BackoffMultiplier: defaultBackoffMult,
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
			t.Parallel()

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
		{"unauthorized", http.StatusUnauthorized, ErrUnauthorized},
		{"bad request", http.StatusBadRequest, ErrBadRequest},
		{"payment required", http.StatusPaymentRequired, ErrPaymentRequired},
		{"rate limited", http.StatusTooManyRequests, ErrRateLimited},
		{"server error", serverStatus, fmt.Errorf("%w: %d", ErrServer, serverStatus)},
		{"success", http.StatusOK, nil},
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
		{"rate limit status", http.StatusTooManyRequests, nil, true},
		{"internal server error", serverStatus, nil, true},
		{"bad gateway", http.StatusBadGateway, nil, true},
		{"service unavailable", http.StatusServiceUnavailable, nil, true},
		{"gateway timeout", http.StatusGatewayTimeout, nil, true},
		{"network error", 0, ErrNetwork, true},
		{"success", http.StatusOK, nil, false},
		{"not found", http.StatusNotFound, nil, false},
		{"bad request", http.StatusBadRequest, nil, false},
		{"unauthorized", http.StatusUnauthorized, nil, false},
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

func TestDo(t *testing.T) {
	t.Parallel()

	t.Run("context cancellation", func(t *testing.T) {
		t.Parallel()

		client := NewClient("test-key", DefaultRetryOptions())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := client.Do(ctx, client.httpClient.R())
		if err == nil {
			t.Error("Expected error on cancelled context")
		}
	})

	t.Run("max retries exceeded", func(t *testing.T) {
		t.Parallel()

		// Create a test server that returns 401 Unauthorized
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		// Create a client with the test server URL
		client := NewClient("invalid-key", RetryOptions{
			MaxRetries:        1,
			BaseDelay:         1 * time.Millisecond,
			MaxDelay:          10 * time.Millisecond,
			BackoffMultiplier: 1.0,
		})
		client.httpClient.SetBaseURL(server.URL)

		ctx := context.Background()

		_, err := client.Do(ctx, client.httpClient.R())
		if err == nil {
			t.Error("Expected error with invalid API key")
		}
	})
}

func TestDebugf(t *testing.T) {
	t.Parallel()

	t.Run("debug enabled", func(t *testing.T) {
		t.Parallel()

		client := NewClient("test-key", RetryOptions{Debug: true})
		client.debugf("test message")
		// Just verify it doesn't panic
	})

	t.Run("debug disabled", func(t *testing.T) {
		t.Parallel()

		client := NewClient("test-key", RetryOptions{Debug: false})
		client.debugf("test message")
		// Just verify it doesn't panic
	})
}

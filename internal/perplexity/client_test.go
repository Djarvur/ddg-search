//nolint:testpackage
package perplexity

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestErrors(t *testing.T) {
	t.Parallel()

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

	if apiBaseURL != "https://api.perplexity.ai/chat/completions" {
		t.Errorf("Expected apiBaseURL 'https://api.perplexity.ai/chat/completions', got %q", apiBaseURL)
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	apiKey := "test-api-key"

	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected apiKey %s, got %s", apiKey, client.apiKey)
	}
}

func TestCheckAPIError(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key")

	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{"unauthorized", http.StatusUnauthorized, ErrUnauthorized},
		{"bad request", http.StatusBadRequest, ErrBadRequest},
		{"payment required", http.StatusPaymentRequired, ErrPaymentRequired},
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

func TestDo(t *testing.T) {
	t.Parallel()

	t.Run("context cancellation", func(t *testing.T) {
		t.Parallel()

		client := NewClient("test-key")
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := client.Do(ctx, client.httpClient.R(), "GET")
		if err == nil {
			t.Error("Expected error on cancelled context")
		}
	})

	t.Run("unauthorized error", func(t *testing.T) {
		t.Parallel()

		// Create a test server that returns 401 Unauthorized
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		// Create a client with the test server URL
		client := NewClient("invalid-key")
		client.httpClient.SetBaseURL(server.URL)

		ctx := context.Background()

		_, err := client.Do(ctx, client.httpClient.R(), "GET")
		if err == nil {
			t.Error("Expected error with invalid API key")
		}
	})

	t.Run("server error", func(t *testing.T) {
		t.Parallel()

		// Create a test server that returns 500 Internal Server Error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		// Create a client with the test server URL
		client := NewClient("test-key")
		client.httpClient.SetBaseURL(server.URL)

		ctx := context.Background()

		_, err := client.Do(ctx, client.httpClient.R(), "GET")
		if err == nil {
			t.Error("Expected error with server error")
		}
	})
}

// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// Errors returned by the Perplexity client.
var (
	ErrUnauthorized    = errors.New("invalid API key")
	ErrBadRequest      = errors.New("bad request")
	ErrPaymentRequired = errors.New("payment required - API quota exceeded")
	ErrServer          = errors.New("server error")
	ErrQueryEmpty      = errors.New("query cannot be empty")
	ErrAPI             = errors.New("API error")
	ErrNetwork         = errors.New("network error")
	ErrNoChoices       = errors.New("no choices in API response")
)

// Internal constants.
const (
	apiBaseURL = "https://api.perplexity.ai/chat/completions"

	// Request timeout.
	defaultRequestTimeout = 30 * time.Second

	// API constants.
	serverStatus = 500
)

// Client wraps the HTTP client for Perplexity API.
type Client struct {
	httpClient  *resty.Client
	apiKey      string
	debugWriter io.Writer
}

// NewClient creates a new Perplexity API client with the given API key.
func NewClient(apiKey string) *Client {
	client := resty.New().
		SetBaseURL(apiBaseURL).
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTimeout(defaultRequestTimeout)

	return &Client{
		httpClient:  client,
		apiKey:      apiKey,
		debugWriter: io.Discard,
	}
}

// Do executes an HTTP request without retry logic.
func (c *Client) Do(ctx context.Context, req *resty.Request, method string) (*resty.Response, error) {
	resp, err := req.SetContext(ctx).Execute(method, "")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrNetwork, err)
	}

	// Check for API errors
	apiErr := c.checkAPIError(resp)
	if apiErr != nil {
		return nil, apiErr
	}

	return resp, nil
}

// checkAPIError checks the response for API errors and returns appropriate error.
func (c *Client) checkAPIError(resp *resty.Response) error {
	status := resp.StatusCode()

	switch {
	case status == http.StatusUnauthorized:
		return ErrUnauthorized
	case status == http.StatusBadRequest:
		return ErrBadRequest
	case status == http.StatusPaymentRequired:
		return ErrPaymentRequired
	case status >= serverStatus:
		return fmt.Errorf("%w: %d", ErrServer, status)
	default:
		return nil
	}
}

// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Errors returned by the Perplexity client.
var (
	ErrRateLimited     = errors.New("rate limited by Perplexity API")
	ErrMaxRetries      = errors.New("max retries exceeded")
	ErrUnauthorized    = errors.New("invalid API key")
	ErrBadRequest      = errors.New("bad request")
	ErrPaymentRequired = errors.New("payment required - API quota exceeded")
)

// Internal constants.
const (
	jitterMs   = 500
	apiBaseURL = "https://api.perplexity.ai"
)

// RetryOptions configures the retry behavior for rate-limited requests.
type RetryOptions struct {
	// MaxRetries is the maximum number of retry attempts (default: 3).
	MaxRetries int
	// BaseDelay is the initial retry delay in nanoseconds (default: 1s).
	BaseDelay time.Duration
	// MaxDelay is the maximum retry delay cap in nanoseconds (default: 30s).
	MaxDelay time.Duration
	// BackoffMultiplier is the exponential backoff multiplier (default: 2.0).
	BackoffMultiplier float64
	// Debug enables verbose logging to stderr for troubleshooting.
	Debug bool
}

// DefaultRetryOptions returns the default retry configuration.
func DefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * time.Second,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// Client wraps the HTTP client with retry logic for Perplexity API.
type Client struct {
	httpClient   *resty.Client
	apiKey       string
	retryOptions RetryOptions
	debugWriter  io.Writer
}

// NewClient creates a new Perplexity API client with the given API key and retry options.
func NewClient(apiKey string, retryOptions RetryOptions) *Client {
	client := resty.New().
		SetBaseURL(apiBaseURL).
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTimeout(30 * time.Second)

	debugWriter := io.Discard
	if retryOptions.Debug {
		debugWriter = os.Stderr
	}

	return &Client{
		httpClient:   client,
		apiKey:       apiKey,
		retryOptions: retryOptions,
		debugWriter:  debugWriter,
	}
}

// isRateLimited checks if the response indicates rate limiting.
func isRateLimited(resp *resty.Response, err error) bool {
	if err != nil {
		return true // Network errors are retryable
	}

	status := resp.StatusCode()

	// 429 Too Many Requests is standard rate limit response
	// 5xx server errors are retryable
	return status == http.StatusTooManyRequests || status >= 500
}

// calculateDelay computes the retry delay with exponential backoff and jitter.
func (c *Client) calculateDelay(attempt int) time.Duration {
	// Exponential backoff
	delay := c.retryOptions.BaseDelay * time.Duration(1<<uint(attempt))
	if c.retryOptions.BackoffMultiplier != 2.0 {
		delay = c.retryOptions.BaseDelay * time.Duration(c.retryOptions.BackoffMultiplier*float64(attempt+1))
	}

	// Cap at max delay
	if delay > c.retryOptions.MaxDelay {
		delay = c.retryOptions.MaxDelay
	}

	// Add jitter to avoid thundering herd
	jitter := time.Duration(rand.Int63n(int64(jitterMs) * 1000000))
	delay += jitter

	return delay
}

// Do executes an HTTP request with retry logic.
func (c *Client) Do(ctx context.Context, req *resty.Request) (*resty.Response, error) {
	var (
		lastErr  error
		lastResp *resty.Response
	)

	c.debugf("starting request with max_retries=%d, base_delay=%v, max_delay=%v",
		c.retryOptions.MaxRetries, c.retryOptions.BaseDelay, c.retryOptions.MaxDelay)

	for attempt := 0; attempt <= c.retryOptions.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := c.calculateDelay(attempt - 1)
			c.debugf("attempt %d: waiting %v before retry", attempt+1, delay)

			select {
			case <-ctx.Done():
				c.debugf("context cancelled during retry delay")

				return nil, fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
			}
		}

		resp, err := req.Send()
		if err != nil {
			c.debugf("attempt %d: request error: %v", attempt+1, err)

			lastErr = err
			if !isRateLimited(nil, err) {
				break // Non-retryable error
			}

			continue
		}

		lastResp = resp
		c.debugf("attempt %d: status=%d", attempt+1, resp.StatusCode())

		// Check for API errors
		if err := c.checkAPIError(resp); err != nil {
			c.debugf("attempt %d: API error: %v", attempt+1, err)
			lastErr = err

			// Don't retry on authentication or bad request errors
			if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrBadRequest) || errors.Is(err, ErrPaymentRequired) {
				break
			}

			// Retry on rate limit errors
			if isRateLimited(resp, nil) {
				continue
			}

			break
		}

		// Success
		c.debugf("attempt %d: success", attempt+1)

		return resp, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return lastResp, ErrMaxRetries
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
	case status == http.StatusTooManyRequests:
		return ErrRateLimited
	case status >= 500:
		return fmt.Errorf("server error: %d", status)
	default:
		return nil
	}
}

// debugf writes debug output to the debug writer.
func (c *Client) debugf(format string, args ...any) {
	if c.debugWriter != io.Discard {
		fmt.Fprintf(c.debugWriter, "[perplexity] "+format+"\n", args...)
	}
}

// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	ErrServer          = errors.New("server error")
	ErrQueryEmpty      = errors.New("query cannot be empty")
	ErrNetwork         = errors.New("network error")
	ErrAPI             = errors.New("API error")
)

// Internal constants.
const (
	jitterMs   = 500
	apiBaseURL = "https://api.perplexity.ai"

	// Retry configuration constants.
	defaultMaxRetries     = 3
	defaultBaseDelay      = 1 * time.Second
	defaultMaxDelay       = 30 * time.Second
	defaultBackoffMult    = 2.0
	defaultRequestTimeout = 30 * time.Second

	// API constants.
	maxTokens    = 500
	temperature  = 0.1
	serverStatus = 500

	// Jitter constants.
	nanosPerMs = 1_000_000
	maxAttempt = 30
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
		MaxRetries:        defaultMaxRetries,
		BaseDelay:         defaultBaseDelay,
		MaxDelay:          defaultMaxDelay,
		BackoffMultiplier: defaultBackoffMult,
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
		SetTimeout(defaultRequestTimeout)

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

// Do executes an HTTP request with retry logic.
//
//nolint:gocognit // Retry logic requires multiple condition checks
func (c *Client) Do(ctx context.Context, req *resty.Request) (*resty.Response, error) {
	var (
		lastErr  error
		lastResp *resty.Response
	)

	c.debugf("starting request with max_retries=%d, base_delay=%v, max_delay=%v",
		c.retryOptions.MaxRetries, c.retryOptions.BaseDelay, c.retryOptions.MaxDelay)

	for attempt := 0; attempt <= c.retryOptions.MaxRetries; attempt++ {
		err := c.waitForRetry(ctx, attempt)
		if err != nil {
			return nil, err
		}

		resp, err := req.SetContext(ctx).Send()
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

		apiErr := c.checkAPIError(resp)
		if apiErr != nil {
			c.debugf("attempt %d: API error: %v", attempt+1, apiErr)
			lastErr = apiErr

			if c.shouldBreakOnAPIError(apiErr) {
				break
			}

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

// waitForRetry waits for the retry delay if this is not the first attempt.
func (c *Client) waitForRetry(ctx context.Context, attempt int) error {
	if attempt == 0 {
		return nil
	}

	delay := c.calculateDelay(attempt - 1)
	c.debugf("attempt %d: waiting %v before retry", attempt+1, delay)

	select {
	case <-ctx.Done():
		c.debugf("context cancelled during retry delay")

		return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
	case <-time.After(delay):
		return nil
	}
}

// shouldBreakOnAPIError determines if we should break the retry loop based on the error.
func (c *Client) shouldBreakOnAPIError(err error) bool {
	// Don't retry on authentication or bad request errors
	if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrBadRequest) || errors.Is(err, ErrPaymentRequired) {
		return true
	}

	return false
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
	case status >= serverStatus:
		return fmt.Errorf("%w: %d", ErrServer, status)
	default:
		return nil
	}
}

// debugf writes debug output to the debug writer.
func (c *Client) debugf(format string, args ...any) {
	if c.debugWriter != io.Discard {
		_, _ = fmt.Fprintf(c.debugWriter, "[perplexity] "+format+"\n", args...)
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
	return status == http.StatusTooManyRequests || status >= serverStatus
}

// calculateDelay computes the retry delay with exponential backoff and jitter.
func (c *Client) calculateDelay(attempt int) time.Duration {
	// Exponential backoff - use float64 to avoid integer overflow
	var delay time.Duration

	if attempt < maxAttempt {
		// Use safe shift operation
		if attempt >= 0 && attempt < 63 {
			delay = c.retryOptions.BaseDelay * time.Duration(1<<uint(attempt))
		} else {
			delay = c.retryOptions.MaxDelay
		}
	} else {
		// For very large attempts, use the max delay
		delay = c.retryOptions.MaxDelay
	}

	if c.retryOptions.BackoffMultiplier != defaultBackoffMult {
		delay = c.retryOptions.BaseDelay * time.Duration(c.retryOptions.BackoffMultiplier*float64(attempt+1))
	}

	// Cap at max delay
	if delay > c.retryOptions.MaxDelay {
		delay = c.retryOptions.MaxDelay
	}

	// Add jitter to avoid thundering herd using crypto/rand
	jitter, err := randomJitter()
	if err == nil {
		delay += jitter
	}

	return delay
}

// randomJitter generates a random jitter duration using crypto/rand.
//
//nolint:gosec // Modulo operation is safe for small positive values
func randomJitter() (time.Duration, error) {
	var buf [8]byte

	_, err := rand.Read(buf[:])
	if err != nil {
		return 0, fmt.Errorf("failed to read random bytes: %w", err)
	}
	// Convert to int64 and mod by jitterMs * nanosPerMs (nanoseconds per millisecond)
	// Use safe conversion to avoid overflow
	r := binary.BigEndian.Uint64(buf[:])
	maxJitter := uint64(jitterMs) * uint64(nanosPerMs)

	jitterNs := int64(r % maxJitter)
	if jitterNs < 0 {
		jitterNs = -jitterNs
	}

	return time.Duration(jitterNs), nil
}

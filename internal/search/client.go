// Package search provides DuckDuckGo search functionality.
package search

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
	"github.com/go-resty/resty/v2"
)

// Errors returned by the search client.
var (
	ErrRateLimited = errors.New("rate limited by DuckDuckGo")
	ErrMaxRetries  = errors.New("max retries exceeded")
)

// Internal constants.
const (
	jitterMs        = 500
	http5xxThreshold = 500
)

// Client wraps the HTTP client with retry logic.
type Client struct {
	httpClient   *resty.Client
	retryOptions config.RetryOptions
	debugWriter  io.Writer
}

// NewClient creates a new search client with the given retry options.
func NewClient(retryOptions config.RetryOptions) *Client {
	return NewClientWithBaseURL(retryOptions, "https://html.duckduckgo.com/html")
}

// NewClientWithBaseURL creates a new search client with a custom base URL (for testing).
func NewClientWithBaseURL(retryOptions config.RetryOptions, baseURL string) *Client {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("User-Agent", "Mozilla/5.0 (compatible; ddg-search/1.0)").
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	debugWriter := io.Discard
	if retryOptions.Debug {
		debugWriter = os.Stderr
	}

	return &Client{
		httpClient:   client,
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

	// 202 Accepted is sometimes used by DDG to indicate rate limiting
	// 429 Too Many Requests is standard rate limit response
	// 5xx server errors are retryable
	return status == http.StatusAccepted || status == http.StatusTooManyRequests || status >= http5xxThreshold
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
			err = fmt.Errorf("request failed: %w", err)
			c.debugf("attempt %d: request error: %v", attempt+1, err)
		}

		if !isRateLimited(resp, err) {
			c.debugf("attempt %d: request succeeded (not rate limited)", attempt+1)

			return resp, err
		}

		// Log the reason for rate limit detection
		c.debugRateLimitReason(resp, err, attempt+1)

		lastErr = err
		lastResp = resp

		if err == nil && resp.StatusCode() == http.StatusTooManyRequests {
			lastErr = ErrRateLimited
		}
	}

	c.debugf("all %d attempts exhausted, returning ErrMaxRetries", c.retryOptions.MaxRetries+1)

	if lastErr != nil {
		return lastResp, errors.Join(ErrMaxRetries, lastErr)
	}

	return lastResp, ErrMaxRetries
}

// Close releases client resources.
func (c *Client) Close() {
	c.httpClient.GetClient().CloseIdleConnections()
}

// calculateDelay computes the retry delay with exponential backoff and jitter.
func (c *Client) calculateDelay(attempt int) time.Duration {
	delay := float64(c.retryOptions.BaseDelay)
	for range attempt {
		delay *= c.retryOptions.BackoffMultiplier
	}

	if delay > float64(c.retryOptions.MaxDelay) {
		delay = float64(c.retryOptions.MaxDelay)
	}
	// Add jitter (0-500ms)
	jitter := rand.Float64() * float64(jitterMs*time.Millisecond) //nolint:gosec // weak random is fine for jitter

	return time.Duration(delay + jitter)
}

// debugf logs a debug message if debug mode is enabled.
func (c *Client) debugf(format string, args ...any) {
	if c.debugWriter != io.Discard {
		_, _ = fmt.Fprintf(c.debugWriter, "[DEBUG] "+format+"\n", args...)
	}
}

// debugRateLimitReason logs why a request was detected as rate-limited.
func (c *Client) debugRateLimitReason(resp *resty.Response, err error, attempt int) {
	if c.debugWriter == io.Discard {
		return
	}

	var reason string

	switch {
	case err != nil:
		reason = fmt.Sprintf("network error: %v", err)
	case resp.StatusCode() == http.StatusTooManyRequests:
		reason = "HTTP 429 Too Many Requests"
	case resp.StatusCode() >= http5xxThreshold:
		reason = fmt.Sprintf("HTTP %d server error (%s)", resp.StatusCode(), resp.Status())
	default:
		reason = "unknown reason"
	}

	c.debugf("attempt %d: rate limited detected - %s", attempt, reason)
}

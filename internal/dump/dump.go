// Package dump provides web page fetching and markdown conversion.
package dump

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-resty/resty/v2"
)

// Default configuration values.
const (
	DefaultTimeout     = 30 * time.Second
	DefaultUserAgent   = "page-dump/1.0"
	MaxRedirects       = 10
	HTTPErrorThreshold = 400
)

// Errors returned by the dump package.
var (
	ErrInvalidURL        = errors.New("invalid URL format")
	ErrUnsupportedScheme = errors.New("unsupported URL scheme (only HTTP and HTTPS are supported)")
	ErrHTTPError         = errors.New("HTTP error")
)

// Config holds the configuration for fetching and converting.
type Config struct {
	Timeout   time.Duration
	UserAgent string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Timeout:   DefaultTimeout,
		UserAgent: DefaultUserAgent,
	}
}

// ValidateURL parses and validates the URL.
func ValidateURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	// Check scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, ErrUnsupportedScheme
	}

	// Ensure host is present
	if parsedURL.Host == "" {
		return nil, ErrInvalidURL
	}

	return parsedURL, nil
}

// Fetch retrieves the HTML content from the given URL.
func Fetch(ctx context.Context, parsedURL *url.URL, cfg Config) (string, error) {
	client := resty.New().
		SetTimeout(cfg.Timeout).
		SetHeader("User-Agent", cfg.UserAgent).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(MaxRedirects))

	resp, err := client.R().
		SetContext(ctx).
		Get(parsedURL.String())
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode() >= HTTPErrorThreshold {
		return "", fmt.Errorf("%w: %d %s", ErrHTTPError, resp.StatusCode(), resp.Status())
	}

	return resp.String(), nil
}

// Convert transforms HTML content to markdown.
func Convert(html string) (string, error) {
	markdown, err := md.ConvertString(html)
	if err != nil {
		return "", fmt.Errorf("conversion failed: %w", err)
	}

	return markdown, nil
}

// FetchAndConvert fetches a URL and converts the HTML to markdown.
func FetchAndConvert(ctx context.Context, rawURL string, cfg Config) (string, error) {
	parsedURL, err := ValidateURL(rawURL)
	if err != nil {
		return "", err
	}

	html, err := Fetch(ctx, parsedURL, cfg)
	if err != nil {
		return "", err
	}

	return Convert(html)
}

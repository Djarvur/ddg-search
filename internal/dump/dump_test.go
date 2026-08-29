package dump

import (
	"context"
	"errors"
	"testing"
	"time"
)

// testDumpURL is the placeholder URL shared across the tests in this file.
const testDumpURL = "https://example.com"

func TestValidateURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		wantErr error
	}{
		{
			name:    "valid HTTP URL",
			url:     "http://example.com",
			wantErr: nil,
		},
		{
			name:    "valid HTTPS URL",
			url:     testDumpURL,
			wantErr: nil,
		},
		{
			name:    "valid URL with path",
			url:     "https://example.com/path/to/page",
			wantErr: nil,
		},
		{
			name:    "valid URL with query",
			url:     "https://example.com?query=value",
			wantErr: nil,
		},
		{
			name:    "unsupported scheme FTP",
			url:     "ftp://example.com",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "unsupported scheme file",
			url:     "file:///path/to/file",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "missing scheme",
			url:     "example.com",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "invalid URL characters",
			url:     "http://[invalid",
			wantErr: ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ValidateURL(tt.url)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ValidateURL(%q) = nil, want %v", tt.url, tt.wantErr)

					return
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ValidateURL(%q) = %v, want %v", tt.url, err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("ValidateURL(%q) = %v, want nil", tt.url, err)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name: "heading and paragraph",
			html: `<h1>Title</h1><p>Content</p>`,
			want: "# Title\n\nContent",
		},
		{
			name: "link",
			html: `<a href="https://example.com">Link</a>`,
			want: "[Link](https://example.com)",
		},
		{
			name: "list",
			html: `<ul><li>Item 1</li><li>Item 2</li></ul>`,
			want: "- Item 1\n- Item 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Convert(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("Convert() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Expected Timeout %v, got %v", DefaultTimeout, cfg.Timeout)
	}

	if cfg.UserAgent != DefaultUserAgent {
		t.Errorf("Expected UserAgent %q, got %q", DefaultUserAgent, cfg.UserAgent)
	}
}

func TestFetchAndConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "invalid URL",
			url:     "not-a-url",
			cfg:     DefaultConfig(),
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com",
			cfg:     DefaultConfig(),
			wantErr: true,
		},
		{
			name:    "valid URL with default config",
			url:     testDumpURL,
			cfg:     DefaultConfig(),
			wantErr: false, // URL is valid, network error is expected but not tested here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			_, err := FetchAndConvert(ctx, tt.url, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAndConvert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchAndConvert_WithCustomConfig(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Timeout:   5 * time.Second,
		UserAgent: "custom-agent/1.0",
	}

	ctx := context.Background()

	_, err := FetchAndConvert(ctx, "invalid-url", cfg)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestConvert_EmptyHTML(t *testing.T) {
	t.Parallel()

	got, err := Convert("")
	if err != nil {
		t.Errorf("Convert() error = %v", err)
	}

	if got != "" {
		t.Errorf("Convert() = %q, want empty string", got)
	}
}

func TestConvert_ComplexHTML(t *testing.T) {
	t.Parallel()

	html := `<html><head><title>Test</title></head><body><h1>Heading</h1><p>Paragraph</p></body></html>`

	got, err := Convert(html)
	if err != nil {
		t.Errorf("Convert() error = %v", err)
	}

	if got == "" {
		t.Error("Convert() returned empty string for valid HTML")
	}
}

func TestValidateURL_MissingHost(t *testing.T) {
	t.Parallel()

	_, err := ValidateURL("http://")
	if err == nil {
		t.Error("Expected error for URL without host")
	}

	if !errors.Is(err, ErrInvalidURL) {
		t.Errorf("Expected ErrInvalidURL, got %v", err)
	}
}

func TestValidateURL_ValidURLs(t *testing.T) {
	t.Parallel()

	validURLs := []string{
		"http://example.com",
		testDumpURL,
		"https://example.com/path",
		"https://example.com/path?query=value",
		"https://example.com:8080/path",
		"http://localhost",
		"http://127.0.0.1",
	}

	for _, url := range validURLs {
		t.Run(url, func(t *testing.T) {
			t.Parallel()

			parsed, err := ValidateURL(url)
			if err != nil {
				t.Errorf("ValidateURL(%q) error = %v", url, err)
			}

			if parsed == nil {
				t.Errorf("ValidateURL(%q) returned nil parsed URL", url)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	if DefaultTimeout != 30*time.Second {
		t.Errorf("Expected DefaultTimeout 30s, got %v", DefaultTimeout)
	}

	if DefaultUserAgent != "page-dump/1.0" {
		t.Errorf("Expected DefaultUserAgent 'page-dump/1.0', got %q", DefaultUserAgent)
	}

	if MaxRedirects != 10 {
		t.Errorf("Expected MaxRedirects 10, got %d", MaxRedirects)
	}

	if HTTPErrorThreshold != 400 {
		t.Errorf("Expected HTTPErrorThreshold 400, got %d", HTTPErrorThreshold)
	}
}

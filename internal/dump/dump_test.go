package dump_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/dump"
)

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
			url:     "https://example.com",
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
			wantErr: dump.ErrUnsupportedScheme,
		},
		{
			name:    "unsupported scheme file",
			url:     "file:///path/to/file",
			wantErr: dump.ErrUnsupportedScheme,
		},
		{
			name:    "missing scheme",
			url:     "example.com",
			wantErr: dump.ErrUnsupportedScheme,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: dump.ErrUnsupportedScheme,
		},
		{
			name:    "invalid URL characters",
			url:     "http://[invalid",
			wantErr: dump.ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := dump.ValidateURL(tt.url)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("dump.ValidateURL(%q) = nil, want %v", tt.url, tt.wantErr)

					return
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("dump.ValidateURL(%q) = %v, want %v", tt.url, err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("dump.ValidateURL(%q) = %v, want nil", tt.url, err)
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

			got, err := dump.Convert(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("dump.Convert() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("dump.Convert() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := dump.DefaultConfig()

	if cfg.Timeout != dump.DefaultTimeout {
		t.Errorf("Expected Timeout %v, got %v", dump.DefaultTimeout, cfg.Timeout)
	}

	if cfg.UserAgent != dump.DefaultUserAgent {
		t.Errorf("Expected UserAgent %q, got %q", dump.DefaultUserAgent, cfg.UserAgent)
	}
}

func TestFetchAndConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		cfg     dump.Config
		wantErr bool
	}{
		{
			name:    "invalid URL",
			url:     "not-a-url",
			cfg:     dump.DefaultConfig(),
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com",
			cfg:     dump.DefaultConfig(),
			wantErr: true,
		},
		{
			name:    "valid URL with default config",
			url:     "https://example.com",
			cfg:     dump.DefaultConfig(),
			wantErr: false, // URL is valid, network error is expected but not tested here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			_, err := dump.FetchAndConvert(ctx, tt.url, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("dump.FetchAndConvert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchAndConvert_WithCustomConfig(t *testing.T) {
	t.Parallel()

	cfg := dump.Config{
		Timeout:   5 * time.Second,
		UserAgent: "custom-agent/1.0",
	}

	ctx := context.Background()

	_, err := dump.FetchAndConvert(ctx, "invalid-url", cfg)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestConvert_EmptyHTML(t *testing.T) {
	t.Parallel()

	got, err := dump.Convert("")
	if err != nil {
		t.Errorf("dump.Convert() error = %v", err)
	}

	if got != "" {
		t.Errorf("dump.Convert() = %q, want empty string", got)
	}
}

func TestConvert_ComplexHTML(t *testing.T) {
	t.Parallel()

	html := `<html><head><title>Test</title></head><body><h1>Heading</h1><p>Paragraph</p></body></html>`

	got, err := dump.Convert(html)
	if err != nil {
		t.Errorf("dump.Convert() error = %v", err)
	}

	if got == "" {
		t.Error("dump.Convert() returned empty string for valid HTML")
	}
}

func TestValidateURL_MissingHost(t *testing.T) {
	t.Parallel()

	_, err := dump.ValidateURL("http://")
	if err == nil {
		t.Error("Expected error for URL without host")
	}

	if !errors.Is(err, dump.ErrInvalidURL) {
		t.Errorf("Expected dump.ErrInvalidURL, got %v", err)
	}
}

func TestValidateURL_ValidURLs(t *testing.T) {
	t.Parallel()

	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"https://example.com/path",
		"https://example.com/path?query=value",
		"https://example.com:8080/path",
		"http://localhost",
		"http://127.0.0.1",
	}

	for _, url := range validURLs {
		t.Run(url, func(t *testing.T) {
			t.Parallel()

			parsed, err := dump.ValidateURL(url)
			if err != nil {
				t.Errorf("dump.ValidateURL(%q) error = %v", url, err)
			}

			if parsed == nil {
				t.Errorf("dump.ValidateURL(%q) returned nil parsed URL", url)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	if dump.DefaultTimeout != 30*time.Second {
		t.Errorf("Expected dump.DefaultTimeout 30s, got %v", dump.DefaultTimeout)
	}

	if dump.DefaultUserAgent != "page-dump/1.0" {
		t.Errorf("Expected dump.DefaultUserAgent 'page-dump/1.0', got %q", dump.DefaultUserAgent)
	}

	if dump.MaxRedirects != 10 {
		t.Errorf("Expected MaxRedirects 10, got %d", dump.MaxRedirects)
	}

	if dump.HTTPErrorThreshold != 400 {
		t.Errorf("Expected HTTPErrorThreshold 400, got %d", dump.HTTPErrorThreshold)
	}
}

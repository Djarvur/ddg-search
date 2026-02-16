package dump

import (
	"errors"
	"testing"
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

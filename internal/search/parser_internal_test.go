package search

import (
	"testing"
)

const testHTML = `<!DOCTYPE html>
<html>
<body>
<div class="result__body">
  <h2 class="result__title">
    <a class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com%2F">Example Title</a>
  </h2>
  <a class="result__snippet">Example snippet text</a>
</div>
<div class="result__body">
  <h2 class="result__title">
    <a class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Ftest.com%2F">Test Title</a>
  </h2>
  <a class="result__snippet">Test snippet</a>
</div>
</body>
</html>`

func TestParserParse(t *testing.T) {
	t.Parallel()

	p := NewParser()

	t.Run("parse multiple results", func(t *testing.T) {
		t.Parallel()

		results, err := p.Parse(testHTML, 10)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Parse() got %d results, want 2", len(results))
		}
	})

	t.Run("limit results", func(t *testing.T) {
		t.Parallel()

		results, err := p.Parse(testHTML, 1)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Parse() got %d results, want 1", len(results))
		}

		if results[0].Title != "Example Title" {
			t.Errorf("Parse() title = %q, want %q", results[0].Title, "Example Title")
		}
	})

	t.Run("empty HTML", func(t *testing.T) {
		t.Parallel()

		results, err := p.Parse("<html><body></body></html>", 10)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Parse() got %d results, want 0", len(results))
		}
	})
}

func TestParserExtractURL(t *testing.T) {
	t.Parallel()

	p := NewParser()

	tests := []struct {
		name    string
		input   string
		wantURL string
	}{
		{
			name:    "extract from redirect URL",
			input:   "//duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com%2F",
			wantURL: "https://example.com/",
		},
		{
			name:    "plain URL",
			input:   "https://example.com/path",
			wantURL: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := p.extractURL(tt.input)
			if got != tt.wantURL {
				t.Errorf("extractURL() = %q, want %q", got, tt.wantURL)
			}
		})
	}
}

func TestParserIsRateLimitPage(t *testing.T) {
	t.Parallel()

	p := NewParser()

	tests := []struct {
		name     string
		html     string
		wantBool bool
	}{
		{
			name:     "normal page",
			html:     `<html><body>Search results here</body></html>`,
			wantBool: false,
		},
		{
			name:     "captcha page",
			html:     `<html><body>Please complete the captcha</body></html>`,
			wantBool: true,
		},
		{
			name:     "anomaly page",
			html:     `<html><body>anomaly detection triggered</body></html>`,
			wantBool: true,
		},
		{
			name:     "bots use duckduckgo",
			html:     `<html><body>Unfortunately, bots use DuckDuckGo too</body></html>`,
			wantBool: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := p.IsRateLimitPage(tt.html)
			if got != tt.wantBool {
				t.Errorf("IsRateLimitPage() = %v, want %v", got, tt.wantBool)
			}
		})
	}
}

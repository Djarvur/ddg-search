package search

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
)

func TestSearcherBuildSearchParams(t *testing.T) {
	t.Parallel()

	searcher := NewSearcher(config.DefaultRetryOptions())
	t.Cleanup(searcher.Close)

	tests := []struct {
		name string
		opts config.SearchOptions
		want url.Values
	}{
		{
			name: "query only, safe search off",
			opts: config.SearchOptions{Query: testQuery},
			want: url.Values{"q": {testQuery}, "p": {"-1"}},
		},
		{
			name: "site filter is prefixed onto the query",
			opts: config.SearchOptions{Query: "docker compose", Site: "github.com"},
			want: url.Values{"q": {"site:github.com docker compose"}, "p": {"-1"}},
		},
		{
			name: "region maps to kl",
			opts: config.SearchOptions{Query: testQuery, Region: "uk-en"},
			want: url.Values{"q": {testQuery}, "kl": {"uk-en"}, "p": {"-1"}},
		},
		{
			name: "time filter maps to df",
			opts: config.SearchOptions{Query: testQuery, TimeFilter: "w"},
			want: url.Values{"q": {testQuery}, "df": {"w"}, "p": {"-1"}},
		},
		{
			name: "safe search on maps p to 1",
			opts: config.SearchOptions{Query: testQuery, SafeSearch: true},
			want: url.Values{"q": {testQuery}, "p": {"1"}},
		},
		{
			name: "every option at once",
			opts: config.SearchOptions{
				Query:      "kubernetes",
				Site:       "k8s.io",
				Region:     "us-en",
				TimeFilter: "d",
				SafeSearch: true,
			},
			want: url.Values{
				"q":  {"site:k8s.io kubernetes"},
				"kl": {"us-en"},
				"df": {"d"},
				"p":  {"1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := searcher.buildSearchParams(tt.opts)
			if got.Encode() != tt.want.Encode() {
				t.Errorf("buildSearchParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSearcherCalculateDelay pins the Searcher's own backoff, which is a separate
// implementation from Client.calculateDelay and adds no jitter.
func TestSearcherCalculateDelay(t *testing.T) {
	t.Parallel()

	searcher := NewSearcher(config.RetryOptions{
		MaxRetries:        5,
		BaseDelay:         1 * time.Second,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
	})
	t.Cleanup(searcher.Close)

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{"first attempt is the base delay", 0, 1 * time.Second},
		{"second attempt doubles", 1, 2 * time.Second},
		{"third attempt doubles again", 2, 4 * time.Second},
		{"fourth attempt is capped", 3, 5 * time.Second},
		{"far past the cap stays capped", 20, 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := searcher.calculateDelay(tt.attempt); got != tt.want {
				t.Errorf("calculateDelay(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

// TestSearcherSearchBlankQuery covers the early return that avoids sending a
// request DuckDuckGo would reject anyway.
func TestSearcherSearchBlankQuery(t *testing.T) {
	t.Parallel()

	searcher := NewSearcher(config.DefaultRetryOptions())
	t.Cleanup(searcher.Close)

	for _, query := range []string{"", "   ", "\t\n"} {
		results, err := searcher.Search(context.Background(), config.SearchOptions{Query: query})
		if err != nil {
			t.Errorf("Search(%q) error = %v, want nil", query, err)
		}

		if len(results) != 0 {
			t.Errorf("Search(%q) returned %d results, want 0", query, len(results))
		}
	}
}

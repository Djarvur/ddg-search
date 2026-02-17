// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"testing"
)

func TestSearchResults_Markdown(t *testing.T) {
	tests := []struct {
		name     string
		results  *SearchResults
		expected string
	}{
		{
			name: "with citations",
			results: &SearchResults{
				Query:  "test query",
				Answer: "This is a test answer.",
				Citations: []string{
					"https://example.com/1",
					"https://example.com/2",
				},
				References: []Reference{
					{Index: 1, URL: "https://example.com/1"},
					{Index: 2, URL: "https://example.com/2"},
				},
			},
			expected: "This is a test answer.\n\n## Sources\n\n1. https://example.com/1\n2. https://example.com/2\n",
		},
		{
			name: "without citations",
			results: &SearchResults{
				Query:      "test query",
				Answer:     "This is a test answer.",
				Citations:  []string{},
				References: []Reference{},
			},
			expected: "This is a test answer.\n\n",
		},
		{
			name: "answer with newlines",
			results: &SearchResults{
				Query:  "test query",
				Answer: "Line 1\nLine 2\nLine 3",
				Citations: []string{
					"https://example.com/1",
				},
				References: []Reference{
					{Index: 1, URL: "https://example.com/1"},
				},
			},
			expected: "Line 1\nLine 2\nLine 3\n\n## Sources\n\n1. https://example.com/1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.results.Markdown()
			if got != tt.expected {
				t.Errorf("Markdown() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSearchResults_Markdown_EmptyAnswer(t *testing.T) {
	results := &SearchResults{
		Query:      "test query",
		Answer:     "",
		Citations:  []string{},
		References: []Reference{},
	}

	got := results.Markdown()

	expected := "\n\n"
	if got != expected {
		t.Errorf("Markdown() = %q, want %q", got, expected)
	}
}

func TestSearchResults_Markdown_SingleCitation(t *testing.T) {
	results := &SearchResults{
		Query:  "test query",
		Answer: "Single citation answer.",
		Citations: []string{
			"https://example.com/1",
		},
		References: []Reference{
			{Index: 1, URL: "https://example.com/1"},
		},
	}

	got := results.Markdown()

	expected := "Single citation answer.\n\n## Sources\n\n1. https://example.com/1\n"
	if got != expected {
		t.Errorf("Markdown() = %q, want %q", got, expected)
	}
}

func TestSearchResults_Markdown_MultipleCitations(t *testing.T) {
	results := &SearchResults{
		Query:  "test query",
		Answer: "Multiple citations answer.",
		Citations: []string{
			"https://example.com/1",
			"https://example.com/2",
			"https://example.com/3",
			"https://example.com/4",
			"https://example.com/5",
		},
		References: []Reference{
			{Index: 1, URL: "https://example.com/1"},
			{Index: 2, URL: "https://example.com/2"},
			{Index: 3, URL: "https://example.com/3"},
			{Index: 4, URL: "https://example.com/4"},
			{Index: 5, URL: "https://example.com/5"},
		},
	}

	got := results.Markdown()

	expected := "Multiple citations answer.\n\n## Sources\n\n1. https://example.com/1\n2. https://example.com/2\n3. https://example.com/3\n4. https://example.com/4\n5. https://example.com/5\n"
	if got != expected {
		t.Errorf("Markdown() = %q, want %q", got, expected)
	}
}

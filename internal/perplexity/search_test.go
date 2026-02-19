package perplexity

import (
	"testing"
)

func TestSearchOptions(t *testing.T) {
	t.Parallel()

	opts := SearchOptions{
		Query:      "test query",
		MaxResults: 5,
		Model:      "sonar-medium-online",
	}

	if opts.Query != "test query" {
		t.Errorf("Expected Query 'test query', got %q", opts.Query)
	}

	if opts.MaxResults != 5 {
		t.Errorf("Expected MaxResults 5, got %d", opts.MaxResults)
	}

	if opts.Model != "sonar-medium-online" {
		t.Errorf("Expected Model 'sonar-medium-online', got %q", opts.Model)
	}
}

func TestSearchOptions_ZeroValues(t *testing.T) {
	t.Parallel()

	opts := SearchOptions{}

	if opts.Query != "" {
		t.Errorf("Expected Query empty, got %q", opts.Query)
	}

	if opts.MaxResults != 0 {
		t.Errorf("Expected MaxResults 0, got %d", opts.MaxResults)
	}

	if opts.Model != "" {
		t.Errorf("Expected Model empty, got %q", opts.Model)
	}
}

func TestReference(t *testing.T) {
	t.Parallel()

	ref := Reference{
		Index: 1,
		URL:   "https://example.com",
	}

	if ref.Index != 1 {
		t.Errorf("Expected Index 1, got %d", ref.Index)
	}

	if ref.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got %q", ref.URL)
	}
}

func TestSearchResults(t *testing.T) {
	t.Parallel()

	results := SearchResults{
		Query:  "test query",
		Answer: "test answer",
		Citations: []string{
			"https://example.com/1",
			"https://example.com/2",
		},
		References: []Reference{
			{Index: 1, URL: "https://example.com/1"},
			{Index: 2, URL: "https://example.com/2"},
		},
	}

	if results.Query != "test query" {
		t.Errorf("Expected Query 'test query', got %q", results.Query)
	}

	if results.Answer != "test answer" {
		t.Errorf("Expected Answer 'test answer', got %q", results.Answer)
	}

	if len(results.Citations) != 2 {
		t.Errorf("Expected 2 citations, got %d", len(results.Citations))
	}

	if len(results.References) != 2 {
		t.Errorf("Expected 2 references, got %d", len(results.References))
	}
}

func TestSearchResults_Empty(t *testing.T) {
	t.Parallel()

	results := SearchResults{}

	if results.Query != "" {
		t.Errorf("Expected Query empty, got %q", results.Query)
	}

	if results.Answer != "" {
		t.Errorf("Expected Answer empty, got %q", results.Answer)
	}

	if results.Citations != nil {
		t.Errorf("Expected Citations nil, got %v", results.Citations)
	}

	if results.References != nil {
		t.Errorf("Expected References nil, got %v", results.References)
	}
}

func TestSearchResults_Markdown(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := tt.results.Markdown()
			if got != tt.expected {
				t.Errorf("Markdown() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSearchResults_Markdown_EmptyAnswer(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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

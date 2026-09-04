package perplexity

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Fixtures shared across the tests in this file.
const (
	testQuery = "test query"
	testURL1  = "https://example.com/1"
	testURL2  = "https://example.com/2"
)

// newTestClient returns a client pointed at srv with retries effectively disabled.
func newTestClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()

	client := NewClient("test-key", RetryOptions{MaxRetries: 0, BackoffMultiplier: defaultBackoffMult})
	client.httpClient.SetBaseURL(srv.URL)

	return client
}

// TestSearchSendsPostToChatCompletions guards against the request being dispatched
// with an empty method and path, which resty resolves to GET on the bare base URL.
func TestSearchSendsPostToChatCompletions(t *testing.T) {
	t.Parallel()

	var (
		gotMethod string
		gotPath   string
		gotBody   map[string]any
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path

		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)

		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"answer"}}]}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv).Search(context.Background(), testQuery, 5, "sonar")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want %q", gotMethod, http.MethodPost)
	}

	if gotPath != chatCompletionsPath {
		t.Errorf("path = %q, want %q", gotPath, chatCompletionsPath)
	}

	if gotBody["model"] != "sonar" {
		t.Errorf("model = %v, want %q", gotBody["model"], "sonar")
	}
}

// TestSearchReadsAnswerFromChoices guards against reading the answer from a
// top-level "answer" field, which the API does not return.
func TestSearchReadsAnswerFromChoices(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
			"answer": "wrong field",
			"choices": [{"message": {"role": "assistant", "content": "the real answer"}}],
			"search_results": [{"title": "Example", "url": "https://example.com/1"}]
		}`))
	}))
	defer srv.Close()

	results, err := newTestClient(t, srv).Search(context.Background(), "q", 0, "sonar")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if results.Answer != "the real answer" {
		t.Errorf("Answer = %q, want %q", results.Answer, "the real answer")
	}

	want := "the real answer\n\n## Sources\n\n1. [Example](https://example.com/1)\n"
	if got := results.Markdown(); got != want {
		t.Errorf("Markdown() = %q, want %q", got, want)
	}
}

func TestSearchSourceSelection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		maxResults int
		wantURLs   []string
	}{
		{
			name:       "prefers search_results over deprecated citations",
			body:       `{"choices":[{"message":{"content":"a"}}],"citations":["https://old.example"],"search_results":[{"title":"New","url":"https://new.example"}]}`,
			maxResults: 0,
			wantURLs:   []string{"https://new.example"},
		},
		{
			name:       "falls back to citations when search_results is absent",
			body:       `{"choices":[{"message":{"content":"a"}}],"citations":["https://old.example"]}`,
			maxResults: 0,
			wantURLs:   []string{"https://old.example"},
		},
		{
			name:       "max-results caps the source list",
			body:       `{"choices":[{"message":{"content":"a"}}],"citations":["https://a.example","https://b.example","https://c.example"]}`,
			maxResults: 2,
			wantURLs:   []string{"https://a.example", "https://b.example"},
		},
		{
			name:       "max-results above the source count keeps every source",
			body:       `{"choices":[{"message":{"content":"a"}}],"citations":["https://a.example"]}`,
			maxResults: 10,
			wantURLs:   []string{"https://a.example"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			results, err := newTestClient(t, srv).Search(context.Background(), "q", tt.maxResults, "sonar")
			if err != nil {
				t.Fatalf("Search() error = %v", err)
			}

			assertSources(t, results, tt.wantURLs)
		})
	}
}

// assertSources checks that results carry exactly wantURLs as sources, numbered from 1.
func assertSources(t *testing.T, results *SearchResults, wantURLs []string) {
	t.Helper()

	if len(results.Citations) != len(wantURLs) {
		t.Fatalf("got %d citations, want %d", len(results.Citations), len(wantURLs))
	}

	for i, want := range wantURLs {
		if results.Citations[i] != want {
			t.Errorf("Citations[%d] = %q, want %q", i, results.Citations[i], want)
		}

		if results.References[i].Index != i+1 {
			t.Errorf("References[%d].Index = %d, want %d", i, results.References[i].Index, i+1)
		}
	}
}

func TestSearchNoChoices(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer srv.Close()

	_, err := newTestClient(t, srv).Search(context.Background(), "q", 0, "sonar")
	if !errors.Is(err, ErrNoChoices) {
		t.Errorf("Search() error = %v, want ErrNoChoices", err)
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key", DefaultRetryOptions())

	_, err := client.Search(context.Background(), "", 5, "sonar")
	if !errors.Is(err, ErrQueryEmpty) {
		t.Errorf("Search() error = %v, want ErrQueryEmpty", err)
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
		Query:  testQuery,
		Answer: "test answer",
		Citations: []string{
			testURL1,
			testURL2,
		},
		References: []Reference{
			{Index: 1, URL: testURL1},
			{Index: 2, URL: testURL2},
		},
	}

	if results.Query != testQuery {
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
				Query:  testQuery,
				Answer: "This is a test answer.",
				Citations: []string{
					testURL1,
					testURL2,
				},
				References: []Reference{
					{Index: 1, URL: testURL1},
					{Index: 2, URL: testURL2},
				},
			},
			expected: "This is a test answer.\n\n## Sources\n\n1. https://example.com/1\n2. https://example.com/2\n",
		},
		{
			name: "without citations",
			results: &SearchResults{
				Query:      testQuery,
				Answer:     "This is a test answer.",
				Citations:  []string{},
				References: []Reference{},
			},
			expected: "This is a test answer.\n\n",
		},
		{
			name: "answer with newlines",
			results: &SearchResults{
				Query:  testQuery,
				Answer: "Line 1\nLine 2\nLine 3",
				Citations: []string{
					testURL1,
				},
				References: []Reference{
					{Index: 1, URL: testURL1},
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
		Query:      testQuery,
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
		Query:  testQuery,
		Answer: "Single citation answer.",
		Citations: []string{
			testURL1,
		},
		References: []Reference{
			{Index: 1, URL: testURL1},
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
		Query:  testQuery,
		Answer: "Multiple citations answer.",
		Citations: []string{
			testURL1,
			testURL2,
			"https://example.com/3",
			"https://example.com/4",
			"https://example.com/5",
		},
		References: []Reference{
			{Index: 1, URL: testURL1},
			{Index: 2, URL: testURL2},
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

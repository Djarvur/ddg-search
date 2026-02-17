// Package perplexity provides Perplexity API search functionality.
package perplexity

import (
	"context"
	"os"
	"testing"
)

// TestSearchIntegration performs an integration test with the real Perplexity API.
// This test is skipped unless PERPLEXITY_API_KEY is set and INTEGRATION_TEST is set to "true".
func TestSearchIntegration(t *testing.T) {
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		t.Skip("PERPLEXITY_API_KEY not set, skipping integration test")
	}

	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("INTEGRATION_TEST not set to 'true', skipping integration test")
	}

	client := NewClient(apiKey, DefaultRetryOptions())
	ctx := context.Background()

	// Test successful search
	t.Run("successful search", func(t *testing.T) {
		results, err := client.Search(ctx, "What is Go programming language?", 5, "sonar-medium-online")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}

		if results == nil {
			t.Fatal("Expected non-nil results")
		}

		if results.Answer == "" {
			t.Error("Expected non-empty answer")
		}

		if len(results.Citations) == 0 {
			t.Error("Expected at least one citation")
		}

		markdown := results.Markdown()
		if markdown == "" {
			t.Error("Expected non-empty markdown output")
		}

		t.Logf("Answer: %s", results.Answer)
		t.Logf("Citations: %v", results.Citations)
	})

	// Test error scenarios
	t.Run("error scenarios", func(t *testing.T) {
		// Test with invalid API key
		t.Run("invalid API key", func(t *testing.T) {
			invalidClient := NewClient("invalid-key", DefaultRetryOptions())
			_, err := invalidClient.Search(ctx, "test query", 5, "sonar-medium-online")
			if err == nil {
				t.Error("Expected error with invalid API key")
			}
		})

		// Test with empty query
		t.Run("empty query", func(t *testing.T) {
			_, err := client.Search(ctx, "", 5, "sonar-medium-online")
			if err == nil {
				t.Error("Expected error with empty query")
			}
		})
	})

	// Test with different models
	t.Run("different models", func(t *testing.T) {
		models := []string{
			"sonar-medium-online",
			"sonar-small-online",
		}

		for _, model := range models {
			t.Run(model, func(t *testing.T) {
				results, err := client.Search(ctx, "What is golang?", 3, model)
				if err != nil {
					t.Logf("Warning: Model %s failed: %v", model, err)
					// Don't fail the test, as some models may not be available
					return
				}

				if results.Answer == "" {
					t.Errorf("Expected non-empty answer for model %s", model)
				}
			})
		}
	})
}

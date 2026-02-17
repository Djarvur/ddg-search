// Package main provides the CLI entry point for perplexity-search.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Djarvur/ddg-search/internal/perplexity"
	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags.
const version = "dev"

// Default CLI values.
const (
	defaultMaxResults = 5
	defaultModel      = "sonar-medium-online"
)

// errNoQuery is returned when no search query is provided.
var errNoQuery = errors.New("no search query provided")

// errNoAPIKey is returned when the Perplexity API key is not set.
var errNoAPIKey = errors.New("PERPLEXITY_API_KEY environment variable not set. Please set it in your .env file or shell environment.")

func main() {
	cmd := &cli.Command{
		Name:    "perplexity-search",
		Usage:   "Search the web using Perplexity API",
		Version: version,
		Flags: []cli.Flag{
			// Search options
			&cli.IntFlag{
				Name:  "max-results",
				Usage: "maximum number of results to return",
				Value: defaultMaxResults,
			},
			&cli.StringFlag{
				Name:  "model",
				Usage: "Perplexity model to use (e.g., sonar-medium-online, sonar-pro-online)",
				Value: defaultModel,
			},
			// Debug
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug logging to stderr",
				Value: false,
			},
		},
		Action: runSearch,
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runSearch(ctx context.Context, cmd *cli.Command) error {
	// Get query from arguments
	args := cmd.Args()
	if args.Len() == 0 {
		return errNoQuery
	}

	query := args.First()

	// Get flags
	maxResults := cmd.Int("max-results")
	model := cmd.String("model")
	debug := cmd.Bool("debug")

	// Get API key from environment
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		return errNoAPIKey
	}

	// Create client
	client := perplexity.NewClient(apiKey, perplexity.RetryOptions{
		MaxRetries:        3,
		BaseDelay:         1 * 1000000000, // 1 second
		MaxDelay:          30 * 1000000000, // 30 seconds
		BackoffMultiplier: 2.0,
		Debug:             debug,
	})

	// Perform search
	results, err := client.Search(ctx, query, maxResults, model)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Output results
	fmt.Println(results.Markdown())

	return nil
}

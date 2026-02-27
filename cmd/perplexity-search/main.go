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
	defaultModel = "sonar-medium-online"
)

// errNoQuery is returned when no search query is provided.
var errNoQuery = errors.New("no search query provided")

// errNoAPIKey is returned when the Perplexity API key is not set.
var errNoAPIKey = errors.New(
	"perplexity API key environment variable not set. " +
		"Please set it in your .env file or shell environment",
)

func main() {
	cmd := &cli.Command{
		Name:    "perplexity-search",
		Usage:   "Search the web using Perplexity API",
		Version: version,
		Flags: []cli.Flag{
			// Search options
			&cli.StringFlag{
				Name:  "model",
				Usage: "Perplexity model to use (e.g., sonar-medium-online, sonar-pro-online)",
				Value: defaultModel,
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
	model := cmd.String("model")

	// Get API key from environment
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		return errNoAPIKey
	}

	// Create client
	client := perplexity.NewClient(apiKey)

	// Perform search
	results, err := client.Search(ctx, query, model)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Output results
	_, _ = fmt.Fprint(os.Stdout, results.Markdown())

	return nil
}

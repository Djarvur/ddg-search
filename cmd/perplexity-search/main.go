// Package main provides the CLI entry point for perplexity-search.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Djarvur/ddg-search/internal/perplexity"
	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags; it must be a var for the linker to patch it.
var version = "dev"

// Default CLI values.
const (
	defaultMaxResults = 5
	defaultModel      = "sonar"

	// Retry configuration.
	defaultMaxRetries        = 3
	defaultMaxDelaySeconds   = 30
	defaultBackoffMultiplier = 2.0
)

// errNoQuery is returned when no search query is provided.
var errNoQuery = errors.New("no search query provided")

// errNoAPIKey is returned when the Perplexity API key is not set.
var errNoAPIKey = errors.New(
	"PERPLEXITY_API_KEY is not set. " +
		"Export it in your shell environment before running perplexity-search",
)

func main() {
	cmd := &cli.Command{
		Name:    "perplexity-search",
		Usage:   "Search the web using Perplexity API",
		Version: version,
		Flags: []cli.Flag{
			// Search options
			&cli.IntFlag{
				Name:  "max-results",
				Usage: "maximum number of sources to list alongside the answer",
				Value: defaultMaxResults,
			},
			&cli.StringFlag{
				Name:  "model",
				Usage: "Perplexity model to use (e.g., sonar, sonar-pro, sonar-reasoning)",
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

	query := strings.Join(args.Slice(), " ")

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
		MaxRetries:        defaultMaxRetries,
		BaseDelay:         time.Second,
		MaxDelay:          time.Duration(defaultMaxDelaySeconds) * time.Second,
		BackoffMultiplier: defaultBackoffMultiplier,
		Debug:             debug,
	})

	// Perform search
	results, err := client.Search(ctx, query, maxResults, model)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Output results
	_, _ = fmt.Fprint(os.Stdout, results.Markdown())

	return nil
}

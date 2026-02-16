// Package main provides the CLI entry point for ddg-search.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Djarvur/ddg-search/internal/config"
	"github.com/Djarvur/ddg-search/internal/search"
	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags.
const version = "dev"

// Default CLI values.
const (
	defaultMaxResults   = 10
	defaultMaxRetries   = 3
	defaultMaxDelaySecs = 30
	backoffMultiplier   = 2.0
)

// errNoQuery is returned when no search query is provided.
var errNoQuery = errors.New("no search query provided")

func main() {
	cmd := &cli.Command{
		Name:    "ddg-search",
		Usage:   "DuckDuckGo search from the command line",
		Version: version,
		Flags: []cli.Flag{
			// Search options
			&cli.IntFlag{
				Name:  "max-results",
				Usage: "maximum number of results to return",
				Value: defaultMaxResults,
			},
			&cli.StringFlag{
				Name:  "site",
				Usage: "filter results to a specific domain",
			},
			&cli.StringFlag{
				Name:  "region",
				Usage: "search region (e.g., us-en, uk-en)",
				Value: "us-en",
			},
			&cli.StringFlag{
				Name:  "time",
				Usage: "time filter: d (day), w (week), m (month), y (year)",
			},
			&cli.BoolFlag{
				Name:  "safe-search",
				Usage: "enable safe search",
				Value: false,
			},
			// Retry options
			&cli.IntFlag{
				Name:  "max-retries",
				Usage: "maximum retry attempts on rate limiting",
				Value: defaultMaxRetries,
			},
			&cli.DurationFlag{
				Name:  "retry-delay",
				Usage: "initial retry delay",
				Value: 1 * time.Second,
			},
			&cli.DurationFlag{
				Name:  "max-retry-delay",
				Usage: "maximum retry delay cap",
				Value: defaultMaxDelaySecs * time.Second,
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

	var querySb82 strings.Builder
	for i := 1; i < args.Len(); i++ {
		querySb82.WriteString(" " + args.Get(i))
	}

	query += querySb82.String()

	// Build search options
	searchOpts := config.SearchOptions{
		Query:      query,
		MaxResults: cmd.Int("max-results"),
		Site:       cmd.String("site"),
		Region:     cmd.String("region"),
		TimeFilter: cmd.String("time"),
		SafeSearch: cmd.Bool("safe-search"),
	}

	// Build retry options
	retryOpts := config.RetryOptions{
		MaxRetries:        cmd.Int("max-retries"),
		BaseDelay:         cmd.Duration("retry-delay"),
		MaxDelay:          cmd.Duration("max-retry-delay"),
		BackoffMultiplier: backoffMultiplier,
		Debug:             cmd.Bool("debug"),
	}

	// Perform search
	searcher := search.NewSearcher(retryOpts)
	defer searcher.Close()

	output, err := searcher.SearchMarkdown(ctx, searchOpts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Output results
	_, err = fmt.Fprintln(os.Stdout, output)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

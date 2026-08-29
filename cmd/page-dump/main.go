// Package main provides the CLI entry point for page-dump.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Djarvur/ddg-search/internal/dump"
	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags; it must be a var for the linker to patch it.
var version = "dev"

// Error for missing URL argument.
var errNoURL = errors.New("no URL provided")

func main() {
	cmd := &cli.Command{
		Name:    "page-dump",
		Usage:   "Fetch a web page and convert to markdown",
		Version: version,
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "request timeout",
				Value: dump.DefaultTimeout,
			},
			&cli.StringFlag{
				Name:  "user-agent",
				Usage: "custom user agent string",
				Value: dump.DefaultUserAgent,
			},
		},
		Action: runDump,
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runDump(ctx context.Context, cmd *cli.Command) error {
	// Get URL from arguments
	args := cmd.Args()
	if args.Len() == 0 {
		return errNoURL
	}

	rawURL := args.First()

	// Build config from flags
	cfg := dump.Config{
		Timeout:   cmd.Duration("timeout"),
		UserAgent: cmd.String("user-agent"),
	}

	// Fetch and convert
	markdown, err := dump.FetchAndConvert(ctx, rawURL, cfg)
	if err != nil {
		return fmt.Errorf("failed to fetch and convert: %w", err)
	}

	// Output markdown to stdout
	_, _ = os.Stdout.WriteString(markdown + "\n")

	return nil
}

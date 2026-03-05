// Package mcplog provides logging utilities for the MCP server.
package mcplog

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

var (
	// ErrUnknownLogLevel is returned when an unknown log level is provided.
	ErrUnknownLogLevel = errors.New("unknown log level")
)

// NewLogger creates a new logger with the specified level.
func NewLogger(level string) (*slog.Logger, error) {
	logLevel, err := parseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)

	return logger, nil
}

// parseLevel parses a log level string.
func parseLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("%w: %s", ErrUnknownLogLevel, level)
	}
}

package mcp_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Djarvur/ddg-search/internal/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &mcp.Config{
		Name:    "test-server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(cfg, logger)
	require.NotNil(t, server)
}

func TestRegisterTool(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &mcp.Config{
		Name:    "test-server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(cfg, logger)

	tool := mcpgo.Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: mcpgo.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"input": map[string]any{
					"type": "string",
				},
			},
			Required: []string{"input"},
		},
	}

	handler := func(
		_ context.Context,
		_ any,
		_ mcpgo.CallToolRequest,
	) (*mcpgo.CallToolResult, error) {
		return mcpgo.NewToolResultText("test response"), nil
	}

	// This should not panic
	server.RegisterTool(tool, handler)
}

func TestToolHandler(t *testing.T) {
	t.Parallel()

	handler := func(
		_ context.Context,
		_ any,
		_ mcpgo.CallToolRequest,
	) (*mcpgo.CallToolResult, error) { //nolint:unparam
		return mcpgo.NewToolResultText("test response"), nil
	}

	result, err := handler(context.Background(), nil, mcpgo.CallToolRequest{})
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestToolHandlerWithError(t *testing.T) {
	t.Parallel()

	handler := func(
		_ context.Context,
		_ any,
		_ mcpgo.CallToolRequest,
	) (*mcpgo.CallToolResult, error) {
		return mcpgo.NewToolResultError("test error"), nil
	}

	result, err := handler(context.Background(), nil, mcpgo.CallToolRequest{})
	require.NoError(t, err)
	require.NotNil(t, result)
}

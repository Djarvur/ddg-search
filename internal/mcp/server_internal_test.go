package mcp_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/mcp"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestServerIntegration(t *testing.T) {
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

	handler := func(_ context.Context, _ any, _ mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		return mcpgo.NewToolResultText("test response"), nil
	}

	server.RegisterTool(tool, handler)
}

func TestServerShutdown(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &mcp.Config{
		Name:    "test-server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(cfg, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := server.Shutdown(ctx)
	require.NoError(t, err)
}

func TestServerServe(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &mcp.Config{
		Name:    "test-server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(cfg, logger)

	// Note: We can't actually test Serve() because it blocks on stdio
	// This test just verifies that server can be created
	require.NotNil(t, server)
}

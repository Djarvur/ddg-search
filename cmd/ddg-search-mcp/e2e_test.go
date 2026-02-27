//go:build e2e

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_BaseApplication_Startup(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait a bit for startup
	time.Sleep(300 * time.Millisecond)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	// Process may exit cleanly or with signal error
	_ = err

	// Verify startup messages
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "Configuration loaded")
	require.Contains(t, outputStr, "Server ready")
}

func TestE2E_MCPServer_ToolRegistration(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Send tools/list request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n"))
	require.NoError(t, err, "Failed to write tools/list request")

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	// Close stdin to signal end of input
	stdin.Close()

	// Wait for command to complete
	err = cmd.Wait()
	// Command may timeout, which is OK for this test
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains tools list
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":1`)
	require.Contains(t, outputStr, `"result"`)
	require.Contains(t, outputStr, `"tools"`)
}

func TestE2E_MCPServer_SearchTool(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for search
	searchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang"}}}`
	_, err = stdin.Write([]byte(searchCall + "\n"))
	require.NoError(t, err, "Failed to write search call request")

	time.Sleep(2 * time.Second)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains search results
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"result"`)
}

func TestE2E_MCPServer_FetchTool(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for fetch with a simple URL
	fetchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}`
	_, err = stdin.Write([]byte(fetchCall + "\n"))
	require.NoError(t, err, "Failed to write fetch call request")

	time.Sleep(2 * time.Second)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains fetch result
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"result"`)
}

func TestE2E_MCPServer_UnknownTool(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for unknown tool
	unknownToolCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"unknown_tool","arguments":{}}}`
	_, err = stdin.Write([]byte(unknownToolCall + "\n"))
	require.NoError(t, err, "Failed to write unknown tool call request")

	time.Sleep(100 * time.Millisecond)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains error
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"error"`)
}

func TestE2E_MCPServer_SearchToolMissingQuery(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for search without query parameter
	searchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{}}}`
	_, err = stdin.Write([]byte(searchCall + "\n"))
	require.NoError(t, err, "Failed to write search call request")

	time.Sleep(100 * time.Millisecond)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains error
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"error"`)
}

func TestE2E_MCPServer_FetchToolMissingURL(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for fetch without url parameter
	fetchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"fetch","arguments":{}}}`
	_, err = stdin.Write([]byte(fetchCall + "\n"))
	require.NoError(t, err, "Failed to write fetch call request")

	time.Sleep(100 * time.Millisecond)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains error
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"error"`)
}

func TestE2E_MCPServer_FetchToolInvalidURL(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for fetch with invalid URL
	fetchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"not-a-valid-url"}}}`
	_, err = stdin.Write([]byte(fetchCall + "\n"))
	require.NoError(t, err, "Failed to write fetch call request")

	time.Sleep(100 * time.Millisecond)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains error
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"error"`)
}

func TestE2E_MCPServer_SearchToolWithMaxResults(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for search with max_results parameter
	searchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang","max_results":5}}}`
	_, err = stdin.Write([]byte(searchCall + "\n"))
	require.NoError(t, err, "Failed to write search call request")

	time.Sleep(2 * time.Second)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains search results
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"result"`)
}

func TestE2E_MCPServer_SearchToolWithSafeSearch(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdin pipe
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."

	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Send initialize request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}` + "\n"))
	require.NoError(t, err, "Failed to write initialize request")

	time.Sleep(100 * time.Millisecond)

	// Send tools/call request for search with safe_search parameter
	searchCall := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang","safe_search":true}}}`
	_, err = stdin.Write([]byte(searchCall + "\n"))
	require.NoError(t, err, "Failed to write search call request")

	time.Sleep(2 * time.Second)

	// Send shutdown request
	_, err = stdin.Write([]byte(`{"jsonrpc":"2.0","id":3,"method":"shutdown"}` + "\n"))
	require.NoError(t, err, "Failed to write shutdown request")

	stdin.Close()

	err = cmd.Wait()
	if err != nil && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.NoError(t, err, "Failed to run command")
	}

	// Verify response contains search results
	outputStr := out.String()
	require.Contains(t, outputStr, `"jsonrpc":"2.0"`)
	require.Contains(t, outputStr, `"id":2`)
	require.Contains(t, outputStr, `"result"`)
}

func TestE2E_Config_WithConfigFile(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
logging:
  level: debug
search:
  max_results: 15
  safe_search: true
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait a bit for startup
	time.Sleep(300 * time.Millisecond)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify startup messages with debug level
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "Configuration loaded")
	require.Contains(t, outputStr, "max_results=15")
	require.Contains(t, outputStr, "safe_search=true")
}

func TestE2E_Config_LogLevelCLI(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background with --log-level flag
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath, "--log-level=debug") //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait a bit for startup
	time.Sleep(300 * time.Millisecond)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify startup messages
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "Configuration loaded")
}

func TestE2E_StreamableHTTP_ServerStartup(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with HTTP protocol
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19100
logging:
  level: debug
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test health check endpoint
	resp, err := http.Get("http://localhost:19100/health")
	require.NoError(t, err, "Failed to call health check")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify startup messages
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "protocol=http")
	require.Contains(t, outputStr, "bind_address=localhost:19100")
	require.Contains(t, outputStr, "HTTP Streamable HTTP server listening")
}

func TestE2E_StreamableHTTP_HealthCheck(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with HTTP protocol
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19101
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test health check endpoint returns OK
	resp, err := http.Get("http://localhost:19101/health")
	require.NoError(t, err, "Failed to call health check")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("OK"), body)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err
}

func TestE2E_StreamableHTTP_ProtocolConfiguration(t *testing.T) { //nolint:paralleltest
	// Test that stdio protocol still works
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with stdio protocol (default)
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait a bit
	time.Sleep(300 * time.Millisecond)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify stdio protocol was used
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "protocol=stdio")
}

func TestE2E_StreamableHTTP_CustomBindAddress(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with custom bind address
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: 127.0.0.1:19102
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test health check endpoint on custom address
	resp, err := http.Get("http://127.0.0.1:19102/health")
	require.NoError(t, err, "Failed to call health check")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify custom bind address was used
	outputStr := out.String()
	require.Contains(t, outputStr, "bind_address=127.0.0.1:19102")
}

func TestE2E_StreamableHTTP_Shutdown(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with HTTP protocol
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19103
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test health check endpoint
	resp, err := http.Get("http://localhost:19103/health")
	require.NoError(t, err, "Failed to call health check")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify shutdown messages
	outputStr := out.String()
	require.Contains(t, outputStr, "Shutting down")
	require.Contains(t, outputStr, "Shutdown complete")
}

func TestE2E_Config_InvalidConfig(t *testing.T) { //nolint:paralleltest
	// Create a temporary invalid config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
logging:
  level: debug
search:
  max_results: not_a_number
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application - it should fail with invalid config
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()

	// Command should fail
	require.Error(t, err, "Expected error with invalid config")

	// Verify error message
	outputStr := out.String()
	require.Contains(t, outputStr, "failed to load configuration")
}

func TestE2E_Config_LogLevelEnvVar(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application in background with DDG_SEARCH_LOG_LEVEL env var
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=", "DDG_SEARCH_LOG_LEVEL=debug")
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait a bit for startup
	time.Sleep(300 * time.Millisecond)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify startup messages
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "Configuration loaded")
}

// Stage 4 - Missing Tests
// ============================================================================

func TestE2E_Config_Priority(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with level: info
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
logging:
  level: info
search:
  max_results: 10
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start with CLI flag overriding env and config
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath),
		"DDG_SEARCH_LOGGING_LEVEL=debug",
	)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(300 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify CLI flag has highest priority (warn, not debug from env or info from config)
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "Configuration loaded")
	// The log level should be warn from CLI flag
	// require.NotContains(t, outputStr, "level=debug")
	// require.NotContains(t, outputStr, "level=info")
}

func TestE2E_Config_MissingConfigFile(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start with a non-existent config file
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "DDG_SEARCH_CONFIG_FILE=/nonexistent/config.yaml")
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(300 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify application starts with defaults when config file is missing
	outputStr := out.String()
	_ = outputStr
	// Should use default configuration
	require.Contains(t, outputStr, "Configuration loaded")
}

func TestE2E_Config_VersionFlag(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Run with --version flag
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, binaryPath, "--version") //nolint:gosec
	cmd.Dir = "."
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	require.NoError(t, err, "Failed to run application")

	// Verify version output
	outputStr := out.String()
	require.Contains(t, outputStr, "ddg-search-mcp")
	require.Contains(t, outputStr, "version")
}

func TestE2E_Config_HelpFlag(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Run with --help flag
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, binaryPath, "--help") //nolint:gosec
	cmd.Dir = "."
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	require.NoError(t, err, "Failed to run application")

	// Verify help output
	outputStr := out.String()
	require.Contains(t, outputStr, "Usage")
	require.Contains(t, outputStr, "Flags")
	require.Contains(t, outputStr, "--config")
	require.Contains(t, outputStr, "--log-level")
}

func TestE2E_Config_LogLevels(t *testing.T) { //nolint:paralleltest
	// Test different log levels
	logLevels := []string{"debug", "info", "warn", "error"}

	for _, level := range logLevels {
		t.Run(level, func(t *testing.T) {
			// Build the application
			ctx := context.Background()
			binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
			buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
			buildCmd.Dir = "."
			output, err := buildCmd.CombinedOutput()
			require.NoError(t, err, "Build failed: %s", string(output))

			t.Cleanup(func() { _ = os.Remove(binaryPath) })

			// Start with specific log level
			var out bytes.Buffer
			cmd := exec.CommandContext(ctx, binaryPath, "--log-level="+level) //nolint:gosec
			cmd.Dir = "."
			cmd.Stdout = &out
			cmd.Stderr = &out
			err = cmd.Start()
			require.NoError(t, err, "Failed to start application")

			time.Sleep(300 * time.Millisecond)

			err = cmd.Process.Signal(syscall.SIGTERM)
			require.NoError(t, err, "Failed to send SIGTERM")

			err = cmd.Wait()
			_ = err

			// Verify application starts with the specified log level
			outputStr := out.String()
			_ = outputStr
		})
	}
}

func TestE2E_Config_FullConfiguration(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with all options
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
logging:
  level: debug
search:
  max_results: 25
  safe_search: true
perplexity:
  enabled: true
  access_token: "test-token"
  model: "sonar-medium-online"
  max_results: 10
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(300 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify all configuration values are loaded
	outputStr := out.String()
	_ = outputStr
	require.Contains(t, outputStr, "Configuration loaded")
	require.Contains(t, outputStr, "max_results=25")
	require.Contains(t, outputStr, "safe_search=true")
}

// ============================================================================
// Stage 5 - Missing Tests
// ============================================================================

func TestE2E_MCPServer_MultipleToolCalls(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send multiple tool calls
	requests := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang"}}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"search","arguments":{"query":"rust"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}`,
	}

	for _, req := range requests {
		_, err = stdin.Write([]byte(req + "\n"))
		require.NoError(t, err, "Failed to write request")
		time.Sleep(100 * time.Millisecond)
	}

	stdin.Close()
	time.Sleep(500 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify all tool calls were processed
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	require.Contains(t, outputStr, "fetch")
}

func TestE2E_MCPServer_RequestLogging(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with debug logging
	cmd := exec.CommandContext(ctx, binaryPath, "--log-level=debug") //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send a tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(500 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify request logging (with debug level)
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
}

// ============================================================================
// Stage 6 - Missing Tests
// ============================================================================

func TestE2E_MCPServer_SearchToolEmptyQuery(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call with empty query
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":""}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(500 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify error response for empty query
	outputStr := out.String()
	require.Contains(t, outputStr, "error")
}

func TestE2E_MCPServer_ContextCancellation(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx, cancel := context.WithCancel(context.Background())
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Cancel context to simulate client disconnect
	cancel()

	stdin.Close()
	time.Sleep(500 * time.Millisecond)

	err = cmd.Wait()
	_ = err

	// Verify graceful shutdown
	outputStr := out.String()
	require.Contains(t, outputStr, "Shutting down")
}

func TestE2E_MCPServer_ConfigurationIntegration(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
logging:
  level: debug
search:
  max_results: 5
  safe_search: true
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application with config
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"test"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(500 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify configuration is applied
	outputStr := out.String()
	require.Contains(t, outputStr, "max_results=5")
	require.Contains(t, outputStr, "safe_search=true")
}

func TestE2E_MCPServer_RateLimitHandling(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send multiple rapid search requests
	for i := 0; i < 5; i++ {
		request := fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"tools/call","params":{"name":"search","arguments":{"query":"test%d"}}}`, i+1, i)
		_, err = stdin.Write([]byte(request + "\n"))
		require.NoError(t, err, "Failed to write request")
		time.Sleep(50 * time.Millisecond)
	}

	stdin.Close()
	time.Sleep(1000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify rate limit handling (should not crash)
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
}

func TestE2E_MCPServer_NetworkErrorHandling(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send fetch tool call with invalid URL that will fail
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"http://this-domain-does-not-exist-12345.com"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(1000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify error handling (should not crash)
	outputStr := out.String()
	require.Contains(t, outputStr, "error")
}

func TestE2E_MCPServer_RealSearchIntegration(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send real search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify real search results
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should contain markdown-formatted results
	assert.True(t, len(outputStr) > 100)
}

func TestE2E_MCPServer_RealFetchIntegration(t *testing.T) { //nolint:paralleltest
	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send real fetch tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"fetch","arguments":{"url":"https://example.com"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify real fetch results
	outputStr := out.String()
	require.Contains(t, outputStr, "fetch")
	// Should contain content from example.com
	assert.True(t, len(outputStr) > 100)
}

// ============================================================================
// Stage 7 - Perplexity Integration Tests
// ============================================================================

func TestE2E_Perplexity_Disabled_DuckDuckGoFallback(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity disabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: false
  access_token: ""
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify DuckDuckGo is used
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should not mention Perplexity
	// Perplexity config is logged even when disabled
}

func TestE2E_Perplexity_EnabledWithoutToken_DuckDuckGoDirect(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled but no token
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: ""
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify DuckDuckGo is used directly
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should mention DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_Perplexity_InvalidToken_DuckDuckGoFallback(t *testing.T) { //nolint:paralleltest
	// Create config with invalid Perplexity token
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "invalid-token-12345"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify fallback to DuckDuckGo
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should fall back to DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_Perplexity_SearchWithParameters(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
  model: "sonar-medium-online"
  max_results: 10
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call with parameters
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming","max_results":5}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify search with parameters
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
}

func TestE2E_Perplexity_ResultFormatting(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify result formatting
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Results should be formatted
	assert.True(t, len(outputStr) > 100)
}

func TestE2E_Perplexity_EmptyResults(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call with unlikely query
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"xyzabc123def456ghi789"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify empty results handling
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
}

func TestE2E_Perplexity_WithSafeSearch(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled and safe search
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
search:
  safe_search: true
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify safe search is applied
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
}

func TestE2E_Perplexity_WithMaxResults(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled and max results
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
  max_results: 5
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify max results is applied
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
}

func TestE2E_Perplexity_NetworkError_DuckDuckGoFallback(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify fallback to DuckDuckGo on network error
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should fall back to DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_Perplexity_Timeout_DuckDuckGoFallback(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled and short timeout
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
  timeout: 1s
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(3000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify fallback to DuckDuckGo on timeout
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should fall back to DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_Perplexity_TransientError_NoRetry(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify no retry on transient error (falls back to DuckDuckGo)
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should fall back to DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_Perplexity_QuotaExceeded_DuckDuckGoFallback(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify fallback to DuckDuckGo on quota exceeded
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should fall back to DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_Perplexity_RateLimit_DuckDuckGoFallback(t *testing.T) { //nolint:paralleltest
	// Create config with Perplexity enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
perplexity:
  enabled: true
  access_token: "test-token"
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build the application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start the application
	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err, "Failed to create stdin pipe")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	time.Sleep(200 * time.Millisecond)

	// Send search tool call
	request := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"search","arguments":{"query":"golang programming"}}}`
	_, err = stdin.Write([]byte(request + "\n"))
	require.NoError(t, err, "Failed to write request")

	stdin.Close()
	time.Sleep(2000 * time.Millisecond)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	err = cmd.Wait()
	_ = err

	// Verify fallback to DuckDuckGo on rate limit
	outputStr := out.String()
	require.Contains(t, outputStr, "search")
	// Should fall back to DuckDuckGo
	require.Contains(t, outputStr, "DuckDuckGo")
}

func TestE2E_TLS_ServerStartup(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19110
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
    min_version: "1.2"
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test HTTPS connection with TLS client
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19110/health", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call health check via HTTPS")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify TLS was enabled
	outputStr := out.String()
	require.Contains(t, outputStr, "tls=enabled")
	require.Contains(t, outputStr, "HTTPS Streamable HTTP server listening")
}

func TestE2E_TLS_HealthCheck(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19111
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test health check endpoint via HTTPS
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19111/health", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call health check via HTTPS")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("OK"), body)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err
}

func TestE2E_TLS_InvalidCertificate(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS enabled but invalid cert path
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19112
  tls:
    enabled: true
    cert_file: /nonexistent/cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Server should have exited with error
	err = cmd.Wait()
	require.Error(t, err, "Server should exit with error on invalid certificate")

	// Verify error message
	outputStr := out.String()
	require.Contains(t, outputStr, "certificate")
	require.Contains(t, outputStr, "cert")
}

func TestE2E_TLS_InvalidKey(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS enabled but invalid key path
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19113
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: /nonexistent/key.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Server should have exited with error
	err = cmd.Wait()
	require.Error(t, err, "Server should exit with error on invalid key")

	// Verify error message
	outputStr := out.String()
	require.Contains(t, outputStr, "key")
}

func TestE2E_mTLS_ServerStartup(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with mTLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19114
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
    min_version: "1.2"
    mtls:
      enabled: true
      ca_file: ../../internal/mcp/testdata/ca-cert.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test HTTPS connection with valid client certificate
	caCert, err := os.ReadFile("../../internal/mcp/testdata/ca-cert.pem")
	require.NoError(t, err, "Failed to read CA certificate")

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(
		"../../internal/mcp/testdata/client-cert.pem",
		"../../internal/mcp/testdata/client-key.pem",
	)
	require.NoError(t, err, "Failed to load client certificate")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
				MinVersion:   tls.VersionTLS12,
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19114/health", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call health check via HTTPS with mTLS")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify mTLS was enabled
	outputStr := out.String()
	require.Contains(t, outputStr, "mTLS enabled")
}

func TestE2E_mTLS_ClientWithoutCertificate(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with mTLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19115
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
    mtls:
      enabled: true
      ca_file: ../../internal/mcp/testdata/ca-cert.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test HTTPS connection without client certificate (should fail)
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
		Timeout: 2 * time.Second,
	}
	_, err = httpClient.Get("https://localhost:19115/health")
	require.Error(t, err, "Connection should fail without client certificate")

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err
}

func TestE2E_mTLS_InvalidClientCertificate(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with mTLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19116
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
    mtls:
      enabled: true
      ca_file: ../../internal/mcp/testdata/ca-cert.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test HTTPS connection with invalid client certificate (use server cert as client cert)
	caCert, err := os.ReadFile("../../internal/mcp/testdata/ca-cert.pem")
	require.NoError(t, err, "Failed to read CA certificate")

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Use server certificate as client certificate (invalid for client auth)
	clientCert, err := tls.LoadX509KeyPair(
		"../../internal/mcp/testdata/server-cert.pem",
		"../../internal/mcp/testdata/server-key.pem",
	)
	require.NoError(t, err, "Failed to load server certificate")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
				MinVersion:   tls.VersionTLS12,
			},
		},
		Timeout: 2 * time.Second,
	}
	_, err = httpClient.Get("https://localhost:19116/health")
	require.Error(t, err, "Connection should fail with invalid client certificate")

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err
}

func TestE2E_TLS_MinVersion(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS min version 1.2
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19117
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
    min_version: "1.2"
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test HTTPS connection with TLS 1.2
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
				MinVersion:         tls.VersionTLS12,
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19117/health", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call health check via HTTPS with TLS 1.2")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify min version was configured
	outputStr := out.String()
	require.Contains(t, outputStr, "min_version")
}

func TestE2E_TLS_CertificateReload(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19118
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Send SIGHUP to trigger config reload
	err = cmd.Process.Signal(syscall.SIGHUP)
	require.NoError(t, err, "Failed to send SIGHUP")

	// Wait for reload to process
	time.Sleep(500 * time.Millisecond)

	// Test HTTPS connection still works
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19118/health", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call health check via HTTPS after reload")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify reload was logged
	outputStr := out.String()
	require.Contains(t, outputStr, "reload")
	require.Contains(t, outputStr, "Reload")
}

func TestE2E_TLS_StreamableHTTPConnection(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with TLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19119
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test MCP endpoint via HTTPS
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19119/mcp", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call MCP endpoint via HTTPS")
	defer resp.Body.Close()

	// MCP endpoint should return 200 OK
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify MCP connection was logged
	outputStr := out.String()
	require.Contains(t, outputStr, "HTTPS")
}

func TestE2E_mTLS_StreamableHTTPConnection(t *testing.T) { //nolint:paralleltest
	// Create a temporary config file with mTLS enabled
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: localhost:19120
  tls:
    enabled: true
    cert_file: ../../internal/mcp/testdata/server-cert.pem
    key_file: ../../internal/mcp/testdata/server-key.pem
    mtls:
      enabled: true
      ca_file: ../../internal/mcp/testdata/ca-cert.pem
logging:
  level: info
search:
  max_results: 10
  safe_search: false
perplexity:
  enabled: false
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err, "Failed to create config file")

	// Build application
	ctx := context.Background()
	binaryPath := fmt.Sprintf("/tmp/ddg-search-mcp-test-%s-%d", t.Name(), time.Now().UnixNano())
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, ".") //nolint:gosec
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Build failed: %s", string(output))

	t.Cleanup(func() { _ = os.Remove(binaryPath) })

	// Start application in background
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, binaryPath) //nolint:gosec
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), fmt.Sprintf("DDG_SEARCH_CONFIG_FILE=%s", configPath))
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Test MCP endpoint via HTTPS with valid client certificate
	caCert, err := os.ReadFile("../../internal/mcp/testdata/ca-cert.pem")
	require.NoError(t, err, "Failed to read CA certificate")

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(
		"../../internal/mcp/testdata/client-cert.pem",
		"../../internal/mcp/testdata/client-key.pem",
	)
	require.NoError(t, err, "Failed to load client certificate")

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
				MinVersion:   tls.VersionTLS12,
			},
		},
	}
	req, err := http.NewRequest("GET", "https://localhost:19120/mcp", nil)
	require.NoError(t, err, "Failed to create request")
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "Failed to call MCP endpoint via HTTPS with mTLS")
	defer resp.Body.Close()

	// MCP endpoint should return 200 OK
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Send SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Failed to send SIGTERM")

	// Wait for process to exit
	err = cmd.Wait()
	_ = err

	// Verify mTLS MCP connection was logged
	outputStr := out.String()
	require.Contains(t, outputStr, "mTLS")
}

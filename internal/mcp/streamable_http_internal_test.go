package mcp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/mcplog"
	"github.com/stretchr/testify/require"
)

const (
	errContextCanceled = "MCP server error: context canceled"
)

func TestServer_serveHTTP(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "http", "localhost:0")
		serverErr <- err
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Cancel to stop server
	cancel()

	select {
	case err := <-serverErr:
		// nil is expected when context is canceled
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for server")
	}
}

func TestServer_serveStdio(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start stdio server in a goroutine
	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "stdio", "")
		if err != nil {
			serverErr <- err
		}
	}()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Cancel context
	cancel()

	// Stdio server should stop
	select {
	case err := <-serverErr:
		// Context canceled is expected
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(1 * time.Second):
		// OK, server stopped
	}
}

func TestServer_healthCheckHandler(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set shutdown context
	server.shutdownCtx = ctx

	// Test health check when server is running
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.healthCheckHandler(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, []byte("OK"), body)

	// Test health check when server is shutting down
	cancel()

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	w = httptest.NewRecorder()

	server.healthCheckHandler(w, req)

	resp = w.Result()
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	require.Equal(t, []byte("Service Unavailable"), body)
}

func TestServer_loggingMiddleware(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	// Wrap with logging middleware
	middleware := server.loggingMiddleware(testHandler)

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Referer", "http://example.com")
	req.Header.Set("User-Agent", "test-agent/1.0")

	w := httptest.NewRecorder()

	// Call middleware
	middleware.ServeHTTP(w, req)

	// Verify response
	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, []byte("test response"), body)
}

func TestLoggingResponseWriter_Write(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	lrw := &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Write data
	data := []byte("test data")
	n, err := lrw.Write(data)
	require.NoError(t, err)
	require.Equal(t, len(data), n)
	require.Equal(t, len(data), lrw.bytesWritten)

	// Write more data
	moreData := []byte(" more data")
	n, err = lrw.Write(moreData)
	require.NoError(t, err)
	require.Equal(t, len(moreData), n)
	require.Equal(t, len(data)+len(moreData), lrw.bytesWritten)

	// Verify response
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, append(data, moreData...), body)
}

func TestLoggingResponseWriter_WriteHeader(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	lrw := &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Write header
	lrw.WriteHeader(http.StatusCreated)
	require.Equal(t, http.StatusCreated, lrw.statusCode)

	// Verify response
	resp := w.Result()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestLoggingResponseWriter_DefaultStatusCode(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	lrw := &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	// Don't call WriteHeader, should default to OK
	data := []byte("test")
	_, err := lrw.Write(data)
	require.NoError(t, err)

	// Status should still be OK (default)
	require.Equal(t, http.StatusOK, lrw.statusCode)
}

func TestServer_UnsupportedProtocol(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx := t.Context()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "unsupported", "localhost:0")
		if err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported protocol")
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for error")
	}
}

func TestServer_HTTPServerCreation(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "http", "localhost:0")
		serverErr <- err
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Verify httpServer is created
	require.NotNil(t, server.GetHTTPServer())

	// Cancel to stop server
	cancel()

	select {
	case err := <-serverErr:
		// Context canceled is expected
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for server")
	}
}

func TestServer_StreamableHTTPServerCreation(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "http", "localhost:0")
		serverErr <- err
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Verify streamableHTTPServer is created
	require.NotNil(t, server.GetStreamableHTTPServer())

	// Cancel to stop server
	cancel()

	select {
	case err := <-serverErr:
		// Context canceled is expected
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for server")
	}
}

func TestServer_Shutdown_HTTP(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "http", "localhost:0")
		serverErr <- err
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Shutdown server
	shutdownErr := make(chan error, 1)

	go func() {
		shutdownErr <- server.Shutdown(ctx)
	}()

	select {
	case err := <-shutdownErr:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for shutdown")
	}

	// Cancel context
	cancel()

	// Server should have stopped
	select {
	case err := <-serverErr:
		// Context canceled is expected
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(1 * time.Second):
		// OK, server stopped
	}
}

func TestServer_Shutdown_Stdio(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	ctx, cancel := context.WithCancel(context.Background())

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, "stdio", "")
		if err != nil {
			serverErr <- err
		}
	}()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Shutdown server
	shutdownErr := make(chan error, 1)

	go func() {
		shutdownErr <- server.Shutdown(ctx)
	}()

	select {
	case err := <-shutdownErr:
		require.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for shutdown")
	}

	// Cancel context
	cancel()

	// Server should have stopped
	select {
	case err := <-serverErr:
		// Context canceled is expected
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(1 * time.Second):
		// OK, server stopped
	}
}

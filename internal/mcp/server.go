// Package mcp provides the MCP server implementation and tool handlers.
package mcp

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	// shutdownTimeout is the timeout for graceful shutdown.
	shutdownTimeout = 5 * time.Second
	// readHeaderTimeout is the timeout for reading HTTP headers.
	readHeaderTimeout = 10 * time.Second
)

var (
	// ErrUnsupportedProtocol is returned when an unsupported protocol is specified.
	ErrUnsupportedProtocol = errors.New("unsupported protocol")
	// ErrMTLSCANotSpecified is returned when mTLS is enabled but CA file is not specified.
	ErrMTLSCANotSpecified = errors.New("mTLS enabled but CA file not specified")
	// ErrFailedToParseCACertificate is returned when CA certificate parsing fails.
	ErrFailedToParseCACertificate = errors.New("failed to parse CA certificate")
	// ErrResponseWriterNoHijacking is returned when response writer doesn't support hijacking.
	ErrResponseWriterNoHijacking = errors.New("response writer does not support hijacking")
	// ErrTLSConfigNoGetCertFile is returned when TLS config does not implement GetCertFile.
	ErrTLSConfigNoGetCertFile = errors.New("TLS config does not implement GetCertFile")
	// ErrTLSConfigNoGetKeyFile is returned when TLS config does not implement GetKeyFile.
	ErrTLSConfigNoGetKeyFile = errors.New("TLS config does not implement GetKeyFile")
	// ErrCAPathNotAbsolute is returned when CA certificate path is not absolute.
	ErrCAPathNotAbsolute = errors.New("CA certificate path must be absolute")
)

// Server wraps the MCP server with additional functionality.
type Server struct {
	mcpServer            *server.MCPServer
	streamableHTTPServer *server.StreamableHTTPServer
	logger               *slog.Logger
	config               *Config
	appConfig            any // Holds full application configuration (mcpconfig.Config)
	httpServer           *http.Server
	shutdownCtx          context.Context //nolint:containedctx // Stored for health check during shutdown
	mu                   sync.RWMutex    // Protects httpServer and streamableHTTPServer fields
}

// Config holds the MCP server configuration.
type Config struct {
	Name    string
	Version string
}

// NewServer creates a new MCP server instance.
func NewServer(cfg *Config, logger *slog.Logger) *Server {
	mcpServer := server.NewMCPServer(
		cfg.Name,
		cfg.Version,
		server.WithToolCapabilities(true),
	)

	return &Server{
		mcpServer: mcpServer,
		logger:    logger,
		config:    cfg,
		appConfig: nil, // Will be set via SetAppConfig
	}
}

// SetAppConfig sets the full application configuration.
func (s *Server) SetAppConfig(appConfig any) {
	s.appConfig = appConfig
}

// GetHTTPServer returns the HTTP server (for testing purposes).
func (s *Server) GetHTTPServer() *http.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.httpServer
}

// GetStreamableHTTPServer returns the Streamable HTTP server (for testing purposes).
func (s *Server) GetStreamableHTTPServer() *server.StreamableHTTPServer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.streamableHTTPServer
}

// RegisterTool registers a tool with the server.
func (s *Server) RegisterTool(tool mcp.Tool, handler ToolHandler) {
	s.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		startTime := time.Now()
		toolName := request.Params.Name

		// Extract arguments as map
		var args map[string]any

		if request.Params.Arguments != nil {
			if argsMap, ok := request.Params.Arguments.(map[string]any); ok {
				args = argsMap
			} else {
				// Try to unmarshal if it's a JSON string
				if jsonBytes, ok := request.Params.Arguments.(string); ok {
					_ = json.Unmarshal([]byte(jsonBytes), &args)
				}
			}
		}

		// Log the tool call
		s.logToolCall(toolName, args, startTime, 0, nil)

		// Call the handler
		result, err := handler(ctx, s.appConfig, request)

		// Log the result
		status := 200
		if err != nil {
			status = 500
		}

		s.logToolCall(toolName, args, startTime, status, err)

		return result, err
	})
}

// ToolHandler is a function that handles tool calls.
type ToolHandler func(
	ctx context.Context,
	appConfig any,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error)

// Serve starts the MCP server with the configured transport (stdio or HTTP).
func (s *Server) Serve(ctx context.Context, protocol, bindAddress string) error {
	s.shutdownCtx = ctx
	s.logger.Info("Starting MCP server", "name", s.config.Name, "version", s.config.Version, "protocol", protocol)

	switch protocol {
	case "stdio":
		return s.serveStdio()
	case "http":
		return s.serveHTTP(ctx, bindAddress)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedProtocol, protocol)
	}
}

// ReloadTLS reloads TLS certificates and returns an error if the new configuration is invalid.
func (s *Server) ReloadTLS() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.httpServer == nil {
		return nil // No HTTP server running, nothing to reload
	}

	// Get current TLS config from appConfig
	cfg, ok := s.appConfig.(interface{ GetTLSConfig() any })
	if !ok {
		return nil // No TLS config available
	}

	tlsConfig := cfg.GetTLSConfig()
	if tlsConfig == nil {
		return nil // TLS not enabled
	}

	// Check if TLS is enabled
	tlsFields, ok := tlsConfig.(interface{ GetEnabled() bool })
	if !ok {
		return nil // TLS config not in expected format
	}

	if !tlsFields.GetEnabled() {
		return nil // TLS disabled, nothing to reload
	}

	// Load new certificates
	certGetter, ok := tlsConfig.(interface{ GetCertFile() string })
	if !ok {
		return ErrTLSConfigNoGetCertFile
	}

	keyGetter, ok := tlsConfig.(interface{ GetKeyFile() string })
	if !ok {
		return ErrTLSConfigNoGetKeyFile
	}

	cert, err := tls.LoadX509KeyPair(certGetter.GetCertFile(), keyGetter.GetKeyFile())
	if err != nil {
		return fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	// Update server TLS config
	s.httpServer.TLSConfig.Certificates = []tls.Certificate{cert}
	s.logger.Info("TLS certificates reloaded")

	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down MCP server")

	// Shutdown HTTP server if running
	if s.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("HTTP server shutdown error: %w", err)
		}
	}

	// Shutdown Streamable HTTP server if running
	if s.streamableHTTPServer != nil {
		err := s.streamableHTTPServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("streamable HTTP server shutdown error: %w", err)
		}
	}

	return nil
}

// serveStdio starts the MCP server with stdio transport.
func (s *Server) serveStdio() error {
	err := server.ServeStdio(s.mcpServer)
	if err != nil {
		return fmt.Errorf("MCP server error: %w", err)
	}

	return nil
}

// startHTTPServer starts the HTTP server in a goroutine and returns the error channel.
func (s *Server) startHTTPServer(ctx context.Context, tlsConfig *tls.Config, bindAddress string) chan error {
	serverErr := make(chan error, 1)

	go func() {
		var err error

		if tlsConfig != nil {
			// Create a custom listener to log connection attempts
			lc := &net.ListenConfig{}

			listener, err := lc.Listen(ctx, "tcp", bindAddress)
			if err != nil {
				serverErr <- fmt.Errorf("failed to create listener: %w", err)

				return
			}

			// Update server address with actual bound address (protected by mutex)
			s.mu.Lock()
			s.httpServer.Addr = listener.Addr().String()
			s.mu.Unlock()

			// Wrap listener with logging
			loggingListener := &loggingListener{
				Listener: listener,
				logger:   s.logger,
			}

			s.logger.Info("Starting HTTPS server with custom listener")
			_ = s.httpServer.ServeTLS(loggingListener, "", "") // Certificates loaded from TLSConfig
		} else {
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	return serverErr
}

// logServerStartup logs the server startup information.
func (s *Server) logServerStartup(tlsConfig *tls.Config, bindAddress string) {
	if tlsConfig != nil {
		mtlsEnabled := tlsConfig.ClientAuth == tls.RequireAndVerifyClientCert
		s.logger.Info("HTTPS Streamable HTTP server listening",
			"address", bindAddress,
			"endpoint", "/mcp",
			"tls", "enabled",
			"min_version", tlsConfig.MinVersion,
			"mtls", mtlsEnabled,
		)
	} else {
		s.logger.Info("HTTP Streamable HTTP server listening",
			"address", bindAddress,
			"endpoint", "/mcp",
		)
	}
}

// serveHTTP starts the MCP server with Streamable HTTP transport.
func (s *Server) serveHTTP(ctx context.Context, bindAddress string) error {
	// Get TLS configuration from appConfig
	tlsConfig, err := s.buildTLSConfig()
	if err != nil {
		return fmt.Errorf("failed to build TLS config: %w", err)
	}

	// Create Streamable HTTP server
	s.mu.Lock()
	s.streamableHTTPServer = server.NewStreamableHTTPServer(s.mcpServer)

	// Create HTTP server with logging middleware
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.healthCheckHandler)

	// MCP endpoint (handles both POST and GET methods)
	mux.Handle("/mcp", s.streamableHTTPServer)

	// Wrap with logging middleware
	handler := s.loggingMiddleware(mux)

	s.httpServer = &http.Server{
		Addr:              bindAddress,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		TLSConfig:         tlsConfig,
	}
	s.mu.Unlock()

	// Log server startup
	s.logServerStartup(tlsConfig, bindAddress)

	// Start server in a goroutine
	serverErr := s.startHTTPServer(ctx, tlsConfig, bindAddress)

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		s.logger.Info("Shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		//nolint:contextcheck // Using background context for shutdown timeout as parent context is canceled
		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("HTTP server shutdown error: %w", err)
		}

		return nil
	case err := <-serverErr:
		return err
	}
}

// configureMTLS configures mutual TLS (mTLS) for the server.
func (s *Server) configureMTLS(config *tls.Config, caFile string) error {
	// Validate CA file path is absolute
	if !filepath.IsAbs(caFile) {
		s.logger.Error("mTLS CA certificate path must be absolute", "path", caFile)

		return fmt.Errorf("%w: %s", ErrCAPathNotAbsolute, caFile)
	}

	// Load CA certificate
	//nolint:gosec // G304 - path is validated as absolute above
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		s.logger.Error("failed to read mTLS CA certificate", "path", caFile, "error", err)

		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		s.logger.Error("failed to parse mTLS CA certificate", "path", caFile)

		return ErrFailedToParseCACertificate
	}

	// Configure client certificate verification
	config.ClientCAs = caCertPool
	config.ClientAuth = tls.RequireAndVerifyClientCert

	s.logger.Info("mTLS enabled", "ca_file", caFile)

	return nil
}

// configureTLSLogging sets up TLS handshake and connection logging callbacks.
func (s *Server) configureTLSLogging(config *tls.Config) {
	// Add TLS handshake logging callback
	// Note: This is set after config creation to avoid circular reference
	tlsConfig := config
	config.GetConfigForClient = func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
		s.logger.Debug("TLS handshake attempt",
			"client", hello.Conn.RemoteAddr(),
			"server_name", hello.ServerName,
			"supported_versions", hello.SupportedVersions,
			"cipher_suites", hello.CipherSuites,
			"min_version", tlsConfig.MinVersion,
			"max_version", tlsConfig.MaxVersion,
		)

		return tlsConfig, nil
	}

	// Add connection state callback to log successful TLS connections
	config.VerifyConnection = func(state tls.ConnectionState) error {
		s.logger.Info("TLS connection established",
			"version", tlsVersionToString(state.Version),
			"cipher", tls.CipherSuiteName(state.CipherSuite),
			"peer_certs", len(state.PeerCertificates),
		)

		return nil
	}
}

// buildTLSConfig builds TLS configuration from appConfig.
func (s *Server) buildTLSConfig() (*tls.Config, error) {
	// Extract TLS config from appConfig
	type tlsConfigGetter interface {
		GetTLSConfig() any
	}

	cfgGetter, ok := s.appConfig.(tlsConfigGetter)
	if !ok {
		return nil, nil //nolint:nilnil // No TLS config available
	}

	tlsCfg := cfgGetter.GetTLSConfig()
	if tlsCfg == nil {
		return nil, nil //nolint:nilnil // TLS not enabled
	}

	// Extract TLS fields
	type tlsFields interface {
		GetEnabled() bool
		GetCertFile() string
		GetKeyFile() string
		GetMinVersion() string
		GetMTLSEnabled() bool
		GetMTLSCAFile() string
	}

	fields, ok := tlsCfg.(tlsFields)
	if !ok {
		return nil, nil //nolint:nilnil // TLS config not in expected format
	}

	if !fields.GetEnabled() {
		return nil, nil //nolint:nilnil // TLS disabled
	}

	// Load certificate and key
	cert, err := tls.LoadX509KeyPair(fields.GetCertFile(), fields.GetKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate and key: %w", err)
	}

	// Build TLS config
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Minimum TLS 1.2
	}

	// Configure TLS logging
	s.configureTLSLogging(config)

	// Configure mTLS if enabled
	if fields.GetMTLSEnabled() {
		caFile := fields.GetMTLSCAFile()
		if caFile == "" {
			s.logger.Error("mTLS enabled but CA file not specified")

			return nil, ErrMTLSCANotSpecified
		}

		err := s.configureMTLS(config, caFile)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// healthCheckHandler handles health check requests.
func (s *Server) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	select {
	case <-s.shutdownCtx.Done():
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("Service Unavailable"))
	default:
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}

// loggingMiddleware logs HTTP requests in Combined Log Format.
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log TLS connection details if available
		if r.TLS != nil {
			state := r.TLS
			s.logger.Debug("TLS connection established",
				"version", tlsVersionToString(state.Version),
				"cipher", tls.CipherSuiteName(state.CipherSuite),
				"client", r.RemoteAddr,
			)

			// Log client certificate for mTLS
			if len(state.PeerCertificates) > 0 {
				cert := state.PeerCertificates[0]
				s.logger.Debug("mTLS client certificate",
					"subject", cert.Subject.String(),
					"issuer", cert.Issuer.String(),
					"client", r.RemoteAddr,
				)
			}
		}

		// Wrap response writer to capture status code and bytes written
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(lrw, r)

		// Log the request
		timestamp := startTime.Format("02/Jan/2006:15:04:05 -0700")
		client := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		status := lrw.statusCode
		bytes := lrw.bytesWritten

		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "-"
		}

		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			userAgent = "-"
		}

		logLine := fmt.Sprintf("[%s] [%s] [%s] [%s] [%d] [%d] [%s] [%s]",
			timestamp, client, method, path, status, bytes, referer, userAgent)

		s.logger.Debug(logLine)
	})
}

// tlsVersionToString converts TLS version constant to string.
func tlsVersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", version)
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code and bytes written.
type loggingResponseWriter struct {
	http.ResponseWriter

	statusCode   int
	bytesWritten int
}

// WriteHeader captures the status code.
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the number of bytes written.
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += n

	return n, err //nolint:wrapcheck // Propagating Write error from underlying ResponseWriter
}

// Flush implements http.Flusher interface for SSE support.
func (lrw *loggingResponseWriter) Flush() {
	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements http.Hijacker interface for WebSocket support.
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := lrw.ResponseWriter.(http.Hijacker); ok {
		conn, rw, err := hijacker.Hijack()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to hijack connection: %w", err)
		}

		return conn, rw, nil
	}

	return nil, nil, ErrResponseWriterNoHijacking
}

// logToolCall logs a tool call in Combined Log Format.
func (s *Server) logToolCall(toolName string, arguments map[string]any, startTime time.Time, status int, err error) {
	argsJSON, jsonErr := json.Marshal(arguments)
	if jsonErr != nil {
		argsJSON = []byte("{}")
	}

	bytes := len(argsJSON)

	timestamp := startTime.Format("02/Jan/2006:15:04:05 -0700")
	client := "stdio"
	method := "CALL"
	path := toolName
	referer := "-"
	userAgent := "mcp-client"

	logLine := fmt.Sprintf("[%s] [%s] [%s] [%s] [%d] [%d] [%s] [%s]",
		timestamp, client, method, path, status, bytes, referer, userAgent)

	if err != nil {
		s.logger.Debug(logLine, "error", err.Error())
	} else {
		s.logger.Debug(logLine)
	}
}

// loggingListener wraps net.Listener to log connection attempts.
type loggingListener struct {
	net.Listener

	logger *slog.Logger
}

// Accept logs incoming connections before accepting them.
func (ll *loggingListener) Accept() (net.Conn, error) {
	conn, err := ll.Listener.Accept()
	if err != nil {
		ll.logger.Error("Failed to accept connection", "error", err)

		return nil, fmt.Errorf("accept connection: %w", err)
	}

	ll.logger.Debug("Connection accepted", "remote_addr", conn.RemoteAddr(), "local_addr", conn.LocalAddr())

	// Wrap connection to capture TLS errors
	return &loggingConn{
		Conn:   conn,
		logger: ll.logger,
	}, nil
}

// loggingConn wraps net.Conn to log TLS errors.
type loggingConn struct {
	net.Conn

	logger *slog.Logger
}

// Close logs when connection is closed.
func (lc *loggingConn) Close() error {
	lc.logger.Debug("Connection closed", "remote_addr", lc.RemoteAddr())

	err := lc.Conn.Close()
	if err != nil {
		return fmt.Errorf("close connection: %w", err)
	}

	return nil
}

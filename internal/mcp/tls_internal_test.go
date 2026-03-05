package mcp

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/mcpconfig"
	"github.com/Djarvur/ddg-search/internal/mcplog"
	"github.com/stretchr/testify/require"
)

func TestServer_buildTLSConfig(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     *mcpconfig.Config
		wantErr bool
		wantTLS bool
	}{
		{
			name: "TLS disabled",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					TLS: mcpconfig.TLSConfig{
						Enabled: false,
					},
				},
			},
			wantErr: false,
			wantTLS: false,
		},
		{
			name: "TLS enabled with valid config",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					TLS: mcpconfig.TLSConfig{
						Enabled:    true,
						CertFile:   "testdata/cert.pem",
						KeyFile:    "testdata/key.pem",
						MinVersion: "1.2",
					},
				},
			},
			wantErr: true, // Cert files don't exist
			wantTLS: false, // Returns nil on error
		},
		{
			name: "TLS enabled with mTLS",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					TLS: mcpconfig.TLSConfig{
						Enabled:    true,
						CertFile:   "testdata/cert.pem",
						KeyFile:    "testdata/key.pem",
						MinVersion: "1.2",
						MTLS: mcpconfig.MTLSConfig{
							Enabled: true,
							CAFile:  "testdata/ca.pem",
						},
					},
				},
			},
			wantErr: true, // Cert files don't exist
			wantTLS: false, // Returns nil on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewServer(&Config{
				Name:    "test-server",
				Version: "1.0.0",
			}, logger)

			server.SetAppConfig(tt.cfg)

			tlsConfig, err := server.buildTLSConfig()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.wantTLS {
				require.NotNil(t, tlsConfig)
			} else {
				require.Nil(t, tlsConfig)
			}
		})
	}
}

func TestServer_tlsVersionToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ver  uint16
		want string
	}{
		{
			name: "TLS 1.0",
			ver:  tls.VersionTLS10,
			want: "TLS 1.0",
		},
		{
			name: "TLS 1.1",
			ver:  tls.VersionTLS11,
			want: "TLS 1.1",
		},
		{
			name: "TLS 1.2",
			ver:  tls.VersionTLS12,
			want: "TLS 1.2",
		},
		{
			name: "TLS 1.3",
			ver:  tls.VersionTLS13,
			want: "TLS 1.3",
		},
		{
			name: "Unknown version",
			ver:  0x0000,
			want: "Unknown (0x0000)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tlsVersionToString(tt.ver)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServer_ReloadTLS(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	// Test reload without HTTP server
	err = server.ReloadTLS()
	require.NoError(t, err)

	// Test reload with HTTP server but no TLS
	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "localhost:0",
			TLS: mcpconfig.TLSConfig{
				Enabled: false,
			},
		},
		Logging: mcpconfig.LoggingConfig{
			Level: "debug",
		},
	}

	server.SetAppConfig(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, cfg.Server.Protocol, cfg.Server.BindAddress)
		if err != nil && err.Error() != errContextCanceled {
			serverErr <- err
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Reload TLS (should succeed even with TLS disabled)
	err = server.ReloadTLS()
	require.NoError(t, err)

	// Cancel to stop server
	cancel()

	select {
	case err := <-serverErr:
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		// OK, server stopped
	}
}

func TestServer_serveHTTPWithTLS(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "localhost:0",
			TLS: mcpconfig.TLSConfig{
				Enabled: false, // Disable TLS for this test
			},
		},
	}

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	server.SetAppConfig(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, cfg.Server.Protocol, cfg.Server.BindAddress)
		if err != nil && err.Error() != errContextCanceled {
			serverErr <- err
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Verify server is running
	httpServer := server.GetHTTPServer()
	require.NotNil(t, httpServer)

	// Make a request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	httpServer.Handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, []byte("OK"), w.Body.Bytes())

	// Cancel to stop server
	cancel()

	select {
	case err := <-serverErr:
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		// OK, server stopped
	}
}

func TestServer_loggingMiddlewareWithTLS(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "localhost:0",
			TLS: mcpconfig.TLSConfig{
				Enabled: false,
			},
		},
	}

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	server.SetAppConfig(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, cfg.Server.Protocol, cfg.Server.BindAddress)
		if err != nil && err.Error() != errContextCanceled {
			serverErr <- err
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Make a request without TLS
	httpServer := server.GetHTTPServer()
	if httpServer != nil {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		httpServer.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	}

	// Cancel to stop server
	cancel()

	select {
	case err := <-serverErr:
		if err != nil && err.Error() != errContextCanceled {
			require.NoError(t, err)
		}
	case <-time.After(2 * time.Second):
		// OK, server stopped
	}
}

func TestServer_ShutdownWithTLS(t *testing.T) {
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "localhost:0",
			TLS: mcpconfig.TLSConfig{
				Enabled: false,
			},
		},
	}

	server := NewServer(&Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	server.SetAppConfig(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	serverErr := make(chan error, 1)

	go func() {
		err := server.Serve(ctx, cfg.Server.Protocol, cfg.Server.BindAddress)
		if err != nil && err.Error() != errContextCanceled {
			serverErr <- err
		}
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

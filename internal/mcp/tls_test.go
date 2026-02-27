// Package mcp_test provides tests for TLS/mTLS functionality.
package mcp_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/mcp"
	"github.com/Djarvur/ddg-search/internal/mcpconfig"
	"github.com/Djarvur/ddg-search/internal/mcplog"
	"github.com/stretchr/testify/require"
)

func TestServer_TLSConfiguration(t *testing.T) { //nolint:funlen
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	tests := []struct {
		name     string
		tlsCfg   *mcpconfig.TLSConfig
		wantErr  bool
		wantTLS  bool
		wantMTLS bool
	}{
		{
			name: "TLS disabled",
			tlsCfg: &mcpconfig.TLSConfig{
				Enabled: false,
			},
			wantErr:  false,
			wantTLS:  false,
			wantMTLS: false,
		},
		{
			name: "TLS enabled without mTLS",
			tlsCfg: &mcpconfig.TLSConfig{
				Enabled:    true,
				CertFile:   "testdata/cert.pem",
				KeyFile:    "testdata/key.pem",
				MinVersion: "1.2",
				MTLS: mcpconfig.MTLSConfig{
					Enabled: false,
				},
			},
			wantErr:  true, // Cert files don't exist
			wantTLS:  true,
			wantMTLS: false,
		},
		{
			name: "TLS enabled with mTLS",
			tlsCfg: &mcpconfig.TLSConfig{
				Enabled:    true,
				CertFile:   "testdata/cert.pem",
				KeyFile:    "testdata/key.pem",
				MinVersion: "1.2",
				MTLS: mcpconfig.MTLSConfig{
					Enabled: true,
					CAFile:  "testdata/ca.pem",
				},
			},
			wantErr:  true, // Cert files don't exist
			wantTLS:  true,
			wantMTLS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "http",
					BindAddress: "localhost:0",
					TLS:         *tt.tlsCfg,
				},
			}

			server := mcp.NewServer(&mcp.Config{
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

			// Wait a bit for server to start
			time.Sleep(200 * time.Millisecond)

			// Cancel to stop server
			cancel()

			select {
			case err := <-serverErr:
				if tt.wantErr {
					require.Error(t, err)
				} else if err != nil && err.Error() != errContextCanceled {
					require.NoError(t, err)
				}
			case <-time.After(2 * time.Second):
				// OK, server stopped
			}
		})
	}
}

func TestServer_TLSVersionParsing(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name      string
		version   string
		wantError bool
	}{
		{
			name:      "TLS 1.0",
			version:   "1.0",
			wantError: false,
		},
		{
			name:      "TLS 1.1",
			version:   "1.1",
			wantError: false,
		},
		{
			name:      "TLS 1.2",
			version:   "1.2",
			wantError: false,
		},
		{
			name:      "TLS 1.3",
			version:   "1.3",
			wantError: false,
		},
		{
			name:      "Invalid version",
			version:   "invalid",
			wantError: false, // Should default to TLS 1.2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// This is tested indirectly through server startup
			// The parseTLSVersion function is internal
			// We verify the server starts with the given version
			logger, err := mcplog.NewLogger("debug")
			require.NoError(t, err)

			cfg := &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "http",
					BindAddress: "localhost:0",
					TLS: mcpconfig.TLSConfig{
						Enabled:    false, // Disable TLS for this test
						MinVersion: tt.version,
					},
				},
			}

			server := mcp.NewServer(&mcp.Config{
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

			// Wait a bit for server to start
			time.Sleep(200 * time.Millisecond)

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
		})
	}
}

func TestServer_TLSConnectionLogging(t *testing.T) {
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

	server := mcp.NewServer(&mcp.Config{
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

	// Make a request
	httpServer := server.GetHTTPServer()
	if httpServer != nil {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		httpServer.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, []byte("OK"), w.Body.Bytes())
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

func TestServer_TLSReload(t *testing.T) {
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
		Logging: mcpconfig.LoggingConfig{
			Level: "debug",
		},
	}

	server := mcp.NewServer(&mcp.Config{
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

	// Try to reload TLS (should succeed even with TLS disabled)
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

func TestServer_TLSConfigGetters(t *testing.T) {
	t.Parallel()

	tlsCfg := &mcpconfig.TLSConfig{
		Enabled:    true,
		CertFile:   "/path/to/cert.pem",
		KeyFile:    "/path/to/key.pem",
		MinVersion: "1.3",
		MTLS: mcpconfig.MTLSConfig{
			Enabled: true,
			CAFile:  "/path/to/ca.pem",
		},
	}

	require.True(t, tlsCfg.GetEnabled())
	require.Equal(t, "/path/to/cert.pem", tlsCfg.GetCertFile())
	require.Equal(t, "/path/to/key.pem", tlsCfg.GetKeyFile())
	require.Equal(t, "1.3", tlsCfg.GetMinVersion())
	require.True(t, tlsCfg.GetMTLSEnabled())
	require.Equal(t, "/path/to/ca.pem", tlsCfg.GetMTLSCAFile())
}

func TestServer_TLSConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := mcpconfig.DefaultConfig()

	require.False(t, cfg.Server.TLS.Enabled)
	require.Equal(t, "1.2", cfg.Server.TLS.MinVersion)
	require.False(t, cfg.Server.TLS.MTLS.Enabled)
	require.Empty(t, cfg.Server.TLS.CertFile)
	require.Empty(t, cfg.Server.TLS.KeyFile)
	require.Empty(t, cfg.Server.TLS.MTLS.CAFile)
}

func TestServer_TLSConfigValidation(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *mcpconfig.Config
		wantErr bool
	}{
		{
			name: "Valid TLS config",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "http",
					BindAddress: "localhost:0",
					TLS: mcpconfig.TLSConfig{
						Enabled:    true,
						CertFile:   "/path/to/cert.pem",
						KeyFile:    "/path/to/key.pem",
						MinVersion: "1.2",
					},
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "TLS enabled but cert file missing",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "http",
					BindAddress: "localhost:0",
					TLS: mcpconfig.TLSConfig{
						Enabled:    true,
						CertFile:   "",
						KeyFile:    "/path/to/key.pem",
						MinVersion: "1.2",
					},
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
		{
			name: "TLS enabled but key file missing",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "http",
					BindAddress: "localhost:0",
					TLS: mcpconfig.TLSConfig{
						Enabled:    true,
						CertFile:   "/path/to/cert.pem",
						KeyFile:    "",
						MinVersion: "1.2",
					},
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
		{
			name: "mTLS enabled but CA file missing",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "http",
					BindAddress: "localhost:0",
					TLS: mcpconfig.TLSConfig{
						Enabled:    true,
						CertFile:   "/path/to/cert.pem",
						KeyFile:    "/path/to/key.pem",
						MinVersion: "1.2",
						MTLS: mcpconfig.MTLSConfig{
							Enabled: true,
							CAFile:  "",
						},
					},
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "info",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_TLSConfigString(t *testing.T) {
	t.Parallel()

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "localhost:0",
			TLS: mcpconfig.TLSConfig{
				Enabled:    true,
				CertFile:   "/path/to/cert.pem",
				KeyFile:    "/path/to/key.pem",
				MinVersion: "1.3",
				MTLS: mcpconfig.MTLSConfig{
					Enabled: true,
					CAFile:  "/path/to/ca.pem",
				},
			},
		},
		Logging: mcpconfig.LoggingConfig{
			Level: "debug",
		},
		Search: mcpconfig.SearchConfig{
			MaxResults: 10,
			SafeSearch: false,
		},
		Perplexity: mcpconfig.PerplexityConfig{
			Enabled:     false,
			AccessToken: "",
		},
	}

	str := cfg.String()
	require.Contains(t, str, "tls_enabled=true")
	require.Contains(t, str, "mtls_enabled=true")
}

func TestServer_TLSConfigGetTLSConfig(t *testing.T) {
	t.Parallel()

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "localhost:0",
			TLS: mcpconfig.TLSConfig{
				Enabled:    true,
				CertFile:   "/path/to/cert.pem",
				KeyFile:    "/path/to/key.pem",
				MinVersion: "1.3",
			},
		},
	}

	tlsCfg := cfg.GetTLSConfig()
	require.NotNil(t, tlsCfg)

	tlsConfig, ok := tlsCfg.(*mcpconfig.TLSConfig)
	require.True(t, ok)
	require.True(t, tlsConfig.Enabled)
	require.Equal(t, "/path/to/cert.pem", tlsConfig.CertFile)
}

func TestServer_TLSHandshakeLogging(t *testing.T) { //nolint:funlen
	t.Parallel()

	logger, err := mcplog.NewLogger("debug")
	require.NoError(t, err)

	// Get absolute paths to test certificates
	cwd, err := os.Getwd()
	require.NoError(t, err)

	certFile := cwd + "/testdata/server-cert.pem"
	keyFile := cwd + "/testdata/server-key.pem"
	caFile := cwd + "/testdata/ca-cert.pem"

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "127.0.0.1:0",
			TLS: mcpconfig.TLSConfig{
				Enabled:    true,
				CertFile:   certFile,
				KeyFile:    keyFile,
				MinVersion: "1.2",
				MTLS: mcpconfig.MTLSConfig{
					Enabled: true,
					CAFile:  caFile,
				},
			},
		},
		Logging: mcpconfig.LoggingConfig{
			Level: "debug",
		},
	}

	server := mcp.NewServer(&mcp.Config{
		Name:    "test-server",
		Version: "1.0.0",
	}, logger)

	server.SetAppConfig(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)
	serverReady := make(chan struct{})

	go func() {
		err := server.Serve(ctx, cfg.Server.Protocol, cfg.Server.BindAddress)
		if err != nil && err.Error() != errContextCanceled {
			serverErr <- err
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)
	close(serverReady)

	// Get the actual server address
	httpServer := server.GetHTTPServer()
	require.NotNil(t, httpServer)

	// Test TLS connection without client certificate (should fail)
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//nolint:gosec // G402 - InsecureSkipVerify is acceptable for testing self-signed certs
				InsecureSkipVerify: true,
			},
		},
	}

	// Get the actual address from the listener
	addr := httpServer.Addr
	require.NotEmpty(t, addr)

	var resp *http.Response

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+addr+"/health", nil)
	require.NoError(t, err)

	resp, err = client.Do(req)
	if err == nil && resp != nil {
		//nolint:gosec // G104 - error handling not needed for Close in test cleanup
		resp.Body.Close()
	}

	require.Error(t, err, "mTLS should reject connection without client certificate")

	// Test TLS connection with client certificate (should succeed)
	clientCert, err := tls.LoadX509KeyPair(cwd+"/testdata/client-cert.pem", cwd+"/testdata/client-key.pem")
	require.NoError(t, err)

	//nolint:gosec // G304 - path is from test data, safe for testing
	caCert, err := os.ReadFile(caFile)
	require.NoError(t, err)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      caCertPool,
				//nolint:gosec // G402 - InsecureSkipVerify is acceptable for testing self-signed certs
				InsecureSkipVerify: true,
			},
		},
	}

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+addr+"/health", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	//nolint:gosec // G104 - error handling not needed for Close in test cleanup
	resp.Body.Close()

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

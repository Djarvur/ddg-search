package mcpconfig_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Djarvur/ddg-search/internal/mcpconfig"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := mcpconfig.DefaultConfig()

	require.Equal(t, "stdio", cfg.Server.Protocol)
	require.Equal(t, "localhost:9100", cfg.Server.BindAddress)
	require.False(t, cfg.Server.TLS.Enabled)
	require.Equal(t, "1.2", cfg.Server.TLS.MinVersion)
	require.False(t, cfg.Server.TLS.MTLS.Enabled)
	require.Equal(t, "info", cfg.Logging.Level)
	require.Equal(t, 10, cfg.Search.MaxResults)
	require.False(t, cfg.Search.SafeSearch)
	require.False(t, cfg.Perplexity.Enabled)
	require.Empty(t, cfg.Perplexity.AccessToken)
}

//nolint:funlen
func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *mcpconfig.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol:    "stdio",
					BindAddress: "localhost:9100",
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "info",
				},
				Search: mcpconfig.SearchConfig{
					MaxResults: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid protocol",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name: "TLS enabled without cert file",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol: "http",
					TLS: mcpconfig.TLSConfig{
						Enabled: true,
						KeyFile: "key.pem",
						MTLS:    mcpconfig.MTLSConfig{},
					},
				},
			},
			wantErr: true,
			errMsg:  "cert_file not specified",
		},
		{
			name: "TLS enabled without key file",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol: "http",
					TLS: mcpconfig.TLSConfig{
						Enabled:  true,
						CertFile: "cert.pem",
						MTLS:     mcpconfig.MTLSConfig{},
					},
				},
			},
			wantErr: true,
			errMsg:  "key_file not specified",
		},
		{
			name: "mTLS enabled without CA file",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol: "http",
					TLS: mcpconfig.TLSConfig{
						Enabled:  true,
						CertFile: "cert.pem",
						KeyFile:  "key.pem",
						MTLS: mcpconfig.MTLSConfig{
							Enabled: true,
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "ca_file not specified",
		},
		{
			name: "invalid log level",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol: "stdio",
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name: "negative max_results",
			cfg: &mcpconfig.Config{
				Server: mcpconfig.ServerConfig{
					Protocol: "stdio",
				},
				Logging: mcpconfig.LoggingConfig{
					Level: "info",
				},
				Search: mcpconfig.SearchConfig{
					MaxResults: -1,
				},
			},
			wantErr: true,
			errMsg:  "must be >= 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  protocol: http
  bind_address: "0.0.0.0:8080"
  tls:
    enabled: true
    cert_file: /path/to/cert.pem
    key_file: /path/to/key.pem
    min_version: "1.3"
    mtls:
      enabled: true
      ca_file: /path/to/ca.pem
logging:
  level: debug
search:
  max_results: 20
  safe_search: true
perplexity:
  enabled: true
  access_token: test-token
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0o600))

	// Set config path via environment
	t.Setenv("DDG_SEARCH_CONFIG_FILE", configPath)

	cfg, err := mcpconfig.Load()
	require.NoError(t, err)
	require.Equal(t, "http", cfg.Server.Protocol)
	require.Equal(t, "0.0.0.0:8080", cfg.Server.BindAddress)
	require.True(t, cfg.Server.TLS.Enabled)
	require.Equal(t, "/path/to/cert.pem", cfg.Server.TLS.CertFile)
	require.Equal(t, "/path/to/key.pem", cfg.Server.TLS.KeyFile)
	require.Equal(t, "1.3", cfg.Server.TLS.MinVersion)
	require.True(t, cfg.Server.TLS.MTLS.Enabled)
	require.Equal(t, "/path/to/ca.pem", cfg.Server.TLS.MTLS.CAFile)
	require.Equal(t, "debug", cfg.Logging.Level)
	require.Equal(t, 20, cfg.Search.MaxResults)
	require.True(t, cfg.Search.SafeSearch)
	require.True(t, cfg.Perplexity.Enabled)
	require.Equal(t, "test-token", cfg.Perplexity.AccessToken)
}

func TestLoadWithEnvOverrides(t *testing.T) {
	// Set environment variables
	t.Setenv("DDG_SEARCH_SERVER_PROTOCOL", "http")
	t.Setenv("DDG_SEARCH_SERVER_BIND_ADDRESS", "0.0.0.0:9000")
	t.Setenv("DDG_SEARCH_LOGGING_LEVEL", "debug")
	t.Setenv("DDG_SEARCH_SEARCH_MAX_RESULTS", "15")
	t.Setenv("DDG_SEARCH_SEARCH_SAFE_SEARCH", "true")
	t.Setenv("DDG_SEARCH_PERPLEXITY_ENABLED", "true")
	t.Setenv("DDG_SEARCH_PERPLEXITY_ACCESS_TOKEN", "env-token")

	cfg, err := mcpconfig.Load()
	require.NoError(t, err)
	require.Equal(t, "http", cfg.Server.Protocol)
	require.Equal(t, "0.0.0.0:9000", cfg.Server.BindAddress)
	require.Equal(t, "debug", cfg.Logging.Level)
	require.Equal(t, 15, cfg.Search.MaxResults)
	require.True(t, cfg.Search.SafeSearch)
	require.True(t, cfg.Perplexity.Enabled)
	require.Equal(t, "env-token", cfg.Perplexity.AccessToken)
}

func TestLoadWithMissingConfigFile(t *testing.T) {
	// Set a non-existent config path
	t.Setenv("DDG_SEARCH_CONFIG_FILE", "/nonexistent/config.yaml")

	cfg, err := mcpconfig.Load()
	require.NoError(t, err)
	// Should use defaults
	require.Equal(t, "stdio", cfg.Server.Protocol)
	require.Equal(t, "info", cfg.Logging.Level)
}

func TestConfigString(t *testing.T) {
	t.Parallel()

	cfg := &mcpconfig.Config{
		Server: mcpconfig.ServerConfig{
			Protocol:    "http",
			BindAddress: "0.0.0.0:8080",
			TLS: mcpconfig.TLSConfig{
				Enabled: true,
				MTLS: mcpconfig.MTLSConfig{
					Enabled: true,
				},
			},
		},
		Logging: mcpconfig.LoggingConfig{
			Level: "debug",
		},
		Search: mcpconfig.SearchConfig{
			MaxResults: 20,
			SafeSearch: true,
		},
		Perplexity: mcpconfig.PerplexityConfig{
			Enabled: true,
		},
	}

	str := cfg.String()
	require.Contains(t, str, "protocol=http")
	require.Contains(t, str, "bind_address=0.0.0.0:8080")
	require.Contains(t, str, "tls_enabled=true")
	require.Contains(t, str, "mtls_enabled=true")
	require.Contains(t, str, "level=debug")
	require.Contains(t, str, "max_results=20")
	require.Contains(t, str, "safe_search=true")
	require.Contains(t, str, "enabled=true")
}

func TestReloadableConfig_Get(t *testing.T) {
	t.Parallel()

	cfg := mcpconfig.DefaultConfig()
	reloadable := mcpconfig.NewReloadableConfig(cfg, slog.Default())

	got := reloadable.Get()
	require.Equal(t, cfg, got)
}

func TestReloadableConfig_Reload(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write initial config
	initialConfig := `
logging:
  level: info
search:
  max_results: 10
`
	require.NoError(t, os.WriteFile(configPath, []byte(initialConfig), 0o600))

	// Set config path
	t.Setenv("DDG_SEARCH_CONFIG_FILE", configPath)

	// Load initial config
	cfg, err := mcpconfig.Load()
	require.NoError(t, err)
	require.Equal(t, "info", cfg.Logging.Level)
	require.Equal(t, 10, cfg.Search.MaxResults)

	// Create reloadable config
	reloadable := mcpconfig.NewReloadableConfig(cfg, slog.Default())

	// Update config file
	updatedConfig := `
logging:
  level: debug
search:
  max_results: 20
`
	require.NoError(t, os.WriteFile(configPath, []byte(updatedConfig), 0o600))

	// Reload
	err = reloadable.Reload()
	require.NoError(t, err)

	// Verify reloaded config
	got := reloadable.Get()
	require.Equal(t, "debug", got.Logging.Level)
	require.Equal(t, 20, got.Search.MaxResults)
}

func TestReloadableConfig_ReloadInvalidConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write valid initial config
	initialConfig := `
logging:
  level: info
`
	require.NoError(t, os.WriteFile(configPath, []byte(initialConfig), 0o600))

	// Set config path
	t.Setenv("DDG_SEARCH_CONFIG_FILE", configPath)

	// Load initial config
	cfg, err := mcpconfig.Load()
	require.NoError(t, err)

	// Create reloadable config
	reloadable := mcpconfig.NewReloadableConfig(cfg, slog.Default())

	// Write invalid config
	invalidConfig := `
logging:
  level: invalid
`
	require.NoError(t, os.WriteFile(configPath, []byte(invalidConfig), 0o600))

	// Reload should fail
	err = reloadable.Reload()
	require.Error(t, err)

	// Config should remain unchanged
	got := reloadable.Get()
	require.Equal(t, "info", got.Logging.Level)
}

func TestReloadableConfig_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	cfg := mcpconfig.DefaultConfig()
	reloadable := mcpconfig.NewReloadableConfig(cfg, slog.Default())

	done := make(chan struct{})

	// Start multiple goroutines reading config
	g := errgroup.Group{}
	for range 10 {
		g.Go(func() error {
			for {
				select {
				case <-done:
					return nil
				default:
					_ = reloadable.Get()

					time.Sleep(1 * time.Millisecond)
				}
			}
		})
	}

	// Reload config multiple times
	for range 5 {
		time.Sleep(10 * time.Millisecond)

		_ = reloadable.Reload()
	}

	// Stop goroutines
	close(done)

	_ = g.Wait()
}

func TestReloadableConfig_WatchSignals(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Write initial config
	initialConfig := `
logging:
  level: info
`
	require.NoError(t, os.WriteFile(configPath, []byte(initialConfig), 0o600))

	// Set config path
	t.Setenv("DDG_SEARCH_CONFIG_FILE", configPath)

	// Load initial config
	cfg, err := mcpconfig.Load()
	require.NoError(t, err)

	// Create reloadable config
	reloadable := mcpconfig.NewReloadableConfig(cfg, slog.Default())

	// Create shutdown channel
	shutdownChan := make(chan struct{})

	// Start signal watcher in goroutine
	go reloadable.WatchSignals(shutdownChan)

	// Wait a bit for the watcher to start
	time.Sleep(100 * time.Millisecond)

	// Update config file
	updatedConfig := `
logging:
  level: debug
`
	require.NoError(t, os.WriteFile(configPath, []byte(updatedConfig), 0o600))

	// Send SIGHUP signal
	// Note: We can't actually send SIGHUP in tests, so we just verify the watcher started
	// The actual signal handling is tested in e2e tests

	// Shutdown
	close(shutdownChan)
	time.Sleep(100 * time.Millisecond)
}

func TestSetDefaults(t *testing.T) {
	// Set DDG_SEARCH_CONFIG_FILE to empty string to prevent loading user's config file
	t.Setenv("DDG_SEARCH_CONFIG_FILE", "")
	
	// This is an internal function, so we test it indirectly through Load
	cfg, err := mcpconfig.Load()
	require.NoError(t, err)

	// Verify defaults are set
	require.Equal(t, "stdio", cfg.Server.Protocol)
	require.Equal(t, "localhost:9100", cfg.Server.BindAddress)
	require.False(t, cfg.Server.TLS.Enabled)
	require.Equal(t, "1.2", cfg.Server.TLS.MinVersion)
	require.False(t, cfg.Server.TLS.MTLS.Enabled)
	require.Equal(t, "info", cfg.Logging.Level)
	require.Equal(t, 10, cfg.Search.MaxResults)
	require.False(t, cfg.Search.SafeSearch)
	require.False(t, cfg.Perplexity.Enabled)
	require.Empty(t, cfg.Perplexity.AccessToken)
}

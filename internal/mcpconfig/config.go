// Package mcpconfig provides configuration management for the MCP server.
package mcpconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	// DefaultMaxResults is the default maximum number of search results.
	DefaultMaxResults = 10
)

var (
	// ErrInvalidProtocol is returned when the protocol is not "stdio" or "http".
	ErrInvalidProtocol = errors.New("invalid protocol (must be 'stdio' or 'http')")
	// ErrTLSCertNotSpecified is returned when TLS is enabled but cert_file is not specified.
	ErrTLSCertNotSpecified = errors.New("TLS enabled but cert_file not specified")
	// ErrTLSKeyNotSpecified is returned when TLS is enabled but key_file is not specified.
	ErrTLSKeyNotSpecified = errors.New("TLS enabled but key_file not specified")
	// ErrMTLSCANotSpecified is returned when mTLS is enabled but ca_file is not specified.
	ErrMTLSCANotSpecified = errors.New("mTLS enabled but ca_file not specified")
	// ErrInvalidLogLevel is returned when the log level is invalid.
	ErrInvalidLogLevel = errors.New("invalid log level (must be 'debug', 'info', 'warn', or 'error')")
	// ErrInvalidMaxResults is returned when max_results is negative.
	ErrInvalidMaxResults = errors.New("search.max_results must be >= 0")
)

// Config holds the configuration for the MCP server.
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Search     SearchConfig     `mapstructure:"search"`
	Perplexity PerplexityConfig `mapstructure:"perplexity"`
}

// ServerConfig holds server configuration.
type ServerConfig struct {
	Protocol    string    `mapstructure:"protocol"`
	BindAddress string    `mapstructure:"bind_address"`
	TLS         TLSConfig `mapstructure:"tls"`
}

// TLSConfig holds TLS configuration.
type TLSConfig struct {
	Enabled    bool       `mapstructure:"enabled"`
	CertFile   string     `mapstructure:"cert_file"`
	KeyFile    string     `mapstructure:"key_file"`
	MinVersion string     `mapstructure:"min_version"`
	MTLS       MTLSConfig `mapstructure:"mtls"`
}

// MTLSConfig holds mutual TLS configuration.
type MTLSConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	CAFile  string `mapstructure:"ca_file"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// SearchConfig holds search tool configuration.
type SearchConfig struct {
	MaxResults int  `mapstructure:"max_results"`
	SafeSearch bool `mapstructure:"safe_search"`
}

// PerplexityConfig holds Perplexity search configuration.
type PerplexityConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	AccessToken string `mapstructure:"access_token"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Protocol:    "stdio",
			BindAddress: "localhost:9100",
			TLS: TLSConfig{
				Enabled:    false,
				MinVersion: "1.2",
				MTLS: MTLSConfig{
					Enabled: false,
				},
			},
		},
		Logging: LoggingConfig{
			Level: "info",
		},
		Search: SearchConfig{
			MaxResults: DefaultMaxResults,
			SafeSearch: false,
		},
		Perplexity: PerplexityConfig{
			Enabled:     false,
			AccessToken: "",
		},
	}
}

// Load loads configuration from file, environment variables, and CLI flags.
// Priority: CLI flags > environment variables > config file > defaults.
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configure config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Check for DDG_SEARCH_CONFIG_FILE environment variable
	configFilePath, configFileSet := os.LookupEnv("DDG_SEARCH_CONFIG_FILE")
	skipConfigFile := false

	switch {
	case !configFileSet:
		// DDG_SEARCH_CONFIG_FILE is not set, search in default paths
		configPaths := []string{
			".",
			filepath.Join(os.Getenv("HOME"), ".config", "ddg-search"),
			"/etc/ddg-search",
		}
		for _, path := range configPaths {
			v.AddConfigPath(path)
		}
	case configFilePath == "":
		// DDG_SEARCH_CONFIG_FILE is set to empty string, skip config file loading
		skipConfigFile = true
	default:
		// DDG_SEARCH_CONFIG_FILE is set to a path, use it
		// Check if the file exists
		_, err := os.Stat(configFilePath)
		if err == nil {
			v.SetConfigFile(configFilePath)
		} else {
			// File doesn't exist, skip config file loading
			skipConfigFile = true
		}
	}

	// Read config file (ignore if not found)
	if !skipConfigFile {
		err := v.ReadInConfig()
		if err != nil {
			var notFoundErr viper.ConfigFileNotFoundError
			if !errors.As(err, &notFoundErr) {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	// Configure environment variables
	v.SetEnvPrefix("DDG_SEARCH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal config
	var cfg Config

	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadFromFile loads configuration from a specific file path.
func LoadFromFile(configFilePath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configure to use the specified config file
	v.SetConfigFile(configFilePath)
	v.SetConfigType("yaml")

	// Read config file
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Configure environment variables
	v.SetEnvPrefix("DDG_SEARCH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal config
	var cfg Config

	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values.
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.protocol", "stdio")
	v.SetDefault("server.bind_address", "localhost:9100")
	v.SetDefault("server.tls.enabled", false)
	v.SetDefault("server.tls.min_version", "1.2")
	v.SetDefault("server.tls.mtls.enabled", false)
	v.SetDefault("logging.level", "info")
	v.SetDefault("search.max_results", DefaultMaxResults)
	v.SetDefault("search.safe_search", false)
	v.SetDefault("perplexity.enabled", false)
	v.SetDefault("perplexity.access_token", "")
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	// Validate protocol
	if c.Server.Protocol != "stdio" && c.Server.Protocol != "http" {
		return fmt.Errorf("%w: %s", ErrInvalidProtocol, c.Server.Protocol)
	}

	// Validate TLS configuration
	if c.Server.TLS.Enabled {
		if c.Server.TLS.CertFile == "" {
			return ErrTLSCertNotSpecified
		}

		if c.Server.TLS.KeyFile == "" {
			return ErrTLSKeyNotSpecified
		}

		if c.Server.TLS.MTLS.Enabled && c.Server.TLS.MTLS.CAFile == "" {
			return ErrMTLSCANotSpecified
		}
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("%w: %s", ErrInvalidLogLevel, c.Logging.Level)
	}

	// Validate search max_results
	if c.Search.MaxResults < 0 {
		return ErrInvalidMaxResults
	}

	return nil
}

// String returns a string representation of the configuration.
func (c *Config) String() string {
	return fmt.Sprintf(
		"Server: protocol=%s, bind_address=%s, tls_enabled=%v, mtls_enabled=%v | "+
			"Logging: level=%s | "+
			"Search: max_results=%d, safe_search=%v | "+
			"Perplexity: enabled=%v",
		c.Server.Protocol,
		c.Server.BindAddress,
		c.Server.TLS.Enabled,
		c.Server.TLS.MTLS.Enabled,
		c.Logging.Level,
		c.Search.MaxResults,
		c.Search.SafeSearch,
		c.Perplexity.Enabled,
	)
}

// GetTLSConfig returns the TLS configuration.
func (c *Config) GetTLSConfig() any {
	return &c.Server.TLS
}

// GetEnabled returns whether TLS is enabled.
func (t *TLSConfig) GetEnabled() bool {
	return t.Enabled
}

// GetCertFile returns the TLS certificate file path.
func (t *TLSConfig) GetCertFile() string {
	return t.CertFile
}

// GetKeyFile returns the TLS key file path.
func (t *TLSConfig) GetKeyFile() string {
	return t.KeyFile
}

// GetMinVersion returns the minimum TLS version.
func (t *TLSConfig) GetMinVersion() string {
	return t.MinVersion
}

// GetMTLSEnabled returns whether mTLS is enabled.
func (t *TLSConfig) GetMTLSEnabled() bool {
	return t.MTLS.Enabled
}

// GetMTLSCAFile returns the mTLS CA certificate file path.
func (t *TLSConfig) GetMTLSCAFile() string {
	return t.MTLS.CAFile
}

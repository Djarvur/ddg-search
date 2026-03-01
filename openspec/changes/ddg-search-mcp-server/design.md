# Design Document

## Detailed design decisions and implementation notes.

## MCP Tool Definitions

### Search Tool (search - Anthropic-compatible naming)

```go
// Tool definition - uses snake_case for Claude Code compatibility
&mcp.Tool{
    Name:        "search",
    Description: "Search the web using Perplexity (if enabled) or DuckDuckGo (fallback)",
    InputSchema: mcp.JSONSchema{
        "type": "object",
        "properties": map[string]interface{}{
            "query":         map[string]string{"type": "string"},
            "max_results":   map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 100},
            "provider":      map[string]interface{}{"type": "string", "enum": []string{"auto", "perplexity", "duckduckgo"}},
            "model":         map[string]interface{}{"type": "string", "enum": []string{"sonar-small-online", "sonar-medium-online", "sonar-pro-online"}},
            "site":          map[string]string{"type": "string"},
            "region":        map[string]string{"type": "string"},
            "time":          map[string]interface{}{"type": "string", "enum": []string{"d", "w", "m", "y"}},
            "output_format": map[string]interface{}{"type": "string", "enum": []string{"text", "json"}},
        },
        "required": []string{"query"},
    },
}
```

### Page Dump Tool (fetch - Anthropic-compatible naming)

```go
&mcp.Tool{
    Name:        "fetch",
    Description: "Fetch a web page and convert it to markdown",
    InputSchema: mcp.JSONSchema{
        "type": "object",
        "properties": map[string]interface{}{
            "url":          map[string]string{"type": "string"},
            "timeout":      map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 300},
            "output_format": map[string]interface{}{"type": "string", "enum": []string{"text", "json"}},
        },
        "required": []string{"url"},
    },
}
```

## Configuration Structure

```go
type Config struct {
    Server      ServerConfig      `mapstructure:"server"`
    TLS         TLSConfig        `mapstructure:"tls"`
    Perplexity  PerplexityConfig  `mapstructure:"perplexity"`
    DuckDuckGo  DuckDuckGoConfig  `mapstructure:"duckduckgo"`
    PageDump    PageDumpConfig    `mapstructure:"pageDump"`
    Logging     LoggingConfig     `mapstructure:"logging"`
    Output      OutputConfig      `mapstructure:"output"`
}

type ServerConfig struct {
    Transport string `mapstructure:"transport"` // stdio | http
    Port      int    `mapstructure:"port"`
    Host      string `mapstructure:"host"`
}

type TLSConfig struct {
    Enabled    bool   `mapstructure:"enabled"`
    CertFile  string `mapstructure:"cert_file"`
    KeyFile   string `mapstructure:"key_file"`
    Combined  string `mapstructure:"combined"`
    CACert    string `mapstructure:"ca_cert"`
    ClientAuth string `mapstructure:"client_auth"` // none | request | require
}

type PerplexityConfig struct {
    Enabled    bool   `mapstructure:"enabled"`
    APIKey     string `mapstructure:"api_key"`
    Model      string `mapstructure:"model"`
    MaxResults int    `mapstructure:"max_results"`
    Timeout    string `mapstructure:"timeout"`
}

type DuckDuckGoConfig struct {
    MaxResults int    `mapstructure:"max_results"`
    Timeout    string `mapstructure:"timeout"`
    UserAgent  string `mapstructure:"user_agent"`
    Region     string `mapstructure:"region"`
}

type PageDumpConfig struct {
    Timeout   string `mapstructure:"timeout"`
    UserAgent string `mapstructure:"user_agent"`
}

type LoggingConfig struct {
    Level  string `mapstructure:"level"`  // debug | info | warn | error
    Format string `mapstructure:"format"` // text | json
}

type OutputConfig struct {
    DefaultFormat string `mapstructure:"default_format"` // text | json
}
```

## Provider Fallback Logic

```go
func (s *Server) handleSearch(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
    query := args["query"].(string)
    provider := getProvider(args) // auto | perplexity | duckduckgo

    // Determine provider
    if provider == "auto" {
        if s.cfg.Perplexity.Enabled && s.cfg.Perplexity.APIKey != "" {
            provider = "perplexity"
        } else {
            provider = "duckduckgo"
        }
    }

    // Try primary provider
    result, err := s.searchWithProvider(ctx, query, provider, args)
    if err != nil {
        // Check if we should fallback
        if shouldFallback(err) && provider == "perplexity" {
            s.logger.Warn("perplexity failed, falling back to duckduckgo", "error", err)
            result, err = s.searchWithProvider(ctx, query, "duckduckgo", args)
        }
        if err != nil {
            return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
        }
    }

    return formatResult(result, args["outputFormat"])
}

func shouldFallback(err error) bool {
    // Rate limit (429), Quota exceeded (402), Auth errors (401, 403)
    var apiErr *perplexity.APIError
    if errors.As(err, &apiErr) {
        return apiErr.StatusCode == 429 || apiErr.StatusCode == 401 ||
               apiErr.StatusCode == 402 || apiErr.StatusCode == 403
    }
    return false
}
```

## Output Format Handling

```go
func formatResult(result interface{}, formatArg string) (*mcp.CallToolResult, error) {
    format := "text" // default
    if f, ok := formatArg.(string); ok {
        format = f
    }

    switch format {
    case "json":
        jsonBytes, err := json.Marshal(result)
        if err != nil {
            return nil, err
        }
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                mcp.NewTextContent(string(jsonBytes)),
            },
        }, nil
    default:
        // Text format - convert results to readable text
        text := formatAsText(result)
        return &mcp.CallToolResult{
            Content: []mcp.Content{
                mcp.NewTextContent(text),
            },
        }, nil
    }
}
```

## Signal Handling

```go
func setupSignalHandling(s *Server) {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

    go func() {
        for sig := range sigChan {
            switch sig {
            case syscall.SIGINT, syscall.SIGTERM:
                s.logger.Info("shutting down", "signal", sig)
                s.Cancel()
                return
            case syscall.SIGHUP:
                s.logger.Info("reloading configuration")
                if err := s.ReloadConfig(); err != nil {
                    s.logger.Error("failed to reload config", "error", err)
                }
            }
        }
    }()
}
```

## Test Organization

```
internal/
├── config/
│   ├── config.go
│   ├── config_test.go        // Public API tests
│   └── config_internal_test.go  // Private logic tests
├── mcp/
│   ├── server.go
│   ├── server_test.go       // Public API tests
│   ├── server_internal_test.go  // Private logic tests
│   ├── handlers.go
│   ├── handlers_test.go     // Public API tests
│   ├── handlers_internal_test.go  // Private logic tests
│   ├── transport.go
│   ├── transport_test.go
│   └── integration_test.go  // MCP protocol tests
│
e2e/
│   ├── mcp_server_test.go   // Full server e2e tests
│   └── signal_test.go       // Signal handling tests
```

## Logging with slog (Proxy-style)

```go
// Initialize logger
func initLogger(cfg *LoggingConfig) {
    var level slog.Level
    switch cfg.Level {
    case "debug":
        level = slog.LevelDebug
    case "info":
        level = slog.LevelInfo
    case "warn":
        level = slog.LevelWarn
    case "error":
        level = slog.LevelError
    }

    var handler slog.Handler
    if cfg.Format == "json" {
        handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
    } else {
        handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
    }

    slog.SetDefault(slog.New(handler))
}

// Proxy-style request logger
type RequestLogger struct {
    logger *slog.Logger
}

func NewRequestLogger(logger *slog.Logger) *RequestLogger {
    return &RequestLogger{logger: logger}
}

func (r *RequestLogger) Log(method, toolName string, duration time.Duration, status int, attrs ...slog.Attr) {
    baseAttrs := []slog.Attr{
        slog.String("method", method),
        slog.String("name", toolName),
        slog.Duration("duration", duration),
        slog.Int("status", status),
    }
    baseAttrs = append(baseAttrs, attrs...)

    switch {
    case status >= 500:
        r.logger.Error("request failed", baseAttrs...)
    case status >= 400:
        r.logger.Warn("bad request", baseAttrs...)
    default:
        r.logger.Info("request completed", baseAttrs...)
    }
}

// Example usage in handlers
func (s *Server) handleWebSearch(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResult, error) {
    start := time.Now()

    query := args["query"].(string)
    provider := s.determineProvider(args)

    result, err := s.search(ctx, query, provider, args)
    status := 200
    if err != nil {
        status = s.getErrorStatus(err)
        if status < 500 && provider == "perplexity" {
            // Try fallback
            result, err = s.search(ctx, query, "duckduckgo", args)
            if err != nil {
                status = s.getErrorStatus(err)
            } else {
                status = 200
                s.logger.Warn("perplexity failed, using duckduckgo",
                    slog.String("query", query),
                    slog.String("error", err.Error()))
            }
        }
    }

    duration := time.Since(start)
    s.requestLogger.Log("tools/call", "search", duration, status,
        slog.String("query", query),
        slog.String("provider", provider))

    return result, err
}
```

## Config Precedence Implementation

```go
func loadConfig() (*Config, error) {
    viper.SetConfigType("yaml")
    viper.SetConfigFile(configFilePath) // from --config flag

    // Environment variables (DDG_* prefix, no .env file support)
    viper.SetEnvPrefix("DDG")
    viper.AutomaticEnv()

    // Special handling for Perplexity API key
    // Uses DDG_PERPLEXITY_API_KEY env var (no .env file support)
    if apiKey := os.Getenv("DDG_PERPLEXITY_API_KEY"); apiKey != "" {
        viper.Set("perplexity.api_key", apiKey)
    }

    // CLI flags override (set after env)
    if transport != "" {
        viper.Set("server.transport", transport)
    }
    if port > 0 {
        viper.Set("server.port", port)
    }

    // Load config file
    if err := viper.ReadInConfig(); err != nil && !errors.Is(err, viper.ConfigFileNotFoundError{}) {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

// Package main implements the ddg-search-mcp server, a Model Context Protocol (MCP)
// server that provides web search and content fetching tools to Claude Code.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Djarvur/ddg-search/internal/mcp"
	"github.com/Djarvur/ddg-search/internal/mcpconfig"
	"github.com/Djarvur/ddg-search/internal/mcplog"
	"github.com/Djarvur/ddg-search/internal/mcpsignal"
	"github.com/spf13/cobra"
)

const (
	version = "dev"
)

func main() {
	rootCmd := newRootCmd()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		flagLogLevel   string
		flagConfigFile string
	)

	cmd := &cobra.Command{
		Use:   "ddg-search-mcp",
		Short: "MCP server for web search and content fetching",
		Long: `ddg-search-mcp is a Model Context Protocol (MCP) server that provides ` +
			`web search and content fetching tools to Claude Code.`,
		Version: version,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runServer(flagLogLevel, flagConfigFile)
		},
	}

	cmd.Flags().StringVar(&flagLogLevel, "log-level", "", "Log level (debug, info, warn, error)")
	cmd.Flags().StringVar(&flagConfigFile, "config", "", "Path to configuration file")

	return cmd
}

func runServer(flagLogLevel string, flagConfigFile string) error {
	cfg, err := loadAndValidateConfig(flagLogLevel, flagConfigFile)
	if err != nil {
		return err
	}

	logger, err := setupLogger(cfg)
	if err != nil {
		return err
	}

	shutdown := mcpsignal.NewShutdownHandler(logger)
	shutdown.WatchSignals()

	mcpServer := setupMCPServer(cfg, logger)

	reloadableCfg := mcpconfig.NewReloadableConfig(cfg, logger)
	startConfigWatcher(reloadableCfg, shutdown, mcpServer, logger)

	serverErr := startServer(mcpServer, cfg, shutdown, logger)

	return waitForShutdown(shutdown, serverErr, logger)
}

func loadAndValidateConfig(flagLogLevel, flagConfigFile string) (*mcpconfig.Config, error) {
	cfg, err := mcpconfig.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if flagConfigFile != "" {
		cfg, err = mcpconfig.LoadFromFile(flagConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration from file: %w", err)
		}
	}

	if flagLogLevel != "" {
		cfg.Logging.Level = flagLogLevel
	}

	validateErr := cfg.Validate()
	if validateErr != nil {
		return nil, fmt.Errorf("invalid configuration: %w", validateErr)
	}

	return cfg, nil
}

func setupLogger(cfg *mcpconfig.Config) (*slog.Logger, error) {
	logger, err := mcplog.NewLogger(cfg.Logging.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	slog.SetDefault(logger)
	logger.Info("Starting ddg-search-mcp", "version", version)
	logger.Info("Configuration loaded", "config", cfg.String())

	return logger, nil
}

func setupMCPServer(cfg *mcpconfig.Config, logger *slog.Logger) *mcp.Server {
	mcpServer := mcp.NewServer(&mcp.Config{
		Name:    "ddg-search-mcp",
		Version: version,
	}, logger)

	mcpServer.SetAppConfig(cfg)
	mcpServer.RegisterTool(mcp.SearchTool(), mcp.HandleSearch)
	mcpServer.RegisterTool(mcp.FetchTool(), mcp.HandleFetch)

	return mcpServer
}

func startConfigWatcher(
	reloadableCfg *mcpconfig.ReloadableConfig,
	shutdown *mcpsignal.ShutdownHandler,
	mcpServer *mcp.Server,
	logger *slog.Logger,
) {
	go reloadableCfg.WatchSignalsWithCallback(shutdown.ShutdownChan(), func() {
		currentCfg := reloadableCfg.Get()
		if currentCfg.Server.Protocol == "http" && currentCfg.Server.TLS.Enabled {
			err := mcpServer.ReloadTLS()
			if err != nil {
				logger.Error("Failed to reload TLS certificates", "error", err)
			}
		}
	})
}

func startServer(
	mcpServer *mcp.Server,
	cfg *mcpconfig.Config,
	shutdown *mcpsignal.ShutdownHandler,
	logger *slog.Logger,
) <-chan error {
	serverErr := make(chan error, 1)
	serverCtx, serverCancel := context.WithCancel(context.Background())

	go func() {
		defer serverCancel()

		err := mcpServer.Serve(serverCtx, cfg.Server.Protocol, cfg.Server.BindAddress)
		if err != nil {
			serverErr <- fmt.Errorf("MCP server error: %w", err)
		}

		logger.Info("MCP server stopped, initiating shutdown")
		shutdown.InitiateShutdown()
	}()

	logger.Info("Server ready, waiting for signals...")

	return serverErr
}

func waitForShutdown(shutdown *mcpsignal.ShutdownHandler, serverErr <-chan error, logger *slog.Logger) error {
	select {
	case <-shutdown.ShutdownChan():
		logger.Info("Shutting down...")
	case err := <-serverErr:
		if err != nil && err.Error() != "MCP server error: context canceled" {
			logger.Error("Server error", "error", err)

			return err
		}

		logger.Info("Shutting down...")
	}

	shutdown.Wait()
	logger.Info("Shutdown complete")

	return nil
}

package mcpconfig

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ReloadableConfig manages configuration with reload capability.
type ReloadableConfig struct {
	mu     sync.RWMutex
	config *Config
	logger *slog.Logger
}

// NewReloadableConfig creates a new reloadable configuration.
func NewReloadableConfig(cfg *Config, logger *slog.Logger) *ReloadableConfig {
	return &ReloadableConfig{
		config: cfg,
		logger: logger,
	}
}

// Get returns a copy of the current configuration.
func (r *ReloadableConfig) Get() *Config {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.config
}

// Reload reloads the configuration from file.
func (r *ReloadableConfig) Reload() error {
	r.logger.Info("Reloading configuration")

	newCfg, err := Load()
	if err != nil {
		r.logger.Error("Failed to reload configuration", "error", err)

		return err
	}

	validateErr := newCfg.Validate()
	if validateErr != nil {
		r.logger.Error("Invalid configuration after reload", "error", validateErr)

		return validateErr
	}

	r.mu.Lock()
	r.config = newCfg
	r.mu.Unlock()

	r.logger.Info("Configuration reloaded successfully", "config", newCfg.String())

	return nil
}

// WatchSignals watches for SIGHUP signal and reloads configuration.
func (r *ReloadableConfig) WatchSignals(shutdownChan <-chan struct{}) {
	r.WatchSignalsWithCallback(shutdownChan, nil)
}

// WatchSignalsWithCallback watches for SIGHUP signal and reloads configuration,
// calling the provided callback after successful reload.
func (r *ReloadableConfig) WatchSignalsWithCallback(shutdownChan <-chan struct{}, callback func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)

	for {
		select {
		case <-shutdownChan:
			signal.Stop(sigChan)

			return
		case sig := <-sigChan:
			if sig == syscall.SIGHUP {
				r.logger.Info("Received SIGHUP signal, reloading configuration")

				reloadErr := r.Reload()
				if reloadErr != nil {
					r.logger.Error("Configuration reload failed, keeping previous configuration", "error", reloadErr)
				} else if callback != nil {
					// Call callback after successful reload
					callback()
				}
			}
		}
	}
}

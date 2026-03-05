// Package mcpsignal provides signal handling for graceful shutdown.
package mcpsignal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ShutdownHandler manages graceful shutdown.
type ShutdownHandler struct {
	logger       *slog.Logger
	wg           sync.WaitGroup
	shutdownChan chan struct{}
	mu           sync.Mutex
}

// NewShutdownHandler creates a new shutdown handler.
func NewShutdownHandler(logger *slog.Logger) *ShutdownHandler {
	return &ShutdownHandler{
		logger:       logger,
		shutdownChan: make(chan struct{}),
	}
}

// Context returns the context that is cancelled on shutdown.
func (h *ShutdownHandler) Context() context.Context {
	h.mu.Lock()
	defer h.mu.Unlock()

	select {
	case <-h.shutdownChan:
		return context.Background()
	default:
		return context.Background()
	}
}

// ShutdownChan returns a channel that is closed when shutdown is initiated.
func (h *ShutdownHandler) ShutdownChan() <-chan struct{} {
	return h.shutdownChan
}

// Add adds a goroutine to wait for before shutdown completes.
func (h *ShutdownHandler) Add(delta int) {
	h.wg.Add(delta)
}

// Done marks a goroutine as done.
func (h *ShutdownHandler) Done() {
	h.wg.Done()
}

// Wait waits for all goroutines to complete.
func (h *ShutdownHandler) Wait() {
	h.wg.Wait()
}

// IsShutdown returns true if shutdown has been initiated.
func (h *ShutdownHandler) IsShutdown() bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	select {
	case <-h.shutdownChan:
		return true
	default:
		return false
	}
}

// WatchSignals watches for shutdown signals.
func (h *ShutdownHandler) WatchSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		h.logger.Info("Received shutdown signal", "signal", sig.String())
		h.InitiateShutdown()
	}()
}

// InitiateShutdown initiates the shutdown process.
func (h *ShutdownHandler) InitiateShutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	select {
	case <-h.shutdownChan:
		// Already shutting down
		return
	default:
		close(h.shutdownChan)
	}
}

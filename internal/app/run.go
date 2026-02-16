package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davidkamere/htmx-go-scaffolding/internal/config"
	"github.com/davidkamere/htmx-go-scaffolding/internal/middleware"
	"github.com/davidkamere/htmx-go-scaffolding/internal/server"
)

func Run() error {
	cfg := config.Load()
	logger := log.New(os.Stdout, "", log.LstdFlags)

	router, closer, err := server.NewRouter(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("build router: %w", err)
	}
	defer func() {
		if closeErr := closer.Close(); closeErr != nil {
			logger.Printf("failed to close store: %v", closeErr)
		}
	}()

	handler := middleware.Chain(
		router,
		middleware.SecurityHeaders(),
		middleware.RequestLogging(logger),
		middleware.Recovery(logger),
	)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Printf("server starting on :%s (env=%s log=%s db=%s)", cfg.Port, cfg.AppEnv, cfg.LogLevel, cfg.DBPath)
		if serveErr := httpServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Printf("shutdown signal received: %s", sig)
	case serveErr := <-errCh:
		return fmt.Errorf("server failed: %w", serveErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	logger.Println("server stopped")
	return nil
}

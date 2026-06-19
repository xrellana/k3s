package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k3s-sample/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	version := getenv("APP_VERSION", "dev")
	port := getenv("PORT", "8080")

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.New(version),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server starting", "addr", srv.Addr, "version", version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("server stopping")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

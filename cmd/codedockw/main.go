package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/worker"
)

func main() {
	slog.Info("Starting codedockw...")

	cfg := config.Get()
	token := cfg.Worker.WorkerToken
	if token == "" {
		slog.Error("CODEDOCK_WORKER_TOKEN is required")
		os.Exit(1)
	}

	apiHost := cfg.Server.ApiHost
	if apiHost == "" {
		apiHost = "http://localhost:8080"
	}

	daemon := worker.NewWorkerDaemon(apiHost, token)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := daemon.Start(ctx); err != nil {
		slog.Error("failed to start worker daemon", "err", err)
		os.Exit(1)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	slog.Info("Interrupt received, shutting down")
	cancel()
}

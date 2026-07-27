package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"

	"codedock.run/codedock/cmd/codedockd/commands"
	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/engine/networking"
	"codedock.run/codedock/internal/engine/observability"
	codedockhttp "codedock.run/codedock/internal/http"
	"codedock.run/codedock/internal/services/system"

	"codedock.run/codedock/internal/telemetry"
	"codedock.run/codedock/internal/version"
)

var codedockVersion = version.Version

func main() {
	_ = godotenv.Load(".env")
	commands.Execute(codedockVersion, startServer, runMCP)
}

func startServer() {
	slog.Info("booting daemon", "version", codedockVersion, "os", runtime.GOOS, "arch", runtime.GOARCH)
	dataDir, db, vlt := commands.InitDataDir()
	defer db.Close()

	telemetry.Init()
	defer telemetry.Close()
	telemetry.Track("system", "daemon_start", map[string]any{
		"version": codedockVersion,
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	})

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Warn("Docker daemon connection warning", "err", err, "detail", "container deployment features disabled")
	}

	traefikMgr := networking.NewTraefikManager(dockerClient, config.Get().Security.TLSEmail)
	if err := traefikMgr.EnsureTraefikRunning(context.Background()); err != nil {
		slog.Warn("failed to start Traefik proxy", "err", err)
	}

	tsdbMgr := observability.NewTSDBManager(dockerClient)
	if err := tsdbMgr.EnsureTSDBRunning(context.Background()); err != nil {
		slog.Warn("failed to start TSDB", "err", err)
	}

	lokiMgr := observability.NewLokiManager(dockerClient)
	if err := lokiMgr.EnsureLokiRunning(context.Background()); err != nil {
		slog.Warn("failed to start Loki", "err", err)
	}

	metricsWorker := observability.NewMetricsWorker(dockerClient)
	metricsWorker.Start()

	logWorker := observability.NewLogWorker(dockerClient)
	logWorker.Start(context.Background())

	system.StartTelemetryReporter(db, codedockVersion)

	host := config.Get().Server.Host
	port := fmt.Sprintf("%d", config.Get().Server.Port)
	addr := host + ":" + port

	deployer := engine.NewDeployer(dockerClient, commands.NewDBDeployerStore(db, vlt))
	apiServer, err := codedockhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient, dataDir)
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: apiServer.Handler(),
	}

	go func() {
		slog.Info("codedock server running", "addr", addr, "version", codedockVersion)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}

	slog.Info("server exited")
}

func runMCP() {
	slog.Info("starting MCP stdio server")
	_, db, vlt := commands.InitDataDir()
	defer db.Close()

	dockerClient, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	deployer := engine.NewDeployer(dockerClient, commands.NewDBDeployerStore(db, vlt))
	traefikMgr := networking.NewTraefikManager(dockerClient, config.Get().Security.TLSEmail)
	apiServer, err := codedockhttp.NewServer(db, vlt, deployer, traefikMgr, dockerClient, "")
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		os.Exit(1)
	}

	if err := apiServer.StartMCPStdio(); err != nil {
		slog.Error("MCP server exited", "err", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	cfg "github.com/shayanmkpr/task-pool/config"
	"github.com/shayanmkpr/task-pool/internal/api"
	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

func main() {
	// Top-level panic recovery
	defer func() {
		if err := recover(); err != nil {
			slog.Error("panic in main",
				"error", err,
				"stack", string(debug.Stack()),
			)
			// no os exit
		}
	}()

	lg, err := logger.New("./app.log")
	if err != nil {
		panic(err)
	}
	defer lg.Close()

	lg.Info("Application started")

	config := cfg.Load()

	memoryStore := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(config.PoolSize, memoryStore)
	workerManager := taskpool.NewWorkerManager(config.WorkerCount, memoryStore)

	workerManager.InitiateWorkers(pool)
	workerManager.MonitorWorkers(lg)

	// Set up HTTP server
	handler := api.NewHandler(pool, memoryStore, lg)
	mux := http.NewServeMux()
	api.RegisterTaskRoutes(mux, handler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      api.RequestLogger(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server listening on http://localhost%s\n", server.Addr)
		lg.Info("starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
	lg.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		lg.Error("server forced to shutdown", "error", err)
	}

	lg.Info("waiting for workers to complete...")
	workerManager.WaitForCompletion(ctx, lg, 100*time.Millisecond)
	workerManager.ForceStopWorkers()

	lg.Info("Application finished")
}

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdarshAnand26/Student-API/internal/config"
	"github.com/AdarshAnand26/Student-API/internal/http/handlers/student"
	"github.com/AdarshAnand26/Student-API/internal/storage/sqlite"
)

func main() {

	// Load config
	cfg := config.MustLoad()

	// Database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// Setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("Server Started", slog.String("Address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to Start Server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error(
			"Failed to shutdown server",
			slog.String("Error", err.Error()),
		)
	}

	slog.Info("Server Shutdown Successfully!!")
}
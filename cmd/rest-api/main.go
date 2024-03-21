package main

import (
	"awesomeProject/internal/http-server/handlers/redirect"
	"awesomeProject/internal/http-server/handlers/url/delete"
	"awesomeProject/internal/http-server/handlers/url/save"
	"awesomeProject/internal/http-server/handlers/url/update"
	"awesomeProject/internal/storage/postgres"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"

	"awesomeProject/internal/config"
	mwLogger "awesomeProject/internal/http-server/middleware/logger"
	loggerconfig "awesomeProject/internal/lib/logger/config"
	"awesomeProject/internal/lib/logger/sl"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// init logger: slog
	log := loggerconfig.SetupLogger(cfg.Env)

	log.Info("starting rest-api",
		slog.String("env", cfg.Env),
		slog.String("version", "0.1"),
	)
	log.Debug("debug messages are enabled")

	// init storage: sqlx
	storage, err := postgres.New(cfg.DBName, cfg.DBUser, cfg.DBPassword)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// init router: chi, chi render
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Put("/url", update.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/{alias}", delete.New(log, storage))

	// run server
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/config"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/database"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/logger"
)





func main() {
	// # Step 1: Set our custom logger as default logger
	appLogger := logger.New("development")
	slog.SetDefault(appLogger)

	// # Step 2: Load Configuration
	slog.Info(
		logger.MsgConfigLoadStart,
		"action", logger.ActionConfigLoadStart,
	)

	cfg, err := config.Load()
	if err != nil {
		slog.Error(
			logger.MsgConfigLoadError,
			"action", logger.ActionConfigLoadError,
			"error", err,
		)

		os.Exit(1) // exit the entire go application completely if config is not loaded
	}

	slog.Info(
		logger.MsgConfigLoaded,
		"action", logger.ActionConfigLoaded,
		"app_env", cfg.AppEnv,
		"port", cfg.Port,
	)

	port := cfg.Port
	appEnv := cfg.AppEnv

	// # Step 3: Connect to Database
	db, err := database.ConnectDB(cfg.Database)
	if err != nil {
		slog.Error(
			logger.MsgDatabaseConnectionError,
			"action", logger.ActionDatabaseConnectionError,
			"error", err,
		)

		os.Exit(1) // exit the entire go application completely if db connection is not established
	}
	defer db.Close()


	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		slog.Info(
			logger.MsgHTTPRequest,
			"action", logger.ActionHTTPRequest,
			"method", r.Method,
			"path", r.URL.Path,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"status":  "healthy",
			"message": "API is running",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			slog.Error(
				logger.MsgHTTPRequestError,
				"action", logger.ActionHTTPRequestError,
				"error", err,
			)
		}
	})

	slog.Info(
		logger.MsgServerStarted,
		"action", logger.ActionServerStarted,
		"port", port,
		"environment", appEnv,
	)

	if err := http.ListenAndServe(port, nil); err != nil {
		slog.Error(
			logger.MsgServerStopped,
			"action", logger.ActionServerStopped,
			"error", err,
		)

		os.Exit(1)
	}
}


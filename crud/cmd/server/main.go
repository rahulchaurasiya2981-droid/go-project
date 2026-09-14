package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/config"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/database"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/health"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/logger"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user"
)

func main() {
	// =========================================================
	// STEP 1: Initialize Logger
	// =========================================================
	appLogger := logger.New("development")
	slog.SetDefault(appLogger)

	// =========================================================
	// STEP 2: Load Configuration
	// =========================================================
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

	// =========================================================
	// STEP 3: Connect to Database
	// =========================================================
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

	healthHandler := health.NewHandler(db)
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// =========================================================================
	// STEP 4: Create HTTP Router & Register Routes with respective handlers
	// =========================================================================
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz/live", healthHandler.Liveness)
	mux.HandleFunc("GET /healthz/ready", healthHandler.Readiness)
	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	// mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)
	// mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)

	// =========================================================
	// STEP 5: Create HTTP Server
	// =========================================================

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// =========================================================
	// STEP 6: Start HTTP Server
	// =========================================================
	slog.Info(
		logger.MsgServerStarted,
		"action", logger.ActionServerStarted,
		"port", port,
		"environment", appEnv,
	)

	err = server.ListenAndServe()
	if err != nil {
		slog.Error(
			logger.MsgServerStopped,
			"action", logger.ActionServerStopped,
			"error", err,
		)

		os.Exit(1)
	}
}

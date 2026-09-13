package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/config"
	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/logger"
)

func ConnectDB(dbConfig config.DatabaseConfig) (*sql.DB, error) {

	db, err := sql.Open(dbConfig.DriverName, dbConfig.ConnectionString)
	if err != nil {
		slog.Error(
			logger.MsgDatabaseConnectionError,
			"action", logger.ActionDatabaseConnectionError,
			"driver", dbConfig.DriverName,
			"host", dbConfig.Host,
			"error", err,
		)
		return nil, fmt.Errorf("open connection: %w", err)
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)       // Maximum number of open DB connections
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)       // Maximum number of idle DB connections
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime) // Maximum lifetime of a DB connection
	db.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime) // Maximum time a connection can stay idle

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		slog.Error(
			logger.MsgDatabaseHealthFailed,
			"action", logger.ActionDatabaseHealthFailed,
			"host", dbConfig.Host,
			"database", dbConfig.Name,
			"error", err,
		)
		return nil, fmt.Errorf("ping database: %w", err)
	}

	slog.Info(
		logger.MsgDatabaseConnected,
		"action", logger.ActionDatabaseConnected,
		"host", dbConfig.Host,
		"database", dbConfig.Name,
	)

	return db, nil
}

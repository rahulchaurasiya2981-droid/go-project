package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	LogLevel string
	Port     string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	DriverName       string
	Host             string
	Port             string
	User             string
	Password         string
	Name             string
	SSLMode          string
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	ConnMaxIdleTime  time.Duration
}

func (d DatabaseConfig) BuildConnectionString() string {
	return (&url.URL{
		Scheme: os.Getenv("DATABASE_SCHEME"),
		User:   url.UserPassword(d.User, d.Password),
		Host:   d.Host + ":" + d.Port,
		Path:   d.Name,
		RawQuery: url.Values{
			"sslmode": []string{d.SSLMode},
		}.Encode(),
	}).String()
}

func Load() (*Config, error) {

	slog.Info(
		"Starting environment configuration load",
		"action", "ENV_LOAD_START",
	)

	if err := godotenv.Load(); err != nil {
		slog.Error(
			"Failed to load environment file",
			"action", "ENV_LOAD_ERROR",
			"error", err,
		)

		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	slog.Info(
		"Environment file loaded successfully",
		"action", "ENV_LOADED",
	)

	dbConfig := DatabaseConfig{
		DriverName:      os.Getenv("DATABASE_DRIVER"),
		Host:            os.Getenv("DATABASE_HOST"),
		Port:            os.Getenv("DATABASE_PORT"),
		User:            os.Getenv("DATABASE_USER"),
		Password:        os.Getenv("DATABASE_PASSWORD"),
		Name:            os.Getenv("DATABASE_NAME"),
		SSLMode:         os.Getenv("DATABASE_SSLMODE"),
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	dbConfig.ConnectionString = dbConfig.BuildConnectionString()

	cfg := &Config{
		AppEnv:   os.Getenv("APP_ENV"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		Port:     os.Getenv("PORT"),
		Database: dbConfig,
	}

	return cfg, nil
}

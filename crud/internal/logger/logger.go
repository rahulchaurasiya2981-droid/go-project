package logger

import (
	"log/slog"
	"os"
	"strings"
)

func parseLevel(value string) slog.Level {
	switch strings.ToUpper(value) {
	case "DEBUG":
		return slog.LevelDebug

	case "INFO":
		return slog.LevelInfo

	case "WARN":
		return slog.LevelWarn

	case "ERROR":
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}

func New(appEnv string) *slog.Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	options := &slog.HandlerOptions{
		Level: level,
	}
	var handler slog.Handler
	handler = slog.NewJSONHandler(os.Stdout, options)
	// switch strings.ToLower(appEnv) {
	// case "development":
	// 	handler = slog.NewJSONHandler(os.Stdout, options)

	// case "production":
	// 	handler = slog.NewJSONHandler(os.Stdout, options)

	// default:
	// 	handler = slog.NewJSONHandler(os.Stdout, options)
	// }

	return slog.New(handler)
}

package log

import (
	"log/slog"
	"os"
	"yuima-agent/internal/config"
)

func New(cfg config.Config) *slog.Logger {
	// Провера уровня по конфигу
	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	// Выбор формата вывода по конфигу
	var handler slog.Handler
	if cfg.LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	
	return slog.New(handler)
}
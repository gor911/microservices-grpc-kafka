package logger

import (
	"context"
	"log/slog"
	"os"
)

type Loglevel string

const (
	Debug Loglevel = "debug"
	Info  Loglevel = "info"
	Warn  Loglevel = "warn"
	Error Loglevel = "error"
)

type Config struct {
	Level Loglevel
}

func New(ctx context.Context, config Config) *slog.Logger {
	var log *slog.Logger

	switch config.Level {
	case Debug:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case Info:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case Warn:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	case Error:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
	}

	return log
}

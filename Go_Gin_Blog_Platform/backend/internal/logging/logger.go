package logging

import (
	"log/slog"
	"os"
)

func NewJSONLogger(appEnv string) *slog.Logger {
	level := slog.LevelInfo
	if appEnv == "local" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

package utils

import (
	"log/slog"
	"os"
)

// Loggers: Debug, Info, Warning, Error

func Init(isProduction bool) {
	// Define the logger
	var handler slog.Handler

	// Configuracion del handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	}

	if isProduction {
		opts.Level = slog.LevelWarn
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)

	}

	// Logger global
	slog.SetDefault(slog.New(handler))
}
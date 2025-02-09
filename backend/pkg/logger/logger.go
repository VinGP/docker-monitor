package logger

import (
	"backend/pkg/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

func NewLogger(prod bool) *slog.Logger {
	var log *slog.Logger

	switch prod {
	case false:
		log = setupPrettySlog()
	case true:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

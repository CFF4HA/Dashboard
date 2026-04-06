package core

import (
	"log"
	"log/slog"
)

var (
	Logger         = slog.New(slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	BackendAddress = "http://localhost:8081"
)

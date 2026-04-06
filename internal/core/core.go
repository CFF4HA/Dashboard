package core

import (
	"log"
	"log/slog"
)

var (
	Logger         = slog.New(slog.NewJSONHandler(log.Writer(), nil))
	BackendAddress = "http://localhost:8081"
)

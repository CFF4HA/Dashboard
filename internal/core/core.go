package core

import (
	"log"
	"log/slog"
)

var (
	Logger = slog.New(slog.NewJSONHandler(log.Writer(), nil))
)

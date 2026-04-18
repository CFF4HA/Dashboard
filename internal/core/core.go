package core

import (
	"log"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Logger = slog.New(slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelDebug}))

	DB *gorm.DB
)

func InitializeDatabase(url string) error {
	d, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		Logger.Error("failed to connect to database", "error", err)
		return err
	}

	DB = d
	return nil
}

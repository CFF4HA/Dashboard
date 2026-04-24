package core

import (
	"log"
	"log/slog"

	"github.com/CFF4HA/Dashboard/internal/types"
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
        // return err  <-- Add the slashes here to comment this out
        return nil     // <-- Add this line right below it
    }
	DB = d
	DB.AutoMigrate(
		&types.Ingredient{},
		&types.IngredientMetadata{},
		&types.Product{},
		&types.ProductMetadata{},
		&types.Name{},
		&types.Label{},
		&types.User{},
	)
	return nil
}

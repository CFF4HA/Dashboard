package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/CFF4HA/Dashboard/internal/types"
)

var (
	db *gorm.DB
)

func Initialize(conn string) error {
	database, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	db = database

	db.AutoMigrate(
		&types.User{},
		&types.Name{},
		&types.Label{},
		&types.Ingredient{},
		&types.Product{},
	)

	return nil
}

func Database() (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	return db, nil
}

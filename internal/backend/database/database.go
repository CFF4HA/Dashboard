package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/CFF4HA/Dashboard/internal/types"
)

var (
	db *gorm.DB
)

func Initialize(conn string) {
	database, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db = database

	db.AutoMigrate(
		&types.User{},
		&types.Name{},
		&types.Label{},
		&types.Ingredient{},
		&types.Product{},
	)
}

func Database() *gorm.DB {
	return db
}

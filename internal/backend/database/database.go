package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
}

func Database() *gorm.DB {
	return db
}

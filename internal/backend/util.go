package backend

import (
	"gorm.io/gorm"
)

func WithOrder(field string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(field)
	}
}

func WithSearch(field, term string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" LIKE ?", term+"%")
	}
}

func WithPreload(associations ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, a := range associations {
			db = db.Preload(a)
		}
		return db
	}
}

func WithAnyOf(field string, ids []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" IN ?", ids)
	}
}

func WithCursor(cursor string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if cursor == "" {
			return db
		}
		return db.Where("id > ?", cursor)
	}
}

func WithLimit(n int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(n + 1) // fetch one extra to know if there's a next page
	}
}

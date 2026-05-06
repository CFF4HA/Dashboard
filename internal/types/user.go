package types

import (
	"github.com/google/uuid"
)

type User struct {
	Model
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"password_hash" gorm:"not null"`

	// The following lists are "Favorite" products and
	// ingredients.
	Products    []Product    `json:"products" gorm:"many2many:user_products;"`
	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:user_ingredients;"`
}

type Session struct {
	Session uuid.UUID `json:"session"`
}

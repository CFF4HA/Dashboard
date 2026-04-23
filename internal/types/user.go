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

	// The following fields will be used to determine a user's
	// ingredient sensitivities. A "Good" product is one that
	// doesn't cause the user any reactions.
	GoodProducts   []Product    `json:"good_products" gorm:"many2many:user_good_products;"`
	GoodIngredient []Ingredient `json:"good_ingredients" gorm:"many2many:user_good_ingredients;"`

	//A "Bad" product is one that cases the user some reaction.
	BadProducts     []Product    `json:"bad_products" gorm:"many2many:user_bad_products;"`
	BadIngreidients []Ingredient `json:"bad_ingredients" gorm:"many2many:user_bad_ingredients;"`
}

type Session struct {
	Session uuid.UUID `json:"session"`
}

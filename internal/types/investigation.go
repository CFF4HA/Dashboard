package types

import (
	"github.com/google/uuid"
)

// An investigation contains a list of products that the
// user is comparing. There are a series of attributes that
// are generated from the list:
// - Shared Ingredients
// - Shared Hazards
// - Shared Effects
// - Shared Symptoms
// - Shared Regulations
// - Ingredients Unique to Each Product
// - Hazards Unique to Each Product
// - Effects Unique to Each Product
// - Symptoms Unique to Each Product
// - Regulations Unique to Each Product
type Investigation struct {
	Model

	ProductList []Product
	Shared      struct {
		Ingredients []Ingredient
		Hazards     []string
		Effects     []string
		Symptoms    []string
		Regulations []string
	}

	Unique map[uuid.UUID]struct {
		Ingredients []Ingredient
		Hazards     []string
		Effects     []string
		Symptoms    []string
		Regulations []string
	}

	// The following is going to be used to investigate
	// potential root cause for harmful effects.
	GoodProducts   []uuid.UUID
	BadProducts    []uuid.UUID
	BadIngredients []uuid.UUID
}

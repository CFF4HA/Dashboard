package types

import ()

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

	SharedIngredients []Ingredient
	SharedHazards     []Label
	SHaredEffects     []Label
	SharedSymptoms    []Label
	SharedRegulations []Label
}

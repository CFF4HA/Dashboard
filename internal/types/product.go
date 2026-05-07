package types

import "github.com/google/uuid"

// A product can be thought of as a list of ingredients with relevant metrics that
// accompany it, all stemming from the list of ingredients.
type Product struct {
	// this is a basic set of database model columns like ID, Created, Updated, and
	// Deleted.
	Model
	UserId *uuid.UUID `json:"user_id" gorm:"type:uuid;"`

	// The following is a name for the product, its origin (e.g. country of origin), and a list of ingredients
	// (the list of ingredients is a many-to-many relationship, and it is managed via GORM).
	IsPublic    bool         `json:"is_public" gorm:"default:false"`
	Name        string       `json:"name" gorm:"type:varchar(255);not null;"`
	Origin      *string      `json:"origin" gorm:"type:varchar(255);default:null"`
	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:product_ingredients;"`
	Tags        []Tag        `json:"tags" gorm:"many2many:product_tags;"`
}

type ProductDraftAutomated struct {
	Name        string   `json:"name"`
	Origin      string   `json:"origin,omitempty"`
	Ingredients []string `json:"ingredients"`
}

type ProductComparison struct {
	SharedIngredients []Ingredient        `json:"shared_ingredients"`
	SharedLabels      map[string][]string `json:"shared_labels"`

	UniqueIngToProduct map[string][]Ingredient `json:"unique_to_product"`
	UniqueLabToProduct map[string][]string     `json:"unique_labels_to_product"`
}

type InvestigationShared struct {
	Ingredients []Ingredient
	Hazards     []string
	Effects     []string
	Symptoms    []string
	Regulations []string
}

type InvestigationUnique struct {
	Ingredients []Ingredient
	Hazards     []string
	Effects     []string
	Symptoms    []string
	Regulations []string
}


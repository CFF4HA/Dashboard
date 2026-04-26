package types

import (
	"github.com/google/uuid"
)

// A product can be thought of as a list of ingredients with relevant metrics that
// accompany it, all stemming from the list of ingredients.
type Product struct {
	// this is a basic set of database model columns like ID, Created, Updated, and
	// Deleted.
	Model

	// The following is a name for the product, its origin (e.g. country of origin), and a list of ingredients
	// (the list of ingredients is a many-to-many relationship, and it is managed via GORM).
	Name        string       `json:"name" gorm:"type:varchar(255);not null;"`
	Origin      *string      `json:"origin" gorm:"type:varchar(255);default:null"`
	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:product_ingredients;"`
	Tags        []Tag        `json:"tags" gorm:"many2many:product_tags;"`

	// The product metadata is another object that contains everything from
	// metrics like number of renders, to AI information blurbs, count of
	// hazards, etc.
	Metadata ProductMetadata `json:"metadata" gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ProductMetadata struct {
	Model
	ProductId uuid.UUID `json:"product_id" gorm:"type:uuid;not null;uniqueIndex;"`

	NumHazard     int `json:"hazard_count" gorm:"type:int;default:0"`
	NumEffect     int `json:"effect_count" gorm:"type:int;default:0"`
	NumSymptom    int `json:"symptom_count" gorm:"type:int;default:0"`
	NumReglations int `json:"regulation_count" gorm:"type:int;default:0"`

	AiBlurb *string `json:"ai_blurb" gorm:"type:text;default:null"`
}

package types

import (
	"github.com/google/uuid"
)

// An Ingredient can have multiple names, multiple hazards, etc. We will
// generally access ingredients by looking up the mapping in a name table.
type Ingredient struct {
	Model
	Failed      bool   `json:"failed" gorm:"type:boolean;not null;default:false"`
	PrimaryName string `json:"primary_name" gorm:"type:text;not null;index;"`

	Labels   []Label            `json:"labels" gorm:"many2many:ingredient_labels;"`
	Tags     []Tag              `json:"tags" gorm:"many2many:ingredient_tags;"`
	Names    []Name             `json:"names" gorm:"many2many:ingredient_names"`
	Metadata IngredientMetadata `json:"metadata" gorm:"foreignKey:IngredientId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type IngredientMetadata struct {
	Model
	IngredientId uuid.UUID `json:"ingredient_id" gorm:"type:uuid;not null;uniqueIndex;"`

	NumHazards     int `json:"num_hazards" gorm:"type:int;default:0"`
	NumEffects     int `json:"num_effects" gorm:"type:int;default:0"`
	NumSymptoms    int `json:"num_symptoms" gorm:"type:int;default:0"`
	NumRegulations int `json:"num_regulations" gorm:"type:int;default:0"`

	AiBlurb *string `json:"ai_blurb" gorm:"type:text;default:null"`
}

// This is the Name table, which is how we will access ingredients.
type Name struct {
	Model

	Text        string       `json:"text" gorm:"type:text;not null;unique;"`
	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:ingredient_names;"`
}

// We should use labels as a way to describe a hazard, symptom, effect, regulation, etc.
// If we need more information, we can add it to a metadata table. This generalization
// should suffice for now.
type Label struct {
	Model

	Type    string  `json:"type" gorm:"type:text;not null;check:type IN ('hazard', 'symptom', 'general', 'effect', 'regulation');default:'general';"`
	Payload string  `json:"name" gorm:"type:text;not null"`
	Origin  *string `json:"origin" gorm:"type:text;default:null"`

	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:ingredient_labels;"`
}

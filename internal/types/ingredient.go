package types

import (
	"time"

	"github.com/google/uuid"
)

// An Ingredient can have multiple names, multiple hazards, etc. We will
// generally access ingredients by looking up the mapping in a name table.
type Ingredient struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Created time.Time `json:"created" gorm:"type:timestamp;not null;default:current_timestamp"`
	Updated time.Time `json:"updated" gorm:"type:timestamp;not null;default:current_timestamp"`

	Labels []Label `json:"labels" gorm:"many2many:ingredient_labels;"`
}

// This is the Name table, which is how we will access ingredients.
type Name struct {
	Text       string     `json:"text" gorm:"type:text;not null;primaryKey;unique;"`
	Ingredient Ingredient `json:"ingredient" gorm:"foreignKey:IngredientId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// We should use labels as a way to describe a hazard, symptom, effect, regulation, etc.
// If we need more information, we can add it to a metadata table. This generalization
// should suffice for now.
type Label struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Type    string    `json:"type" gorm:"type:varchar(255);not null;check:type IN ('hazard', 'symptom', 'general', 'effect', 'regulation');default:'general';"`
	Payload string    `json:"name" gorm:"type:varchar(255);not null"`
	Origin  string    `json:"origin" gorm:"type:varchar(255);not null"`

	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:ingredient_labels;"`
}

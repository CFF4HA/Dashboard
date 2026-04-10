package types

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	ingredient = Ingredient{
		Id:      uuid.New(),
		Created: time.Now().Add(-time.Hour),
		Updated: time.Now(),
		Names: []Name{
			{
				Text: "Formaldehyde",
			},
			{
				Text: "Methanal",
			},
		},
		Labels: []Label{
			{
				Id:      uuid.New(),
				Type:    "hazard",
				Payload: "Carcinogenic",
				Origin:  "http://example.com/carcinogenic",
			},
			{
				Id:      uuid.New(),
				Type:    "symptom",
				Payload: "Skin Irritation",
				Origin:  "http://example.com/skin-irritation",
			},
			{
				Id:      uuid.New(),
				Type:    "effect",
				Payload: "May cause cancer",
				Origin:  "http://example.com/may-cause-cancer",
			},
		},
	}
)

// An Ingredient can have multiple names, multiple hazards, etc. We will
// generally access ingredients by looking up the mapping in a name table.
type Ingredient struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Created time.Time `json:"created" gorm:"type:timestamp;not null;default:current_timestamp"`
	Updated time.Time `json:"updated" gorm:"type:timestamp;not null;default:current_timestamp"`
	Failed  bool      `json:"failed" gorm:"type:boolean;not null;default:false"`

	Labels []Label `json:"labels" gorm:"many2many:ingredient_labels;"`
	Names  []Name  `json:"names" gorm:"foreignKey:IngredientId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// This is the Name table, which is how we will access ingredients.
type Name struct {
	Text string `json:"text" gorm:"type:text;not null;primaryKey;unique;"`

	IngredientId uuid.UUID  `json:"ingredient_id" gorm:"type:uuid;not null;index;"`
	Ingredient   Ingredient `json:"ingredient" gorm:"foreignKey:IngredientId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// We should use labels as a way to describe a hazard, symptom, effect, regulation, etc.
// If we need more information, we can add it to a metadata table. This generalization
// should suffice for now.
type Label struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type    string    `json:"type" gorm:"type:text;not null;check:type IN ('hazard', 'symptom', 'general', 'effect', 'regulation');default:'general';"`
	Payload string    `json:"name" gorm:"type:text;not null"`
	Origin  string    `json:"origin" gorm:"type:text;not null"`

	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:ingredient_labels;"`
}

func (l Label) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	return Label{
		Id:          uuid.New(),
		Type:        "hazard",
		Payload:     "Carcinogenic",
		Origin:      "http://example.com/carcinogenic",
		Ingredients: []Ingredient{ingredient, ingredient, ingredient},
	}, nil
}

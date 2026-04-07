package types

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ProductDraftIngredient struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Exists bool      `json:"exists"`
	Failed bool      `json:"failed" gorm:"type:boolean;not null;default:false"`
}

type ProductDraft struct {
	Ingredients []ProductDraftIngredient `json:"ingredients"`
}

type Product struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Created time.Time `json:"created" gorm:"type:timestamp;not null;default:current_timestamp"`
	Updated time.Time `json:"updated" gorm:"type:timestamp;not null;default:current_timestamp"`
	Name    string    `json:"name" gorm:"type:varchar(255);not null;"`
	Origin  string    `json:"origin" gorm:"type:varchar(255);default:'None'"`

	UserUid     uuid.UUID    `json:"user_uid" gorm:"type:uuid;"`
	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:product_ingredients;"`
}

func (p Product) Data(w http.ResponseWriter, r *http.Request) (any, error) {

	return &Product{
		Id:      uuid.New(),
		Created: time.Now().Add(-time.Hour),
		Updated: time.Now(),
		Name:    "Example Product",
		Origin:  "Example Origin",
		UserUid: uuid.Nil,
		Ingredients: []Ingredient{
			ingredient,
		},
	}, nil
}

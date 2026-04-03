package types

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Created time.Time `json:"created" gorm:"type:timestamp;not null;default:current_timestamp"`
	Updated time.Time `json:"updated" gorm:"type:timestamp;not null;default:current_timestamp"`
	Name    string    `json:"name" gorm:"type:varchar(255);not null;"`
	Origin  string    `json:"origin" gorm:"type:varchar(255);default:'None'"`

	UserUid     uuid.UUID    `json:"user_uid" gorm:"type:uuid;"`
	Ingredients []Ingredient `json:"ingredients" gorm:"many2many:product_ingredients;"`
}

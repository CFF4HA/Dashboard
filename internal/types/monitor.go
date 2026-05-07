package types

import (
	"time"
)

type Monitor struct {
	Model

	IngredientId          string    `json:"ingredient_id"`
	UserId                string    `json:"user_id"`
	LastKnownUpdatedField time.Time `json:"last_known_updated_field"`
}

package ingredient

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func RetrieveIngredientHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	name := strings.ToLower(strings.TrimSpace(r.FormValue("name")))
	if name == "" {
		return errors.New("ingredient name is required")
	}

	// Return early if ingredient already exists by primary name.
	var existing types.Ingredient
	if tx := core.DB.Where("primary_name ILIKE ?", name).Preload("Labels").Preload("Names").First(&existing); tx.Error == nil {
		return nil
	}

	_, err := RetrieveIngredient(name)
	return err
}

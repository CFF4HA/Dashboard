package ingredient

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/tagging"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
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

	ing, err := RetrieveIngredient(name)
	if err != nil {
		return err
	}

	// Run the requesting user's tagging rules against the new ingredient.
	if ing != nil && !ing.Failed {
		if u, userErr := user.GetUserFromRequest(w, r); userErr == nil {
			tagging.TagNewIngredient(*ing, u.Model.Id)
		}
	}

	return nil
}

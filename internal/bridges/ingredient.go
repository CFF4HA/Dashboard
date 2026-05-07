package bridges

import (
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
)

type IngredientsByName struct {
}

func (i IngredientsByName) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	ingredients, err := backend.GetIngredientsByPrimaryName(r.FormValue("name"), r.FormValue("cursor"))
	if err != nil {
		return nil, errors.New("failed to get ingredients by name: " + err.Error())
	}

	return ingredients, nil
}

func (i IngredientsByName) Name() string {
	return "Ingredients"
}

package bridges

import (
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
)

var (
	IngredientById = DataBridge{Key: "Ingredient", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		return backend.GetIngredientById(r.FormValue("id"))
	}}

	IngredientNotesByIngredientId = DataBridge{Key: "IngredientNotesByIngredientId", Func: func(w http.ResponseWriter, r *http.Request) (any, error) {
		u, err := user.GetUserFromRequestNoRedirect(w, r)
		if err != nil || u == nil {
			return nil, errors.New("user not authenticated")
		}

		return backend.GetIngredientNotes(r.FormValue("id"), u.Id, r.FormValue("cursor"))
	}}
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

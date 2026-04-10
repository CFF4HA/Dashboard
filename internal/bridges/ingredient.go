package bridges

import (
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
)

func IngredientByNameProvider(r *http.Request, model map[string]any) (any, error) {
	name := strings.Trim(r.FormValue("name"), " ")
	force := strings.Trim(r.FormValue("force"), " ") != ""

	core.Logger.Info("ingredients bridge", "force", force)
	v, err := database.PollForIngredient(name, 10, 2, force)
	if err != nil {
		return nil, err
	}

	return v, nil
}

package bridges

import (
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func IngredientByNameProvider(r *http.Request, model map[string]any) (any, error) {
	name := strings.Trim(r.FormValue("name"), " ")
	force := strings.Trim(r.FormValue("force"), " ") != ""

	core.Logger.Info("ingredients bridge", "force", force)
	v, err := database.PollForIngredient(name, 10, 2, force)
	if err != nil {
		return nil, err
	}

	model["Ingredient"] = v
	return v, nil
}

func IngredientMetadataProvider(r *http.Request, model map[string]any) (any, error) {
	var meta types.IngredientMetadata

	v, ok := model["Ingredient"]
	if !ok {
		return nil, nil
	}

	ingrdient, ok := v.(*types.Ingredient)
	if !ok {
		return nil, nil
	}

	counts := map[string]int{}
	for _, label := range ingrdient.Labels {
		switch label.Type {
		case "hazard":
			counts["hazard"]++
		case "effect":
			counts["effect"]++
		case "symptom":
			counts["symptom"]++
		case "general":
			counts["general"]++
		case "regulation":
			counts["regulation"]++
		}
	}

	meta.CountHazards = counts["hazard"]
	meta.CountEffects = counts["effect"]
	meta.CountSymptoms = counts["symptom"]
	meta.CountGeneral = counts["general"]
	meta.CountRegulations = counts["regulation"]

	return meta, nil
}

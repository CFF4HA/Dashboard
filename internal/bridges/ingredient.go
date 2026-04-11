package bridges

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func SyncDatabaseWithIngredient(name string) {
	// this will do a request to the python backend that should trigger a
	// scrape of PubChem
	request, err := http.NewRequest("GET", core.BackendAddress+"/ingredient?name="+url.QueryEscape(name), nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	core.Logger.Info("synced database with ingredient", "name", name)
}

// this ingredient function will serve as the bridge between the
// server and the database. It expects the existence of an "id",
// or "name" query parameter.
func Ingredient(r *http.Request) (any, error) {
	ingredient := &types.Ingredient{}
	db := database.Database()
	id := strings.Trim(r.FormValue("id"), " ")
	name := strings.ToLower(strings.Trim(r.FormValue("name"), " "))

	// this is the search by ID path and takes
	// precence over the name search path
	if id != "" {
		err := db.Preload("Names").Preload("Labels").First(ingredient, "id = ?", id).Error
		if err != nil {
			go SyncDatabaseWithIngredient(name)

			core.Logger.Info("ingredient not found in database, syncing with backend", "name", name)
			return nil, err
		}

		return ingredient, nil
	}

	// the following assumes that you are searching by primary name only
	err := db.Preload("Names").Preload("Labels").First(ingredient, "primary_name = ?", name).Error
	if err != nil {
		go SyncDatabaseWithIngredient(name)

		core.Logger.Info("ingredient not found in database, syncing with backend", "name", name)
		return nil, err
	}

	core.Logger.Info("ingredient found in database", "name", name)
	return ingredient, nil
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

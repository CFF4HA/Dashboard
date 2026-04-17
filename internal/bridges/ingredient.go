package bridges

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

// SyncDatabaseWithIngredient triggers a PubChem scrape on the Python backend
// for the given ingredient name, populating the local database.
func SyncDatabaseWithIngredient(name string) {
	req, err := http.NewRequest("GET", core.BackendAddress+"/ingredient?name="+url.QueryEscape(name), nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	core.Logger.Info("synced database with ingredient", "name", name)
}

// Ingredient is the acquisition policy used by the cached ingredient bridge.
// It looks up an ingredient by ID (preferred) or primary_name, triggering a
// backend sync if the record is not found locally.
func Ingredient(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	db, err := database.Database()
	if err != nil {
		return nil, nil
	}

	id := strings.Trim(r.FormValue("id"), " ")
	name := strings.ToLower(strings.Trim(r.FormValue("primary_name"), " "))

	ingredient := &types.Ingredient{}

	// ID takes precedence over name when both are present.
	if id != "" {
		if err := db.Preload("Names").Preload("Labels").First(ingredient, "id = ?", id).Error; err != nil {
			SyncDatabaseWithIngredient(name)
			core.Logger.Info("ingredient not found by ID, syncing with backend", "id", id)
			return nil, err
		}
		return ingredient, nil
	}

	// Fall back to primary_name lookup.
	if err := db.Preload("Names").Preload("Labels").First(ingredient, "primary_name = ?", name).Error; err != nil {
		SyncDatabaseWithIngredient(name)
		core.Logger.Info("ingredient not found by name, syncing with backend", "name", name)
		return nil, err
	}

	core.Logger.Info("ingredient found in database", "name", name)
	return ingredient, nil
}

// IngredientMetadataProvider counts label types on the ingredient stored under
// KeyIngredient in the model map, returning an IngredientMetadata summary.
func IngredientMetadataProvider(r *http.Request, model map[string]any) (any, error) {
	v, ok := model[KeyIngredient]
	if !ok {
		return nil, nil
	}

	ingredient, ok := v.(*types.Ingredient)
	if !ok {
		return nil, nil
	}

	counts := map[string]int{}
	for _, label := range ingredient.Labels {
		counts[label.Type]++
	}

	return types.IngredientMetadata{
		CountHazards:     counts["hazard"],
		CountEffects:     counts["effect"],
		CountSymptoms:    counts["symptom"],
		CountGeneral:     counts["general"],
		CountRegulations: counts["regulation"],
	}, nil
}

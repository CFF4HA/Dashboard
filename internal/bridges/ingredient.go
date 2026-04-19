package bridges

import (
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
)

type IngredientSearchResult struct {
	Ingredients   []types.Ingredient
	Query         string
	ResultsTarget string // CSS selector for the container that holds these results
}

var IngredientSearchBridge = verb.Map("Search", func(r *http.Request, m map[string]any) (any, error) {
	q := strings.TrimSpace(r.FormValue("query"))

	var ingredients []types.Ingredient

	if q == "" {
		tx := core.DB.Preload("Metadata").Preload("Names").Preload("Labels").
			Order("created DESC").Limit(20).Find(&ingredients)
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else {
		// Two-step: get matching IDs via JOIN, then preload associations separately
		// to avoid GROUP BY conflicts with GORM's Preload.
		var ids []string
		tx := core.DB.Model(&types.Ingredient{}).
			Select("DISTINCT ingredients.id").
			Joins("LEFT JOIN ingredient_names ing_n ON ing_n.ingredient_id = ingredients.id").
			Joins("LEFT JOIN names n ON n.id = ing_n.name_id").
			Where("ingredients.primary_name ILIKE ? OR n.text ILIKE ?", "%"+q+"%", "%"+q+"%").
			Pluck("id", &ids)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if len(ids) > 0 {
			tx = core.DB.Preload("Metadata").Preload("Names").Preload("Labels").
				Where("id IN ?", ids).Find(&ingredients)
			if tx.Error != nil {
				return nil, tx.Error
			}
		}
	}

	target := r.FormValue("results_target")
	if target == "" {
		target = "#search-results"
	}

	return IngredientSearchResult{Ingredients: ingredients, Query: q, ResultsTarget: target}, nil
})

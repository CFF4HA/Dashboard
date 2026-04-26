package bridges

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	"github.com/google/uuid"
)

type CategorizedIngredients struct {
}

func (c CategorizedIngredients) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
	var Payload struct {
		Safe   []uuid.UUID `json:"safe"`
		Unsafe []uuid.UUID `json:"unsafe"`
	}

	user, ok := m["User"]
	if !ok {
		core.Logger.Warn("No user in context")
		return nil, nil
	} else if _, ok := user.(*types.User); !ok {
		core.Logger.Warn("User in context is not of type User")
		return nil, nil
	}
	user_object := user.(*types.User)

	if core.DB.Preload("GoodIngredient").Preload("BadIngreidients").First(&user_object, "id = ?", user_object.Model.Id).Error != nil {
		return nil, errors.New("failed to load user with ingredients")
	}
	for _, user_good := range user_object.GoodIngredient {
		Payload.Safe = append(Payload.Safe, user_good.Model.Id)
	}

	for _, user_bad := range user_object.BadIngreidients {
		Payload.Unsafe = append(Payload.Unsafe, user_bad.Model.Id)
	}

	return Payload, nil
}

func (c CategorizedIngredients) Name() string {
	return "Categories"
}

type IngredientSearchResult struct {
	Ingredients   []types.Ingredient
	Query         string
	ResultsTarget string // CSS selector for the container that holds these results
}

var IngredientSearchBridge = verb.Map("Search", func(r *http.Request, m map[string]any) (any, error) {
	q := strings.TrimSpace(r.FormValue("query"))
	labelOp := strings.TrimSpace(r.FormValue("label_filter_op"))
	labelText := strings.TrimSpace(r.FormValue("label_filter_text"))
	tagFilters := r.Form["tag_filter"]

	var ingredients []types.Ingredient
	preload := core.DB.Preload("Tags").Preload("Metadata").Preload("Names").Preload("Labels")

	// Fast path — no filters at all, return most recent.
	if q == "" && labelText == "" && len(tagFilters) == 0 {
		if tx := preload.Order("failed ASC").Order("created DESC").Limit(20).Find(&ingredients); tx.Error != nil {
			return nil, tx.Error
		}
		target := r.FormValue("results_target")
		if target == "" {
			target = "#search-results"
		}
		return IngredientSearchResult{Ingredients: ingredients, Query: q, ResultsTarget: target}, nil
	}

	// Two-step (see AGENTS.md): DISTINCT IDs via JOIN, then Preload with IN.
	idQuery := core.DB.Model(&types.Ingredient{}).Select("DISTINCT ingredients.id")

	if q != "" {
		idQuery = idQuery.
			Joins("LEFT JOIN ingredient_names ing_n ON ing_n.ingredient_id = ingredients.id").
			Joins("LEFT JOIN names n ON n.id = ing_n.name_id").
			Where("ingredients.primary_name ILIKE ? OR n.text ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if labelText != "" {
		// Sub-query: ingredient IDs carrying at least one matching label.
		labelSub := core.DB.Table("ingredient_labels").
			Select("DISTINCT ingredient_labels.ingredient_id").
			Joins("JOIN labels l ON l.id = ingredient_labels.label_id").
			Where("l.payload ILIKE ?", "%"+labelText+"%")

		switch labelOp {
		case "excludes":
			idQuery = idQuery.Where("ingredients.id NOT IN (?)", labelSub)
		default: // "contains"
			idQuery = idQuery.Where("ingredients.id IN (?)", labelSub)
		}
	}

	if len(tagFilters) > 0 {
		tagSub := core.DB.Table("ingredient_tags").
			Select("DISTINCT ingredient_tags.ingredient_id").
			Where("ingredient_tags.tag_id IN ?", tagFilters)
		idQuery = idQuery.Where("ingredients.id IN (?)", tagSub)
	}

	var ids []string
	if tx := idQuery.Pluck("id", &ids); tx.Error != nil {
		return nil, tx.Error
	}

	if len(ids) > 0 {
		if tx := preload.Where("id IN ?", ids).Order("failed ASC").Order("created DESC").Find(&ingredients); tx.Error != nil {
			return nil, tx.Error
		}
	}

	target := r.FormValue("results_target")
	if target == "" {
		target = "#search-results"
	}

	return IngredientSearchResult{Ingredients: ingredients, Query: q, ResultsTarget: target}, nil
})

var IngredientDetailBridge = verb.Map("Ingredient", func(r *http.Request, m map[string]any) (any, error) {
	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid ingredient id")
	}

	var ingredient types.Ingredient
	tx := core.DB.Preload("Metadata").Preload("Tags").Preload("Names").Preload("Labels").First(&ingredient, "id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ingredient, nil
})

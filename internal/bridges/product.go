package bridges

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	verbs_gorm "github.com/DAlba-sudo/verbs/gorm"
	"github.com/google/uuid"
)

type CategorizedProducts struct{}

func (c CategorizedProducts) Data(w http.ResponseWriter, r *http.Request, m map[string]any) (any, error) {
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

	if core.DB.Preload("GoodProducts").Preload("BadProducts").Preload("BadIngreidients").First(&user_object, "id = ?", user_object.Model.Id).Error != nil {
		return nil, errors.New("failed to load user with products")
	}

	// Explicit bad products.
	for _, p := range user_object.BadProducts {
		Payload.Unsafe = append(Payload.Unsafe, p.Model.Id)
	}

	// Products that contain at least one of the user's bad ingredients are also unsafe.
	var badIngIDs []uuid.UUID
	for _, ing := range user_object.BadIngreidients {
		badIngIDs = append(badIngIDs, ing.Model.Id)
	}
	if len(badIngIDs) > 0 {
		var derived []uuid.UUID
		core.DB.Table("product_ingredients").
			Select("DISTINCT product_id").
			Where("ingredient_id IN ?", badIngIDs).
			Pluck("product_id", &derived)
		for _, pid := range derived {
			if !slices.Contains(Payload.Unsafe, pid) {
				Payload.Unsafe = append(Payload.Unsafe, pid)
			}
		}
	}

	// Explicit good products, excluding any now derived as unsafe.
	for _, p := range user_object.GoodProducts {
		if !slices.Contains(Payload.Unsafe, p.Model.Id) {
			Payload.Safe = append(Payload.Safe, p.Model.Id)
		}
	}

	return Payload, nil
}

func (c CategorizedProducts) Name() string {
	return "Categories"
}

func ProductBridge(limit int, modifiers map[string]string) verb.Bridge {
	return verbs_gorm.GORM("Products", types.Product{}, types.Product{}, core.DB, verbs_gorm.GormOptions{
		Preload:      []string{"Ingredients", "Metadata"},
		Limit:        limit,
		KeyModifiers: modifiers,
	})
}

type ProductSearchResult struct {
	Products []types.Product
	Query    string
}

var ProductSearchBridge = verb.Map("Search", func(r *http.Request, m map[string]any) (any, error) {
	q := strings.TrimSpace(r.FormValue("query"))
	labelOp := strings.TrimSpace(r.FormValue("label_filter_op"))
	labelText := strings.TrimSpace(r.FormValue("label_filter_text"))

	var products []types.Product
	db := core.DB.Preload("Ingredients").Preload("Metadata")

	// Fast path — no filters at all, return most recent.
	if q == "" && labelText == "" {
		if tx := db.Order("created DESC").Limit(20).Find(&products); tx.Error != nil {
			return nil, tx.Error
		}
		return ProductSearchResult{Products: products, Query: q}, nil
	}

	// Two-step pattern (see AGENTS.md): DISTINCT IDs via JOIN, then Preload with IN.
	idQuery := core.DB.Model(&types.Product{}).Select("DISTINCT products.id")

	if q != "" {
		idQuery = idQuery.Where("products.name ILIKE ?", "%"+q+"%")
	}

	if labelText != "" {
		// Sub-query: the set of product IDs that have at least one ingredient
		// whose label payload matches the user-supplied text.
		labelSub := core.DB.Table("product_ingredients").
			Select("DISTINCT product_ingredients.product_id").
			Joins("JOIN ingredient_labels il ON il.ingredient_id = product_ingredients.ingredient_id").
			Joins("JOIN labels l ON l.id = il.label_id").
			Where("l.payload ILIKE ?", "%"+labelText+"%")

		switch labelOp {
		case "excludes":
			idQuery = idQuery.Where("products.id NOT IN (?)", labelSub)
		default: // "contains" (default behaviour)
			idQuery = idQuery.Where("products.id IN (?)", labelSub)
		}
	}

	var ids []string
	if tx := idQuery.Pluck("id", &ids); tx.Error != nil {
		return nil, tx.Error
	}

	if len(ids) > 0 {
		if tx := db.Where("id IN ?", ids).Find(&products); tx.Error != nil {
			return nil, tx.Error
		}
	}

	return ProductSearchResult{Products: products, Query: q}, nil
})

var ProductDetailBridge = verb.Map("Product", func(r *http.Request, m map[string]any) (any, error) {
	idStr := strings.TrimSpace(r.FormValue("id"))
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid product id")
	}

	var product types.Product
	tx := core.DB.
		Preload("Metadata").
		Preload("Ingredients").
		Preload("Ingredients.Labels").
		Preload("Ingredients.Names").
		Preload("Ingredients.Metadata").
		First(&product, "id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return product, nil
})

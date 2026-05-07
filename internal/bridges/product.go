package bridges

import (
	"errors"
	"net/http"
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

var ProductSearchBridgeV2 = verb.Map("SearchV2", func(r *http.Request, m map[string]any) (any, error) {
	product_name := strings.TrimSpace(r.FormValue("product_name"))
	product_tags := r.Form["product_tags"]

	var ingredients []string
	for _, s := range strings.Split(r.FormValue("ingredients"), ",") {
		if s = strings.TrimSpace(s); s != "" {
			ingredients = append(ingredients, s)
		}
	}

	var products []types.Product
	db := core.DB.Preload("Ingredients").Preload("Metadata").Preload("Tags")

	// Fast path — no filters, return most recent.
	if product_name == "" && len(ingredients) == 0 && len(product_tags) == 0 {
		if tx := db.Order("created DESC").Limit(20).Find(&products); tx.Error != nil {
			return nil, tx.Error
		}
		return ProductSearchResult{Products: products, Query: product_name}, nil
	}

	// Two-step pattern: DISTINCT IDs via JOIN, then Preload with IN.
	idQuery := core.DB.Model(&types.Product{}).Select("DISTINCT products.id")

	if product_name != "" {
		idQuery = idQuery.Where("products.name ILIKE ?", "%"+product_name+"%")
	}

	if len(ingredients) > 0 {
		// Match ingredients by primary_name OR any synonym (names table).
		ingSub := core.DB.Table("product_ingredients").
			Select("DISTINCT product_ingredients.product_id").
			Joins("JOIN ingredients ON ingredients.id = product_ingredients.ingredient_id").
			Joins("LEFT JOIN ingredient_names ON ingredient_names.ingredient_id = ingredients.id").
			Joins("LEFT JOIN names ON names.id = ingredient_names.name_id")

		conds := make([]string, 0, len(ingredients))
		args := make([]interface{}, 0, len(ingredients)*2)
		for _, ing := range ingredients {
			conds = append(conds, "(ingredients.primary_name ILIKE ? OR names.text ILIKE ?)")
			args = append(args, "%"+ing+"%", "%"+ing+"%")
		}
		ingSub = ingSub.Where(strings.Join(conds, " OR "), args...)
		idQuery = idQuery.Where("products.id IN (?)", ingSub)
	}

	if len(product_tags) > 0 {
		tagSub := core.DB.Table("product_tags").
			Select("DISTINCT product_tags.product_id").
			Where("product_tags.tag_id IN ?", product_tags)
		idQuery = idQuery.Where("products.id IN (?)", tagSub)
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

	return ProductSearchResult{Products: products, Query: product_name}, nil
})

var ProductSearchBridge = verb.Map("Search", func(r *http.Request, m map[string]any) (any, error) {
	q := strings.TrimSpace(r.FormValue("query"))
	labelOp := strings.TrimSpace(r.FormValue("label_filter_op"))
	labelText := strings.TrimSpace(r.FormValue("label_filter_text"))
	tagFilters := r.Form["tag_filter"]

	var products []types.Product
	db := core.DB.Preload("Ingredients").Preload("Metadata").Preload("Tags")

	// Fast path — no filters at all, return most recent.
	if q == "" && labelText == "" && len(tagFilters) == 0 {
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

	if len(tagFilters) > 0 {
		tagSub := core.DB.Table("product_tags").
			Select("DISTINCT product_tags.product_id").
			Where("product_tags.tag_id IN ?", tagFilters)
		idQuery = idQuery.Where("products.id IN (?)", tagSub)
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
		Preload("Tags").
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

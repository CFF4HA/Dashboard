package bridges

import (
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	verbs_gorm "github.com/DAlba-sudo/verbs/gorm"
)

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

	var products []types.Product

	db := core.DB.Preload("Ingredients").Preload("Metadata")

	if q == "" {
		if tx := db.Order("created DESC").Limit(20).Find(&products); tx.Error != nil {
			return nil, tx.Error
		}
	} else {
		if tx := db.Where("name ILIKE ?", "%"+q+"%").Find(&products); tx.Error != nil {
			return nil, tx.Error
		}
	}

	return ProductSearchResult{Products: products, Query: q}, nil
})

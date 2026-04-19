package bridges

import (
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verbs/gorm"
)

var ()

func ProductBridge(limit int, modifiers map[string]string) verb.Bridge {
	return gorm.GORM("Products", types.Product{}, types.Product{}, core.DB, gorm.GormOptions{
		Preload:      []string{"Ingredients", "Metadata"},
		Limit:        limit,
		KeyModifiers: modifiers,
	})
}

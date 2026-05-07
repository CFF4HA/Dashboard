package ingredient

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func DeleteIngredientHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")

	// Explicitly remove join table rows before deleting the ingredient so FK
	// constraints don't block the delete on databases without cascade DDL.
	core.DB.Exec("DELETE FROM ingredient_labels WHERE ingredient_id = ?", id)
	core.DB.Exec("DELETE FROM ingredient_names WHERE ingredient_id = ?", id)
	core.DB.Exec("DELETE FROM ingredient_tags WHERE ingredient_id = ?", id)
	core.DB.Exec("DELETE FROM product_ingredients WHERE ingredient_id = ?", id)

	// IngredientMetadata cascades via the FK constraint defined on the struct.
	return core.DB.Unscoped().Delete(&types.Ingredient{}, "id = ?", id).Error
}

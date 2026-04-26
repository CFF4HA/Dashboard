package ingredient

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

// SyncIngredientHandler deletes the existing ingredient record and re-fetches
// it from PubChem, preserving any product associations the ingredient had.
func SyncIngredientHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")

	var existing types.Ingredient
	if tx := core.DB.First(&existing, "id = ?", id); tx.Error != nil {
		return tx.Error
	}

	// Capture which products contained this ingredient so we can restore them.
	var productIDs []string
	core.DB.Table("product_ingredients").
		Select("product_id").
		Where("ingredient_id = ?", id).
		Pluck("product_id", &productIDs)

	// Clean up join tables and delete the ingredient.
	core.DB.Exec("DELETE FROM ingredient_labels WHERE ingredient_id = ?", id)
	core.DB.Exec("DELETE FROM ingredient_names WHERE ingredient_id = ?", id)
	core.DB.Exec("DELETE FROM ingredient_tags WHERE ingredient_id = ?", id)
	core.DB.Exec("DELETE FROM product_ingredients WHERE ingredient_id = ?", id)
	core.DB.Unscoped().Delete(&existing)

	// Re-fetch from PubChem under the same primary name.
	newIng, err := RetrieveIngredient(existing.PrimaryName)
	if err != nil {
		return err
	}

	// Restore product associations using the new ingredient's ID.
	for _, pid := range productIDs {
		core.DB.Exec(
			"INSERT INTO product_ingredients (product_id, ingredient_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			pid, newIng.Id,
		)
	}

	return nil
}

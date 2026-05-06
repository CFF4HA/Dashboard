package product

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/ingredient"
	"github.com/CFF4HA/Dashboard/internal/types"
)

// UpdateProductHandler handles POST /product — updates an existing product's
// name, origin, and ingredient list, recalculating all metadata counts.
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	idStr := strings.TrimSpace(r.FormValue("id"))
	name := strings.TrimSpace(r.FormValue("name"))
	origin := strings.TrimSpace(r.FormValue("origin"))
	ingredientNames := r.Form["ingredient_names"]

	if idStr == "" {
		return errors.New("product id is required")
	}
	if len(name) == 0 {
		return errors.New("product name is required")
	}
	if len(ingredientNames) == 0 {
		return errors.New("at least one ingredient name is required")
	}

	var product types.Product
	if tx := core.DB.Preload("Metadata").First(&product, "id = ?", idStr); tx.Error != nil {
		return tx.Error
	}

	// Update scalar fields.
	product.Name = name
	if origin != "" {
		product.Origin = &origin
	} else {
		product.Origin = nil
	}

	// Resolve each ingredient, fetching from PubChem if not found locally.
	var newIngredients []types.Ingredient

	for _, ingName := range ingredientNames {
		ingName = strings.ToLower(strings.TrimSpace(ingName))
		if ingName == "" {
			continue
		}

		var ing types.Ingredient
		tx := core.DB.Model(&types.Ingredient{}).
			Preload("Metadata").Preload("Labels").Preload("Names").
			Where("primary_name ~* ?", ingName).First(&ing)
		if tx.Error != nil {
			if tx.Error.Error() != "record not found" {
				return tx.Error
			}
			i, err := ingredient.RetrieveIngredient(ingName)
			if err != nil {
				return err
			}
			ing = *i
		}

		newIngredients = append(newIngredients, ing)
	}

	// Replace the product_ingredients join table rows entirely.
	if err := core.DB.Model(&product).Association("Ingredients").Replace(newIngredients); err != nil {
		return err
	}

	if err := core.DB.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

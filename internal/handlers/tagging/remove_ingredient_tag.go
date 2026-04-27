package tagging

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// RemoveIngredientTagHandler handles DELETE /ingredient/tag.
//
// It removes a tag from a single ingredient, then cascades the removal to
// any product that contains that ingredient — but only removes the tag from
// a product if no other ingredient in that product still carries the tag.
func RemoveIngredientTagHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	ingredientID, err := uuid.Parse(strings.TrimSpace(r.FormValue("ingredient_id")))
	if err != nil {
		return errors.New("invalid ingredient_id")
	}

	tagID, err := uuid.Parse(strings.TrimSpace(r.FormValue("tag_id")))
	if err != nil {
		return errors.New("invalid tag_id")
	}

	// Load both records so GORM can resolve the association correctly.
	var tag types.Tag
	if err := core.DB.First(&tag, "id = ?", tagID).Error; err != nil {
		return err
	}

	var ing types.Ingredient
	if err := core.DB.First(&ing, "id = ?", ingredientID).Error; err != nil {
		return err
	}

	// Remove the ingredient → tag link.
	if err := core.DB.Model(&tag).Association("Ingredients").Delete(&ing); err != nil {
		return err
	}

	// Find every product that:
	//   (a) contains this ingredient, AND
	//   (b) is currently tagged with this tag.
	var affectedProducts []types.Product
	core.DB.
		Joins("JOIN product_ingredients pi ON pi.product_id = products.id").
		Joins("JOIN product_tags pt ON pt.product_id = products.id").
		Where("pi.ingredient_id = ? AND pt.tag_id = ?", ingredientID, tagID).
		Find(&affectedProducts)

	for _, product := range affectedProducts {
		// Count other ingredients in this product that still carry the tag.
		// If none remain the product tag is no longer justified.
		var remaining int64
		core.DB.Table("ingredient_tags").
			Joins("JOIN product_ingredients pi ON pi.ingredient_id = ingredient_tags.ingredient_id").
			Where(
				"pi.product_id = ? AND ingredient_tags.tag_id = ? AND pi.ingredient_id != ?",
				product.Id, tagID, ingredientID,
			).
			Count(&remaining)

		if remaining == 0 {
			if err := core.DB.Model(&tag).Association("Products").Delete(&product); err != nil {
				core.Logger.Error("failed to cascade tag removal to product",
					"product_id", product.Id, "tag_id", tagID, "error", err)
			}
		}
	}

	core.Logger.Info("removed tag from ingredient (and cascade)",
		"ingredient_id", ingredientID, "tag_id", tagID,
		"products_checked", len(affectedProducts))
	return nil
}

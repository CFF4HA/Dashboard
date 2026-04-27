package tagging

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// HandleTagBadIngredients accepts a list of ingredient IDs via ingredient_ids[],
// finds or creates a "bad" tag, applies it to every ingredient, and cascades
// the tag to any product that contains at least one of those ingredients.
func HandleTagBadIngredients(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	rawIDs := r.Form["ingredient_ids"]
	if len(rawIDs) == 0 {
		return errors.New("no ingredient ids provided")
	}

	var ids []uuid.UUID
	for _, raw := range rawIDs {
		id, err := uuid.Parse(strings.TrimSpace(raw))
		if err != nil {
			core.Logger.Warn("invalid ingredient id in tag-bad request, skipping", "id", raw)
			continue
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return errors.New("no valid ingredient ids provided")
	}

	// Find or create the "bad" tag.
	var tag types.Tag
	tx := core.DB.Where("name = ?", "bad").First(&tag)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			return tx.Error
		}
		tag = types.Tag{
			Model: types.Model{
				Id:      uuid.New(),
				Created: time.Now(),
				Updated: time.Now(),
			},
			Name:        "bad",
			Description: "Ingredients flagged as potentially harmful via product comparison.",
		}
		if err := core.DB.Create(&tag).Error; err != nil {
			return err
		}
	}

	// Load the ingredient records and apply the tag.
	var ingredients []types.Ingredient
	if err := core.DB.Where("id IN ?", ids).Find(&ingredients).Error; err != nil {
		return err
	}

	if err := core.DB.Model(&tag).Association("Ingredients").Append(&ingredients); err != nil {
		core.Logger.Error("failed to tag bad ingredients", "count", len(ingredients), "error", err)
		return err
	}

	core.Logger.Info("tagged ingredients as bad", "count", len(ingredients))

	// Cascade: tag every product that contains at least one of these ingredients.
	var products []types.Product
	core.DB.
		Joins("JOIN product_ingredients ON product_ingredients.product_id = products.id").
		Where("product_ingredients.ingredient_id IN ?", ids).
		Find(&products)

	if len(products) > 0 {
		if err := core.DB.Model(&tag).Association("Products").Append(&products); err != nil {
			core.Logger.Error("failed to cascade bad tag to products", "count", len(products), "error", err)
			return err
		}
		core.Logger.Info("cascaded bad tag to products", "count", len(products))
	}

	return nil
}

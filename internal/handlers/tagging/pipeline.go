package tagging

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// TagNewIngredient runs every enabled tagging rule owned by userID against a
// single ingredient that was just created, applying matching tags immediately.
func TagNewIngredient(ing types.Ingredient, userID uuid.UUID) {
	var rules []types.TaggingRule
	core.DB.Where("user_id = ? AND enabled = true", userID).Preload("Tag").Find(&rules)

	for _, rule := range rules {
		var matchCount int64
		core.DB.Model(&types.Label{}).
			Joins("JOIN ingredient_labels ON ingredient_labels.label_id = labels.id").
			Where("ingredient_labels.ingredient_id = ? AND payload ~* ?", ing.Id, rule.Regex).
			Count(&matchCount)

		hasMatch := matchCount > 0
		if (rule.Contains && hasMatch) || (!rule.Contains && !hasMatch) {
			if err := core.DB.Model(&rule.Tag).Association("Ingredients").Append(&ing); err != nil {
				core.Logger.Error("failed to tag new ingredient", "ingredient_id", ing.Id, "tag_id", rule.Tag.Id, "error", err)
			}
		}
	}
}

// TagNewProduct runs every enabled tagging rule owned by userID against the
// ingredients of a newly created product, tagging matching ingredients and
// bubbling each applied tag up to the product itself.
func TagNewProduct(product types.Product, userID uuid.UUID) {
	if len(product.Ingredients) == 0 {
		return
	}

	var rules []types.TaggingRule
	core.DB.Where("user_id = ? AND enabled = true", userID).Preload("Tag").Find(&rules)

	var ingIDs []uuid.UUID
	for _, ing := range product.Ingredients {
		ingIDs = append(ingIDs, ing.Id)
	}

	for _, rule := range rules {
		matchingSub := core.DB.Model(&types.Label{}).
			Select("ingredient_labels.ingredient_id").
			Joins("JOIN ingredient_labels ON ingredient_labels.label_id = labels.id").
			Where("ingredient_labels.ingredient_id IN ? AND payload ~* ?", ingIDs, rule.Regex)

		var matchingIngredients []types.Ingredient
		if rule.Contains {
			core.DB.Where("id IN (?)", matchingSub).Find(&matchingIngredients)
		} else {
			core.DB.Where("id IN ? AND id NOT IN (?)", ingIDs, matchingSub).Find(&matchingIngredients)
		}

		if len(matchingIngredients) == 0 {
			continue
		}

		for _, ing := range matchingIngredients {
			if err := core.DB.Model(&rule.Tag).Association("Ingredients").Append(&ing); err != nil {
				core.Logger.Error("failed to tag product ingredient", "ingredient_id", ing.Id, "tag_id", rule.Tag.Id, "error", err)
			}
		}

		if err := core.DB.Model(&rule.Tag).Association("Products").Append(&product); err != nil {
			core.Logger.Error("failed to bubble tag to product", "product_id", product.Id, "tag_id", rule.Tag.Id, "error", err)
		}
	}
}

func HandleTaggingRuleRun(w http.ResponseWriter, r *http.Request) error {
	rule_id := r.FormValue("rule_id")

	var rule types.TaggingRule
	if tx := core.DB.Preload("Tag").First(&rule, "id = ?", rule_id); tx.Error != nil {
		core.Logger.Error("failed to find tagging rule", "rule_id", rule_id, "error", tx.Error)
		return tx.Error
	}

	return TagIngredients(rule)
}

func TagIngredients(job types.TaggingRule) error {
	// Build a subquery that resolves to the set of ingredient IDs which have
	// at least one label whose payload matches the regex. Both branches reuse
	// this to keep the predicate at the ingredient level, which is the only
	// correct granularity: an ingredient may carry multiple labels, so
	// filtering on individual labels and traversing back to ingredients gives
	// wrong results for the NOT case (and duplicate appends for the IN case).
	matchingIngredientIDs := core.DB.Model(&types.Label{}).
		Select("ingredient_labels.ingredient_id").
		Joins("JOIN ingredient_labels ON ingredient_labels.label_id = labels.id").
		Where("payload ~* ?", job.Regex)

	var ingredients []types.Ingredient
	if job.Contains {
		core.DB.Where("id IN (?)", matchingIngredientIDs).Find(&ingredients)
	} else {
		// Exclude any ingredient that has even one label matching the regex.
		core.DB.Where("id NOT IN (?)", matchingIngredientIDs).Find(&ingredients)
	}

	for _, ingredient := range ingredients {
		// Place the ingredient in the tag's many2many association table.
		err := core.DB.Model(&job.Tag).Association("Ingredients").Append(&ingredient)
		if err != nil {
			core.Logger.Error("failed to tag ingredient", "ingredient_id", ingredient.Id, "tag_id", job.Tag.Id, "error", err)
			return err
		}
	}

	core.Logger.Info("successfully tagged ingredients", "number_ingredients_tagged", len(ingredients), "tag_id", job.Tag.Id)

	var ingredientIDs []uuid.UUID
	for _, ing := range ingredients {
		ingredientIDs = append(ingredientIDs, ing.Id)
	}
	if len(ingredientIDs) > 0 {
		var products []types.Product
		core.DB.
			Joins("JOIN product_ingredients ON product_ingredients.product_id = products.id").
			Where("product_ingredients.ingredient_id IN ?", ingredientIDs).
			Find(&products)

		for _, product := range products {
			err := core.DB.Model(&job.Tag).Association("Products").Append(&product)
			if err != nil {
				core.Logger.Error("failed to tag product", "product_id", product.Id, "tag_id", job.Tag.Id, "error", err)
			}
		}
		core.Logger.Info("tagged products via ingredient cascade", "count", len(products), "tag_id", job.Tag.Id)
	}

	return nil
}

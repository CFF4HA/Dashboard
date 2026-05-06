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
		var matchingIngredients []types.Ingredient

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
	return nil
}

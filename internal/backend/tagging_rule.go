package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/CFF4HA/Dashboard/pkg/boolee"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TagAllIngredientsForTaggingRule(rule *types.TaggingRule) error {
	// we want to first pull all the ingredients that match the pattern.
	var ingredient_ids []string

	core.Logger.Debug("tagging rule ingredient query", "sql", core.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var ids []string
		return tx.Model(&types.Ingredient{}).
			Scopes(boolee.WithBooleeAny("primary_name", "payload", rule.Pattern, "ILIKE")).
			Joins("LEFT JOIN labels ON ingredients.id=labels.ingredient_id").
			Distinct("ingredients.id").
			Pluck("ingredients.id", &ids)
	}))

	tx := core.DB.Model(&types.Ingredient{}).
		Scopes(boolee.WithBooleeAny("primary_name", "payload", rule.Pattern, "ILIKE")).
		Joins("LEFT JOIN labels ON ingredients.id=labels.ingredient_id").
		Distinct("ingredients.id").
		Pluck("ingredients.id", &ingredient_ids)

	if tx.Error != nil {
		core.Logger.Error("failed to retrieve ingredients matching pattern for tagging rule", "tagging_rule_id", rule.Id, "pattern", rule.Pattern, "error", tx.Error)
		return errors.New("failed to retrieve ingredients matching pattern, try again later.")
	}

	core.Logger.Info("tagging ingredients for tagging rule", "tagging_rule_id", rule.Id, "pattern", rule.Pattern, "enabled", rule.Enabled, "num_ingredients_to_check", len(ingredient_ids))

	if len(ingredient_ids) == 0 {
		return nil
	}

	type ingredientTag struct {
		IngredientID uuid.UUID `gorm:"column:ingredient_id;primaryKey"`
		TagID        uuid.UUID `gorm:"column:tag_id;primaryKey"`
	}

	rows := make([]ingredientTag, 0, len(ingredient_ids))
	for _, idStr := range ingredient_ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			core.Logger.Warn("skipping ingredient with unparseable id", "id", idStr, "error", err)
			continue
		}
		rows = append(rows, ingredientTag{IngredientID: id, TagID: rule.TagID})
	}

	if len(rows) == 0 {
		return nil
	}

	tx = core.DB.Table("ingredient_tags").Clauses(clause.OnConflict{DoNothing: true}).Create(&rows)
	if tx.Error != nil {
		core.Logger.Error("failed to apply tag to ingredients", "tagging_rule_id", rule.Id, "tag_id", rule.TagID, "error", tx.Error)
		return errors.New("failed to apply tag to matched ingredients, try again later.")
	}

	core.Logger.Info("successfully tagged ingredients", "tagging_rule_id", rule.Id, "tag_id", rule.TagID, "num_tagged", len(rows))

	taggedIngredientIDs := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		taggedIngredientIDs[i] = row.IngredientID
	}

	var product_ids []string
	tx = core.DB.Table("product_ingredients").
		Where("ingredient_id IN ?", taggedIngredientIDs).
		Distinct("product_id").
		Pluck("product_id", &product_ids)

	if tx.Error != nil {
		core.Logger.Error("failed to retrieve products for tagging rule", "tagging_rule_id", rule.Id, "error", tx.Error)
		return errors.New("failed to retrieve products for tagging, try again later.")
	}

	if len(product_ids) == 0 {
		return nil
	}

	type productTag struct {
		ProductID uuid.UUID `gorm:"column:product_id;primaryKey"`
		TagID     uuid.UUID `gorm:"column:tag_id;primaryKey"`
	}

	productRows := make([]productTag, 0, len(product_ids))
	for _, idStr := range product_ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			core.Logger.Warn("skipping product with unparseable id", "id", idStr, "error", err)
			continue
		}
		productRows = append(productRows, productTag{ProductID: id, TagID: rule.TagID})
	}

	if len(productRows) == 0 {
		return nil
	}

	tx = core.DB.Table("product_tags").Clauses(clause.OnConflict{DoNothing: true}).Create(&productRows)
	if tx.Error != nil {
		core.Logger.Error("failed to apply tag to products", "tagging_rule_id", rule.Id, "tag_id", rule.TagID, "error", tx.Error)
		return errors.New("failed to apply tag to matched products, try again later.")
	}

	core.Logger.Info("successfully tagged products", "tagging_rule_id", rule.Id, "tag_id", rule.TagID, "num_tagged", len(productRows))
	return nil
}

// ------------------------------
// Generic Utility Functions, not Route Related
//
// The following are utility functions that do not assume the
// method the data was retrieved with.
// ------------------------------
func InsertTaggingRule(userId string, tagId string, taggingSetId string, enabled bool, pattern string) (*types.TaggingRule, error) {
	rule := &types.TaggingRule{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Pattern: pattern,
		Enabled: enabled,
	}

	taggingSetIdParsed, err := uuid.Parse(taggingSetId)
	if err != nil {
		core.Logger.Error("failed to parse tagging set id", "tagging_set_id", taggingSetId, "error", err)
		return nil, errors.New("a valid tagging set is required to perform this operation, please try again after selecting a tagging set")
	}

	tagIdParsed, err := uuid.Parse(tagId)
	if err != nil {
		core.Logger.Error("failed to parse tag id", "tag_id", tagId, "error", err)
		return nil, errors.New("a valid tag is required to perform this operation, please try again after selecting a tag")
	}

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		core.Logger.Error("failed to parse user id", "user_id", userId, "error", err)
		return nil, errors.New("a valid user is required to perform this operation, please try again after signing in")
	}

	rule.TagID = tagIdParsed
	rule.UserId = userIdParsed
	rule.TaggingSetID = &taggingSetIdParsed

	tx := core.DB.Create(rule)
	if tx.Error != nil {
		core.Logger.Error("failed to insert tagging rule into database", "error", tx.Error)
		return nil, errors.New("failed to insert tagging rule into database, try again later.")
	}

	// TODO: Now we want to run it and tag ingredients that fall into the relevant category
	// as done here.
	if TagAllIngredientsForTaggingRule(rule) != nil {
		core.Logger.Error("failed to apply tagging rule to existing ingredients", "tagging_rule_id", rule.Id)
	}

	return rule, nil
}

func RunTaggingRule(id string) error {
	var rule *types.TaggingRule

	tx := core.DB.First(&rule, "id = ?", id)
	if tx.Error != nil {
		return errors.New("failed to retrieve tagging rule from database, try again later.")
	}

	if !rule.Enabled {
		return errors.New("this tagging rule is currently disabled, enable it before running.")
	}

	return TagAllIngredientsForTaggingRule(rule)
}

func DeleteTaggingRule(id string) error {
	tx := core.DB.Delete(&types.TaggingRule{}, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to delete tagging rule from database", "tagging_rule_id", id, "error", tx.Error)
		return errors.New("failed to delete tagging rule from database, try again later.")
	}

	return nil
}

func GetTaggingRuleById(id string) (*types.TaggingRule, error) {
	var rule *types.TaggingRule

	tx := core.DB.Scopes(WithPreload("Tag", "TaggingSet", "User")).First(&rule, "id = ?", id)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging rule from database", "tagging_rule_id", id, "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging rule from database, try again later.")
	}

	return rule, nil
}

func GetTaggingRulesForTag(tagId string, cursor string) ([]*types.TaggingRule, error) {
	var rules []*types.TaggingRule

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithPreload("Tag", "TaggingSet", "User")).Find(&rules, "tag_id = ?", tagId)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging rules from database", "tag_id", tagId, "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging rules from database, try again later.")
	}

	return rules, nil
}

func GetTaggingRules(cursor string) ([]types.TaggingRule, error) {
	var rules []types.TaggingRule

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithPreload("Tag", "TaggingSet", "User")).Find(&rules)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging rules from database", "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging rules from database, try again later.")
	}

	return rules, nil
}

func GetTaggingRulesForUser(userId string, cursor string) ([]*types.TaggingRule, error) {
	var rules []*types.TaggingRule

	tx := core.DB.Scopes(WithCursor(cursor), WithOrder("id"), WithPreload("Tag", "TaggingSet", "User")).Find(&rules, "user_id = ?", userId)
	if tx.Error != nil {
		core.Logger.Error("failed to retrieve tagging rules from database", "user_id", userId, "error", tx.Error)
		return nil, errors.New("failed to retrieve tagging rules from database, try again later.")
	}

	return rules, nil
}

// ------------------------------
// Route Related
//
// These set of endpoints are used in the API verification stage.
// ------------------------------
func RouteInsertTaggingRule(w http.ResponseWriter, r *http.Request) error {
	user, err := user.GetUserFromRequest(w, r)
	if err != nil {
		return errors.New("a valid user is required to perform this operation, please try again after signing in")
	}

	rule, err := InsertTaggingRule(user.Id.String(), r.FormValue("tag_id"), r.FormValue("tagging_set_id"), r.FormValue("enabled") == "true", r.FormValue("pattern"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(rule)
}

func RouteDeleteTaggingRule(w http.ResponseWriter, r *http.Request) error {
	return DeleteTaggingRule(r.FormValue("id"))
}

func RouteGetTaggingRulesForUser(w http.ResponseWriter, r *http.Request) error {
	user, err := user.GetUserFromRequest(w, r)
	if err != nil {
		return errors.New("a valid user is required to perform this operation, please try again after signing in")
	}

	rules, err := GetTaggingRulesForUser(user.Id.String(), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(rules)
}

func RouteGetTaggingRules(w http.ResponseWriter, r *http.Request) error {
	rules, err := GetTaggingRules(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(rules)
}

func RouteGetTaggingRuleById(w http.ResponseWriter, r *http.Request) error {
	rule, err := GetTaggingRuleById(r.FormValue("id"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(rule)
}

func RouteGetTaggingRulesForTag(w http.ResponseWriter, r *http.Request) error {
	rules, err := GetTaggingRulesForTag(r.FormValue("tag_id"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(rules)
}

func RouteRunTaggingRule(w http.ResponseWriter, r *http.Request) error {
	return RunTaggingRule(r.FormValue("id"))
}

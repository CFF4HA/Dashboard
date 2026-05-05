package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

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

	return rule, nil
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

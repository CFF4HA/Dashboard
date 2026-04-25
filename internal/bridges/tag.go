package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	"github.com/google/uuid"
)

var (
	TagsByName = verb.DataBridge{
		Key: "Tags",
		Provide: func(r *http.Request, model map[string]any) (any, error) {
			model["Name"] = r.FormValue("name")
			return GetTagsByName(r.FormValue("name"))
		},
	}
)

type TaggingRules struct{}

func (t TaggingRules) Name() string {
	return "Rules"
}

func (t TaggingRules) Data(w http.ResponseWriter, r *http.Request, model map[string]any) (any, error) {
	user, err := user.GetUserFromRequest(w, r)
	if err != nil {
		return nil, err
	}

	return GetTaggingRules(user.Model.Id)
}

func GetTaggingRules(user uuid.UUID) ([]types.TaggingRule, error) {
	var tagrules []types.TaggingRule
	if err := core.DB.Where("user_id = ?", user).Preload("Tag").Find(&tagrules).Error; err != nil {
		if err.Error() == "record not found" {
			return tagrules, nil
		}

		return nil, err
	}

	return tagrules, nil
}

func GetTagsByName(name string) ([]types.Tag, error) {
	var tags []types.Tag
	if err := core.DB.Where("name ILIKE ?", name+"%").Find(&tags).Error; err != nil {
		if err.Error() == "record not found" {
			return tags, nil
		}

		core.Logger.Warn("failed to get tags by name", "name", name)
		return nil, err
	}

	return tags, nil
}

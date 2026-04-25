package bridges

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
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

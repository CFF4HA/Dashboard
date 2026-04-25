package tagging

import (
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

// this function will delete a tag rule by it's id.
func HandleTagRuleDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return errors.New("id is required")
	}

	tx := core.DB.Begin()
	if err := tx.Delete(&types.TaggingRule{}, "id = ?", id).Error; err != nil {
		core.Logger.Error("failed to delete tagging rule", "id", id, "err", err)
		return errors.New("failed to delete tagging rule")
	}

	tx.Commit()
	return nil
}

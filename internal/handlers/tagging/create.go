package tagging

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// This handler will be used to create a tagging rule. The request parameters
// will be used as follows:
//
// "regex": this is the pattern that will be compared against the ingredient
// labels. It is required, and will be placed 1-1 into the relevatne field in the database.
//
// "operation": this must be one of two values: "contains" or "dncontains". It is
// used to set the database fild to true or false, also required.
//
// "tag": this is the tag that will be applied to the ingredients matching the
// regex and operation combination. The tag must be created if it doesn't exist
// and the tag must be associated with the user via the tag user many2many table.
//
// Should return 200 when the tagging rule is created.
// TODO: Re-check that this function is actually correct.
func HandleTagRuleCreate(w http.ResponseWriter, r *http.Request) error {
	regex := r.FormValue("regex")
	operation := r.FormValue("operation")
	tag := strings.ToLower(strings.TrimSpace(r.FormValue("tag")))

	// Section: Input Validation
	if regex == "" {
		return errors.New("regex is required")
	} else if operation != "contains" && operation != "dncontains" {
		return errors.New("operation must be either 'contains' or 'dncontains'")
	} else if tag == "" {
		return errors.New("tag is required")
	}

	// Section: Ensure that there is a user
	user, err := user.GetUserFromRequest(w, r)
	if err != nil {
		return err
	}

	// Section Create or retrieve the tag
	// The tag existence check will be done by checking the tag table against
	// the name field, which is unique.
	var tagv types.Tag
	tx := core.DB.Where("name = ?", tag).First(&tagv)
	if tx.Error.Error() == "record not found" {
		// here we want to create the tag.
		tagv.Model.Id = uuid.New()
		tagv.Model.Created = time.Now()
		tagv.Model.Updated = time.Now()
		tagv.Name = tag

		if err := core.DB.Create(&tagv).Error; err != nil {
			// the tag could not be created so we return an error.
			return err
		}
	} else if tx.Error != nil {
		return tx.Error
	}

	// pre-fills the tagging rule with the known variables
	tagging_rule := &types.TaggingRule{
		Contains: operation == "contains",
		Regex:    regex,
		Enabled:  true,
		UserId:   user.Model.Id,
		TagID:    tagv.Model.Id,
	}

	if err := core.DB.Create(tagging_rule).Error; err != nil {
		return err
	}

	core.Logger.Info("created a tagging rule", "regex", regex, "operation", operation, "tag", tag)
	return nil
}

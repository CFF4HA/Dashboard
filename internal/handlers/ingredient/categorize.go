package ingredient

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/google/uuid"
)

// This function expects an ingredient ID and a "status" that
// reports whether it is "safe" or "unsafe" for the user.
// It requires a signed-in user.
func Categorize(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	id, ok := r.Form["id"]
	if !ok {
		return errors.New("missing ingredient id")
	}

	status, ok := r.Form["status"]
	if !ok {
		return errors.New("missing ingredient status")
	} else if len(status) != 1 {
		return errors.New("invalid ingredient status")
	}

	c, err := r.Cookie("session")
	if err != nil {
		return errors.New("missing session cookie")
	}

	session, err := uuid.Parse(c.Value)
	if err != nil {
		return errors.New("invalid session cookie")
	}

	uid, ok := user.GetSession(session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return errors.New("invalid session")
	}

	// We now want to insert the id and the user into the relevant
	// many2many table. If the status is "safe", we want to insert it into
	// "good_ingredients". If the status is "unsafe", we want to insert it into
	// the "bad_ingredients" table.
	var table, opposing_table string
	switch status[0] {
	case "safe":
		table = "good_ingredients"
		opposing_table = "bad_ingredients"
	case "unsafe":
		table = "bad_ingredients"
		opposing_table = "good_ingredients"
	default:
		return errors.New("invalid ingredient status")
	}

	tx := core.DB.Begin()
	for _, _id := range id {
		ingredient, err := uuid.Parse(_id)
		if err != nil {
			core.Logger.Warn("skipping an ingredient due to malformed ID", "ingredient", _id)
			continue
		}

		var returned_uid uuid.UUID
		query := fmt.Sprintf("SELECT user_id FROM user_%s WHERE user_id = ? AND ingredient_id = ?", table)
		if core.DB.Exec(query, uid, ingredient).Scan(&returned_uid).Error == nil {
			// no error here means that there existed a record with this category already
			// so we will go ahead and skip.
			core.Logger.Warn("skipping an ingredient because it was already categorized", "user_id", uid, "ingredient_id", ingredient, "category", table)
			continue
		}

		query = fmt.Sprintf("DELETE FROM user_%s WHERE user_id = ? AND ingredient_id = ?", opposing_table)
		if core.DB.Exec(query, uid, ingredient).Error != nil {
			core.Logger.Warn("failed to remove ingredient from opposing category", "user_id", uid, "ingredient_id", ingredient, "category", opposing_table)
			continue
		}

		query = fmt.Sprintf("INSERT INTO user_%s VALUES (?, ?)", table)
		if err := tx.Exec(query, uid, ingredient).Error; err != nil {
			tx.Rollback()
			core.Logger.Error("failed to categorize ingredient", "error", err)
			continue
		}
	}

	tx.Commit()
	core.Logger.Info("successfully categorized ingredients", "user_id", uid, "category", table, "count", len(id))
	return nil
}

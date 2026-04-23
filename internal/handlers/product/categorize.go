package product

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/google/uuid"
)

func Categorize(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	id, ok := r.Form["id"]
	if !ok {
		return errors.New("missing product id")
	}

	status, ok := r.Form["status"]
	if !ok {
		return errors.New("missing product status")
	} else if len(status) != 1 {
		return errors.New("invalid product status")
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

	var table, opposing_table string
	switch status[0] {
	case "safe":
		table = "good_products"
		opposing_table = "bad_products"
	case "unsafe":
		table = "bad_products"
		opposing_table = "good_products"
	default:
		return errors.New("invalid product status")
	}

	tx := core.DB.Begin()
	for _, _id := range id {
		prod, err := uuid.Parse(_id)
		if err != nil {
			core.Logger.Warn("skipping a product due to malformed ID", "product", _id)
			continue
		}

		var returned_uid uuid.UUID
		query := fmt.Sprintf("SELECT user_id FROM user_%s WHERE user_id = ? AND product_id = ?", table)
		if core.DB.Exec(query, uid, prod).Scan(&returned_uid).Error == nil {
			core.Logger.Warn("skipping a product because it was already categorized", "user_id", uid, "product_id", prod, "category", table)
			continue
		}

		query = fmt.Sprintf("DELETE FROM user_%s WHERE user_id = ? AND product_id = ?", opposing_table)
		if core.DB.Exec(query, uid, prod).Error != nil {
			core.Logger.Warn("failed to remove product from opposing category", "user_id", uid, "product_id", prod, "category", opposing_table)
			continue
		}

		query = fmt.Sprintf("INSERT INTO user_%s VALUES (?, ?)", table)
		if err := tx.Exec(query, uid, prod).Error; err != nil {
			tx.Rollback()
			core.Logger.Error("failed to categorize product", "error", err)
			continue
		}

		if table == "good_products" {
			var ingredientIDs []uuid.UUID
			if err := core.DB.Table("product_ingredients").
				Select("ingredient_id").
				Where("product_id = ?", prod).
				Pluck("ingredient_id", &ingredientIDs).Error; err != nil {
				core.Logger.Warn("failed to load product ingredients for safe cascade", "product_id", prod, "error", err)
			} else {
				for _, ingID := range ingredientIDs {
					if err := tx.Exec("INSERT INTO user_good_ingredients (user_id, ingredient_id) VALUES (?, ?) ON CONFLICT DO NOTHING", uid, ingID).Error; err != nil {
						core.Logger.Warn("failed to cascade ingredient to safe list", "ingredient_id", ingID, "error", err)
					}
				}
			}
		}
	}

	tx.Commit()
	core.Logger.Info("successfully categorized products", "user_id", uid, "category", table, "count", len(id))
	return nil
}

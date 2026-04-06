package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

func ProductIngredientListGET(w http.ResponseWriter, r *http.Request) error {
	db := database.Database()
	draft := &types.ProductDraft{}

	for _, ingredient := range strings.Split(r.FormValue("ingredients"), ",") {
		if strings.Trim(ingredient, " \t\r") == "" {
			continue
		}

		var count int64
		name := strings.Trim(ingredient, " \t\r\n")
		tx := db.Model(&types.Name{}).Where("text ILIKE ?", name).Count(&count)
		if tx.Error != nil {
			return tx.Error
		}

		if count == 0 {
			go syncDatabaseWithIngredient(name)
		}

		draft.Ingredients = append(draft.Ingredients, types.ProductDraftIngredient{
			Name:   ingredient,
			Exists: count > 0,
		})
	}

	return json.NewEncoder(w).Encode(draft)
}

func ProductGET(w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	id := r.FormValue("id")

	db := database.Database()
	if id != "" {
		// in this scenarios, we want to get a specific product by id.
		var product *types.Product

		product_id, err := uuid.Parse(id)
		if err != nil {
			return err
		}

		if tx := db.Where("id = ?", product_id).First(&product); tx.Error != nil {
			return tx.Error
		}

		return json.NewDecoder(r.Body).Decode(&product)
	}

	// in this scenario we want to get all products matching a given name.
	var products []*types.Product
	if tx := db.Where("name LIKE ?", "%"+name+"%").Find(&products); tx.Error != nil {
		core.Logger.Error("failed to query for products", "error", tx.Error, "name", name)
		return tx.Error
	}

	return json.NewDecoder(r.Body).Decode(&products)
}

func ProductPUT(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("name") == "" || r.FormValue("origin") == "" {
		return errors.New("name and origin are required to create a product")
	}
	db := database.Database()

	product := &types.Product{
		Name:   r.FormValue("name"),
		Origin: r.FormValue("origin"),
	}

	if tx := db.Create(product); tx.Error != nil {
		core.Logger.Error("failed to create product", "error", tx.Error)
		return tx.Error
	}

	core.Logger.Info("PUT /product")
	return nil
}

func ProductPATCH(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")
	if id == "" {
		return errors.New("id is required to update a product")
	}

	product_id, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	db := database.Database()

	product := &types.Product{
		Id: product_id,
	}

	if name := r.FormValue("name"); name != "" {
		product.Name = name
	}
	if origin := r.FormValue("origin"); origin != "" {
		product.Origin = origin
	}
	product.Updated = time.Now()

	return db.Model(&product).Where("id = ?", product_id).Updates(product).Error
}

func ProductDELETE(w http.ResponseWriter, r *http.Request) error {
	if r.FormValue("id") == "" {
		return errors.New("id is required to delete a product")
	}
	id := r.FormValue("id")

	db := database.Database()
	if tx := db.Delete(&types.Product{}, "id = ?", id); tx.Error != nil {
		return tx.Error
	}

	core.Logger.Info("DELETE /product", "id", id)
	return nil
}

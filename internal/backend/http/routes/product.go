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

		var ingredient_name_first types.Name
		name := strings.Trim(ingredient, " \t\r\n")
		tx := db.Model(&types.Name{}).Where("text ILIKE ?", name).Preload("Ingredient").First(&ingredient_name_first)
		if tx.Error != nil {
			if tx.Error.Error() == "record not found" {
				syncDatabaseWithIngredient(name)
				db.Model(&types.Name{}).Where("text ILIKE ?", name).Preload("Ingredient").First(&ingredient_name_first)
			} else {
				return tx.Error
			}
		}

		draft.Ingredients = append(draft.Ingredients, types.ProductDraftIngredient{
			Id:     ingredient_name_first.Ingredient.Id,
			Name:   ingredient,
			Exists: true,
			Failed: ingredient_name_first.Ingredient.Failed,
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
		var product types.Product

		product_id, err := uuid.Parse(id)
		if err != nil {
			return err
		}

		if tx := db.Where("id = ?", product_id).Preload("Ingredients").First(&product); tx.Error != nil {
			return tx.Error
		}

		for i, _ := range product.Ingredients {
			db.Model(&product.Ingredients[i]).Association("Names").Find(&product.Ingredients[i].Names)
			db.Model(&product.Ingredients[i]).Association("Labels").Find(&product.Ingredients[i].Labels)
		}

		return json.NewEncoder(w).Encode(&product)
	}

	// in this scenario we want to get all products matching a given name.
	var products []*types.Product
	if tx := db.Where("name LIKE ?", "%"+name+"%").Preload("Ingredients").Find(&products); tx.Error != nil {
		core.Logger.Error("failed to query for products", "error", tx.Error, "name", name)
		return tx.Error
	}

	return json.NewEncoder(w).Encode(&products)
}

func ProductPUT(w http.ResponseWriter, r *http.Request) error {
	var request struct {
		Name        string   `json:"name"`
		Origin      string   `json:"origin"`
		Ingredients []string `json:"ingredients"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	db := database.Database()

	product := &types.Product{
		Name:   request.Name,
		Origin: request.Origin,
	}

	for _, name := range request.Ingredients {
		// we first want to get the ingredient that matches the name
		var ingredientName types.Name
		tx := db.Model(&types.Name{}).Where("text ILIKE ?", name).Preload("Ingredient").First(&ingredientName)
		if tx.Error != nil {
			return tx.Error
		}

		// we now want to add this ingredient to the list of the product
		// ingredients.
		product.Ingredients = append(product.Ingredients, ingredientName.Ingredient)
	}

	if tx := db.Create(product); tx.Error != nil {
		core.Logger.Error("failed to create product", "error", tx.Error)
		return tx.Error
	}

	core.Logger.Info("PUT /product")
	w.WriteHeader(http.StatusCreated)
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

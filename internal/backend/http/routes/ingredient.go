package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

const (
	DB_QUERY_PERIOD    = 3
	DB_TIMEOUT_SECONDS = 20
)

var (
	IngredientScrapingBackend = "http://localhost:8082"
)

func syncDatabaseWithIngredient(name string) {
	request, err := http.NewRequest("GET", IngredientScrapingBackend+"/ingredient?name="+url.QueryEscape(name), nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

// This route will be used to search based on names and labels.
// This is different than GET since GET is the route used to
// get an exact name. This will always return a list of ingredients.
func IngredientSEARCH(w http.ResponseWriter, r *http.Request) error {
	db := database.Database()
	var ingredient []types.Ingredient

	chemical_name := r.FormValue("name")
	if chemical_name == "" {
		return json.NewEncoder(w).Encode(ingredient)
	}

	// chemical_name = "Acid"
	searchTerm := "%" + chemical_name + "%"

	err := db.Model(&types.Ingredient{}).
		// 1. Join the Names table to filter by name text
		Joins("JOIN names ON names.ingredient_id = ingredients.id").
		// 2. Filter using ILIKE
		Where("names.text ILIKE ?", searchTerm).
		// 3. Ensure we only get unique Ingredient records
		Distinct("ingredients.*").
		// 4. Limit to top 5
		Limit(5).
		// 5. Optionally Preload all names for those ingredients
		Preload("Names").
		Find(&ingredient).Error

	if err != nil {
		return json.NewEncoder(w).Encode(ingredient)
	}

	return json.NewEncoder(w).Encode(ingredient)
}

func IngredientGET(w http.ResponseWriter, r *http.Request) error {
	var ingredient *types.Ingredient
	name := new(types.Name)
	db := database.Database()

	if formValue := r.FormValue("id"); formValue != "" {
		var ingredient types.Ingredient
		if tx := db.Preload("Labels").Preload("Names").First(&ingredient, "id = ?", formValue); tx.Error != nil {
			return tx.Error
		}

		return json.NewEncoder(w).Encode(ingredient)
	}

	// HANDLE USER AUTH RELATED PIECES HERE IF NEEDED
	// we assume that the user exists, and is in the system.
	chemical_name := r.FormValue("name")
	found := false
	timeout := false

	// setup the infrastructure to query for the ingredient.
	ticker := time.NewTicker(time.Duration(DB_QUERY_PERIOD) * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DB_TIMEOUT_SECONDS)*time.Second)
	defer cancel()

	go syncDatabaseWithIngredient(chemical_name)

	// Ping the database, or timeout.
	for !found && !timeout {
		select {
		case <-ticker.C:
			core.Logger.Debug("checking for ingredient", "chemical_name", chemical_name)

			if tx := db.Model(&types.Name{}).Where("text ILIKE ?", "%"+chemical_name+"%").Preload("Ingredient").First(name); tx.Error != nil {
				continue
			}

			ingredient = &name.Ingredient
			db.Preload("Labels").Preload("Names").Find(ingredient)
			found = true
		case <-ctx.Done():
			core.Logger.Debug("ingredient search timed out", "chemical_name", chemical_name)
			timeout = true
		}
	}

	if !found {
		return errors.New("ingredient not found")
	}

	core.Logger.Info("GET /ingredient", "chemical", name.Text, "ingredient", ingredient)
	return json.NewEncoder(w).Encode(ingredient)
}

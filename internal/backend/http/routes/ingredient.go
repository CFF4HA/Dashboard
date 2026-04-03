package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

const (
	DB_QUERY_PERIOD    = 3
	DB_TIMEOUT_SECONDS = 20
)

// This route will be used to search based on names and labels.
// This is different than GET since GET is the route used to
// get an exact name. This will always return a list of ingredients.
func IngredientSEARCH(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func IngredientGET(w http.ResponseWriter, r *http.Request) error {
	var ingredient *types.Ingredient
	name := new(types.Name)

	// HANDLE USER AUTH RELATED PIECES HERE IF NEEDED
	// we assume that the user exists, and is in the system.
	db := database.Database()
	chemical_name := r.FormValue("name")
	found := false
	timeout := false

	// setup the infrastructure to query for the ingredient.
	ticker := time.NewTicker(time.Duration(DB_QUERY_PERIOD) * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DB_TIMEOUT_SECONDS)*time.Second)
	defer cancel()

	// Ping the database, or timeout.
	for !found && !timeout {
		select {
		case <-ticker.C:
			core.Logger.Debug("checking for ingredient", "chemical_name", chemical_name)

			if tx := db.Model(&types.Name{}).Where("text = ?", chemical_name).Preload("Ingredient").First(name); tx.Error != nil {
				continue
			}

			ingredient = &name.Ingredient
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

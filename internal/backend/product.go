package backend

import (
	"errors"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// ------------------------------
// Generic Utility Functions, not Route Related
//
// These functions provide the functionality the routes need to perform
// actions.
// ------------------------------
func InsertProduct(name string, origin string, ingredient_list []string) (*types.Product, error) {
	// Ingredient list is a list of ingredient names, must be
	// retrieved from the data sources on demand or via database (or cache).
	ingredients := []types.Ingredient{}
	for _, i := range ingredient_list {
		ingredient, err := RetrieveIngredientByPrimaryName(i)
		if err != nil {
			// we were not able to retrieve the ingredient, so we should not insert
			// the product.
			return nil, errors.New("failed to insert product due to failed ingredient retrieval, pubchem or some other data source may be down, try again later.")
		}

		ingredients = append(ingredients, *ingredient)
	}

	prod := &types.Product{
		Model: types.Model{
			Id:      uuid.New(),
			Created: time.Now(),
			Updated: time.Now(),
		},
		Name:        name,
		Origin:      &origin,
		Ingredients: ingredients,
		Metadata:    types.ProductMetadata{},
	}
	prod.Metadata.Model.Id = prod.Id

	// TODO: Add auto-tagging based on the user's rules
	// and the system rules.

	// TODO: Add metadata information for the number of hazards, etc.

	tx := core.DB.Begin()
	if tx.Create(prod); tx.Error != nil {
		core.Logger.Error("failed to insert product into database", "error", tx.Error)
		return nil, errors.New("failed to insert product into database, try again later.")
	}

	return prod, tx.Commit().Error
}

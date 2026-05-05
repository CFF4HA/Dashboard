package backend

import (
	"errors"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/CFF4HA/Dashboard/pkg/pubchem"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

// ------------------------------
// Generic Utility Functions, not Route Related
//
// These functions provide the functionality the routes need to perform
// actions.
// ------------------------------

// This function is the function responsible for pulling ingredient information from the data sources.
func pullIngredientByName(name string) (*types.Ingredient, error) {
	// Will always upload an ingredient, marks it as failed if it did not
	// work as expected. This is to prevent repeated attempts to pull the same ingredient.
	name = strings.ToLower(strings.TrimSpace(name))

	cid, err := pubchem.GetCompoundId(name)
	if err != nil {
		core.Logger.Error("failed to retrieve compound id from pubchem", "name", name, "error", err)
		return nil, errors.New("failed to retrieve ingredient information, try again later.")
	}

	compound, err := pubchem.GetCompoundAsJSON(cid)
	if err != nil {
		core.Logger.Error("failed to retrieve compound information from pubchem", "name", name, "cid", cid, "error", err)
		return nil, errors.New("failed to retrieve ingredient information, try again later.")
	}

	// TODO: Parse information from the configurable database backed pubchem config
	// set by the system.
	allLabelConfigs := []types.PubChemLabelConfig{}
	if tx := core.DB.Model(&types.PubChemLabelConfig{}).Find(allLabelConfigs); tx.Error != nil {
		core.Logger.Error("failed to retrieve pubchem label configs from database", "error", tx.Error)
		return nil, errors.New("failed to retrieve ingredient information, try again later.")
	}

	for _, cfg := range allLabelConfigs {
		// TODO: Finish implementation of the gjson parsing
		// and marshalling them into arrays of strings.

		// TODO: For each string in the array of strings, create a new label and match it
		// with the created ingredient.
	}

	// TODO: Do the tagging based on the rules in the system.

	// TODO: Return the ingredient.
	return nil, nil
}

func RetrieveIngredientByPrimaryName(name string) (*types.Ingredient, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	tx := core.DB.Begin()

	ing := &types.Ingredient{}
	if tx.Where("primary_name = ?", name).First(ing); tx.Error != nil {
		// if the error is a record not found error we need to perform a search for
		// the same chemical using the NAMES table, but it must be an
		// exact match.
		if tx.Error == gorm.ErrRecordNotFound {
			var synonym types.Name
			if tx.Scopes(WithPreload("Ingredients")).Where("name = ?", name).First(&synonym).Error != nil {
				// there is no match, this has never been searched before, so we
				// perform a search for the same chemical.
				return pullIngredientByName(name)
			}

			// returning the first ingredient because there's no good way to demultiplex this
			// right now.
			//
			// TODO: Potentially improve this functionality in the future.
			return &synonym.Ingredients[0], nil
		}

		return ing, nil
	}

	return ing, nil
}

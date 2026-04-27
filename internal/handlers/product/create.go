package product

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/ai"
	"github.com/CFF4HA/Dashboard/internal/handlers/ingredient"
	"github.com/CFF4HA/Dashboard/internal/handlers/tagging"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// this file serves as the primary location
// for all product creation related HTTP handlers.

func CreateProduct(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	ingredient_names := r.Form["ingredient_names"]
	origin := strings.TrimSpace(r.FormValue("origin"))
	name := strings.TrimSpace(r.FormValue("name"))

	// here is our validation step, we will check to make sure there is a name
	// and an origin.
	if len(name) == 0 {
		return errors.New("product name is required")
	} else if len(ingredient_names) == 0 {
		return errors.New("at least one ingredient name is required")
	}

	product := types.Product{
		Model: types.Model{
			Id: uuid.New(),
		},
		Name: name,
		Metadata: types.ProductMetadata{
			Model: types.Model{
				Id: uuid.New(),
			},
		},
	}

	// sets the foreign key relationship between the product and its metadata.
	product.Metadata.ProductId = product.Id

	// sets the origin if it exists.
	if origin != "" {
		product.Origin = new(string)
		*product.Origin = origin
	}

	// begins the ingredient resolution process, and if they don't
	// exist then we simply create a failed ingredient in the database
	// so resolutions don't continue.

	for _, ingredient_name := range ingredient_names {
		ingredient_name = strings.ToLower(strings.TrimSpace(ingredient_name))
		var ing types.Ingredient
		tx := core.DB.Model(&types.Ingredient{}).Preload("Metadata").Preload("Labels").Preload("Names").Where("primary_name ~* ?", ingredient_name).First(&ing)
		if tx.Error != nil {
			if tx.Error.Error() != "record not found" {
				return tx.Error
			}

			i, err := ingredient.RetrieveIngredient(ingredient_name)
			if err != nil {
				return err
			}

			ing = *i
		}

		product.Ingredients = append(product.Ingredients, ing)
		product.Metadata.NumHazard += ing.Metadata.NumHazards
		product.Metadata.NumEffect += ing.Metadata.NumEffects
		product.Metadata.NumSymptom += ing.Metadata.NumSymptoms
		product.Metadata.NumReglations += ing.Metadata.NumRegulations
	}

	// saves the product to the database, which will also save the metadata and the
	tx := core.DB.Create(&product)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	// Run the requesting user's tagging rules against the product's ingredients
	// and bubble any matching tags up to the product itself.
	if u, userErr := user.GetUserFromRequest(w, r); userErr == nil {
		tagging.TagNewProduct(product, u.Model.Id)
	}

	w.Header().Set("HX-Trigger", fmt.Sprintf(`{"productCreated":{"productId":%q}}`, product.Id.String()))
	return ai.ProductBlurb(r.Context(), product)
}

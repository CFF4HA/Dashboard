package product

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
)

// this file serves as the primary location
// for all product creation related HTTP handlers.

func CreateProduct(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	ingredient_names := core.RequestValuesAsSlice(r, "ingredient_names")
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

	// TODO

	// saves the product to the database, which will also save the metadata and the
	tx := core.DB.Create(&product)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return tx.Error
}

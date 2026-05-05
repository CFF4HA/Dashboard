package backend

import (
	"errors"
	"net/http"
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

		core.Logger.Debug("successfully retrieved ingredient information", "name", i, "ingredient_id", ingredient.Id)
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
	}

	// TODO: Add auto-tagging based on the user's rules
	// and the system rules.

	tx := core.DB.Begin()
	if tx.Create(prod); tx.Error != nil {
		core.Logger.Error("failed to insert product into database", "error", tx.Error)
		return nil, errors.New("failed to insert product into database, try again later.")
	}

	return prod, tx.Commit().Error
}

func GetProducts(cursor string) ([]types.Product, error) {
	var prods []types.Product
	tx := core.DB.Scopes(WithPreload("Ingredients"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
		Find(&prods)
	return prods, tx.Error
}

func GetProductById(id string) (*types.Product, error) {
	prod := &types.Product{}
	tx := core.DB.Scopes(WithPreload("Ingredients")).First(prod, "id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return prod, nil
}

func GetProductsByName(name string, cursor string) ([]types.Product, error) {
	var products []types.Product

	// we first do a search by primary name
	tx := core.DB.Scopes(WithPreload("Ingredients"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
		Where("name ~* ?", name).Find(&products)
	if tx.Error != nil {
		core.Logger.Error("failed to search for products by name", "name", name, "error", tx.Error)
		// TODO: Implement auxiliary name search.

		// we then would want to do a search by labels of type 'name' that
		// are in the system (but this is not immediately clear how to do).
		return products, nil
	}

	return products, nil
}

func DeleteProductById(id string) error {
	// TODO: we first want to delete the product_ingredients join table rows,
	// and then we can delete the product itself.
	tx := core.DB.Begin()
	if tx.Exec("DELETE FROM product_ingredients WHERE product_id = ?", id).Error != nil {
		core.Logger.Error("failed to delete product ingredients associations", "product_id", id, "error", tx.Error)
		return errors.New("failed to delete product ingredients associations, try again later.")
	}

	if tx.Exec("DELETE FROM product_tags WHERE product_id = ?", id).Error != nil {
		core.Logger.Error("failed to delete product ingredients associations", "product_id", id, "error", tx.Error)
		return errors.New("failed to delete product tags associations, try again later.")
	}

	if tx.Delete(&types.Product{}, "id = ?", id).Error != nil {
		core.Logger.Error("failed to delete product", "product_id", id, "error", tx.Error)
		return errors.New("failed to delete product, try again later.")
	}

	tx.Commit()
	core.Logger.Debug("successfully deleted product", "product_id", id)
	return nil
}

// ------------------------------
// Routing Related Functions
// ------------------------------
func RouteProductPUT(w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	origin := r.FormValue("origin")
	ingredientList := r.Form["ingredient"]

	_, err := InsertProduct(name, origin, ingredientList)
	if err != nil {
		return err
	}

	return nil
}

func RouteProductGET(w http.ResponseWriter, r *http.Request) error {
	products, err := GetProducts(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully retrieved products", "count", len(products), "products", products)
	return nil
}

func RouteGetProductsByName(w http.ResponseWriter, r *http.Request) error {
	products, err := GetProductsByName(r.FormValue("name"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully retrieved products by name", "name", r.FormValue("name"), "count", len(products), "products", products)
	return nil
}

func RouteGetProductById(w http.ResponseWriter, r *http.Request) error {
	product, err := GetProductById(r.FormValue("id"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully retrieved product by id", "id", r.FormValue("id"), "product", product)
	return nil
}

func RouteDeleteProductById(w http.ResponseWriter, r *http.Request) error {
	return DeleteProductById(r.FormValue("id"))
}

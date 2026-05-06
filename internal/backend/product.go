package backend

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ------------------------------
// Generic Utility Functions, not Route Related
//
// These functions provide the functionality the routes need to perform
// actions.
// ------------------------------
func InsertProduct(name string, origin string, ingredient_list []string, isPublic bool, owner *uuid.UUID) (*types.Product, error) {
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
		IsPublic:    isPublic,
		Origin:      &origin,
		Ingredients: ingredients,
	}

	if owner != nil {
		prod.UserId = owner
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
	tx := core.DB.Scopes(WithPreload("Ingredients"), WithCursor(cursor), WithLimit(20), WithOrder("id"),
		WithPreload("Tags"),
	).
		Find(&prods)
	return prods, tx.Error
}

func GetProductById(id string) (*types.Product, error) {
	prod := &types.Product{}
	tx := core.DB.Scopes(WithPreload("Ingredients", "Tags")).First(prod, "id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return prod, nil
}

func GetProductsByName(name string, cursor string) ([]types.Product, error) {
	var products []types.Product

	// we first do a search by primary name
	tx := core.DB.Scopes(WithPreload("Ingredients", "Tags"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
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

func TagProduct(product_id string, tag_id string) error {
	tx := core.DB.Exec("INSERT INTO product_tags (product_id, tag_id) VALUES (?, ?)", product_id, tag_id)
	if tx.Error != nil {
		core.Logger.Error("failed to tag product", "product_id", product_id, "tag_id", tag_id, "error", tx.Error)
		return errors.New("failed to tag product, try again later.")
	}

	return nil
}

func RemoveTagFromProduct(product_id string, tag_id string) error {
	tx := core.DB.Exec("DELETE FROM product_tags WHERE product_id = ? AND tag_id = ?", product_id, tag_id)
	if tx.Error != nil {
		core.Logger.Error("failed to remove tag from product", "product_id", product_id, "tag_id", tag_id, "error", tx.Error)
		return errors.New("failed to remove tag from product, try again later.")
	}

	// TODO: Add to a blacklist to prevent automated tagging from overriding.

	return nil
}

func GetProductsByTag(tag_id string, cursor string) ([]types.Product, error) {
	var products []types.Product
	tx := core.DB.Scopes(WithPreload("Ingredients", "Tags"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
		Where("id IN (SELECT product_id FROM product_tags WHERE tag_id = ?)", tag_id).Find(&products)

	if tx.Error != nil {
		core.Logger.Error("failed to retrieve products by tag", "tag_id", tag_id, "error", tx.Error)
		return nil, errors.New("failed to retrieve products by tag, try again later.")
	}

	core.Logger.Debug("successfully retrieved products by tag", "tag_id", tag_id, "count", len(products), "products", products)
	return products, nil
}

func GetProductsByIngredient(ingredient_ids string, cursor string) ([]types.Product, error) {
	var products []types.Product

	tx := core.DB.Scopes(WithPreload("Ingredients"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
		Where("id IN (SELECT product_id FROM product_ingredients WHERE ingredient_id IN (?))", ingredient_ids).Find(&products)

	if tx.Error != nil {
		core.Logger.Error("failed to retrieve products by ingredient", "ingredient_ids", ingredient_ids, "error", tx.Error)
		return nil, errors.New("failed to retrieve products by ingredient, try again later.")
	}

	core.Logger.Debug("successfully retrieved products by ingredient", "ingredient_ids", ingredient_ids, "count", len(products), "products", products)
	return products, nil
}

func GetProductsByUser(user_id *uuid.UUID, cursor string) ([]types.Product, error) {
	var products []types.Product

	var tx *gorm.DB
	if user_id == nil {
		tx = core.DB.Scopes(WithPreload("Ingredients"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
			Where("user_id IS NULL").Find(&products)
	} else {
		tx = core.DB.Scopes(WithPreload("Ingredients"), WithCursor(cursor), WithLimit(20), WithOrder("id")).
			Where("user_id = ?", user_id).Find(&products)
	}

	if tx.Error != nil {
		core.Logger.Error("failed to retrieve products by user", "user_id", user_id, "error", tx.Error)
		return nil, errors.New("failed to retrieve products by user, try again later.")
	}

	core.Logger.Debug("successfully retrieved products by user", "user_id", user_id, "count", len(products), "products", products)
	return products, nil
}

// ------------------------------
// Routing Related Functions
// ------------------------------
func RouteProductPUT(w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	origin := r.FormValue("origin")
	ingredientList := r.Form["ingredient"]
	isPublic := r.FormValue("is_public") == "true"
	var owner_id *uuid.UUID

	owner, err := user.GetUserFromRequestNoRedirect(w, r)
	if err == nil && owner != nil {
		owner_id = &owner.Id
	}

	prod, err := InsertProduct(name, origin, ingredientList, isPublic, owner_id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(prod)
}

func RouteProductGET(w http.ResponseWriter, r *http.Request) error {
	products, err := GetProducts(r.FormValue("cursor"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully retrieved products", "count", len(products), "products", products)
	return json.NewEncoder(w).Encode(products)
}

func RouteGetProductsByName(w http.ResponseWriter, r *http.Request) error {
	products, err := GetProductsByName(r.FormValue("name"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully retrieved products by name", "name", r.FormValue("name"), "count", len(products), "products", products)
	return json.NewEncoder(w).Encode(products)
}

func RouteGetProductById(w http.ResponseWriter, r *http.Request) error {
	product, err := GetProductById(r.FormValue("id"))
	if err != nil {
		return err
	}

	core.Logger.Debug("successfully retrieved product by id", "id", r.FormValue("id"), "product", product)
	return json.NewEncoder(w).Encode(product)
}

func RouteDeleteProductById(w http.ResponseWriter, r *http.Request) error {
	return DeleteProductById(r.FormValue("id"))
}

func RouteGetProductsByIngredient(w http.ResponseWriter, r *http.Request) error {
	products, err := GetProductsByIngredient(r.FormValue("ingredient_ids"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(products)
}

func RouteGetProductsByTag(w http.ResponseWriter, r *http.Request) error {
	products, err := GetProductsByTag(r.FormValue("tag_id"), r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(products)
}

func RouteGetProductsByUser(w http.ResponseWriter, r *http.Request) error {
	var userId *uuid.UUID

	u, err := user.GetUserFromRequestNoRedirect(w, r)
	if err == nil && u != nil {
		userId = &u.Id
	}

	products, err := GetProductsByUser(userId, r.FormValue("cursor"))
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(products)
}

func RouteProductTag(w http.ResponseWriter, r *http.Request) error {
	return TagProduct(r.FormValue("product_id"), r.FormValue("tag_id"))
}

func RouteProductTagRemove(w http.ResponseWriter, r *http.Request) error {
	return RemoveTagFromProduct(r.FormValue("product_id"), r.FormValue("tag_id"))
}

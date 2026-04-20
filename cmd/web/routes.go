package main

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/handlers/ingredient"
	"github.com/CFF4HA/Dashboard/internal/handlers/product"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
)

func SearchBars(v *verb.Verb) {
	v.Component("v2/searchbars/searchbar-manual_search.html", htmx.Div().Classes("p-2")).
		Bridge(bridges.ManualSearchMultiplexer)
}

func Index(v *verb.Verb) {
	v.Page("", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
}

func Ingredients(v *verb.Verb) {
	v.Page("/ingredients", "v2/pages/ingredients.html")

	v.Component("v2/components/ingredient-search_results.html", htmx.Div()).
		Bridge(bridges.IngredientSearchBridge)

	v.Action(http.MethodPut, "/ingredient", ingredient.RetrieveIngredientHandler)
}

func Products(v *verb.Verb) {
	// this will instantiate the create product route
	v.Action(http.MethodPut, "/product", product.CreateProduct)

	// this is the components that are rendered as part of the
	// product page.
	v.Component("v2/forms/form-product_create.html", htmx.Div())

	// search results component used by the shared searchbar
	v.Component("v2/components/product-search_results.html", htmx.Div()).
		Bridge(bridges.ProductSearchBridge)

	// this will be the primary products page.
	v.Page("/products", "v2/pages/products.html").
		Bridge(bridges.ProductBridge(20, map[string]string{
			"name": " ILIKE ",
		}))
}

func User(v *verb.Verb) {
	// handles the user registration part of the user flow.
	v.Action(http.MethodPut, "/user", user.HandleUserPUT)
	v.Action(http.MethodPost, "/user/login", user.HandleUserPOST)

	v.Component("v2/forms/form-user_signup.html", htmx.Div())

	v.Page("/login", "v2/pages/login.html")
}

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

// SearchBars handles the HTMX component for the manual search bar
func SearchBars(v *verb.Verb) {
	v.Component("v2/searchbars/searchbar-manual_search.html", htmx.Div().Classes("p-2")).
		Bridge(bridges.ManualSearchMultiplexer)
}

// Index handles the main dashboard and core navigation pages
func Index(v *verb.Verb) {
	// Primary Dashboard
	v.Page("", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
	v.Page("/", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
}

// Ingredients handles ingredient-specific pages and search results
func Ingredients(v *verb.Verb) {
	v.Page("/ingredients", "v2/pages/ingredients.html").
		Bridge(bridges.UserSessionRequired{}).
		Bridge(bridges.UserFavoriteIngredientsBridge)

	v.Component("v2/components/ingredient-search_results.html", htmx.Div()).
		Bridge(bridges.IngredientSearchBridge).
		Bridge(bridges.UserFavoriteIngredientsBridge)

	v.Action(http.MethodPut, "/ingredient", ingredient.RetrieveIngredientHandler)
}

// Products handles product-specific pages, creation, and search results
func Products(v *verb.Verb) {
	v.Action(http.MethodPut, "/product", product.CreateProduct)

	v.Component("v2/forms/form-product_create.html", htmx.Div())

	v.Component("v2/components/product-search_results.html", htmx.Div()).
		Bridge(bridges.ProductSearchBridge).
		Bridge(bridges.UserFavoriteProductsBridge)

	// Tab components
	v.Component("v2/components/tab-products.html", htmx.Div())
	v.Component("v2/components/tab-ingredients.html", htmx.Div())
	v.Component("v2/components/tab-compare.html", htmx.Div())
	v.Component("v2/components/tab-user_dashboard.html", htmx.Div())

	// Product section components
	v.Component("v2/components/btn-add_product.html", htmx.Div())

	v.Component("v2/components/section-your_products.html", htmx.Div()).
		Bridge(bridges.UserFavoriteProductsBridge)

	v.Component("v2/components/section-recommended_products.html", htmx.Div())
}

func Compare(v *verb.Verb) {
	v.Component("v2/components/searchbar-compare.html", htmx.Div())
	v.Component("v2/components/table-compare.html", htmx.Div()).
		Bridge(bridges.ProductSearchBridge).
		Bridge(bridges.UserFavoriteProductsBridge)

	v.Component("v2/components/compare-search_results.html", htmx.Div()).
	Bridge(bridges.ProductSearchBridge)
}

// User handles authentication, registration, and questionnaire flows
func User(v *verb.Verb) {
	// User Actions
	v.Action(http.MethodPut, "/user", user.HandleUserPUT)
	v.Action(http.MethodPost, "/user/login", user.HandleUserPOST)
	v.Action(http.MethodPost, "/user/product", user.HandleAddUserProduct)
	v.Action(http.MethodPost, "/user/ingredient", user.HandleAddUserIngredient)
	v.Action(http.MethodGet, "/guest", user.HandleGuestLogin)

	// User Components
	v.Component("v2/forms/form-user_signup.html", htmx.Div())
	v.Component("v2/components/table-tagging_rules.html", htmx.Div()).
    	Bridge(bridges.UserSessionRequired{})

	// Static Pages
	v.Page("/login", "v2/pages/login.html")
	v.Page("/create-login", "v2/pages/createLogin.html")

	// Questionnaire Flow
	v.Page("/questionnaire", "v2/pages/questionnaire1.html")
	v.Page("/questionnaire/step-2", "v2/pages/questionnaire2.html")
	v.Page("/questionnaire/step-3", "v2/pages/questionnaire3.html")

	// Redirect/Landing Helpers
	v.Page("/guest", "v2/pages/index.html")
	v.Page("/user", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
	v.Page("/user/login-success", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
	v.Page("/user/login", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
	v.Page("/user/dashboard", "v2/pages/userDashboard.html").Bridge(bridges.UserSessionRequired{})
	v.Page("/products", "v2/pages/products.html").Bridge(bridges.UserSessionRequired{})
}
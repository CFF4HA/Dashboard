package main

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
)

func SearchBars(v *verb.Verb) {
	v.Component("v2/searchbars/searchbar-manual_search.html", htmx.Div().Classes("p-2")).
		Bridge(bridges.ManualSearchMultiplexer)
}

func Test(v *verb.Verb) {
	v.Page("/test", "v2/pages/test.html")
}

func Index(v *verb.Verb) {
	v.Page("", "v2/pages/index.html").Bridge(bridges.UserSessionRequired{})
}

func Compare(v *verb.Verb) {
	v.Page("/compare", "v2/pages/compare.html")
	v.Component("v2/components/compare-product_search.html", htmx.Div()).
		Bridge(bridges.ProductSearchBridge)
	v.Component("v2/components/investigation.html", htmx.Div()).
		Bridge(bridges.InvestigationBridge{})

	v.Component("v2/components/compare-product_card.html", htmx.Div()).
		Bridge(bridges.ProductDetailBridge)
}

func Ingredients(v *verb.Verb) {
	v.ActionClassic(http.MethodGet, "/ingredient/get", backend.RouteGetIngredients)
	v.ActionClassic(http.MethodGet, "/ingredient/get/name", backend.RouteGetIngredientsByPrimaryName)
	v.ActionClassic(http.MethodPost, "/ingredient/tag", backend.RouteTagIngredient)
	v.ActionClassic(http.MethodDelete, "/ingredient/tag/remove", backend.RouteRemoveTagFromIngredient)
	v.ActionClassic(http.MethodGet, "/ingredient/retrieve", backend.RouteRetrieveIngredientByPrimaryName)
	v.ActionClassic(http.MethodDelete, "/ingredient/remove", backend.RouteRemoveIngredientById)

	v.Page("/ingredients", "v2/pages/ingredients.html").
		Bridge(bridges.CategorizedIngredients{}).
		Bridge(bridges.UserFavoriteIngredientsBridge)

	v.Component("v2/components/ingredient-search_results.html", htmx.Div()).
		Bridge(bridges.CategorizedIngredients{}).
		Bridge(bridges.IngredientSearchBridge).
		Bridge(bridges.UserFavoriteIngredientsBridge)

	v.Component("v2/components/ingredient-detail.html", htmx.Div()).
		Bridge(bridges.IngredientDetailBridge)
}

func Products(v *verb.Verb) {
	// this will instantiate the create product route
	v.ActionClassic(http.MethodPut, "/product/create", backend.RouteProductPUT)
	v.ActionClassic(http.MethodGet, "/product/get", backend.RouteProductGET)
	v.ActionClassic(http.MethodGet, "/product/get/name", backend.RouteGetProductsByName)
	v.ActionClassic(http.MethodGet, "/product/get/id", backend.RouteGetProductById)
	v.ActionClassic(http.MethodGet, "/product/get/ingredients", backend.RouteGetProductsByIngredient)
	v.ActionClassic(http.MethodGet, "/product/get/tag", backend.RouteGetProductsByTag)
	v.ActionClassic(http.MethodDelete, "/product/remove", backend.RouteDeleteProductById)
	v.ActionClassic(http.MethodDelete, "/product/tag/remove", backend.RouteProductTagRemove)
	v.ActionClassic(http.MethodPost, "/product/tag", backend.RouteProductTag)

	// this is the components that are rendered as part of the
	// product page.
	v.Component("v2/forms/form-product_create.html", htmx.Div())

	// search results component used by the shared searchbar
	v.Component("v2/components/product-search_results.html", htmx.Div()).
		Bridge(bridges.CategorizedProducts{}).
		Bridge(bridges.ProductSearchBridge).
		Bridge(bridges.UserFavoriteProductsBridge)

	v.Component("v2/components/product-detail.html", htmx.Div()).
		Bridge(bridges.ProductDetailBridge)

	// this will be the primary products page.
	v.Page("/products", "v2/pages/products.html").
		Bridge(bridges.CategorizedProducts{}).
		Bridge(bridges.ProductBridge(20, map[string]string{
			"name": " ILIKE ",
		})).
		Bridge(bridges.UserFavoriteProductsBridge)
}

func User(v *verb.Verb) {
	// handles the user registration part of the user flow.

	v.Component("v2/forms/form-user_signup.html", htmx.Div())

	v.Page("/login", "v2/pages/login.html")
	v.Page("/user/dashboard", "v2/pages/user.html")
}

func Tagging(v *verb.Verb) {
	v.Component("v2/forms/form-tag_rule_create.html", htmx.Div())
	v.Component("v2/tables/table-tagging_rules.html", htmx.Div()).
		Bridge(bridges.TaggingRules{})
	v.Component("v2/components/tagging-search_results_sm.html", htmx.Div()).
		Bridge(bridges.TagsByName)

	v.Component("v2/components/tag-filter.html", htmx.Div()).
		Bridge(bridges.AllTagsBridge)
}

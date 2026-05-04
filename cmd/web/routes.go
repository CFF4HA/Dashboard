package main

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/handlers/ingredient"
	"github.com/CFF4HA/Dashboard/internal/handlers/product"
	"github.com/CFF4HA/Dashboard/internal/handlers/tagging"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
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
	v.Page("", "v2/pages/index.html")
	v.Page("/v2", "v3/pages/landing.html")
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
	v.Page("/ingredients", "v2/pages/ingredients.html").
		Bridge(bridges.CategorizedIngredients{}).
		Bridge(bridges.UserFavoriteIngredientsBridge)

	v.Component("v2/components/ingredient-search_results.html", htmx.Div()).
		Bridge(bridges.CategorizedIngredients{}).
		Bridge(bridges.IngredientSearchBridge).
		Bridge(bridges.UserFavoriteIngredientsBridge)

	v.Component("v2/components/ingredient-detail.html", htmx.Div()).
		Bridge(bridges.IngredientDetailBridge)

	v.Action(http.MethodPut, "/ingredient", ingredient.RetrieveIngredientHandler)
	v.Action(http.MethodDelete, "/ingredient/delete", ingredient.DeleteIngredientHandler)
	v.Action(http.MethodPost, "/ingredient/sync", ingredient.SyncIngredientHandler)
	v.Action(http.MethodPost, "/ingredient/categorize", ingredient.Categorize)
}

func Products(v *verb.Verb) {
	// this will instantiate the create product route
	v.Action(http.MethodPut, "/product", product.CreateProduct)
	v.Action(http.MethodPost, "/product/update", product.UpdateProductHandler)
	v.Action(http.MethodDelete, "/product/delete", product.DeleteProductHandler)
	v.Action(http.MethodPost, "/product/categorize", product.Categorize)

	// this is the components that are rendered as part of the
	// product page.
	v.Component("v2/forms/form-product_create.html", htmx.Div())

	// search results component used by the shared searchbar
	v.Component("v2/components/product-search_results.html", htmx.Div()).
		Bridge(bridges.CategorizedProducts{}).
		Bridge(bridges.ProductSearchBridgeV2).
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
	v.Action(http.MethodPut, "/user", user.HandleUserPUT)
	v.Action(http.MethodPost, "/user/login", user.HandleUserPOST)
	v.Action(http.MethodPost, "/user/product", user.HandleAddUserProduct)
	v.Action(http.MethodPost, "/user/ingredient", user.HandleAddUserIngredient)

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

	v.Action(http.MethodPut, "/tag/rule/create", tagging.HandleTagRuleCreate)
	v.Action(http.MethodDelete, "/tag/rule/delete", tagging.HandleTagRuleDelete)
	v.Action(http.MethodPost, "/tag/rule/run", tagging.HandleTaggingRuleRun)
	v.Action(http.MethodPost, "/tag/bad", tagging.HandleTagBadIngredients)
	v.Action(http.MethodDelete, "/ingredient/tag", tagging.RemoveIngredientTagHandler)
}

func FormsV2(v *verb.Verb) {
	v.Component("v3/forms/product-create.html", htmx.Div())
	v.Component("v3/forms/product-create_manual.html", htmx.Div())
}

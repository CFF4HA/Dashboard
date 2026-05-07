package main

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend"
	"github.com/CFF4HA/Dashboard/internal/bridges"
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

func Notifications(v *verb.Verb) {
	v.ActionClassic(http.MethodPut, "/notification/create", backend.RouteInsertNotification)
	v.ActionClassic(http.MethodDelete, "/notification/remove", backend.RouteDeleteNotification)
	v.ActionClassic(http.MethodGet, "/notification/get", backend.RouteGetNotifications)
	v.ActionClassic(http.MethodGet, "/notification/get/enabled", backend.RouteGetNotificationsByEnabled)
}

func Ingredients(v *verb.Verb) {
	v.ActionClassic(http.MethodGet, "/ingredient/get", backend.RouteGetIngredients)
	v.ActionClassic(http.MethodGet, "/ingredient/get/name", backend.RouteGetIngredientsByPrimaryName)
	v.ActionClassic(http.MethodGet, "/ingredient/get/id", backend.RouteGetIngredientById)
	v.ActionClassic(http.MethodPost, "/ingredient/tag", backend.RouteTagIngredient)
	v.ActionClassic(http.MethodDelete, "/ingredient/tag/remove", backend.RouteRemoveTagFromIngredient)
	v.ActionClassic(http.MethodGet, "/ingredient/retrieve", backend.RouteRetrieveIngredientByPrimaryName)
	v.ActionClassic(http.MethodDelete, "/ingredient/remove", backend.RouteRemoveIngredientById)
	v.ActionClassic(http.MethodPost, "/ingredient/note", backend.RouteIngredientAddNote)
	v.ActionClassic(http.MethodGet, "/ingredient/note/get", backend.RouteGetIngredientNotes)
	v.ActionClassic(http.MethodDelete, "/ingredient/note/remove", backend.RouteDeleteIngredientNote)
	v.ActionClassic(http.MethodPut, "/ingredient/scrape/config/create", backend.RouteInsertPubChemLabelConfig)
	v.ActionClassic(http.MethodDelete, "/ingredient/scrape/config/remove", backend.RouteDeletePubChemLabelConfig)
	v.ActionClassic(http.MethodGet, "/ingredient/scrape/config/get", backend.RouteGetPubChemLabelConfigs)

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

func Usage(v *verb.Verb) {
	v.ActionClassic(http.MethodGet, "/usage/metrics/latest", backend.RouteGetLatestUsageMetric)
}

func Products(v *verb.Verb) {
	// this will instantiate the create product route
	v.ActionClassic(http.MethodPut, "/product/create", backend.RouteProductPUT)
	v.ActionClassic(http.MethodGet, "/product/get", backend.RouteProductGET)
	v.ActionClassic(http.MethodGet, "/product/get/name", backend.RouteGetProductsByName)
	v.ActionClassic(http.MethodGet, "/product/get/user", backend.RouteGetProductsByUser)
	v.ActionClassic(http.MethodGet, "/product/get/id", backend.RouteGetProductById)
	v.ActionClassic(http.MethodGet, "/product/get/ingredients", backend.RouteGetProductsByIngredient)
	v.ActionClassic(http.MethodGet, "/product/get/tag", backend.RouteGetProductsByTag)
	v.ActionClassic(http.MethodDelete, "/product/remove", backend.RouteDeleteProductById)
	v.ActionClassic(http.MethodDelete, "/product/tag/remove", backend.RouteProductTagRemove)
	v.ActionClassic(http.MethodPost, "/product/tag", backend.RouteProductTag)
	v.ActionClassic(http.MethodGet, "/product/get/operation", backend.RouteGetProductByOperation)
	v.ActionClassic(http.MethodGet, "/product/scrape", backend.RouteInsertProductAutomatedLLM)

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
	// Add the classic actions here

	// handles the user registration part of the user flow.
	v.Action(http.MethodPut, "/user", user.HandleUserPUT)
	v.Action(http.MethodPost, "/user/login", user.HandleUserPOST)
	v.Action(http.MethodPost, "/user/product", user.HandleAddUserProduct)
	v.Action(http.MethodPost, "/user/ingredient", user.HandleAddUserIngredient)

	v.ActionClassic(http.MethodPut, "/v2/user", backend.RouteInsertUser)
	v.ActionClassic(http.MethodPost, "/v2/user/login", backend.RouteUserLogin)
	v.ActionClassic(http.MethodDelete, "/v2/user", backend.RouteDeleteUser)
	v.ActionClassic(http.MethodGet, "/v2/user/get", backend.RouteGetUsers)
	v.ActionClassic(http.MethodPut, "/v2/user/role", backend.RouteUserRoleAdd)
	v.ActionClassic(http.MethodDelete, "/v2/user/role", backend.RouteUserRoleRemove)

	v.Component("v2/forms/form-user_signup.html", htmx.Div())

	v.Page("/login", "v2/pages/login.html")
	v.Page("/user/dashboard", "v2/pages/user.html")
}

func Tagging(v *verb.Verb) {
	v.ActionClassic(http.MethodGet, "/tag/get", backend.RouteTagGet)
	v.ActionClassic(http.MethodGet, "/tag/get/name", backend.RouteTagGetByName)
	v.ActionClassic(http.MethodGet, "/tag/get/id", backend.RouteTagGetById)
	v.ActionClassic(http.MethodDelete, "/tag/remove", backend.RouteTagDelete)
	v.ActionClassic(http.MethodPut, "/tag/create", backend.RouteTagCreate)

	v.ActionClassic(http.MethodPut, "/tag/rule/create", backend.RouteInsertTaggingRule)
	v.ActionClassic(http.MethodDelete, "/tag/rule/remove", backend.RouteDeleteTaggingRule)
	v.ActionClassic(http.MethodGet, "/tag/rule/get/id", backend.RouteGetTaggingRuleById)
	v.ActionClassic(http.MethodGet, "/tag/rule/get/user", backend.RouteGetTaggingRulesForUser)
	v.ActionClassic(http.MethodGet, "/tag/rule/get", backend.RouteGetTaggingRules)

	v.ActionClassic(http.MethodPut, "/tag/set/create", backend.RouteInsertTaggingSet)
	v.ActionClassic(http.MethodDelete, "/tag/set/remove", backend.RouteDeleteTaggingSet)
	v.ActionClassic(http.MethodGet, "/tag/set/get", backend.RouteGetTaggingSets)
	v.ActionClassic(http.MethodGet, "/tag/set/get/id", backend.RouteGetTaggingSetById)
	v.ActionClassic(http.MethodGet, "/tag/set/get/name", backend.RouteGetTaggingSetByName)

	v.Component("v3/components/tag/tag-set-form.html", htmx.Div())
	v.Component("v3/components/tag/tag-set-mini-view.html", htmx.Div()).
		Bridge(verb.Map("TagSets", func(r *http.Request, m map[string]any) (any, error) {
			sets, err := backend.GetTaggingSets(r.FormValue("cursor"))
			if err != nil {
				return nil, err
			}

			return sets, nil
		})).
		Bridge(verb.Map("Tags", func(r *http.Request, m map[string]any) (any, error) {
			tags, err := backend.GetTags(r.FormValue("cursor"))
			if err != nil {
				return nil, err
			}

			return tags, nil
		}))

	v.Component("v3/components/tag/tag-set-detail-view.html", htmx.Div())

	v.Component("v2/forms/form-tag_rule_create.html", htmx.Div())
	v.Component("v2/tables/table-tagging_rules.html", htmx.Div()).
		Bridge(bridges.TaggingRules{})
	v.Component("v2/components/tagging-search_results_sm.html", htmx.Div()).
		Bridge(bridges.TagsByName)

	v.Component("v2/components/tag-filter.html", htmx.Div()).
		Bridge(bridges.AllTagsBridge)
}

func FormsV2(v *verb.Verb) {
	v.Component("v3/forms/product-create.html", htmx.Div())
	v.Component("v3/forms/product-create_manual.html", htmx.Div())
}

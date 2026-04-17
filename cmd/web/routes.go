package main

import (
	"gorm.io/gorm"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
	"github.com/DAlba-sudo/verbs"
	verbsgorm "github.com/DAlba-sudo/verbs/gorm"
)

// registerRoutes builds all shared bridges and wires up pages and components.
func registerRoutes(v *verb.Verb, cfg Config, dbconn *gorm.DB) {
	// ingredientWithCache retries acquisition via bridges.Ingredient on a
	// cache miss, capped at 10 concurrent acquisitions, keyed by id or primary_name.
	ingredientWithCache := verbs.QueryCachedResource(bridges.KeyIngredient, verbs.QueryCachedResourceOptions{
		AcquisitionPolicy:         bridges.Ingredient,
		MaxConcurrentAcquisitions: 10,
		Fingerprint:               []string{"id", "primary_name"},
	})

	// ingredient queries a single Ingredient row with regex matching on primary_name.
	// The concrete verbsgorm type is kept (not cast to verb.Bridge) so that
	// SetModifiers can be called further below.
	ingredient := verbsgorm.GORM(
		bridges.KeyIngredient,
		&types.Ingredient{},
		&types.Ingredient{},
		dbconn,
		verbsgorm.GormOptions{
			Select:  "*",
			Preload: []string{"Names", "Labels"},
			Limit:   1,
			KeyModifiers: map[string]string{
				"primary_name": " ~* ",
			},
		},
	)

	// names queries Name rows with regex matching on the text field.
	names := verbsgorm.GORM(
		bridges.KeyNames,
		&types.Name{},
		types.Name{},
		dbconn,
		verbsgorm.GormOptions{
			Select:  "*",
			Preload: []string{"Ingredients"},
			KeyModifiers: map[string]string{
				"text": " ~* ",
			},
		},
	)

	// --- Pages ---

	v.Page("", "v2/pages/index.html")
	v.Page("/products", "v2/pages/products.html")
	v.Page("/ingredients", "v2/pages/ingredients.html").
		Bridge(ingredient)

	// --- Components ---

	// Search bar with live health checks for database and LLM service.
	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div")).
		Bridge(bridges.DatabaseChecker).
		Bridge(bridges.NewLLMChecker(cfg.LLMAddr))

	// AI-powered search router: classifies user input into an intent and
	// decodes the intent-specific payload for the template.
	v.Component("v2/components/ai/router.html", htmx.Div()).
		Bridge(bridges.Aux(cfg.LLMAddr, cfg.CFClientID, cfg.CFClientSecret)).
		Bridge(bridges.AuxRouter)

	// Ingredient search: list view.
	// On a cache miss the component re-polls itself every 5 seconds while
	// a background sync completes.
	retryList := htmx.Div().
		GET("/htmx/ingredient-search_result").
		SelfEncodeRequest().
		Trigger("load delay:5s").
		Swap("outerHTML")

	v.Component("v2/components/ingredient/ingredient-search_result.html", htmx.Create("div")).
		Bridge(ingredient.SetModifiers(map[string]string{"primary_name": " ~* "})).
		OnError(retryList).
		OnError(ingredientWithCache)

	// Ingredient search: single result with label metadata.
	// On a cache miss the component re-polls itself every 3 seconds.
	retrySingle := htmx.Div().
		GET("/htmx/ingredient-search_result_single").
		SelfEncodeRequest().
		Trigger("load delay:3s").
		Swap("outerHTML")

	v.Component("v2/components/ingredient/ingredient-search_result_single.html", htmx.Create("div")).
		Bridge(ingredient.SetModifiers(map[string]string{"primary_name": " ~* "})).
		OnError(retrySingle).
		OnError(ingredientWithCache).
		Bridge(verb.Map(bridges.KeyMetadata, bridges.IngredientMetadataProvider))

	// Ingredient manual search bar and its result list.
	v.Component("v2/components/ingredient/ingredient-search_bar_manual.html", htmx.Create("div"))
	v.Component("v2/components/ingredient/ingredient-search_bar_manual_res.html", htmx.Create("div")).
		Bridge(names)

	// Ingredient sync status indicator.
	v.Component("v2/components/ingredient/_ingredient_sync.html", htmx.Create("div")).
		Bridge(ingredientWithCache)

	// Product creation form (no bridges — purely client-driven).
	v.Component("v2/components/product/product-creation_form.html", htmx.Create("div"))
}

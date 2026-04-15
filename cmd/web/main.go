package main

import (
	"flag"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"

	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
	"github.com/DAlba-sudo/verbs"
	"github.com/DAlba-sudo/verbs/gorm"
)

// This is all basic boiler plate, as the frontend you will not have to touch this.
func main() {
	templateDir := flag.String("template", "templates", "the directory for the templates")
	staticDir := flag.String("static", "static", "the directory for the static files")
	address := flag.String("address", "127.0.0.1", "the address to bind the server to")
	port := flag.Int("port", 8080, "the port to bind the server to")
	live := flag.Bool("reload", true, "whether to do live reloads on the template files")
	backend := flag.String("backend", "http://localhost:8081", "the address of the backend server")
	llm := flag.String("llm", "http://localhost:8082", "the address of the LLM server")
	db := flag.String("db", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "the connection string for the database")
	cf_client_id := flag.String("cf_client_id", "", "the client id for the cloudflare instance")
	cf_client_secret := flag.String("cf_client_secret", "", "the client secret for the cloudflare instance")
	flag.Parse()

	bridges.LLM = llm
	core.BackendAddress = strings.TrimRight(*backend, "/")
	if err := database.Initialize(*db); err != nil {
		core.Logger.Error("failed to initialize database", "error", err)
	}

	dbconn, err := database.Database()
	if err != nil {
		core.Logger.Error("failed to connect to database", "error", err)
	}

	v := verb.New(*address, *port, verb.Settings{
		Templates:  *templateDir,
		Static:     *staticDir,
		LiveReload: *live,
		Bridges: []verb.Bridge{
			verbs.Request{},
		},
	})
	v.Func("lower", func(s string) string {
		return strings.ToLower(s)
	})
	v.Func("default", func(a any, b any) any {
		if b == nil {
			return a
		}

		return b
	})
	v.Func("split", func(a any, b string) []string {
		return strings.Split(a.(string), b)
	})
	v.Func("trim", func(a any) string {
		if _, ok := a.(string); !ok {
			return ""
		}
		return strings.TrimSpace(a.(string))
	})

	// Global Bridges (i.e., where state is important)
	IngredientWithCache := verbs.QueryCachedResource("Ingredient", verbs.QueryCachedResourceOptions{
		AcquisitionPolicy:         bridges.Ingredient,
		MaxConcurrentAcquisitions: 10,
		Fingerprint:               []string{"id", "primary_name"},
	})

	Ingredient := gorm.GORM(
		"Ingredient",
		&types.Ingredient{},
		&types.Ingredient{},
		dbconn,
		gorm.GormOptions{
			Select:  "*",
			Preload: []string{"Names", "Labels"},
			Limit:   1,
			KeyModifiers: map[string]string{
				"primary_name": " ~* ",
			},
		},
	)

	Names := gorm.GORM(
		"Names",
		&types.Name{},
		types.Name{},
		dbconn,
		gorm.GormOptions{
			Select:  "*",
			Preload: []string{"Ingredients"},
			KeyModifiers: map[string]string{
				"text": " ~* ",
			},
		},
	)

	// Index Page
	v.Page("", "v2/pages/index.html")
	v.Page("/products", "v2/pages/products.html")
	v.Page("/ingredients", "v2/pages/ingredients.html").
		Bridge(Ingredient)

	// Aux Router Search Bar
	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div")).
		Bridge(bridges.DatabaseChecker).
		Bridge(bridges.LMMChecker)

	v.Component("v2/components/ai/router.html", htmx.Div()).
		Bridge(bridges.Aux(*llm, *cf_client_id, *cf_client_secret)).
		Bridge(bridges.AuxRouter)

	// Ingredient Search Result, Individual
	v.Component("v2/components/ingredient/ingredient-search_result.html", htmx.Create("div")).
		Bridge(Ingredient.SetModifiers(map[string]string{
			"primary_name": " ~* ",
		})).
		OnError(htmx.Div().
			GET("/htmx/ingredient-search_result").
			SelfEncodeRequest().
			Trigger("load delay:5s").
			Swap("outerHTML")).
		OnError(IngredientWithCache)

	v.Component("v2/components/ingredient/ingredient-search_result_single.html", htmx.Create("div")).
		Bridge(Ingredient.SetModifiers(map[string]string{
			"primary_name": " ~* ",
		})).
		OnError(htmx.Div().
			GET("/htmx/ingredient-search_result_single").
			SelfEncodeRequest().
			Trigger("load delay:3s").
			Swap("outerHTML")).
		OnError(IngredientWithCache).
		Bridge(verb.Map("Metadata", bridges.IngredientMetadataProvider))

	v.Component("v2/components/ingredient/ingredient-search_bar_manual.html", htmx.Create("div"))
	v.Component("v2/components/ingredient/ingredient-search_bar_manual_res.html", htmx.Create("div")).
		Bridge(Names)

	v.Component("v2/components/ingredient/_ingredient_sync.html", htmx.Create("div")).
		Bridge(IngredientWithCache)

	// Product
	v.Component("v2/components/product/product-creation_form.html", htmx.Create("div"))

	if err := v.Serve(); err != nil {
		panic(err)
	}
}

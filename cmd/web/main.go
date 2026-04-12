package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/core"
	auxrouter "github.com/CFF4HA/Dashboard/pkg/aux-router"

	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
	"github.com/DAlba-sudo/verbs"
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
	flag.Parse()

	core.BackendAddress = strings.TrimRight(*backend, "/")
	database.Initialize(*db)

	v := verb.New(*address, *port, verb.Settings{
		Templates:  *templateDir,
		Static:     *staticDir,
		LiveReload: *live,
		Bridges:    []verb.Bridge{verbs.Request{}},
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

	// Global Bridges (i.e., where state is important)
	IngredientWithCache := verbs.QueryCachedResource("Ingredient", verbs.QueryCachedResourceOptions{
		AcquisitionPolicy:         bridges.Ingredient,
		MaxConcurrentAcquisitions: 10,
		Fingerprint:               []string{"id", "query"},
	})

	IngredientNames := verbs.QueryParameter("query", &verbs.QueryParameterOptions{
		Name: "Ingredients",
	})
	Name := verbs.QueryParameter("name", &verbs.QueryParameterOptions{
		Name:  "Name",
		First: true,
	})

	// Index Page
	v.Page("", "v2/pages/index.html")

	// Aux Router Search Bar
	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div"))
	v.Component("v2/components/ai/router.html", htmx.Div()).
		Bridge(bridges.Aux(*llm)).
		Bridge(verb.Map("Payload", func(r *http.Request, m map[string]any) (any, error) {
			_, ok := m["Aux"]
			if !ok {
				return nil, nil
			}

			aux_response, ok := m["Aux"].(auxrouter.Response)
			if !ok {
				return nil, nil
			}

			switch aux_response.Intent {
			case "single_ingredient_search":
				var Ingredient struct {
					Ingredient string `json:"ingredient"`
				}

				if json.Unmarshal(aux_response.Response, &Ingredient) != nil {
					return nil, nil
				}

				return Ingredient, nil

			case "multi_ingredient_search":
				var Ingredients struct {
					Ingredients []string `json:"ingredients"`
				}
				if json.Unmarshal(aux_response.Response, &Ingredients) != nil {
					return nil, nil
				}

				return Ingredients, nil
			case "product_create":
				var Product struct {
					Product     string   `json:"product"`
					Ingredients []string `json:"ingredients"`
				}

				if json.Unmarshal(aux_response.Response, &Product) != nil {
					return nil, nil
				}

				return Product, nil
			}

			return nil, nil
		}))

	// Ingredient Search Result, Individual
	// Query Parameters: `id`, `query`
	v.Component("v2/components/ingredient/ingredient-search_result.html", htmx.Create("div")).
		Bridge(IngredientWithCache).
		OnError(htmx.Div().
			GET("/htmx/ingredient-search_result").
			SelfEncodeRequest().
			Trigger("load delay:5s").
			Swap("outerHTML"))

	// Query Parameters: `id`, `query`
	v.Component("v2/components/ingredient/ingredient-search_result_single.html", htmx.Create("div")).
		OnError(htmx.Div().
			GET("/htmx/ingredient-search_result_single").
			SelfEncodeRequest().
			Trigger("load delay:3s").
			Swap("outerHTML")).
		Bridge(IngredientWithCache).
		Bridge(verb.Map("Metadata", bridges.IngredientMetadataProvider))

	// Query Parameters: `name`, `ingredients`
	v.Component("v2/components/ingredient/ingredient-list_stats.html", htmx.Div().GET("/htmx/ingredient-list_stats").
		Trigger("load delay:7s").
		SelfEncodeRequest()).
		Bridge(verbs.QueryParameter("ingredients", &verbs.QueryParameterOptions{Name: "IngredientList"})).
		Bridge(verb.Map("Ingredients", bridges.Ingredients)).
		Bridge(verb.Map("Counts", bridges.CountLabelsForIngredients))

	// Product
	v.Component("v2/components/product/product-creation_form.html", htmx.Create("div")).
		Bridge(Name).
		Bridge(IngredientNames)

	if err := v.Serve(); err != nil {
		panic(err)
	}
}

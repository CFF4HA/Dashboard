package main

import (
	"flag"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
	"github.com/google/uuid"
)

var (
	ai = bridges.AIRouter{
		Sessions: make(map[string]bridges.View),
		Lock:     &sync.RWMutex{},
	}
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
	flag.Parse()

	core.BackendAddress = strings.TrimRight(*backend, "/")
	core.LLMAddress = strings.TrimRight(*llm, "/")

	v := verb.New(*address, *port, verb.Settings{
		Templates:  *templateDir,
		Static:     *staticDir,
		LiveReload: *live,
	})

	v.Page("", "v2/pages/index.html").
		Bridge(verb.DataBridge{
			Key: "TotalIngredients",
			Provide: func(r *http.Request) (any, error) {
				return 3, nil
			},
		}).
		Bridge(verb.DataBridge{
			Key: "TotalProducts",
			Provide: func(r *http.Request) (any, error) {
				return 25, nil
			},
		})
	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div"))
	v.Component("v2/components/product/product-search_result.html", htmx.Create("div")).
		Bridge(verb.DataBridge{
			Key: "Product",
			Provide: func(r *http.Request) (any, error) {
				return types.Product{
					Id:      uuid.New(),
					Name:    "Example Product",
					Origin:  "https://www.vanicream.com/product/free-and-clear-shampoo",
					Created: time.Now().Add(-time.Hour),
					Updated: time.Now(),
					Ingredients: []types.Ingredient{
						{
							Names:  []types.Name{{Text: "Example Name 1"}, {Text: "Example Name 2"}},
							Labels: []types.Label{{Type: "hazard", Payload: "Example Hazard", Origin: "http://example.com/hazard"}},
						},
					},
				}, nil
			},
		})

	v.Component("v2/components/ingredient/ingredient-search_result.html", htmx.Create("div")).
		Bridge(verb.DataBridge{
			Key:     "Ingredient",
			Provide: bridges.IngredientSearchProvider,
		}).
		Bridge(verb.DataBridge{
			Key: "Htmx",
			Provide: func(r *http.Request) (any, error) {
				if _, ok := core.IngredientsCache[r.FormValue("name")]; ok {
					return htmx.Div(), nil
				}

				name := r.FormValue("name")
				return htmx.Div().Trigger("load delay:2s").GET("/htmx/ingredient-search_result?name=" + url.QueryEscape(name)).Swap("outerHTML"), nil
			},
		})

	v.Component("v2/components/ai/router.html", htmx.Div()).
		Bridge(ai)

	if err := v.Serve(); err != nil {
		panic(err)
	}
}

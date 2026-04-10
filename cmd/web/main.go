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
	})

	v.Page("", "v2/pages/index.html")

	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div"))

	v.Component("v2/components/ingredient/ingredient-search_result.html", htmx.Create("div")).
		Bridge(verb.Map("Ingredient", bridges.IngredientByNameProvider))

	v.Component("v2/components/ai/router.html", htmx.Div()).
		Bridge(bridges.Aux(*llm)).
		Bridge(verb.Map("Ingredient", func(r *http.Request, m map[string]any) (any, error) {
			_, ok := m["Aux"]
			if !ok {
				return nil, nil
			}

			aux_response, ok := m["Aux"].(auxrouter.Response)
			if !ok {
				return nil, nil
			}

			var Ingredient struct {
				Ingredient string `json:"ingredient"`
			}

			if aux_response.Intent == "single_ingredient_search" {
				if json.Unmarshal(aux_response.Response, &Ingredient) != nil {
					return nil, nil
				}

				return Ingredient, nil
			}

			return nil, nil
		}))

	if err := v.Serve(); err != nil {
		panic(err)
	}
}

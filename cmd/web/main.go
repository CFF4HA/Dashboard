package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"

	"github.com/CFF4HA/Dashboard/pkg/aux-router"
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
	flag.Parse()

	core.BackendAddress = strings.TrimRight(*backend, "/")

	v := verb.New(*address, *port, verb.Settings{
		Templates:  *templateDir,
		Static:     *staticDir,
		LiveReload: *live,
	})

	v.Page("", "v2/pages/index.html")
	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div"))

	aux := auxrouter.Aux{LlmServerEndpoint: *llm}
	aux.Intent("single_ingredient_search", "user wishes to search for single ingredient",
		auxrouter.IntentExample{
			Prompt:         "Water",
			ResponseSchema: `{"ingredient": "Water"}`,
		},
		auxrouter.IntentExample{
			Prompt:         "What is Citral?",
			ResponseSchema: `{"ingredient": "Citral"}`,
		},
	)
	aux.Intent("unknown", "the intent is unknown or cannot be determined with certainty (note: return this if user input seems questionable)",
		auxrouter.IntentExample{
			Prompt:         "What do you know about the universe?",
			ResponseSchema: `{}`,
		})
	aux.Intent("multi_ingredient_search", "user wishes to search for multiple ingredients",
		auxrouter.IntentExample{
			Prompt:         "Water, Citral, and Limonene",
			ResponseSchema: `{"ingredients": ["Water", "Citral", "Limonene"]}`,
		},
		auxrouter.IntentExample{
			Prompt:         "What do you know about Water, Citral, and Limonene?",
			ResponseSchema: `{"ingredients": ["Water", "Citral", "Limonene"]}`,
		},
	)
	aux.Intent("product_search_name", "user wishes to search for a specific product by name",
		auxrouter.IntentExample{
			Prompt:         "Tell me about the ingredient in CeraVe Hydrating Cleanser",
			ResponseSchema: `{"product": "CeraVe Hydrating Cleanser"}`,
		},
		auxrouter.IntentExample{
			Prompt:         "VaniCream Shampoo",
			ResponseSchema: `{"product": "VaniCream Shampoo"}`,
		},
	)
	v.Component("v2/components/ai/router.html", htmx.Div()).
		Bridge(aux).
		Bridge(verb.Map("Ingredient", func(r *http.Request, m map[string]any) (any, error) {
			_, ok := m[aux.Name()]
			if !ok {
				return nil, nil
			}

			aux_response, ok := m[aux.Name()].(auxrouter.Response)
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

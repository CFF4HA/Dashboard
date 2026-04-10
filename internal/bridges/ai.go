package bridges

import (
	"github.com/CFF4HA/Dashboard/pkg/aux-router"
)

var ()

func Aux(llm string) auxrouter.Aux {
	aux := auxrouter.Aux{LlmServerEndpoint: llm}

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

	return aux
}

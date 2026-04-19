package bridges

import (
	"net/http"
	"strings"

	"github.com/DAlba-sudo/verb"
)

var (
	ManualSearchMultiplexer = verb.Map("Search", func(r *http.Request, m map[string]any) (any, error) {
		// the key is the search mode that you'd like to show as "active"
		ingredient := &ManualSearchPayoad{
			Target:           "/htmx/ingredient-search_results",
			ActiveSearchMode: "Ingredient",
			Description:      "Search for ingredients in our database. This will return a list of ingredients that match your search query.",
		}
		product := &ManualSearchPayoad{
			Target:           "/htmx/product-search_results",
			ActiveSearchMode: "Product",
			Description:      "Search for products in our database. This will return a list of products that match your search query.",
		}

		var payload struct {
			Active *ManualSearchPayoad
			Others []*ManualSearchPayoad
		}

		modes := []*ManualSearchPayoad{ingredient, product}
		for _, mode := range modes {
			if mode.ActiveSearchMode == strings.TrimSpace(r.FormValue("search_mode")) {
				payload.Active = mode
				continue
			}

			payload.Others = append(payload.Others, mode)
		}

		if payload.Active == nil {
			payload.Active = ingredient
			payload.Others = []*ManualSearchPayoad{product}
		}

		return payload, nil
	})
)

type ManualSearchPayoad struct {
	// This is the search bar data destination
	// (i.e., where the data will be sent)
	Target string `json:"target"`

	// This is the human readable name of the search mode that
	// is currently being used.
	ActiveSearchMode string `json:"activeSearchMode"`

	// This is the description that will go in the bottom of the search bar.
	Description string `json:"description,omitempty"`
}

package bridges

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

type DraftProductIngredient struct{}

func (d DraftProductIngredient) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	name := strings.Trim(r.FormValue("name"), " \t\r\n")
	if name == "" {
		return nil, nil
	}

	req, err := http.NewRequest("GET", core.BackendAddress+"/product/draft?ingredients="+url.QueryEscape(name), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var ProductDraft types.ProductDraft
	if err := json.NewDecoder(resp.Body).Decode(&ProductDraft); err != nil {
		return nil, err
	}

	return ProductDraft.Ingredients[0], nil
}

type DraftProduct struct{}

func (d DraftProduct) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	ingredient_as_query := url.QueryEscape(r.FormValue("ingredients"))
	req, err := http.NewRequest("GET", core.BackendAddress+"/product/draft?ingredients="+ingredient_as_query, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var ProductDraft types.ProductDraft
	if err := json.NewDecoder(resp.Body).Decode(&ProductDraft); err != nil {
		return nil, err
	}

	return ProductDraft, nil
}

type CreateProduct struct{}

func (c CreateProduct) Data(w http.ResponseWriter, r *http.Request) (any, error) {

	return nil, nil
}

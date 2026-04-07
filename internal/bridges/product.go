package bridges

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

type ProductCount struct{}

func (p ProductCount) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	var products []types.Product

	ids := strings.Trim(r.FormValue("ids"), " \t\r\n")
	for _, id := range strings.Split(ids, ",") {
		var p types.Product
		req, err := http.NewRequest("GET", core.BackendAddress+"/product?id="+url.QueryEscape(id), nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			continue
		}

		products = append(products, p)
	}

	var CountMetadata struct {
		Hazards  []int
		Symptoms []int
		Effects  []int
	}

	for _, p := range products {
		hazards := 0
		symptoms := 0
		effects := 0

		for _, ingredient := range p.Ingredients {

			for _, label := range ingredient.Labels {
				if label.Type == "hazard" {
					hazards++
				} else if label.Type == "symptom" {
					symptoms++
				} else if label.Type == "effect" {
					effects++
				}
			}

		}
		CountMetadata.Hazards = append(CountMetadata.Hazards, hazards)
		CountMetadata.Symptoms = append(CountMetadata.Symptoms, symptoms)
		CountMetadata.Effects = append(CountMetadata.Effects, effects)

	}

	return CountMetadata, nil
}

type ProductComparator struct{}

func (p ProductComparator) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	ids := strings.Trim(r.FormValue("ids"), " \t\r\n")

	return strings.Split(ids, ","), nil
}

type ProductBridge struct{}

func (p ProductBridge) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	id := strings.Trim(r.FormValue("id"), " \t\r\n")
	req, err := http.NewRequest("GET", core.BackendAddress+"/product?id="+url.QueryEscape(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var prod types.Product
	if err := json.NewDecoder(resp.Body).Decode(&prod); err != nil {
		return nil, err
	}

	return prod, nil
}

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

type ProductList struct{}

func (p ProductList) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	var products []types.Product

	req, err := http.NewRequest("GET", core.BackendAddress+"/product", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil
}

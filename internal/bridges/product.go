package bridges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

// backendGet sends a GET request to the backend service at the given path
// and decodes the JSON response body into target.
func backendGet(path string, target any) error {
	resp, err := http.DefaultClient.Get(core.BackendAddress + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// ProductCount fetches each product in a comma-separated "ids" form value and
// tallies hazard, symptom, and effect label counts across all of them.
type ProductCount struct{}

func (p ProductCount) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	ids := strings.Trim(r.FormValue("ids"), " \t\r\n")

	var products []types.Product
	for _, id := range strings.Split(ids, ",") {
		var prod types.Product
		if err := backendGet("/product?id="+url.QueryEscape(id), &prod); err != nil {
			continue
		}
		products = append(products, prod)
	}

	var counts struct {
		Hazards  []int
		Symptoms []int
		Effects  []int
	}
	for _, prod := range products {
		h, s, e := 0, 0, 0
		for _, ingredient := range prod.Ingredients {
			for _, label := range ingredient.Labels {
				switch label.Type {
				case "hazard":
					h++
				case "symptom":
					s++
				case "effect":
					e++
				}
			}
		}
		counts.Hazards = append(counts.Hazards, h)
		counts.Symptoms = append(counts.Symptoms, s)
		counts.Effects = append(counts.Effects, e)
	}

	return counts, nil
}

// ProductComparator parses a comma-separated "ids" form value into a slice
// of ID strings for use in comparison views.
type ProductComparator struct{}

func (p ProductComparator) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	ids := strings.Trim(r.FormValue("ids"), " \t\r\n")
	return strings.Split(ids, ","), nil
}

// ProductSearch searches for products by name.
type ProductSearch struct{}

func (p ProductSearch) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	name := strings.Trim(r.FormValue("name"), " \t\r\n")
	var products []types.Product
	if err := backendGet("/product/search?name="+url.QueryEscape(name), &products); err != nil {
		return nil, err
	}
	return products, nil
}

// ProductBridge fetches a single product by ID.
type ProductBridge struct{}

func (p ProductBridge) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	id := strings.Trim(r.FormValue("id"), " \t\r\n")
	var prod types.Product
	if err := backendGet("/product?id="+url.QueryEscape(id), &prod); err != nil {
		return nil, err
	}
	return prod, nil
}

// DraftProductIngredient resolves a single ingredient name against the backend
// draft endpoint, returning the first matched ingredient.
type DraftProductIngredient struct{}

func (d DraftProductIngredient) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	name := strings.Trim(r.FormValue("name"), " \t\r\n")
	if name == "" {
		return nil, nil
	}

	var draft types.ProductDraft
	if err := backendGet("/product/draft?ingredients="+url.QueryEscape(name), &draft); err != nil {
		return nil, err
	}

	if len(draft.Ingredients) == 0 {
		return nil, fmt.Errorf("no ingredients returned from draft for %q", name)
	}
	return draft.Ingredients[0], nil
}

// DraftProduct resolves a full comma-separated ingredient list against the
// backend draft endpoint, returning the complete ProductDraft.
type DraftProduct struct{}

func (d DraftProduct) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	query := "/product/draft?ingredients=" + url.QueryEscape(r.FormValue("ingredients"))
	var draft types.ProductDraft
	if err := backendGet(query, &draft); err != nil {
		return nil, err
	}
	return draft, nil
}

// CreateProduct is a placeholder for the product creation bridge.
// TODO: implement product creation via the backend API.
type CreateProduct struct{}

func (c CreateProduct) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	return nil, nil
}

// ProductList fetches all products from the backend.
type ProductList struct{}

func (p ProductList) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	var products []types.Product
	if err := backendGet("/product", &products); err != nil {
		return nil, err
	}
	return products, nil
}

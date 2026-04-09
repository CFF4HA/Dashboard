package bridges

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func IngredientSearchProvider(r *http.Request) (any, error) {
	var ingredients types.Ingredient
	name := strings.Trim(r.FormValue("name"), " \r\n\t")
	if name == "" {
		return nil, nil
	} else if strings.Contains(name, ",") {
		return nil, nil
	}

	req, err := http.NewRequest("GET", core.BackendAddress+"/ingredient/search?name="+url.QueryEscape(name), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&ingredients); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	core.IngredientsCache[name] = ingredients
	return ingredients, nil
}

type SearchIngredients struct{}

func (s SearchIngredients) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	var ingredients []types.Ingredient
	name := strings.Trim(r.FormValue("name"), " \r\n\t")
	if name == "" {
		return nil, nil
	} else if strings.Contains(name, ",") {
		return nil, nil
	}

	req, err := http.NewRequest("GET", core.BackendAddress+"/ingredient/search?name="+url.QueryEscape(name), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ingredients); err != nil {
		return nil, err
	}

	return ingredients, nil
}

type IngredientBridge struct{}

func (i IngredientBridge) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	id := strings.Trim(r.FormValue("id"), " \r\n\t")

	req, err := http.NewRequest("GET", core.BackendAddress+"/ingredient?id="+url.QueryEscape(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var Ingredient types.Ingredient
	if err := json.NewDecoder(resp.Body).Decode(&Ingredient); err != nil {
		return nil, err
	}

	return Ingredient, nil
}

package bridges

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

type SearchIngredients struct{}

func (s SearchIngredients) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	var ingredients []types.Ingredient
	name := r.FormValue("name")
	if name == "" {

	} else if strings.Contains(name, ",") {
		return nil, nil
	}

	req, err := http.NewRequest("GET", core.BackendAddress+"/ingredient/search?name="+r.FormValue("name"), nil)
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

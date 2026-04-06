package bridges

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

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

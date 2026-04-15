package bridges

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/DAlba-sudo/auxrouter"
	"github.com/DAlba-sudo/verb"
)

var (
	LLM *string

	DatabaseChecker = verb.Map("DatabaseErr", func(r *http.Request, m map[string]any) (any, error) {
		_, err := database.Database()
		if err != nil {
			return err.Error(), nil
		}

		return nil, nil
	})

	LMMChecker = verb.Map("LLMErr", func(r *http.Request, m map[string]any) (any, error) {
		resp, err := http.DefaultClient.Get(*LLM)
		if err != nil {
			return err.Error(), nil
		} else if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return errors.New("LLM server returned non-200 status code").Error(), nil
		}
		defer resp.Body.Close()

		return nil, nil
	})

	AuxRouter = verb.Map("Payload", func(r *http.Request, m map[string]any) (any, error) {
		_, ok := m["Aux"]
		if !ok {
			return nil, nil
		}

		aux_response, ok := m["Aux"].(auxrouter.Response)
		if !ok {
			return nil, nil
		}

		switch aux_response.Intent {
		case "single_ingredient_search":
			var Ingredient struct {
				Ingredient string `json:"ingredient"`
			}

			if json.Unmarshal(aux_response.Response, &Ingredient) != nil {
				return nil, nil
			}

			return Ingredient, nil

		case "multi_ingredient_search":
			var Ingredients struct {
				Ingredients []string `json:"ingredients"`
			}
			if json.Unmarshal(aux_response.Response, &Ingredients) != nil {
				return nil, nil
			}

			return Ingredients, nil
		case "product_create":
			var Product struct {
				Product     string   `json:"product"`
				Ingredients []string `json:"ingredients"`
			}

			if json.Unmarshal(aux_response.Response, &Product) != nil {
				return nil, nil
			}

			return Product, nil
		}

		return nil, nil
	})
)

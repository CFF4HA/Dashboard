package bridges

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/DAlba-sudo/auxrouter"
	"github.com/DAlba-sudo/verb"
)

// DatabaseChecker is a bridge that reports database health.
// It exposes any connection error as a string under KeyDatabaseErr.
// A nil value means the database is reachable.
var DatabaseChecker = verb.Map(KeyDatabaseErr, func(r *http.Request, m map[string]any) (any, error) {
	_, err := database.Database()
	if err != nil {
		return err.Error(), nil
	}
	return nil, nil
})

// NewLLMChecker returns a bridge that reports LLM service health.
// It exposes any reachability error as a string under KeyLLMErr.
// A nil value means the service is reachable and returned HTTP 200.
func NewLLMChecker(addr string) verb.Bridge {
	return verb.Map(KeyLLMErr, func(r *http.Request, m map[string]any) (any, error) {
		resp, err := http.DefaultClient.Get(addr)
		if err != nil {
			return err.Error(), nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.New("LLM server returned non-200 status code").Error(), nil
		}
		return nil, nil
	})
}

// AuxRouter reads the response from the Aux bridge and decodes the
// intent-specific payload, exposing it under KeyPayload.
var AuxRouter = verb.Map(KeyPayload, func(r *http.Request, m map[string]any) (any, error) {
	raw, ok := m[KeyAux]
	if !ok {
		return nil, nil
	}

	auxResp, ok := raw.(auxrouter.Response)
	if !ok {
		return nil, nil
	}

	switch auxResp.Intent {
	case "single_ingredient_search":
		var payload struct {
			Ingredient string `json:"ingredient"`
		}
		if json.Unmarshal(auxResp.Response, &payload) != nil {
			return nil, nil
		}
		return payload, nil

	case "multi_ingredient_search":
		var payload struct {
			Ingredients []string `json:"ingredients"`
		}
		if json.Unmarshal(auxResp.Response, &payload) != nil {
			return nil, nil
		}
		return payload, nil

	case "product_create":
		var payload struct {
			Product     string   `json:"product"`
			Ingredients []string `json:"ingredients"`
		}
		if json.Unmarshal(auxResp.Response, &payload) != nil {
			return nil, nil
		}
		return payload, nil
	}

	return nil, nil
})

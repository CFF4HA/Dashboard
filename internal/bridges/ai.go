package bridges

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/CFF4HA/Dashboard/internal/core"
)

type ModelRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Format  string         `json:"format"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options"`
}

type View struct {
	View string
}

type AIRouter struct {
	Sessions map[string]View
	Lock     *sync.RWMutex
}

func (a AIRouter) Data(w http.ResponseWriter, r *http.Request) (any, error) {
	query := strings.Trim(r.FormValue("router-search_query"), " ")
	if query == "" {
		return nil, nil
	}

	// here we would want to contact the backend server, or the LLM server
	// itself to get a response from the query.
	prompt :=
		`Instructions: Act as an intent and entity extractor. Return JSON with keys "intent", "entities", and "summary". 
Intents: IngredientSearch (one chemical), ProductSearch (one brand), ProductCreation (comma list). 

Example 1: 
Input: Linalool 
Output: {"intent": "IngredientSearch", "entities": ["Linalool"], "summary": "Searching for the ingredient Linalool."}

Example 2:
Input: Cerave Cream
Output: {"intent": "ProductSearch", "entities": ["Cerave Cream"], "summary": "Looking up the product Cerave Cream."}

Example 3:
Input: Glycerin, Water
Output: {"intent": "ProductCreation", "entities": ["Glycerin", "Water"], "summary": "Creating a new product with Glycerin and Water."}

Input: %s`

	payload := ModelRequest{
		Model:  "gemma-text",
		Prompt: fmt.Sprintf(prompt, query),
		Format: "json",
		Stream: false,
		Options: map[string]any{
			"think": false,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/generate", core.LLMAddress), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		core.Logger.Error("LLM Server returned non-200 status code", "StatusCode", resp.StatusCode, "status", resp.Status)
		resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	var Response struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		return nil, err
	}

	var ResponseData struct {
		Intent   string   `json:"intent"`
		Entities []string `json:"entities"`
		Summary  string   `json:"summary"`
	}
	json.Unmarshal([]byte(Response.Response), &ResponseData)

	return ResponseData, nil
}

func (a AIRouter) Name() string {
	return "AI"
}

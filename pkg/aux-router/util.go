package auxrouter

import (
	"bytes"
	"context"
	"io"
	"strings"

	"encoding/json"
	"net/http"
)

type OllamaRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Format  string         `json:"format"`
	System  string         `json:"system"`
	Stream  bool           `json:"stream"`
	Think   bool           `json:"think"`
	Options map[string]any `json:"options"`
}

type OllamaResponse struct {
	Response json.RawMessage `json:"response"`
}

func (a *Aux) LLM(ctx context.Context, prompt string) (*OllamaResponse, error) {
	ollamaRequest := OllamaRequest{
		Model:  "gemma-text",
		Prompt: prompt,
		Format: "json",
		System: `You are a specialized intent classifier. You will be provided with a list of system intents and a user query. Each intent contains a unique key, a description, and a response schema. Your goal is to map the user query to the most relevant intent and return a JSON object. This object must contain three fields: "intent" for the key, "certainty" for the confidence score as a float, and "response" for a nested JSON object populated with data extracted from the query. You must ensure the "response" field is a native JSON object that follows the provided examples, a nested JSON object (not a string) populated with data extracted from the query according to the intent's schema. Do not include any conversational text, explanations, or markdown code blocks outside of the JSON output.`,
		Stream: false,
		Think:  false,
		Options: map[string]any{
			"num_ctx": 8192,
		},
	}

	body, err := json.Marshal(ollamaRequest)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(a.LlmServerEndpoint, "/") + "/api/generate"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ollama := &OllamaResponse{}
	if err := json.Unmarshal(data, ollama); err != nil {
		return nil, err
	}

	return ollama, nil
}

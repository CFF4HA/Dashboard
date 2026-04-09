package auxrouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	IntentKeyNoIntents   = "no_intents"
	IntentKeyNoLLMServer = "no_llm_server"
	IntentKeyNoQuery     = "no_query"
	IntentLLMError       = "llm_error"

	llmPromptTemplate = `Supported Intents:
		%s

		User Input:
		%s`
)

var ()

// the Aux struct is the main entry point for the LLM based UI
// router.
type Aux struct {
	// This is a user supplied list of "intents" (i.e.,
	// system supported views, actions).
	Intents []Intent `json:"intents"`

	// the http endpoint to use to contact the llm.
	LlmServerEndpoint string `json:"llm_server_endpoint"`
}

func (a Aux) Data(w http.ResponseWriter, r *http.Request, model map[string]any) (any, error) {
	var response Response
	if a.LlmServerEndpoint == "" {
		response.Certainty = 1
		response.Intent = IntentKeyNoLLMServer

		return response, nil
	} else if len(a.Intents) == 0 {
		response.Certainty = 1
		response.Intent = IntentKeyNoIntents

		return response, nil
	}

	query := r.FormValue("query")
	if query == "" {
		response.Certainty = 1
		response.Intent = IntentKeyNoQuery

		return response, nil
	}

	// we contact the llm using the specified endpoint and we
	// provide it with our relevant prompt(s)
	data, err := json.Marshal(a.Intents)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(strings.Trim(llmPromptTemplate, " "), string(data), query)

	// we send directly to the llm now
	ollama, err := a.LLM(prompt)
	if err != nil {
		response.Certainty = 1
		response.Intent = IntentLLMError

		return response, nil
	}

	var unescapedString string
	if err := json.Unmarshal(ollama.Response, &unescapedString); err != nil {
		return nil, fmt.Errorf("failed to unescape JSON string: %w", err)
	}

	if err := json.Unmarshal([]byte(unescapedString), &response); err != nil {

		response.Certainty = 1
		response.Intent = IntentLLMError
		return response, nil
	}

	return response, nil
}

func (a Aux) Name() string {
	return "Aux"
}

func (a *Aux) Intent(key string, desc string, examples ...IntentExample) {
	a.Intents = append(a.Intents, Intent{
		Key:         key,
		Description: desc,
		Example:     examples,
	})
}

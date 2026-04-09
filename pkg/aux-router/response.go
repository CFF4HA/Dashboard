package auxrouter

import (
	"encoding/json"
)

type Response struct {
	// the key value matching a given intent
	Intent string `json:"intent"`

	// the certainty with which the LLM believes the prompt
	// matches the intent
	Certainty float64 `json:"certainty"`

	// a JSON encoded string representing the expected result
	// schema
	Response json.RawMessage `json:"response"`
}

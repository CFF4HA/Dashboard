package auxrouter

type Intent struct {
	// This is a unique string that the LLM can output
	// to say that this was the intent.
	Key string `json:"key"`

	// This is a description describing to the LLM what
	// user actions are meant to apply to this intent.
	Description string `json:"description"`

	// This provides an example prompt and return spec.
	Example []IntentExample `json:"example"`
}

type IntentExample struct {
	Prompt         string `json:"prompt"`
	ResponseSchema any    `json:"response_schema"`
}

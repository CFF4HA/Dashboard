package bridges

// Bridge data keys are the string identifiers used to pass data between
// bridges and templates via the model map. Using named constants here
// prevents typos and makes it easy to search for all usages of a key.
const (
	KeyIngredient  = "Ingredient"
	KeyNames       = "Names"
	KeyMetadata    = "Metadata"
	KeyDatabaseErr = "DatabaseErr"
	KeyLLMErr      = "LLMErr"
	KeyPayload     = "Payload"
	KeyAux         = "Aux"
)

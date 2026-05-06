package types

// The following types are related to usability data for the site
// like data pulled from pubchem, sessions created, product counts,
// ingredient counts, label counts, etc.

type UsageMetric struct {
	Model
	ProductsCount    int64 `json:"products_count"`
	IngredientsCount int64 `json:"ingredients_count"`
	LabelsCount      int64 `json:"labels_count"`

	AutomatedScrapesCount int `json:"automated_scrapes_count"`
	PubChemLookupsCount   int `json:"pub_chem_lookups_count"`
	SessionsCount         int `json:"sessions_count"`

	TagsCount        int64 `json:"tags_count"`
	TaggingRuleCount int64 `json:"tagging_rule_count"`
}

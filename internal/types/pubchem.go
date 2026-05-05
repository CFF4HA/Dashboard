package types

type PubChemLabelConfig struct {
	GsonField string `json:"gson_field"`
	LabelType string `json:"label_type" gorm:"type:text;not null;default:'general';"`
}

package types

type Notification struct {
	Model

	BackgroundColor string `json:"background_color" gorm:"type:text;not null;default:'#000000';"`
	Content         string `json:"content" gorm:"type:text;not null;"`
	Enabled         bool   `json:"enabled" gorm:"type:boolean;not null;default:false;"`
}

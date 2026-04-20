package types

import (
	"github.com/google/uuid"
)

type User struct {
	Model
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"password_hash" gorm:"not null"`
}

type Session struct {
	Session uuid.UUID `json:"session"`
}

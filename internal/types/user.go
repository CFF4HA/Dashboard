package types

import (
	"github.com/google/uuid"
)

type User struct {
	Uid          uuid.UUID `json:"uid" gorm:"type:uuid;primaryKey;not null;default:gen_random_uuid();"`
	Username     string    `json:"username" gorm:"unique;not null"`
	PasswordHash string    `json:"password_hash" gorm:"not null"`
}

type Session struct {
	Session uuid.UUID `json:"session"`
}

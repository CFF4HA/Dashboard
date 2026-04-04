package types

import (
	"github.com/google/uuid"
)

type User struct {
	Uid          uuid.UUID `json:"uid" gorm:"type:uuid;primaryKey;not null"`
	Username     string    `json:"username" gorm:"unique;not null"`
	PasswordHash string    `json:"password_hash" gorm:"not null"`

	Products []Product `json:"products" gorm:"foreignKey:UserUid;constraint:OnDelete:CASCADE"`
}

type Session struct {
	Session uuid.UUID `json:"session"`
}

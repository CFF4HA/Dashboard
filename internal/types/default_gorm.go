package types

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id      uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Created time.Time  `json:"created" gorm:"type:timestamp;not null;default:current_timestamp"`
	Updated time.Time  `json:"updated" gorm:"type:timestamp;not null;default:current_timestamp"`
	Deleted *time.Time `json:"deleted" gorm:"type:timestamp;default:null"`
}

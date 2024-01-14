package models

import "github.com/google/uuid"

type Donation struct {
	BaseModel

	GUID    uuid.UUID `gorm:"uniqueIndex,size:36" json:"guid"`
	User    uint      `json:"-"`
	Issue   uint      `json:"-"`
	Status  string    `json:"status"`
	Amount  float32   `json:"amount"`
	Deleted bool      `json:"-"`
}

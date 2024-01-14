package models

import "github.com/google/uuid"

type Issue struct {
	BaseModel

	GUID    uuid.UUID `gorm:"uniqueIndex,size:36" json:"guid"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Pool    float32   `json:"pool"`
}

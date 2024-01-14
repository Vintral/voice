package models

import "github.com/google/uuid"

type Issue struct {
	BaseModel

	GUID    uuid.UUID `json:"guid"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Pool    float32   `json:"pool"`
}

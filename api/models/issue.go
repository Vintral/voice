package models

import "github.com/google/uuid"

type Issue struct {
	BaseModel

	GUID    uuid.UUID
	Title   string
	Summary string
	Pool    float32
}

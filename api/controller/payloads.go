package controller

import "github.com/google/uuid"

type donationPayload struct {
	Issue  uuid.UUID `json:"issue"`
	Amount float32   `json:"amount"`
}

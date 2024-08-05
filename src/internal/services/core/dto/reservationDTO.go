package dto

import "github.com/google/uuid"

type ReserveBookDTO struct {
	ReaderID uuid.UUID
	BookID   uuid.UUID
}

package intf

import (
	"BookSmart-services/core/models"
	"context"
	"github.com/google/uuid"
)

type IReservationService interface {
	Create(ctx context.Context, readerID, bookID uuid.UUID) error
	Update(ctx context.Context, reservation *models.ReservationModel) error
	GetAllReservationsByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error)
	GetByID(ctx context.Context, reservationID uuid.UUID) (*models.ReservationModel, error)
}

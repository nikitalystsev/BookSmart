package repositories

import (
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=IReservationRepo.go -destination=../tests/unitTests/mocks/mockIReservationRepo.go

type IReservationRepo interface {
	Create(ctx context.Context, reservation *models.ReservationModel) error
	GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error)
	GetByID(ctx context.Context, reservationID uuid.UUID) (*models.ReservationModel, error)
	Update(ctx context.Context, reservation *models.ReservationModel) error
	GetOverdueByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error)
	GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error)
}

package interfaces

import (
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

type IReservationService interface {
	Create(ctx context.Context, readerID, bookID uuid.UUID) error
	Update(ctx context.Context, reservation *models.ReservationModel) error
}

package interfaces

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

type IReaderService interface {
	Register(ctx context.Context, reader *models.ReaderModel) error
	Login(ctx context.Context, reader *dto.ReaderLoginDTO) error
	ReserveBook(ctx context.Context, readerID, bookID uuid.UUID) error
	ExtendBook(ctx context.Context, reservation models.ReservationModel) error
}

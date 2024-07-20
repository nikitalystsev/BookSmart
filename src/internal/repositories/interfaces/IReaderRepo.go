package interfaces

import (
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=IReaderRepo.go -destination=../tests/unitTests/mocks/mockIReaderRepo.go

type IReaderRepo interface {
	Create(ctx context.Context, reader *models.ReaderModel) error
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.ReaderModel, error)
}

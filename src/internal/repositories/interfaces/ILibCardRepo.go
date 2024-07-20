package interfaces

import (
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=ILibCardRepo.go -destination=../tests/unitTests/mocks/mockILibCardRepo.go

type ILibCardRepo interface {
	Create(ctx context.Context, libCard *models.LibCardModel) error
	GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error)
	GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error)
}

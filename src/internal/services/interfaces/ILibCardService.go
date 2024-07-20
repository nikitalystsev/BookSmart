package interfaces

import (
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

type ILibCardService interface {
	Create(ctx context.Context, readerID uuid.UUID) error
	Update(ctx context.Context, libCard *models.LibCardModel) error
	IsValidLibCard(libCard *models.LibCardModel) bool
}

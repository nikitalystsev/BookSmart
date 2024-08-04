package intf

import (
	"BookSmart-services/models"
	"context"
	"github.com/google/uuid"
)

type ILibCardService interface {
	Create(ctx context.Context, readerID uuid.UUID) error
	Update(ctx context.Context, libCard *models.LibCardModel) error
	GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error)
}

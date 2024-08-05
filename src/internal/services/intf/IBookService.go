package intf

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"context"
	"github.com/google/uuid"
)

type IBookService interface {
	Create(ctx context.Context, book *models.BookModel) error
	Delete(ctx context.Context, ID uuid.UUID) error
	GetByID(ctx context.Context, ID uuid.UUID) (*models.BookModel, error)
	GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error)
}

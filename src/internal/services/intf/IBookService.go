package intf

import (
	"BookSmart-services/dto"
	"BookSmart-services/models"
	"context"
	"github.com/google/uuid"
)

type IBookService interface {
	Create(ctx context.Context, book *models.BookModel) error
	Delete(ctx context.Context, bookID uuid.UUID) error
	GetByID(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error)
	GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error)
}

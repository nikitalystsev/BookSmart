package intfRepo

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=IBookRepo.go -destination=../../tests/unitTests/serviceTests/mocks/mockBookRepo.go --package=mocks

type IBookRepo interface {
	Create(ctx context.Context, book *models.BookModel) error
	GetByID(ctx context.Context, ID uuid.UUID) (*models.BookModel, error)
	GetByTitle(ctx context.Context, title string) (*models.BookModel, error)
	Delete(ctx context.Context, ID uuid.UUID) error
	Update(ctx context.Context, book *models.BookModel) error
	GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error)
}
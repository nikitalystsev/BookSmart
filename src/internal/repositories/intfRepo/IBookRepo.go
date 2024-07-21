package intfRepo

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=IBookRepo.go -destination=../../tests/unitTests/mocks/mockIBookRepo.go

type IBookRepo interface {
	Create(ctx context.Context, book *models.BookModel) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookModel, error)
	GetByTitle(ctx context.Context, title string) (*models.BookModel, error)
	DeleteByTitle(ctx context.Context, title string) error
	Update(ctx context.Context, book *models.BookModel) error
	GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error)
}

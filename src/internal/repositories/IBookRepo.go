package repositories

import (
	"BookSmart/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
)

// ErrNotFound my repository errors
var ErrNotFound = errors.New("[-] ERROR! Book was not found")

//go:generate mockgen -source=IBookRepo.go -destination=../tests/unitTests/mocks/mockIBookRepo.go

type IBookRepo interface {
	Create(ctx context.Context, book *models.BookModel) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookModel, error)
	GetByTitle(ctx context.Context, title string) (*models.BookModel, error)
	DeleteByTitle(ctx context.Context, title string) error
	Update(ctx context.Context, book *models.BookModel) error
}

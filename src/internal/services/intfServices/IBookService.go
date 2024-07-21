package intfServices

import (
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

type IBookService interface {
	Create(ctx context.Context, book *models.BookModel) error
	Delete(ctx context.Context, book *models.BookModel) error
	GetByID(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error)
}

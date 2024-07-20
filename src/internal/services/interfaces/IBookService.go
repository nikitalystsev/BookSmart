package interfaces

import (
	"BookSmart/internal/models"
	"context"
)

type IBookService interface {
	Create(ctx context.Context, book *models.BookModel) error
	Delete(ctx context.Context, book *models.BookModel) error
}

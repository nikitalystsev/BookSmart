package services

import (
	"BookSmart/internal/models"
	"context"
)

type IBookService interface {
	Create(ctx context.Context, book *models.BookModel) error
	DeleteByTitle(ctx context.Context, title string) error
}

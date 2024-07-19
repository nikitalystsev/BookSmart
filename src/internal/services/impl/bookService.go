package impl

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories"
	"context"
	"errors"
	"fmt"
)

const (
	BookRarityCommon = "Обычная"
	BookRarityRare   = "Редкая"
	BookRarityUnique = "Уникальная"
)

type BookService struct {
	bookRepo repositories.IBookRepo
}

func NewBookService(bookRepo repositories.IBookRepo) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (bs *BookService) Create(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByTitle(ctx, book.Title)
	if err != nil && !errors.Is(err, errors.New("[!] ERROR! Object not found")) {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook != nil {
		return errors.New("[!] ERROR! Book with this title already exists")
	}

	err = bs.bookRepo.Create(ctx, book)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating book: %v", err)
	}

	return nil
}

func (bs *BookService) DeleteByTitle(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByTitle(ctx, book.Title)
	if err != nil && !errors.Is(err, errors.New("[!] ERROR! Object not found")) {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook == nil {
		return errors.New("[!] ERROR! Book with this title does not exist")
	}

	err = bs.bookRepo.DeleteByTitle(ctx, book.Title)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error deleting book: %v", err)
	}

	return nil
}

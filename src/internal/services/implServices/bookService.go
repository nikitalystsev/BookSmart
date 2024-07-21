package implServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const (
	BookRarityCommon = "Обычная"
	BookRarityRare   = "Редкая"
	BookRarityUnique = "Уникальная"
)

type BookService struct {
	bookRepo intfRepo.IBookRepo
}

func NewBookService(bookRepo intfRepo.IBookRepo) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (bs *BookService) Create(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByTitle(ctx, book.Title)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
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

func (bs *BookService) Delete(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByTitle(ctx, book.Title)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
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

func (bs *BookService) GetByID(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error) {
	book, err := bs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("[!] ERROR! Error retrieving book information: %v", err)
	}

	if book == nil {
		return nil, errors.New("[!] ERROR! Book with this title does not exist")
	}

	return book, nil
}

func (bs *BookService) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	books, err := bs.bookRepo.GetByParams(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("[!] ERROR! Error searching for books: %v", err)
	}

	return books, nil
}

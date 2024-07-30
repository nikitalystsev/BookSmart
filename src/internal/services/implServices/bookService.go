package implServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/internal/services/intfServices"
	"BookSmart/pkg/logging"
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
	logger   logging.Logger
}

func NewBookService(bookRepo intfRepo.IBookRepo, logger logging.Logger) intfServices.IBookService {
	return &BookService{bookRepo: bookRepo, logger: logger}
}

func (bs *BookService) Create(ctx context.Context, book *models.BookModel) error {
	err := bs.baseValidation(ctx, book)
	if err != nil {
		return err
	}

	err = bs.bookRepo.Create(ctx, book)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating book: %v", err)
	}

	return nil
}

func (bs *BookService) Delete(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByID(ctx, book.ID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook == nil {
		return errors.New("[!] ERROR! Book with this title does not exist")
	}

	err = bs.bookRepo.Delete(ctx, book.ID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error deleting book: %v", err)
	}

	return nil
}

func (bs *BookService) GetByID(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error) {
	book, err := bs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		return nil, fmt.Errorf("[!] ERROR! Error retrieving book information: %v", err)
	}

	if book == nil {
		return nil, errors.New("[!] ERROR! Book with this ID does not exist")
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

func (bs *BookService) baseValidation(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByID(ctx, book.ID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook != nil {
		return errors.New("[!] ERROR! Book with this title already exists")
	}

	if book.Title == "" {
		return errors.New("[!] ERROR! Empty book title")
	}

	if book.Author == "" {
		return errors.New("[!] ERROR! Empty book author")
	}

	if book.Rarity == "" {
		return errors.New("[!] ERROR! Empty book rarity")
	}

	if book.CopiesNumber <= 0 {
		return errors.New("[!] ERROR! Invalid book copies number")
	}

	return nil
}

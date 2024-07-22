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
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

const (
	MaxBooksPerReader    = 5
	ReaderPhoneNumberLen = 11
)

type ReaderService struct {
	readerRepo intfRepo.IReaderRepo
	bookRepo   intfRepo.IBookRepo
}

func NewReaderService(
	readerRepo intfRepo.IReaderRepo,
	bookRepo intfRepo.IBookRepo,
) *ReaderService {
	return &ReaderService{
		readerRepo: readerRepo,
		bookRepo:   bookRepo,
	}
}

func (rs *ReaderService) Register(ctx context.Context, reader *models.ReaderModel) error {
	err := rs.baseValidation(ctx, reader)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reader.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error hashing password: %v", err)
	}

	reader.Password = string(hashedPassword)

	err = rs.readerRepo.Create(ctx, reader)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating reader: %v", err)
	}

	return nil
}

func (rs *ReaderService) Login(ctx context.Context, reader *dto.ReaderLoginDTO) error {
	exitingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}

	if exitingReader == nil {
		return fmt.Errorf("[!] ERROR! Reader with this phoneNumbers does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(exitingReader.Password), []byte(reader.Password))
	if err != nil {
		return fmt.Errorf("[!] ERROR! Wrong password")
	}

	return nil
}

func (rs *ReaderService) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}
	if existingReader == nil {
		return fmt.Errorf("[!] ERROR! Reader with this ID does not exist")
	}

	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}
	if existingBook == nil {
		return fmt.Errorf("[!] ERROR! Book with this ID does not exist")
	}

	isFavorite, err := rs.readerRepo.IsFavorite(ctx, readerID, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking if book is already a favorite: %v", err)
	}
	if isFavorite {
		return fmt.Errorf("[!] ERROR! Book is already in favorites")
	}
	err = rs.readerRepo.AddToFavorites(ctx, readerID, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error adding book to favorites: %v", err)
	}

	return nil
}

func (rs *ReaderService) baseValidation(ctx context.Context, reader *models.ReaderModel) error {
	existingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}

	if existingReader != nil {
		return errors.New("[!] ERROR! Reader with this phoneNumbers already exists")
	}

	if reader.Fio == "" {
		return errors.New("[!] ERROR! Field Fio is required")
	}
	if reader.PhoneNumber == "" {
		return errors.New("[!] ERROR! Field PhoneNumber is required")
	}

	if reader.Age <= 0 {
		return errors.New("[!] ERROR! Field Age is required")
	}

	if len(reader.PhoneNumber) != ReaderPhoneNumberLen {
		return errors.New("[!] ERROR! Reader phoneNumbers len")
	}

	_, err = strconv.Atoi(reader.PhoneNumber)
	if err != nil {
		return errors.New("[!] ERROR! Reader phoneNumbers incorrect format")
	}

	return nil
}

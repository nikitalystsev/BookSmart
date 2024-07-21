package implServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const MaxBooksPerReader = 5

type ReaderService struct {
	readerRepo intfRepo.IReaderRepo
}

func CreateNewReaderService(
	readerRepo intfRepo.IReaderRepo,
) *ReaderService {
	return &ReaderService{
		readerRepo: readerRepo,
	}
}

func (rs *ReaderService) Register(ctx context.Context, reader *models.ReaderModel) error {
	existingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}

	if existingReader != nil {
		return errors.New("[!] ERROR! Reader with this phoneNumbers already exists")
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

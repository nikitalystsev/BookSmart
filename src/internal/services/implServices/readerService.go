package implServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/internal/services/intfServices"
	"BookSmart/pkg/auth"
	hash2 "BookSmart/pkg/hash"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

const (
	MaxBooksPerReader    = 5
	ReaderPhoneNumberLen = 11
)

var (
	accessTokenTTL  = time.Hour * 2       // В минутах
	refreshTokenTTL = time.Hour * 24 * 30 // В минутах (30 дней)
)

type ReaderService struct {
	readerRepo   intfRepo.IReaderRepo
	bookRepo     intfRepo.IBookRepo
	tokenManager auth.ITokenManager
	hasher       hash2.IPasswordHasher
}

func NewReaderService(
	readerRepo intfRepo.IReaderRepo,
	bookRepo intfRepo.IBookRepo,
	tokenManager auth.ITokenManager,
	hasher hash2.IPasswordHasher,
) intfServices.IReaderService {
	return &ReaderService{
		readerRepo:   readerRepo,
		bookRepo:     bookRepo,
		tokenManager: tokenManager,
		hasher:       hasher,
	}
}

// SignUp Зарегистрироваться
func (rs *ReaderService) SignUp(ctx context.Context, reader *models.ReaderModel) error {
	err := rs.baseValidation(ctx, reader)
	if err != nil {
		return err
	}

	hashedPassword, err := rs.hasher.Hash(reader.Password)
	if err != nil {
		return err
	}

	reader.Password = hashedPassword

	err = rs.readerRepo.Create(ctx, reader)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating reader: %v", err)
	}

	return nil
}

// SignIn Войти
func (rs *ReaderService) SignIn(ctx context.Context, reader *dto.ReaderLoginDTO) (intfServices.Tokens, error) {
	exitingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return intfServices.Tokens{}, fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}

	if exitingReader == nil {
		return intfServices.Tokens{}, fmt.Errorf("[!] ERROR! Reader with this phoneNumbers does not exist")
	}

	err = rs.hasher.Compare(exitingReader.Password, reader.Password)
	if err != nil {
		return intfServices.Tokens{}, err
	}

	return rs.createTokens(ctx, exitingReader.ID)
}

func (rs *ReaderService) RefreshTokens(ctx context.Context, refreshToken string) (intfServices.Tokens, error) {
	existingReader, err := rs.readerRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return intfServices.Tokens{}, err
	}

	return rs.createTokens(ctx, existingReader.ID)
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

func (rs *ReaderService) createTokens(ctx context.Context, readerID uuid.UUID) (intfServices.Tokens, error) {
	var (
		res intfServices.Tokens
		err error
	)

	res.AccessToken, err = rs.tokenManager.NewJWT(readerID, accessTokenTTL)
	if err != nil {
		return res, fmt.Errorf("[!] ERROR! Error generating access token: %v", err)
	}

	res.RefreshToken, err = rs.tokenManager.NewRefreshToken()
	if err != nil {
		return res, fmt.Errorf("[!] ERROR! Error generating refresh token: %v", err)
	}

	err = rs.readerRepo.SaveRefreshToken(ctx, readerID, res.RefreshToken, refreshTokenTTL)
	if err != nil {
		return res, fmt.Errorf("[!] ERROR! Error saving refresh token: %v", err)
	}

	return res, nil
}

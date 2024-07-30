package implServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/internal/services/errsService"
	"BookSmart/internal/services/intfServices"
	"BookSmart/pkg/auth"
	hash2 "BookSmart/pkg/hash"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logger       *logrus.Entry
}

func NewReaderService(
	readerRepo intfRepo.IReaderRepo,
	bookRepo intfRepo.IBookRepo,
	tokenManager auth.ITokenManager,
	hasher hash2.IPasswordHasher,
	logger *logrus.Entry,
) intfServices.IReaderService {
	return &ReaderService{
		readerRepo:   readerRepo,
		bookRepo:     bookRepo,
		tokenManager: tokenManager,
		hasher:       hasher,
		logger:       logger,
	}
}

// SignUp Зарегистрироваться
func (rs *ReaderService) SignUp(ctx context.Context, reader *models.ReaderModel) error {
	rs.logger.Info("starting sign up process")

	err := rs.baseValidation(ctx, reader)
	if err != nil {
		rs.logger.Errorf("reader validation failed: %v", err)
		return err
	}

	rs.logger.Info("hashing password")

	hashedPassword, err := rs.hasher.Hash(reader.Password)
	if err != nil {
		rs.logger.Errorf("hashing failed: %v", err)
		return err
	}

	reader.Password = hashedPassword

	rs.logger.Infof("creating reader in repository: %+v", reader)

	err = rs.readerRepo.Create(ctx, reader)
	if err != nil {
		rs.logger.Errorf("error creating reader: %v", err)
		return err
	}

	rs.logger.Info("book creation successful")

	return nil
}

// SignIn Войти
func (rs *ReaderService) SignIn(ctx context.Context, reader *dto.ReaderLoginDTO) (intfServices.Tokens, error) {
	rs.logger.Infof("attempting sign in with phoneNumber: %s", reader.PhoneNumber)

	exitingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking reader existence: %v", err)
		return intfServices.Tokens{}, err
	}

	if exitingReader == nil {
		rs.logger.Warn("reader with this phoneNumber does not exist")
		return intfServices.Tokens{}, errsService.ErrReaderDoesNotExists
	}

	rs.logger.Info("compare password with hashing password")

	err = rs.hasher.Compare(exitingReader.Password, reader.Password)
	if err != nil {
		rs.logger.Errorf("compare password with hashing password failed: %v", err)
		return intfServices.Tokens{}, err
	}

	return rs.createTokens(ctx, exitingReader.ID)
}

func (rs *ReaderService) RefreshTokens(ctx context.Context, refreshToken string) (intfServices.Tokens, error) {
	rs.logger.Info("attempting refresh tokens")

	existingReader, err := rs.readerRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		rs.logger.Errorf("error checking reader existence: %v", err)
		return intfServices.Tokens{}, err
	}

	return rs.createTokens(ctx, existingReader.ID)
}

func (rs *ReaderService) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	rs.logger.Info("attempting to add book to favorites")

	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil {
		rs.logger.Errorf("error checking reader existence: %v", err)
		return err
	}
	if existingReader == nil {
		rs.logger.Warn("reader with this ID does not exist")
		return errsService.ErrReaderDoesNotExists
	}

	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		rs.logger.Errorf("error checking book existence: %v", err)
		return err
	}
	if existingBook == nil {
		rs.logger.Warn("book with this ID does not exist")
		return errsService.ErrBookDoesNotExists
	}

	isFavorite, err := rs.readerRepo.IsFavorite(ctx, readerID, bookID)
	if err != nil {
		rs.logger.Errorf("error checking favorite book: %v", err)
		return err
	}
	if isFavorite {
		rs.logger.Errorf("book with this ID already is a favorite")
		return errsService.ErrBookAlreadyIsFavorite
	}

	err = rs.readerRepo.AddToFavorites(ctx, readerID, bookID)
	if err != nil {
		rs.logger.Errorf("error adding book to favorites: %v", err)
		return err
	}

	rs.logger.Infof("book successfully added in favorites")

	return nil
}

func (rs *ReaderService) baseValidation(ctx context.Context, reader *models.ReaderModel) error {
	existingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)

	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking reader existence: %v", err)
		return err
	}

	if existingReader != nil {
		rs.logger.Warn("reader with this phoneNumbers already exists")
		return errsService.ErrReaderAlreadyExist
	}

	if reader.Fio == "" {
		rs.logger.Warn("empty reader fio")
		return errsService.ErrEmptyReaderFio
	}

	if reader.PhoneNumber == "" {
		rs.logger.Warn("empty reader phoneNumber")
		return errsService.ErrEmptyReaderPhoneNumber
	}

	if reader.Age <= 0 {
		rs.logger.Warn("invalid reader age")
		return errsService.ErrInvalidReaderAge
	}

	if len(reader.PhoneNumber) != ReaderPhoneNumberLen {
		rs.logger.Warn("invalid reader phoneNumber len")
		return errsService.ErrInvalidReaderPhoneNumberLen
	}

	_, err = strconv.Atoi(reader.PhoneNumber)
	if err != nil {
		rs.logger.Warn("invalid reader phoneNumber format")
		return errsService.ErrInvalidReaderPhoneNumberFormat
	}

	rs.logger.Info("reader validation successful")

	return nil
}

func (rs *ReaderService) createTokens(ctx context.Context, readerID uuid.UUID) (intfServices.Tokens, error) {
	rs.logger.Info("attempting to create Tokens")

	var (
		res intfServices.Tokens
		err error
	)

	rs.logger.Info("generate access token")

	res.AccessToken, err = rs.tokenManager.NewJWT(readerID, accessTokenTTL)
	if err != nil {
		rs.logger.Errorf("error generating access token: %v", err)
		return res, err
	}

	rs.logger.Info("generate refresh token")

	res.RefreshToken, err = rs.tokenManager.NewRefreshToken()
	if err != nil {
		rs.logger.Errorf("error generating refresh token: %v", err)
		return res, err
	}

	rs.logger.Info("save refresh token")

	err = rs.readerRepo.SaveRefreshToken(ctx, readerID, res.RefreshToken, refreshTokenTTL)
	if err != nil {
		rs.logger.Errorf("Error saving refresh token: %v", err)
		return res, err
	}

	rs.logger.Info("successfully created tokens")

	return res, nil
}

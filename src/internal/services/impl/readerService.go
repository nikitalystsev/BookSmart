package impl

import (
	errsRepo "BookSmart-repositories/errs"
	intfRepo "BookSmart-repositories/intf"
	"BookSmart-services/dto"
	"BookSmart-services/errs"
	"BookSmart-services/intf"
	"BookSmart-services/models"
	"BookSmart-services/pkg/auth"
	"BookSmart-services/pkg/hash"
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
	ReaderPasswordLen    = 10

	ReaderRole = "Reader"
)

type ReaderService struct {
	readerRepo      intfRepo.IReaderRepo
	bookRepo        intfRepo.IBookRepo
	tokenManager    auth.ITokenManager
	hasher          hash.IPasswordHasher
	logger          *logrus.Entry
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewReaderService(
	readerRepo intfRepo.IReaderRepo,
	bookRepo intfRepo.IBookRepo,
	tokenManager auth.ITokenManager,
	hasher hash.IPasswordHasher,
	logger *logrus.Entry,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) intf.IReaderService {
	return &ReaderService{
		readerRepo:      readerRepo,
		bookRepo:        bookRepo,
		tokenManager:    tokenManager,
		hasher:          hasher,
		logger:          logger,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// SignUp Зарегистрироваться
func (rs *ReaderService) SignUp(ctx context.Context, reader *models.ReaderModel) error {
	if reader == nil {
		rs.logger.Warn("reader object is nil")
		return errs.ErrReaderObjectIsNil
	}

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

	reader.Role = ReaderRole
	reader.Password = hashedPassword

	rs.logger.Infof("creating reader in repository: %+v", reader)

	err = rs.readerRepo.Create(ctx, reader)
	if err != nil {
		rs.logger.Errorf("error creating reader: %v", err)
		return err
	}

	rs.logger.Info("reader successful creation")

	return nil
}

// SignIn Войти
func (rs *ReaderService) SignIn(ctx context.Context, reader *dto.ReaderSignInDTO) (intf.Tokens, error) {
	if reader == nil {
		rs.logger.Warn("reader object is nil")
		return intf.Tokens{}, errs.ErrReaderObjectIsNil
	}

	rs.logger.Infof("attempting sign in with phoneNumber: %s", reader.PhoneNumber)

	exitingReader, err := rs.readerRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking reader existence: %v", err)
		return intf.Tokens{}, err
	}

	if exitingReader == nil {
		rs.logger.Warn("reader with this phoneNumber does not exist")
		return intf.Tokens{}, errs.ErrReaderDoesNotExists
	}

	rs.logger.Info("compare password with hashing password")

	err = rs.hasher.Compare(exitingReader.Password, reader.Password)
	if err != nil {
		rs.logger.Errorf("compare password with hashing password failed: %v", err)
		return intf.Tokens{}, err
	}

	return rs.createTokens(ctx, exitingReader.ID, exitingReader.Role)
}

func (rs *ReaderService) RefreshTokens(ctx context.Context, refreshToken string) (intf.Tokens, error) {
	rs.logger.Info("attempting refresh tokens")

	existingReader, err := rs.readerRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		rs.logger.Errorf("error checking reader existence: %v", err)
		return intf.Tokens{}, err
	}

	return rs.createTokens(ctx, existingReader.ID, existingReader.Role)
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
		return errs.ErrReaderDoesNotExists
	}

	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		rs.logger.Errorf("error checking book existence: %v", err)
		return err
	}
	if existingBook == nil {
		rs.logger.Warn("book with this ID does not exist")
		return errs.ErrBookDoesNotExists
	}

	isFavorite, err := rs.readerRepo.IsFavorite(ctx, readerID, bookID)
	if err != nil {
		rs.logger.Errorf("error checking favorite book: %v", err)
		return err
	}
	if isFavorite {
		rs.logger.Errorf("book with this ID already is a favorite")
		return errs.ErrBookAlreadyIsFavorite
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
		return errs.ErrReaderAlreadyExist
	}

	if reader.Fio == "" {
		rs.logger.Warn("empty reader fio")
		return errs.ErrEmptyReaderFio
	}

	if reader.PhoneNumber == "" {
		rs.logger.Warn("empty reader phoneNumber")
		return errs.ErrEmptyReaderPhoneNumber
	}

	if reader.Password == "" {
		rs.logger.Warn("empty reader password")
		return errs.ErrEmptyReaderPassword
	}

	if len(reader.Password) != ReaderPasswordLen {
		rs.logger.Warn("invalid reader password len")
		return errs.ErrInvalidReaderPasswordLen
	}

	if reader.Age <= 0 {
		rs.logger.Warn("invalid reader age")
		return errs.ErrInvalidReaderAge
	}

	if len(reader.PhoneNumber) != ReaderPhoneNumberLen {
		rs.logger.Warn("invalid reader phoneNumber len")
		return errs.ErrInvalidReaderPhoneNumberLen
	}

	_, err = strconv.Atoi(reader.PhoneNumber)
	if err != nil {
		rs.logger.Warn("invalid reader phoneNumber format")
		return errs.ErrInvalidReaderPhoneNumberFormat
	}

	rs.logger.Info("reader validation successful")

	return nil
}

func (rs *ReaderService) createTokens(ctx context.Context, readerID uuid.UUID, readerRole string) (intf.Tokens, error) {
	rs.logger.Info("attempting to create Tokens")

	var (
		res intf.Tokens
		err error
	)

	rs.logger.Info("generate access token")

	res.AccessToken, err = rs.tokenManager.NewJWT(readerID, readerRole, rs.accessTokenTTL)
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

	err = rs.readerRepo.SaveRefreshToken(ctx, readerID, res.RefreshToken, rs.refreshTokenTTL)
	if err != nil {
		rs.logger.Errorf("Error saving refresh token: %v", err)
		return res, err
	}

	rs.logger.Info("successfully created tokens")

	return res, nil
}
package serviceTests

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"BookSmart-services/impl"
	mockrepo "Booksmart/internal/tests/unitTests/serviceTests/mocks"
	"Booksmart/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReaderService_SignUp(t *testing.T) {
	type args struct {
		reader *models.ReaderModel
	}
	type mockBehaviour func(m *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args)
	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success successful signUp",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password34",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
				h.EXPECT().Hash(args.reader.Password).Return(gomock.Any().String(), nil)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error checking reader existence",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, fmt.Errorf("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error reader already exists",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(args.reader, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReaderAlreadyExist
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error missing fio",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrEmptyReaderFio
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error missing phoneNumber",
			args: args{
				&models.ReaderModel{
					ID:       uuid.New(),
					Fio:      "John Doe",
					Password: "password",
					Age:      25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrEmptyReaderPhoneNumber
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error invalid age",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password78",
					Age:         0,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrInvalidReaderAge
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error invalid phoneNumber len",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "1234567890",
					Password:    "password89",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrInvalidReaderPhoneNumberLen
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error invalid phoneNumber format",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "1234567890a",
					Password:    "password56",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrInvalidReaderPhoneNumberFormat
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error invalid password len",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678903",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrInvalidReaderPasswordLen
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error creating reader",
			args: args{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password54",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
				h.EXPECT().Hash(args.reader.Password).Return(gomock.Any().String(), nil)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("database error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
			readerService := impl.NewReaderService(
				mockReaderRepo, nil, nil,
				mockHasher, logging.GetLoggerForTests(),
				1, 2,
			)

			testCase.mockBehaviour(mockReaderRepo, mockHasher, testCase.args)

			err := readerService.SignUp(context.Background(), testCase.args.reader)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_SignIn(t *testing.T) {
	type args struct {
		readerDTO *dto.ReaderSignInDTO
		reader    *models.ReaderModel
	}
	type mockBehaviour func(
		r *mockrepo.MockIReaderRepo,
		h *mockrepo.MockIPasswordHasher,
		t *mockrepo.MockITokenManager, args args,
	)
	type expectedFunc func(t *testing.T, err error)

	var (
		accessTokenTTL  = time.Hour * 2       // В минутах
		refreshTokenTTL = time.Hour * 24 * 30 // В минутах (30 дней)
	)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success successful signIn",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashed_password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.readerDTO.PhoneNumber).Return(args.reader, nil)
				h.EXPECT().Compare(args.reader.Password, args.readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(args.reader.ID, args.reader.Role, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), args.reader.ID, "refreshToken", refreshTokenTTL).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error reader not found",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.readerDTO.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReaderDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error wrong password",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "wrong_password",
				},
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashed_password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.readerDTO.PhoneNumber).Return(args.reader, nil)
				h.EXPECT().Compare(args.reader.Password, args.readerDTO.Password).Return(fmt.Errorf("[!] ERROR! Wrong password"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Wrong password")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error access token generation error",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashed_password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.readerDTO.PhoneNumber).Return(args.reader, nil)
				h.EXPECT().Compare(args.reader.Password, args.readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(args.reader.ID, args.reader.Role, accessTokenTTL).Return("accessToken", fmt.Errorf("some error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("some error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error refresh token generation error",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashed_password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.readerDTO.PhoneNumber).Return(args.reader, nil)
				h.EXPECT().Compare(args.reader.Password, args.readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(args.reader.ID, args.reader.Role, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", fmt.Errorf("some error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("some error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: save refresh token error",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashed_password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), args.readerDTO.PhoneNumber).Return(args.reader, nil)
				h.EXPECT().Compare(args.reader.Password, args.readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(args.reader.ID, args.reader.Role, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), args.reader.ID, "refreshToken", refreshTokenTTL).Return(fmt.Errorf("some error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("some error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
			mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
			readerService := impl.NewReaderService(
				mockReaderRepo, nil, mockTokenManager,
				mockHasher, logging.GetLoggerForTests(),
				accessTokenTTL, refreshTokenTTL,
			)

			testCase.mockBehaviour(mockReaderRepo, mockHasher, mockTokenManager, testCase.args)

			_, err := readerService.SignIn(context.Background(), testCase.args.readerDTO)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_GetByPhoneNumber(t *testing.T) {
	type args struct {
		existingReader *models.ReaderModel
	}
	type mockBehaviour func(
		m *mockrepo.MockIReaderRepo,
		args args,
	)
	type expectedFunc func(t *testing.T, err error)

	const (
		accessTokenTTL  = time.Hour * 2
		refreshTokenTTL = time.Hour * 24 * 30
	)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success successful get by phone number",
			args: args{
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(m *mockrepo.MockIReaderRepo, args args) {
				m.EXPECT().GetByPhoneNumber(gomock.Any(), args.existingReader.PhoneNumber).Return(args.existingReader, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error reader does not exist",
			args: args{
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(m *mockrepo.MockIReaderRepo, args args) {
				m.EXPECT().GetByPhoneNumber(gomock.Any(), args.existingReader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReaderDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error in database",
			args: args{
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(m *mockrepo.MockIReaderRepo, args args) {
				m.EXPECT().GetByPhoneNumber(gomock.Any(), args.existingReader.PhoneNumber).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			readerService := impl.NewReaderService(
				mockReaderRepo, nil, nil,
				nil, logging.GetLoggerForTests(),
				accessTokenTTL, refreshTokenTTL,
			)

			testCase.mockBehaviour(mockReaderRepo, testCase.args)

			_, err := readerService.GetByPhoneNumber(context.Background(), testCase.args.existingReader.PhoneNumber)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_RefreshTokens(t *testing.T) {
	type args struct {
		refreshToken   string
		existingReader *models.ReaderModel
	}
	type mockBehaviour func(
		m *mockrepo.MockIReaderRepo,
		t *mockrepo.MockITokenManager,
		args args,
	)
	type expectedFunc func(t *testing.T, err error)

	const (
		accessTokenTTL  = time.Hour * 2
		refreshTokenTTL = time.Hour * 24 * 30
	)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success successful token refresh",
			args: args{
				refreshToken: "validRefreshToken",
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), args.refreshToken).Return(args.existingReader, nil)
				t.EXPECT().NewJWT(args.existingReader.ID, args.existingReader.Role, accessTokenTTL).Return("newAccessToken", nil)
				t.EXPECT().NewRefreshToken().Return("newRefreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), args.existingReader.ID, "newRefreshToken", refreshTokenTTL).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error failed to get reader by refresh token",
			args: args{
				refreshToken:   "invalidRefreshToken",
				existingReader: nil,
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), args.refreshToken).Return(nil, errs.ErrReaderDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReaderDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error generating access token",
			args: args{
				refreshToken: "validRefreshToken",
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, t *mockrepo.MockITokenManager, args args) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), args.refreshToken).Return(args.existingReader, nil)
				t.EXPECT().NewJWT(args.existingReader.ID, args.existingReader.Role, accessTokenTTL).Return("", fmt.Errorf("error generating JWT"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("error generating JWT")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
			readerService := impl.NewReaderService(
				mockReaderRepo, nil, mockTokenManager,
				nil, logging.GetLoggerForTests(),
				accessTokenTTL, refreshTokenTTL,
			)

			testCase.mockBehaviour(mockReaderRepo, mockTokenManager, testCase.args)

			_, err := readerService.RefreshTokens(context.Background(), testCase.args.refreshToken)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_AddToFavorites(t *testing.T) {
	type args struct {
		readerID uuid.UUID
		bookID   uuid.UUID
	}
	type mockBehaviour func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, args args)

	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success successful AddToFavorites",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), args.readerID, args.bookID).Return(false, nil)
				r.EXPECT().AddToFavorites(gomock.Any(), args.readerID, args.bookID).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error reader not found",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReaderDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book not found",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error already in favorites",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), args.readerID, args.bookID).Return(true, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookAlreadyIsFavorite
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error adding to favorites",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), args.readerID, args.bookID).Return(false, nil)
				r.EXPECT().AddToFavorites(gomock.Any(), args.readerID, args.bookID).Return(fmt.Errorf("error add"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("error add")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			readerService := impl.NewReaderService(
				mockReaderRepo, mockBookRepo, nil,
				nil, logging.GetLoggerForTests(),
				1, 2,
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, testCase.args)

			err := readerService.AddToFavorites(context.Background(), testCase.args.readerID, testCase.args.bookID)

			testCase.expected(t, err)
		})
	}
}

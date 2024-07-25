package unitTests

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/services/implServices"
	mockrepo "BookSmart/internal/tests/unitTests/mocks"
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

	type mockBehaviour func(m *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel)

	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		reader *models.ReaderModel
	}

	testTable := []struct {
		name          string
		input         inputStruct
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success: successful SignUp",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
				h.EXPECT().Hash(reader.Password).Return(gomock.Any().String(), nil)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: error checking reader existence",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, fmt.Errorf("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error checking reader existence: database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: reader already exists",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(reader, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader with this phoneNumbers already exists")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: missing fio",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Field Fio is required")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: missing PhoneNumber",
			input: inputStruct{
				&models.ReaderModel{
					ID:       uuid.New(),
					Fio:      "John Doe",
					Password: "password",
					Age:      25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Field PhoneNumber is required")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: invalid Age",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         0,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Field Age is required")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: invalid PhoneNumber length",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "1234567890",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader phoneNumbers len")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: invalid PhoneNumber format",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "1234567890a",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader phoneNumbers incorrect format")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error creating reader",
			input: inputStruct{
				&models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "password",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrNotFound)
				h.EXPECT().Hash(reader.Password).Return(gomock.Any().String(), nil)
				r.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error creating reader: database error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
			readerService := implServices.NewReaderService(mockReaderRepo, nil, nil, mockHasher)

			testCase.mockBehaviour(mockReaderRepo, mockHasher, testCase.input.reader)

			err := readerService.SignUp(context.Background(), testCase.input.reader)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_SignIn(t *testing.T) {
	type mockBehaviour func(
		r *mockrepo.MockIReaderRepo,
		h *mockrepo.MockIPasswordHasher,
		t *mockrepo.MockITokenManager,
		readerDTO *dto.ReaderLoginDTO,
		reader *models.ReaderModel,
	)

	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		readerDTO *dto.ReaderLoginDTO
		reader    *models.ReaderModel
	}

	var (
		accessTokenTTL  = time.Hour * 2       // В минутах
		refreshTokenTTL = time.Hour * 24 * 30 // В минутах (30 дней)
	)

	testTable := []struct {
		name          string
		input         inputStruct
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success: successful SignIn",
			input: inputStruct{
				readerDTO: &dto.ReaderLoginDTO{
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, readerDTO *dto.ReaderLoginDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), reader.ID, "refreshToken", refreshTokenTTL).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: reader not found",
			input: inputStruct{
				readerDTO: &dto.ReaderLoginDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, readerDTO *dto.ReaderLoginDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Reader with this phoneNumbers does not exist")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: wrong password",
			input: inputStruct{
				readerDTO: &dto.ReaderLoginDTO{
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, readerDTO *dto.ReaderLoginDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(fmt.Errorf("[!] ERROR! Wrong password"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Wrong password")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: access token generation error",
			input: inputStruct{
				readerDTO: &dto.ReaderLoginDTO{
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, readerDTO *dto.ReaderLoginDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", fmt.Errorf("some error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error generating access token: some error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: refresh token generation error",
			input: inputStruct{
				readerDTO: &dto.ReaderLoginDTO{
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, readerDTO *dto.ReaderLoginDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", fmt.Errorf("some error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error generating refresh token: some error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: save refresh token error",
			input: inputStruct{
				readerDTO: &dto.ReaderLoginDTO{
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, h *mockrepo.MockIPasswordHasher, t *mockrepo.MockITokenManager, readerDTO *dto.ReaderLoginDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), reader.ID, "refreshToken", refreshTokenTTL).Return(fmt.Errorf("some error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error saving refresh token: some error")
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
			readerService := implServices.NewReaderService(mockReaderRepo, nil, mockTokenManager, mockHasher)

			testCase.mockBehaviour(mockReaderRepo, mockHasher, mockTokenManager, testCase.input.readerDTO, testCase.input.reader)

			_, err := readerService.SignIn(context.Background(), testCase.input.readerDTO)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_RefreshTokens(t *testing.T) {
	type mockBehaviour func(
		m *mockrepo.MockIReaderRepo,
		t *mockrepo.MockITokenManager,
		refreshToken string,
		existingReader *models.ReaderModel,
	)

	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		refreshToken   string
		existingReader *models.ReaderModel
	}

	const (
		accessTokenTTL  = time.Hour * 2
		refreshTokenTTL = time.Hour * 24 * 30
	)

	testTable := []struct {
		name          string
		input         inputStruct
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Error: successful token refresh",
			input: inputStruct{
				refreshToken: "validRefreshToken",
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, t *mockrepo.MockITokenManager, refreshToken string, existingReader *models.ReaderModel) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), refreshToken).Return(existingReader, nil)
				t.EXPECT().NewJWT(existingReader.ID, accessTokenTTL).Return("newAccessToken", nil)
				t.EXPECT().NewRefreshToken().Return("newRefreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), existingReader.ID, "newRefreshToken", refreshTokenTTL).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: failed to get reader by refresh token",
			input: inputStruct{
				refreshToken:   "invalidRefreshToken",
				existingReader: nil,
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, t *mockrepo.MockITokenManager, refreshToken string, existingReader *models.ReaderModel) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), refreshToken).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrNotFound
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error generating access token",
			input: inputStruct{
				refreshToken: "validRefreshToken",
				existingReader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "John Doe",
					PhoneNumber: "12345678901",
					Password:    "hashedPassword",
					Age:         25,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, t *mockrepo.MockITokenManager, refreshToken string, existingReader *models.ReaderModel) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), refreshToken).Return(existingReader, nil)
				t.EXPECT().NewJWT(existingReader.ID, accessTokenTTL).Return("", fmt.Errorf("error generating JWT"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error generating access token: error generating JWT")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
			readerService := implServices.NewReaderService(mockReaderRepo, nil, mockTokenManager, nil)

			testCase.mockBehaviour(mockReaderRepo, mockTokenManager, testCase.input.refreshToken, testCase.input.existingReader)

			_, err := readerService.RefreshTokens(context.Background(), testCase.input.refreshToken)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_AddToFavorites(t *testing.T) {
	type mockBehaviour func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, readerID, bookID uuid.UUID)

	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		readerID uuid.UUID
		bookID   uuid.UUID
	}

	testTable := []struct {
		name          string
		input         inputStruct
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success: successful AddToFavorites",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), readerID, bookID).Return(false, nil)
				r.EXPECT().AddToFavorites(gomock.Any(), readerID, bookID).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: reader not found",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader with this ID does not exist")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: book not found",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Book with this ID does not exist")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: already in favorites",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), readerID, bookID).Return(true, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Book is already in favorites")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error adding to favorites",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), readerID, bookID).Return(false, nil)
				r.EXPECT().AddToFavorites(gomock.Any(), readerID, bookID).Return(fmt.Errorf("error add"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Error adding book to favorites: error add")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			readerService := implServices.NewReaderService(mockReaderRepo, mockBookRepo, nil, nil)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, testCase.input.readerID, testCase.input.bookID)

			err := readerService.AddToFavorites(context.Background(), testCase.input.readerID, testCase.input.bookID)

			testCase.expected(t, err)
		})
	}
}

package serviceTests

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/services/errsService"
	"BookSmart/internal/services/implServices"
	"BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"BookSmart/pkg/logging"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReaderService_SignUp(t *testing.T) {
	type mockBehaviour func(m *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		reader *models.ReaderModel
	}

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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
				h.EXPECT().Hash(reader.Password).Return(gomock.Any().String(), nil)
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, fmt.Errorf("database error"))
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(reader, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrReaderAlreadyExist
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrEmptyReaderFio
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrEmptyReaderPhoneNumber
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrInvalidReaderAge
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrInvalidReaderPhoneNumberLen
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrInvalidReaderPhoneNumberFormat
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrInvalidReaderPasswordLen
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
				h.EXPECT().Hash(reader.Password).Return(gomock.Any().String(), nil)
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

			mockReaderRepo := mocks.NewMockIReaderRepo(ctrl)
			mockHasher := mocks.NewMockIPasswordHasher(ctrl)
			readerService := implServices.NewReaderService(
				mockReaderRepo, nil, nil,
				mockHasher, logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockHasher, testCase.args.reader)

			err := readerService.SignUp(context.Background(), testCase.args.reader)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_SignIn(t *testing.T) {
	type mockBehaviour func(
		r *mocks.MockIReaderRepo,
		h *mocks.MockIPasswordHasher,
		t *mocks.MockITokenManager,
		readerDTO *dto.ReaderSignInDTO,
		reader *models.ReaderModel,
	)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		readerDTO *dto.ReaderSignInDTO
		reader    *models.ReaderModel
	}

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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, t *mocks.MockITokenManager, readerDTO *dto.ReaderSignInDTO, reader *models.ReaderModel) {
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
			name: "Error reader not found",
			args: args{
				readerDTO: &dto.ReaderSignInDTO{
					PhoneNumber: "12345678901",
					Password:    "password",
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, t *mocks.MockITokenManager, readerDTO *dto.ReaderSignInDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrReaderDoesNotExists
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, t *mocks.MockITokenManager, readerDTO *dto.ReaderSignInDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(fmt.Errorf("[!] ERROR! Wrong password"))
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, t *mocks.MockITokenManager, readerDTO *dto.ReaderSignInDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", fmt.Errorf("some error"))
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, t *mocks.MockITokenManager, readerDTO *dto.ReaderSignInDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", nil)
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, h *mocks.MockIPasswordHasher, t *mocks.MockITokenManager, readerDTO *dto.ReaderSignInDTO, reader *models.ReaderModel) {
				r.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
				h.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
				t.EXPECT().NewJWT(reader.ID, accessTokenTTL).Return("accessToken", nil)
				t.EXPECT().NewRefreshToken().Return("refreshToken", nil)
				r.EXPECT().SaveRefreshToken(gomock.Any(), reader.ID, "refreshToken", refreshTokenTTL).Return(fmt.Errorf("some error"))
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

			mockReaderRepo := mocks.NewMockIReaderRepo(ctrl)
			mockHasher := mocks.NewMockIPasswordHasher(ctrl)
			mockTokenManager := mocks.NewMockITokenManager(ctrl)
			readerService := implServices.NewReaderService(
				mockReaderRepo, nil, mockTokenManager,
				mockHasher, logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockHasher, mockTokenManager, testCase.args.readerDTO, testCase.args.reader)

			_, err := readerService.SignIn(context.Background(), testCase.args.readerDTO)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_RefreshTokens(t *testing.T) {
	type mockBehaviour func(
		m *mocks.MockIReaderRepo,
		t *mocks.MockITokenManager,
		refreshToken string,
		existingReader *models.ReaderModel,
	)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		refreshToken   string
		existingReader *models.ReaderModel
	}

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
			mockBehaviour: func(r *mocks.MockIReaderRepo, t *mocks.MockITokenManager, refreshToken string, existingReader *models.ReaderModel) {
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
			name: "Error failed to get reader by refresh token",
			args: args{
				refreshToken:   "invalidRefreshToken",
				existingReader: nil,
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, t *mocks.MockITokenManager, refreshToken string, existingReader *models.ReaderModel) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), refreshToken).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsRepo.ErrNotFound
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
			mockBehaviour: func(r *mocks.MockIReaderRepo, t *mocks.MockITokenManager, refreshToken string, existingReader *models.ReaderModel) {
				r.EXPECT().GetByRefreshToken(gomock.Any(), refreshToken).Return(existingReader, nil)
				t.EXPECT().NewJWT(existingReader.ID, accessTokenTTL).Return("", fmt.Errorf("error generating JWT"))
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

			mockReaderRepo := mocks.NewMockIReaderRepo(ctrl)
			mockTokenManager := mocks.NewMockITokenManager(ctrl)
			readerService := implServices.NewReaderService(
				mockReaderRepo, nil, mockTokenManager,
				nil, logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockTokenManager, testCase.args.refreshToken, testCase.args.existingReader)

			_, err := readerService.RefreshTokens(context.Background(), testCase.args.refreshToken)

			testCase.expected(t, err)
		})
	}
}

func TestReaderService_AddToFavorites(t *testing.T) {
	type mockBehaviour func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, readerID, bookID uuid.UUID)

	type expectedFunc func(t *testing.T, err error)
	type args struct {
		readerID uuid.UUID
		bookID   uuid.UUID
	}

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
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, readerID, bookID uuid.UUID) {
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
			name: "Error reader not found",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrReaderDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book not found",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error already in favorites",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), readerID, bookID).Return(true, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrBookAlreadyIsFavorite
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error adding to favorites",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{}, nil)
				r.EXPECT().IsFavorite(gomock.Any(), readerID, bookID).Return(false, nil)
				r.EXPECT().AddToFavorites(gomock.Any(), readerID, bookID).Return(fmt.Errorf("error add"))
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

			mockReaderRepo := mocks.NewMockIReaderRepo(ctrl)
			mockBookRepo := mocks.NewMockIBookRepo(ctrl)
			readerService := implServices.NewReaderService(
				mockReaderRepo, mockBookRepo, nil,
				nil, logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, testCase.args.readerID, testCase.args.bookID)

			err := readerService.AddToFavorites(context.Background(), testCase.args.readerID, testCase.args.bookID)

			testCase.expected(t, err)
		})
	}
}

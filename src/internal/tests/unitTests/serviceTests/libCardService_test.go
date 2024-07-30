package serviceTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/services/errsService"
	"BookSmart/internal/services/implServices"
	mockrepositories "BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"BookSmart/pkg/logging"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLibCardService_Create(t *testing.T) {
	type mockBehavior func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		readerID uuid.UUID
	}

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successful creation",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errsRepo.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error checking existing library card",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error library card already exists",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				existingCard := &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     readerID,
					LibCardNum:   "1234567890123",
					Validity:     365,
					IssueDate:    time.Now().AddDate(0, 0, -10),
					ActionStatus: true,
				}
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(existingCard, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrLibCardAlreadyExist
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error creating library card",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errsRepo.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("create error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLibCardRepo := mockrepositories.NewMockILibCardRepo(ctrl)
			libCardService := implServices.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockLibCardRepo, testCase.args.readerID)

			err := libCardService.Create(context.Background(), testCase.args.readerID)

			testCase.expected(t, err)
		})
	}
}

func TestLibCardService_Update(t *testing.T) {
	type mockBehavior func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		libCard *models.LibCardModel
	}

	testReaderID := uuid.New()

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successful update",
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     testReaderID,
					LibCardNum:   "1234567890123",
					Validity:     implServices.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: false,
				},
			},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel) {
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error checking existing library card",
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     testReaderID,
					LibCardNum:   "1234567890123",
					Validity:     implServices.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel) {
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error library card does not exists",
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     testReaderID,
					LibCardNum:   "1234567890123",
					Validity:     implServices.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel) {
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrLibCardDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error library card is valid",
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     testReaderID,
					LibCardNum:   "1234567890123",
					Validity:     implServices.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel) {
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrLibCardIsValid
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error update library card",
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     testReaderID,
					LibCardNum:   "1234567890123",
					Validity:     implServices.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: false,
				},
			},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel) {
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, nil)
				m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("update error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLibCardRepo := mockrepositories.NewMockILibCardRepo(ctrl)
			libCardService := implServices.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockLibCardRepo, testCase.args.libCard)

			err := libCardService.Update(context.Background(), testCase.args.libCard)

			testCase.expected(t, err)
		})
	}
}

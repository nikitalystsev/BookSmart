package serviceTests

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLibCardService_Create(t *testing.T) {
	type args struct {
		readerID uuid.UUID
	}
	type mockBehavior func(m *mockrepo.MockILibCardRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successful creation",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(nil, errs.ErrLibCardDoesNotExists)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error checking existing library card",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error library card already exists",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				existingCard := &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     args.readerID,
					LibCardNum:   "1234567890123",
					Validity:     365,
					IssueDate:    time.Now().AddDate(0, 0, -10),
					ActionStatus: true,
				}
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(existingCard, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrLibCardAlreadyExist
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error creating library card",
			args: args{readerID: uuid.New()},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(nil, errs.ErrLibCardDoesNotExists)
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

			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
			libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockLibCardRepo, testCase.args)

			err := libCardService.Create(context.Background(), testCase.args.readerID)

			testCase.expected(t, err)
		})
	}
}

func TestLibCardService_Update(t *testing.T) {
	type args struct {
		libCard *models.LibCardModel
	}
	type mockBehavior func(m *mockrepo.MockILibCardRepo, args args)
	type expectedFunc func(t *testing.T, err error)

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
					Validity:     impl.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: false,
				},
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByNum(gomock.Any(), args.libCard.LibCardNum).Return(args.libCard, nil)
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
					Validity:     impl.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByNum(gomock.Any(), args.libCard.LibCardNum).Return(nil, errors.New("database error"))
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
					Validity:     impl.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByNum(gomock.Any(), args.libCard.LibCardNum).Return(nil, errs.ErrLibCardDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrLibCardDoesNotExists
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
					Validity:     impl.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByNum(gomock.Any(), args.libCard.LibCardNum).Return(args.libCard, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrLibCardIsValid
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
					Validity:     impl.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: false,
				},
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByNum(gomock.Any(), args.libCard.LibCardNum).Return(args.libCard, nil)
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

			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
			libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockLibCardRepo, testCase.args)

			err := libCardService.Update(context.Background(), testCase.args.libCard)

			testCase.expected(t, err)
		})
	}
}

func TestLibCardService_GetByReaderID(t *testing.T) {
	type args struct {
		readerID uuid.UUID
	}
	type mockBehavior func(m *mockrepo.MockILibCardRepo, args args)
	type expectedFunc func(t *testing.T, libCard *models.LibCardModel, err error)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successfully getting libCard by readerID",
			args: args{
				readerID: uuid.New(),
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				libCard := &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     args.readerID,
					LibCardNum:   "1234567890123",
					Validity:     impl.LibCardValidityPeriod,
					IssueDate:    time.Now(),
					ActionStatus: true,
				}
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(libCard, nil)
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, libCard)
			},
		},
		{
			name: "Error error checking libCard existence",
			args: args{
				readerID: uuid.New(),
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
				assert.Nil(t, libCard)
			},
		},
		{
			name: "Error library card does not exists",
			args: args{
				readerID: uuid.New(),
			},
			mockBehavior: func(m *mockrepo.MockILibCardRepo, args args) {
				m.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(nil, errs.ErrLibCardDoesNotExists)
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				expectedError := errs.ErrLibCardDoesNotExists
				assert.Equal(t, expectedError, err)
				assert.Nil(t, libCard)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
			libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockLibCardRepo, testCase.args)

			libCard, err := libCardService.GetByReaderID(context.Background(), testCase.args.readerID)

			testCase.expected(t, libCard, err)
		})
	}
}

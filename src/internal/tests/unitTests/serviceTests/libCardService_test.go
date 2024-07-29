package serviceTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/services/implServices"
	mockrepositories "BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLibCardService_Create(t *testing.T) {
	type mockBehavior func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID)
	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		readerID uuid.UUID
	}

	testReaderID := uuid.New()

	testTable := []struct {
		name         string
		input        inputStruct
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name:  "Success: successful creation",
			input: inputStruct{readerID: testReaderID},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errsRepo.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name:  "Error: error checking existing library card",
			input: inputStruct{readerID: testReaderID},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Error checking libCard existence: database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name:  "Error: library card already exists",
			input: inputStruct{readerID: testReaderID},
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
				expectedError := fmt.Errorf("[!] ERROR! User with ID %v already has a library card", testReaderID)
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name:  "Error: error creating library card",
			input: inputStruct{readerID: testReaderID},
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errsRepo.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Error creating libCard: create error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLibCardRepo := mockrepositories.NewMockILibCardRepo(ctrl)
			libCardService := implServices.NewLibCardService(mockLibCardRepo)

			testCase.mockBehavior(mockLibCardRepo, testCase.input.readerID)

			err := libCardService.Create(context.Background(), testCase.input.readerID)

			testCase.expected(t, err)
		})
	}
}

func TestLibCardService_Update(t *testing.T) {
	type mockBehavior func(m *mockrepositories.MockILibCardRepo, libCard *models.LibCardModel)
	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		libCard *models.LibCardModel
	}

	testReaderID := uuid.New()

	testTable := []struct {
		name         string
		input        inputStruct
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success: successful update",
			input: inputStruct{
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
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, errsRepo.ErrNotFound)
				m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: error checking existing library card",
			input: inputStruct{
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
				expectedError := errors.New("[!] ERROR! Error checking libCard existence: database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: library card does not exists",
			input: inputStruct{
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
				expectedError := fmt.Errorf("[!] ERROR! libCard with ID %v does not exist", "1234567890123")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error library card is valid",
			input: inputStruct{
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
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! libCard with ID %v is valid", "1234567890123")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error update library card",
			input: inputStruct{
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
				m.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, errsRepo.ErrNotFound)
				m.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("[!] ERROR! Error updating libCard: update error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLibCardRepo := mockrepositories.NewMockILibCardRepo(ctrl)
			libCardService := implServices.NewLibCardService(mockLibCardRepo)

			testCase.mockBehavior(mockLibCardRepo, testCase.input.libCard)

			err := libCardService.Update(context.Background(), testCase.input.libCard)

			testCase.expected(t, err)
		})
	}
}

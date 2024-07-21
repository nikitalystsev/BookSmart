package unitTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/services/implServices"
	mockrepositories "BookSmart/internal/tests/unitTests/mocks"
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

	testReaderID := uuid.New()

	testTable := []struct {
		name          string
		readerID      uuid.UUID
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:     "successful creation",
			readerID: testReaderID,
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errs.ErrNotFound)

				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "error checking existing library card",
			readerID: testReaderID,
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errors.New("some error"))
			},
			expectedError: errors.New("[!] ERROR! Error checking libCard existence: some error"),
		},
		{
			name:     "library card already exists",
			readerID: testReaderID,
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
			expectedError: fmt.Errorf("[!] ERROR! User with ID %v already has a library card", testReaderID),
		},
		{
			name:     "error creating library card",
			readerID: testReaderID,
			mockBehavior: func(m *mockrepositories.MockILibCardRepo, readerID uuid.UUID) {
				m.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errs.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("create error"))
			},
			expectedError: errors.New("[!] ERROR! Error creating libCard: create error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockLibCardRepo := mockrepositories.NewMockILibCardRepo(ctrl)
			libCardService := implServices.NewLibCardService(mockLibCardRepo)
			readerID := testCase.readerID

			testCase.mockBehavior(mockLibCardRepo, readerID)

			err := libCardService.Create(context.Background(), readerID)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

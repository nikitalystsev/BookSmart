package unitTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/services/impl"
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

func TestReservationService_Create(t *testing.T) {
	type mockBehavior func(
		res *mockrepositories.MockIReservationRepo,
		b *mockrepositories.MockIBookRepo,
		r *mockrepositories.MockIReaderRepo,
		readerID, bookID uuid.UUID,
	)

	testReaderID := uuid.New()
	testBookID := uuid.New()

	testTable := []struct {
		name          string
		readerID      uuid.UUID
		bookID        uuid.UUID
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:     "successful reservation creation",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(&models.ReaderModel{ID: readerID}, nil)
				b.EXPECT().GetByID(context.Background(), bookID).Return(&models.BookModel{ID: bookID}, nil)
				res.EXPECT().GetByReaderAndBook(context.Background(), readerID, bookID).Return(nil, errors.New("[!] ERROR! Object not found"))
				res.EXPECT().Create(context.Background(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "error checking reader existence",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(nil, errors.New("some error"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Error checking reader existence: some error"),
		},
		{
			name:     "reader not found",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(nil, nil)
			},
			expectedError: fmt.Errorf("[!] ERROR! Error checking reader existence: reader does not exist"),
		},
		{
			name:     "error checking book existence",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(&models.ReaderModel{ID: readerID}, nil)
				b.EXPECT().GetByID(context.Background(), bookID).Return(nil, errors.New("some error"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Error checking book existence: some error"),
		},
		{
			name:     "book not found",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(&models.ReaderModel{ID: readerID}, nil)
				b.EXPECT().GetByID(context.Background(), bookID).Return(nil, nil)
			},
			expectedError: fmt.Errorf("[!] ERROR! Error checking book existence: book does not exist"),
		},
		{
			name:     "error checking existing reservation",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(&models.ReaderModel{ID: readerID}, nil)
				b.EXPECT().GetByID(context.Background(), bookID).Return(&models.BookModel{ID: bookID}, nil)
				res.EXPECT().GetByReaderAndBook(context.Background(), readerID, bookID).Return(nil, errors.New("some error"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Error checking existing reservation: some error"),
		},
		{
			name:     "reservation already exists",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(&models.ReaderModel{ID: readerID}, nil)
				b.EXPECT().GetByID(context.Background(), bookID).Return(&models.BookModel{ID: bookID}, nil)
				res.EXPECT().GetByReaderAndBook(context.Background(), readerID, bookID).Return(&models.ReservationModel{ID: uuid.New(), ReaderID: readerID, BookID: bookID}, nil)
			},
			expectedError: fmt.Errorf("[!] ERROR! Book with ID %v is already reserved by reader with ID %v", testBookID, testReaderID),
		},
		{
			name:     "error creating reservation",
			readerID: testReaderID,
			bookID:   testBookID,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, b *mockrepositories.MockIBookRepo, r *mockrepositories.MockIReaderRepo, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(context.Background(), readerID).Return(&models.ReaderModel{ID: readerID}, nil)
				b.EXPECT().GetByID(context.Background(), bookID).Return(&models.BookModel{ID: bookID}, nil)
				res.EXPECT().GetByReaderAndBook(context.Background(), readerID, bookID).Return(nil, errors.New("[!] ERROR! Object not found"))
				res.EXPECT().Create(context.Background(), gomock.Any()).Return(errors.New("error creating reservation"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Error creating reservation: error creating reservation"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepositories.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepositories.NewMockIBookRepo(ctrl)
			mockReservationRepo := mockrepositories.NewMockIReservationRepo(ctrl)
			reservationService := impl.NewReservationService(mockReservationRepo, mockBookRepo, mockReaderRepo)

			readerID := testCase.readerID
			bookID := testCase.bookID

			testCase.mockBehavior(mockReservationRepo, mockBookRepo, mockReaderRepo, readerID, bookID)

			err := reservationService.Create(context.Background(), readerID, bookID)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestReservationService_Update(t *testing.T) {
	type mockBehavior func(
		res *mockrepositories.MockIReservationRepo,
		reservation *models.ReservationModel,
	)

	testReservation := &models.ReservationModel{
		ID:         uuid.New(),
		ReaderID:   uuid.New(),
		BookID:     uuid.New(),
		IssueDate:  time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
		State:      "Активна",
	}

	testTable := []struct {
		name          string
		reservation   *models.ReservationModel
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:        "successful reservation update",
			reservation: testReservation,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, reservation *models.ReservationModel) {
				res.EXPECT().GetByID(context.Background(), reservation.ID).Return(reservation, nil)
				res.EXPECT().Update(context.Background(), reservation).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "reservation not found",
			reservation: testReservation,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, reservation *models.ReservationModel) {
				res.EXPECT().GetByID(context.Background(), reservation.ID).Return(nil, errors.New("[!] ERROR! Object not found"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Reservation with ID %v not found", testReservation.ID),
		},
		{
			name:        "error checking reservation existence",
			reservation: testReservation,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, reservation *models.ReservationModel) {
				res.EXPECT().GetByID(context.Background(), reservation.ID).Return(nil, errors.New("some error"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Error checking reservation existence: some error"),
		},
		{
			name: "reservation is closed",
			reservation: &models.ReservationModel{
				ID:         testReservation.ID,
				ReaderID:   testReservation.ReaderID,
				BookID:     testReservation.BookID,
				IssueDate:  testReservation.IssueDate,
				ReturnDate: testReservation.ReturnDate,
				State:      "Закрыта",
			},
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, reservation *models.ReservationModel) {
				res.EXPECT().GetByID(context.Background(), reservation.ID).Return(reservation, nil)
			},
			expectedError: fmt.Errorf("[!] ERROR! This reservation is already closed"),
		},
		{
			name: "reservation is overdue",
			reservation: &models.ReservationModel{
				ID:         testReservation.ID,
				ReaderID:   testReservation.ReaderID,
				BookID:     testReservation.BookID,
				IssueDate:  testReservation.IssueDate,
				ReturnDate: testReservation.ReturnDate,
				State:      "Просрочена",
			},
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, reservation *models.ReservationModel) {
				res.EXPECT().GetByID(context.Background(), reservation.ID).Return(reservation, nil)
			},
			expectedError: fmt.Errorf("[!] ERROR! This reservation is past its return date"),
		},
		{
			name:        "error updating reservation",
			reservation: testReservation,
			mockBehavior: func(res *mockrepositories.MockIReservationRepo, reservation *models.ReservationModel) {
				res.EXPECT().GetByID(context.Background(), reservation.ID).Return(reservation, nil)
				res.EXPECT().Update(context.Background(), reservation).Return(errors.New("some error"))
			},
			expectedError: fmt.Errorf("[!] ERROR! Error updating reservation: some error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReservationRepo := mockrepositories.NewMockIReservationRepo(ctrl)
			reservationService := impl.NewReservationService(mockReservationRepo, nil, nil)

			testCase.mockBehavior(mockReservationRepo, testCase.reservation)

			err := reservationService.Update(context.Background(), testCase.reservation)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

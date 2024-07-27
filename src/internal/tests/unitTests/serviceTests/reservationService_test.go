package serviceTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/services/implServices"
	"BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReservationService_Create(t *testing.T) {
	type mockBehaviour func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID)

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
			name: "Success: create reservation",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{CopiesNumber: 1}, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				b.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{CopiesNumber: 1}, nil)
				trm.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: reader does not exist",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(nil, fmt.Errorf("reader not found"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Error checking reader existence: reader not found")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: book does not exist",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(nil, fmt.Errorf("book not found"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Error checking book existence: book not found")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: invalid library card",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: false}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader does not have a valid library card")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: overdue books exist",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return([]*models.ReservationModel{{}}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader has overdue books")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: active reservations limit reached",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(make([]*models.ReservationModel, implServices.MaxBooksPerReader), nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader has reached the limit of active reservations")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: no book copies available",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{CopiesNumber: 0}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! No copies of the book are available in the library")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: reader does not meet age requirement",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{Age: 10}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{AgeLimit: 18, CopiesNumber: 1}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Reader does not meet the age requirement for this book")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: book is unique and cannot be reserved",
			input: inputStruct{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{Rarity: implServices.BookRarityUnique, CopiesNumber: 1}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! This book is unique and cannot be reserved")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mocks.NewMockIReaderRepo(ctrl)
			mockBookRepo := mocks.NewMockIBookRepo(ctrl)
			mockLibCardRepo := mocks.NewMockILibCardRepo(ctrl)
			mockReservationRepo := mocks.NewMockIReservationRepo(ctrl)
			mockTransactionManager := mocks.NewMockITransactionManager(ctrl)
			reservationService := implServices.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				mockTransactionManager,
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, mockLibCardRepo, mockReservationRepo, mockTransactionManager, testCase.input.readerID, testCase.input.bookID)

			err := reservationService.Create(context.Background(), testCase.input.readerID, testCase.input.bookID)
			testCase.expected(t, err)
		})
	}
}

func TestReservationService_Update(t *testing.T) {
	type mockBehaviour func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel)

	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		reservation *models.ReservationModel
	}

	testTable := []struct {
		name          string
		input         inputStruct
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success: valid update",
			input: inputStruct{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), reservation.ReaderID).Return(nil, nil)
				b.EXPECT().GetByID(gomock.Any(), reservation.BookID).Return(&models.BookModel{Rarity: implServices.BookRarityCommon}, nil)
				res.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: invalid library card",
			input: inputStruct{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: false}, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.EqualError(t, err, "[!] ERROR! Reader does not have a valid library card")
			},
		},
		{
			name: "Error: overdue books exist",
			input: inputStruct{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), reservation.ReaderID).Return([]*models.ReservationModel{{}}, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.EqualError(t, err, "[!] ERROR! Reader has overdue books")
			},
		},
		{
			name: "Error: reservation already extended",
			input: inputStruct{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationExtended,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), reservation.ReaderID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.EqualError(t, err, "[!] ERROR! This reservation has already been extended")
			},
		},
		{
			name: "Error: unique book not renewable",
			input: inputStruct{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetOverdueByReaderID(gomock.Any(), reservation.ReaderID).Return(nil, nil)
				b.EXPECT().GetByID(gomock.Any(), reservation.BookID).Return(&models.BookModel{Rarity: implServices.BookRarityUnique}, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.EqualError(t, err, "[!] ERROR! This book is not renewed")
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReaderRepo := mocks.NewMockIReaderRepo(ctrl)
			mockBookRepo := mocks.NewMockIBookRepo(ctrl)
			mockLibCardRepo := mocks.NewMockILibCardRepo(ctrl)
			mockReservationRepo := mocks.NewMockIReservationRepo(ctrl)

			reservationService := implServices.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				nil,
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, mockLibCardRepo, mockReservationRepo, testCase.input.reservation)

			err := reservationService.Update(context.Background(), testCase.input.reservation)
			testCase.expected(t, err)
		})
	}
}

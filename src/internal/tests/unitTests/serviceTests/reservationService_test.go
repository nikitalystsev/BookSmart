package serviceTests

import (
	"BookSmart/internal/models"
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
)

func TestReservationService_Create(t *testing.T) {
	type mockBehaviour func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID)
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
			name: "Success create reservation",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{CopiesNumber: 1}, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
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
			name: "Error reader does not exist",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(nil, fmt.Errorf("reader not found"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("reader not found")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book does not exist",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(nil, fmt.Errorf("book not found"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("book not found")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error invalid libCard",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: false}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrLibCardIsInvalid
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error expired books exist",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return([]*models.ReservationModel{{}}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrReaderHasExpiredBooks
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error active reservations limit reached",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(make([]*models.ReservationModel, implServices.MaxBooksPerReader), nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrReservationsLimitExceeded
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error no book copies available",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{CopiesNumber: 0}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrBookNoCopiesNum
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error reader does not meet age requirement",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{Age: 10}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{AgeLimit: 18, CopiesNumber: 1}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrReservationAgeLimit
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book is unique and cannot be reserved",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, trm *mocks.MockITransactionManager, readerID, bookID uuid.UUID) {
				r.EXPECT().GetByID(gomock.Any(), readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), bookID).Return(&models.BookModel{Rarity: implServices.BookRarityUnique, CopiesNumber: 1}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrUniqueBookNotReserved
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
				logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, mockLibCardRepo, mockReservationRepo, mockTransactionManager, testCase.args.readerID, testCase.args.bookID)

			err := reservationService.Create(context.Background(), testCase.args.readerID, testCase.args.bookID)
			testCase.expected(t, err)
		})
	}
}

func TestReservationService_Update(t *testing.T) {
	type mockBehaviour func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		reservation *models.ReservationModel
	}

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success valid update",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), reservation.ReaderID).Return(nil, nil)
				b.EXPECT().GetByID(gomock.Any(), reservation.BookID).Return(&models.BookModel{Rarity: implServices.BookRarityCommon}, nil)
				res.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Error invalid library card",
			args: args{
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
				expiredError := errsService.ErrLibCardIsInvalid
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error expired books exist",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), reservation.ReaderID).Return([]*models.ReservationModel{{}}, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errsService.ErrReaderHasExpiredBooks
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error reservation already extended",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationExtended,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), reservation.ReaderID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errsService.ErrReservationIsAlreadyExtended
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error unique book not renewable",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    implServices.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mocks.MockIReaderRepo, b *mocks.MockIBookRepo, l *mocks.MockILibCardRepo, res *mocks.MockIReservationRepo, reservation *models.ReservationModel) {
				l.EXPECT().GetByReaderID(gomock.Any(), reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), reservation.ReaderID).Return(nil, nil)
				b.EXPECT().GetByID(gomock.Any(), reservation.BookID).Return(&models.BookModel{Rarity: implServices.BookRarityUnique}, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errsService.ErrRareAndUniqueBookNotExtended
				assert.Equal(t, expiredError, err)
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
				logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, mockLibCardRepo, mockReservationRepo, testCase.args.reservation)

			err := reservationService.Update(context.Background(), testCase.args.reservation)
			testCase.expected(t, err)
		})
	}
}

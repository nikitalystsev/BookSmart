package serviceTests

import (
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"BookSmart-services/impl"
	mockrepo "Booksmart/internal/tests/unitTests/serviceTests/mocks"
	"Booksmart/pkg/logging"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReservationService_Create(t *testing.T) {
	type args struct {
		readerID uuid.UUID
		bookID   uuid.UUID
	}
	type mockBehaviour func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args)
	type expectedFunc func(t *testing.T, err error)

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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{CopiesNumber: 1}, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				b.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{CopiesNumber: 1}, nil)
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(nil, fmt.Errorf("reader not found"))
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(nil, fmt.Errorf("book not found"))
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
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(&models.LibCardModel{ActionStatus: false}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrLibCardIsInvalid
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error expired books exist",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return([]*models.ReservationModel{{}}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReaderHasExpiredBooks
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error active reservations limit reached",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(make([]*models.ReservationModel, impl.MaxBooksPerReader), nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReservationsLimitExceeded
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error no book copies available",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{CopiesNumber: 0}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookNoCopiesNum
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error reader does not meet age requirement",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{Age: 10}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{AgeLimit: 18, CopiesNumber: 1}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReservationAgeLimit
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book is unique and cannot be reserved",
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, trm *mockrepo.MockITransactionManager, args args) {
				r.EXPECT().GetByID(gomock.Any(), args.readerID).Return(&models.ReaderModel{}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				res.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, nil)
				l.EXPECT().GetByReaderID(gomock.Any(), args.readerID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{Rarity: impl.BookRarityUnique, CopiesNumber: 1}, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrUniqueBookNotReserved
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
			mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
			mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
			reservationService := impl.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				mockTransactionManager,
				logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, mockLibCardRepo, mockReservationRepo, mockTransactionManager, testCase.args)

			err := reservationService.Create(context.Background(), testCase.args.readerID, testCase.args.bookID)
			testCase.expected(t, err)
		})
	}
}

func TestReservationService_Update(t *testing.T) {
	type args struct {
		reservation *models.ReservationModel
	}
	type mockBehaviour func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, args args)
	type expectedFunc func(t *testing.T, err error)

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
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, args args) {
				l.EXPECT().GetByReaderID(gomock.Any(), args.reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.reservation.ReaderID).Return(nil, nil)
				b.EXPECT().GetByID(gomock.Any(), args.reservation.BookID).Return(&models.BookModel{Rarity: impl.BookRarityCommon}, nil)
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
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, args args) {
				l.EXPECT().GetByReaderID(gomock.Any(), args.reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: false}, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errs.ErrLibCardIsInvalid
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error expired books exist",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, args args) {
				l.EXPECT().GetByReaderID(gomock.Any(), args.reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.reservation.ReaderID).Return([]*models.ReservationModel{{}}, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errs.ErrReaderHasExpiredBooks
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error reservation already extended",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    impl.ReservationExtended,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, args args) {
				l.EXPECT().GetByReaderID(gomock.Any(), args.reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.reservation.ReaderID).Return(nil, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errs.ErrReservationIsAlreadyExtended
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error unique book not renewable",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(r *mockrepo.MockIReaderRepo, b *mockrepo.MockIBookRepo, l *mockrepo.MockILibCardRepo, res *mockrepo.MockIReservationRepo, args args) {
				l.EXPECT().GetByReaderID(gomock.Any(), args.reservation.ReaderID).Return(&models.LibCardModel{ActionStatus: true}, nil)
				res.EXPECT().GetExpiredByReaderID(gomock.Any(), args.reservation.ReaderID).Return(nil, nil)
				b.EXPECT().GetByID(gomock.Any(), args.reservation.BookID).Return(&models.BookModel{Rarity: impl.BookRarityUnique}, nil)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errs.ErrRareAndUniqueBookNotExtended
				assert.Equal(t, expiredError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
			mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)

			reservationService := impl.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				nil,
				logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReaderRepo, mockBookRepo, mockLibCardRepo, mockReservationRepo, testCase.args)

			err := reservationService.Update(context.Background(), testCase.args.reservation)
			testCase.expected(t, err)
		})
	}
}

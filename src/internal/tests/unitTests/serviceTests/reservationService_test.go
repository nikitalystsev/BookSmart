package serviceTests

import (
	mockrepo "Booksmart/internal/tests/unitTests/serviceTests/mocks"
	"Booksmart/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
				b.EXPECT().GetByID(gomock.Any(), args.bookID).Return(&models.BookModel{CopiesNumber: 0}, nil)
				trm.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
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

func TestReaderService_GetByBookID(t *testing.T) {
	type args struct {
		reservation *models.ReservationModel
	}
	type mockBehaviour func(res *mockrepo.MockIReservationRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name          string
		args          args
		mockBehaviour mockBehaviour
		expected      expectedFunc
	}{
		{
			name: "Success success get by book id",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(res *mockrepo.MockIReservationRepo, args args) {
				res.EXPECT().GetByBookID(gomock.Any(), args.reservation.BookID).Return([]*models.ReservationModel{args.reservation}, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Error reservation does not exist",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(res *mockrepo.MockIReservationRepo, args args) {
				res.EXPECT().GetByBookID(gomock.Any(), args.reservation.BookID).Return(nil, errs.ErrReservationDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expiredError := errs.ErrReservationDoesNotExists
				assert.Equal(t, expiredError, err)
			},
		},
		{
			name: "Error database error",
			args: args{
				reservation: &models.ReservationModel{
					ReaderID: uuid.New(),
					BookID:   uuid.New(),
					State:    impl.ReservationIssued,
				},
			},
			mockBehaviour: func(res *mockrepo.MockIReservationRepo, args args) {
				res.EXPECT().GetByBookID(gomock.Any(), args.reservation.BookID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expiredError := errors.New("database error")
				assert.Equal(t, expiredError, err)
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

			reservationService := impl.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				nil,
				logging.GetLoggerForTests(),
			)

			testCase.mockBehaviour(mockReservationRepo, testCase.args)

			_, err := reservationService.GetByBookID(context.Background(), testCase.args.reservation.BookID)
			testCase.expected(t, err)
		})
	}
}

func TestReservationService_GetByID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	type mockBehavior func(m *mockrepo.MockIReservationRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successfully getting reservation by ID",
			args: args{
				id: uuid.New(),
			},
			mockBehavior: func(m *mockrepo.MockIReservationRepo, args args) {
				reservation := &models.ReservationModel{
					ID:         args.id,
					BookID:     uuid.New(),
					ReaderID:   uuid.New(),
					IssueDate:  time.Now(),
					ReturnDate: time.Now().AddDate(0, 0, 14),
					State:      "Issued",
				}
				m.EXPECT().GetByID(gomock.Any(), args.id).Return(reservation, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Error error checking reservation existence",
			args: args{
				id: uuid.New(),
			},
			mockBehavior: func(m *mockrepo.MockIReservationRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.id).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error reservation does not exist",
			args: args{
				id: uuid.New(),
			},
			mockBehavior: func(m *mockrepo.MockIReservationRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.id).Return(nil, errs.ErrReservationDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrReservationDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)

			reservationService := impl.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				nil,
				logging.GetLoggerForTests(),
			)

			testCase.mockBehavior(mockReservationRepo, testCase.args)

			_, err := reservationService.GetByID(context.Background(), testCase.args.id)

			testCase.expected(t, err)
		})
	}
}

func TestReservationService_GetAllReservationsByReaderID(t *testing.T) {
	type args struct {
		readerID uuid.UUID
	}
	type mockBehavior func(m *mockrepo.MockIReservationRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testReaderID := uuid.New()

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successfully getting all reservations by readerID",
			args: args{
				readerID: testReaderID,
			},
			mockBehavior: func(m *mockrepo.MockIReservationRepo, args args) {
				activeReservations := []*models.ReservationModel{
					{
						ID:         uuid.New(),
						BookID:     uuid.New(),
						ReaderID:   args.readerID,
						IssueDate:  time.Now(),
						ReturnDate: time.Now().AddDate(0, 0, 14),
						State:      "Issued",
					},
				}
				expiredReservations := []*models.ReservationModel{
					{
						ID:         uuid.New(),
						BookID:     uuid.New(),
						ReaderID:   args.readerID,
						ReturnDate: time.Now().AddDate(0, -1, 0),
						State:      "Expired",
					},
				}
				m.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(activeReservations, nil)
				m.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(expiredReservations, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Error checking active reservations",
			args: args{
				readerID: testReaderID,
			},
			mockBehavior: func(m *mockrepo.MockIReservationRepo, args args) {
				m.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error checking expired reservations",
			args: args{
				readerID: testReaderID,
			},
			mockBehavior: func(m *mockrepo.MockIReservationRepo, args args) {
				activeReservations := []*models.ReservationModel{
					{
						ID:         uuid.New(),
						BookID:     uuid.New(),
						ReaderID:   args.readerID,
						IssueDate:  time.Now(),
						ReturnDate: time.Now().AddDate(0, 0, 14),
						State:      "Issued",
					},
				}
				m.EXPECT().GetActiveByReaderID(gomock.Any(), args.readerID).Return(activeReservations, nil)
				m.EXPECT().GetExpiredByReaderID(gomock.Any(), args.readerID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
			mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)

			reservationService := impl.NewReservationService(
				mockReservationRepo,
				mockBookRepo,
				mockReaderRepo,
				mockLibCardRepo,
				nil,
				logging.GetLoggerForTests(),
			)
			testCase.mockBehavior(mockReservationRepo, testCase.args)

			_, err := reservationService.GetAllReservationsByReaderID(context.Background(), testCase.args.readerID)

			testCase.expected(t, err)
		})
	}
}

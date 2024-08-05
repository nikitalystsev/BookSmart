package repositoryTests

import (
	"BookSmart-repositories/impl"
	"BookSmart-services/models"
	"Booksmart/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestReservationRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rr := impl.NewReservationRepo(db, logging.GetLoggerForTests())

	type args struct {
		reservation *models.ReservationModel
	}

	testsTable := []struct {
		name         string
		mockBehavior func(args args)
		args         args
		expected     func(t *testing.T, err error)
	}{
		{
			name: "Success create reservation",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO bs.reservation`).
					WithArgs(args.reservation.ID, args.reservation.ReaderID, args.reservation.BookID,
						args.reservation.IssueDate, args.reservation.ReturnDate, args.reservation.State).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				reservation: &models.ReservationModel{
					ID:         uuid.New(),
					ReaderID:   uuid.New(),
					BookID:     uuid.New(),
					IssueDate:  time.Now(),
					ReturnDate: time.Now().Add(14 * 24 * time.Hour),
					State:      "Issued",
				},
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error executing query",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO bs.reservation`).
					WithArgs(args.reservation.ID, args.reservation.ReaderID, args.reservation.BookID, args.reservation.IssueDate, args.reservation.ReturnDate, args.reservation.State).
					WillReturnError(errors.New("insert error"))
			},
			args: args{
				reservation: &models.ReservationModel{
					ID:         uuid.New(),
					ReaderID:   uuid.New(),
					BookID:     uuid.New(),
					IssueDate:  time.Now(),
					ReturnDate: time.Now().Add(14 * 24 * time.Hour),
					State:      "active",
				},
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("insert error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = rr.Create(context.Background(), testCase.args.reservation)

			testCase.expected(t, err)
		})
	}
}

func TestReservationRepo_GetByReaderAndBook(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rr := impl.NewReservationRepo(db, logging.GetLoggerForTests())

	type args struct {
		readerID uuid.UUID
		bookID   uuid.UUID
	}

	testsTable := []struct {
		name         string
		mockBehavior func(args args)
		args         args
		expected     func(t *testing.T, reservation *models.ReservationModel, err error)
	}{
		{
			name: "Success get reservation by reader and book",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
					AddRow(uuid.New(), args.readerID, args.bookID, time.Now(), time.Now().Add(14*24*time.Hour), "active")

				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID, args.bookID).WillReturnRows(rows)
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			expected: func(t *testing.T, reservation *models.ReservationModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reservation)
				assert.Equal(t, "active", reservation.State)
			},
		},
		{
			name: "Error executing query",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID, args.bookID).WillReturnError(errors.New("query error"))
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			expected: func(t *testing.T, reservation *models.ReservationModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reservation)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var reservation *models.ReservationModel
			reservation, err = rr.GetByReaderAndBook(context.Background(), testCase.args.readerID, testCase.args.bookID)

			testCase.expected(t, reservation, err)
		})
	}
}

func TestReservationRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rr := impl.NewReservationRepo(db, logging.GetLoggerForTests())

	type args struct {
		id uuid.UUID
	}

	testsTable := []struct {
		name         string
		mockBehavior func(args args)
		args         args
		expected     func(t *testing.T, reservation *models.ReservationModel, err error)
	}{
		{
			name: "Success get reservation by reader and book",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
					AddRow(args.id, uuid.New(), uuid.New(), time.Now(), time.Now().Add(14*24*time.Hour), "active")

				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.id).WillReturnRows(rows)
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, reservation *models.ReservationModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reservation)
				assert.Equal(t, "active", reservation.State)
			},
		},
		{
			name: "Error executing query",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.id).WillReturnError(errors.New("query error"))
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, reservation *models.ReservationModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reservation)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var reservation *models.ReservationModel
			reservation, err = rr.GetByID(context.Background(), testCase.args.id)

			testCase.expected(t, reservation, err)
		})
	}
}

func TestReservationRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rr := impl.NewReservationRepo(db, logging.GetLoggerForTests())

	type args struct {
		reservation *models.ReservationModel
	}

	testsTable := []struct {
		name         string
		mockBehavior func(args args)
		args         args
		expected     func(t *testing.T, err error)
	}{
		{
			name: "Success update reservation",
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE bs.reservation SET (.+) WHERE (.+)`).
					WithArgs(args.reservation.IssueDate, args.reservation.ReturnDate, args.reservation.State, args.reservation.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				reservation: &models.ReservationModel{
					ID:         uuid.New(),
					IssueDate:  time.Now(),
					ReturnDate: time.Now().AddDate(0, 0, 7),
					State:      "Extended",
				},
			},
			expected: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Error updating reservation",
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE bs.reservation SET (.+) WHERE (.+)`).
					WithArgs(args.reservation.IssueDate, args.reservation.ReturnDate, args.reservation.State, args.reservation.ID).
					WillReturnError(errors.New("update error"))
			},
			args: args{
				reservation: &models.ReservationModel{
					ID:         uuid.New(),
					IssueDate:  time.Now(),
					ReturnDate: time.Now().AddDate(0, 0, 7),
					State:      "Extended",
				},
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("update error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = rr.Update(context.Background(), testCase.args.reservation)

			testCase.expected(t, err)
		})
	}
}

func TestReservationRepo_GetExpiredByReaderID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rr := impl.NewReservationRepo(db, logging.GetLoggerForTests())

	type args struct {
		readerID   uuid.UUID
		returnDate time.Time
	}

	testsTable := []struct {
		name         string
		mockBehavior func(args args)
		args         args
		expected     func(t *testing.T, reservations []*models.ReservationModel, err error)
	}{
		{
			name: "Success get expired reservations",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
					AddRow(uuid.New(), args.readerID, uuid.New(), time.Now().AddDate(0, 0, -10), time.Now().AddDate(0, 0, -5), "Expired")

				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID, time.Now()).WillReturnRows(rows)
			},
			args: args{
				readerID:   uuid.New(),
				returnDate: time.Now(),
			},
			expected: func(t *testing.T, reservations []*models.ReservationModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reservations)
				assert.Equal(t, 1, len(reservations))
				assert.Equal(t, "Expired", reservations[0].State)
			},
		},
		{
			name: "Error query execution fails",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID, time.Now()).WillReturnError(errors.New("query error"))
			},
			args: args{
				readerID:   uuid.New(),
				returnDate: time.Now(),
			},
			expected: func(t *testing.T, reservations []*models.ReservationModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reservations)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			name: "Error row scan fails",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
					AddRow("invalid-uuid", args.readerID, uuid.New(), time.Now().AddDate(0, 0, -10), time.Now().AddDate(0, 0, -5), "overdue")

				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID, time.Now()).WillReturnRows(rows)
			},
			args: args{
				readerID:   uuid.New(),
				returnDate: time.Now(),
			},
			expected: func(t *testing.T, reservations []*models.ReservationModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reservations)
				expectedError := fmt.Errorf("sql: Scan error on column index 0, name \"id\": Scan: invalid UUID length: 12")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var reservations []*models.ReservationModel
			reservations, err = rr.GetExpiredByReaderID(context.Background(), testCase.args.readerID)

			testCase.expected(t, reservations, err)
		})
	}
}

func TestReservationRepo_GetActiveByReaderID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rr := impl.NewReservationRepo(db, logging.GetLoggerForTests())

	type args struct {
		readerID uuid.UUID
	}

	testsTable := []struct {
		name         string
		mockBehavior func(args args)
		args         args
		expected     func(t *testing.T, reservations []*models.ReservationModel, err error)
	}{
		{
			name: "Success get active reservations",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
					AddRow(uuid.New(), args.readerID, uuid.New(), time.Now().AddDate(0, 0, -10), time.Now().AddDate(0, 0, 5), "active").
					AddRow(uuid.New(), args.readerID, uuid.New(), time.Now().AddDate(0, 0, -5), time.Now().AddDate(0, 0, 10), "expired")

				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID).WillReturnRows(rows)
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, reservations []*models.ReservationModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reservations)
				assert.Equal(t, 2, len(reservations))
				assert.Equal(t, "active", reservations[0].State)
				assert.Equal(t, "expired", reservations[1].State)
			},
		},
		{
			name: "Error query execution fails",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID).WillReturnError(errors.New("query error"))
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, reservations []*models.ReservationModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reservations)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			name: "Error row scan fails",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
					AddRow("invalid-uuid", args.readerID, uuid.New(), time.Now().AddDate(0, 0, -10), time.Now().AddDate(0, 0, 5), "active")

				mock.ExpectQuery(`SELECT (.+) FROM bs.reservation WHERE (.+)`).
					WithArgs(args.readerID).WillReturnRows(rows)
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, reservations []*models.ReservationModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reservations)
				expectedError := fmt.Errorf("sql: Scan error on column index 0, name \"id\": Scan: invalid UUID length: 12")
				assert.Contains(t, err.Error(), expectedError.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var reservations []*models.ReservationModel
			reservations, err = rr.GetActiveByReaderID(context.Background(), testCase.args.readerID)

			testCase.expected(t, reservations, err)
		})
	}
}

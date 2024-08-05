package repositoryTests

import (
	errsRepo "BookSmart-repositories/errs"
	"BookSmart-repositories/impl"
	"BookSmart-services/models"
	"Booksmart/pkg/logging"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestLibCardRepo_Create(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	lcr := impl.NewLibCardRepo(db, logging.GetLoggerForTests())

	type args struct {
		libCard *models.LibCardModel
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success create libCard",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO bs.lib_card VALUES`).
					WithArgs(args.libCard.ID, args.libCard.ReaderID, args.libCard.LibCardNum, args.libCard.Validity,
						args.libCard.IssueDate, args.libCard.ActionStatus).
					WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     uuid.New(),
					LibCardNum:   "1234567890123",
					Validity:     12,
					IssueDate:    time.Now(),
					ActionStatus: true,
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
				mock.ExpectExec(`INSERT INTO bs.lib_card VALUES`).
					WithArgs(args.libCard.ID, args.libCard.ReaderID, args.libCard.LibCardNum, args.libCard.Validity,
						args.libCard.IssueDate, args.libCard.ActionStatus).
					WillReturnError(errors.New("insert error"))
			},
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					ReaderID:     uuid.New(),
					LibCardNum:   "1234567890123",
					Validity:     12,
					IssueDate:    time.Now(),
					ActionStatus: true,
				},
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("insert error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = lcr.Create(context.Background(), testCase.args.libCard)

			testCase.expected(t, err)
		})
	}
}

func TestLibCardRepo_GetByReaderID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	lcr := impl.NewLibCardRepo(db, logging.GetLoggerForTests())

	type args struct {
		readerID uuid.UUID
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, libCard *models.LibCardModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success get libCard by readerID",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "reader_id", "lib_card_num", "validity", "issue_date", "action_status"}).
					AddRow(uuid.New(), args.readerID, "1234567890123", 12, time.Now(), true)

				mock.ExpectQuery(`SELECT (.+) FROM bs.lib_card WHERE (.+)`).
					WithArgs(args.readerID).WillReturnRows(rows)
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, libCard)
				assert.Equal(t, "1234567890123", libCard.LibCardNum)
				assert.Equal(t, 12, libCard.Validity)
				assert.True(t, libCard.ActionStatus)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error no rows found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.lib_card WHERE (.+)`).
					WithArgs(args.readerID).WillReturnError(sql.ErrNoRows)
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Nil(t, libCard)
				assert.Equal(t, errsRepo.ErrNotFound, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error query execution",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.lib_card WHERE (.+)`).
					WithArgs(args.readerID).WillReturnError(errors.New("query error"))
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Error(t, err)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
				assert.Nil(t, libCard)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var libCard *models.LibCardModel
			libCard, err = lcr.GetByReaderID(context.Background(), testCase.args.readerID)

			testCase.expected(t, libCard, err)
		})
	}
}

func TestLibCardRepo_GetByNum(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	lcr := impl.NewLibCardRepo(db, logging.GetLoggerForTests())

	type args struct {
		libCardNum string
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, libCard *models.LibCardModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success get libCard by num",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "reader_id", "lib_card_num", "validity", "issue_date", "action_status"}).
					AddRow(uuid.New(), uuid.New(), args.libCardNum, 12, time.Now(), true)

				mock.ExpectQuery(`SELECT (.+) FROM bs.lib_card WHERE (.+)`).
					WithArgs(args.libCardNum).WillReturnRows(rows)
			},
			args: args{
				libCardNum: "1234567890123",
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, libCard)
				assert.Equal(t, "1234567890123", libCard.LibCardNum)
				assert.Equal(t, 12, libCard.Validity)
				assert.True(t, libCard.ActionStatus)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error no rows found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.lib_card WHERE (.+)`).
					WithArgs(args.libCardNum).WillReturnError(sql.ErrNoRows)
			},
			args: args{
				libCardNum: "1234567890123",
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Nil(t, libCard)
				assert.Equal(t, errsRepo.ErrNotFound, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: query execution",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.lib_card WHERE (.+)`).
					WithArgs(args.libCardNum).
					WillReturnError(errors.New("query error"))
			},
			args: args{
				libCardNum: "1234567890123",
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Error(t, err)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
				assert.Nil(t, libCard)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var libCard *models.LibCardModel
			libCard, err = lcr.GetByNum(context.Background(), testCase.args.libCardNum)

			testCase.expected(t, libCard, err)
		})
	}
}

func TestLibCardRepo_Update(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	lcr := impl.NewLibCardRepo(db, logging.GetLoggerForTests())

	type args struct {
		libCard *models.LibCardModel
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success update lib card",
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE bs.lib_card SET (.+) WHERE (.+)`).
					WithArgs(args.libCard.IssueDate, args.libCard.ID).
					WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					IssueDate:    time.Now(),
					ReaderID:     uuid.New(),
					LibCardNum:   "1234567890123",
					Validity:     365,
					ActionStatus: true,
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
				mock.ExpectExec(`UPDATE bs.lib_card SET (.+) WHERE (.+)`).
					WithArgs(args.libCard.IssueDate, args.libCard.ID).
					WillReturnError(errors.New("update error"))
			},
			args: args{
				libCard: &models.LibCardModel{
					ID:           uuid.New(),
					IssueDate:    time.Now(),
					ReaderID:     uuid.New(),
					LibCardNum:   "1234567890123",
					Validity:     365,
					ActionStatus: true,
				},
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("update error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = lcr.Update(context.Background(), testCase.args.libCard)

			testCase.expected(t, err)
		})
	}
}

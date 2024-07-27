package repositoryTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/implRepo/postgres"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	lcr := postgres.NewLibCardRepo(db)

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
			name: "Success: create lib card",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO lib_card VALUES`).
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
			name: "Error: executing query",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO lib_card VALUES`).
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
				expectedError := fmt.Errorf("error inserting libCard: insert error")
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

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	lcr := postgres.NewLibCardRepo(db)

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
			name: "Success: get lib card by reader ID",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "readerid", "libcardnum", "validity", "issuedate", "actionstatus"}).
					AddRow(uuid.New(), args.readerID, "1234567890123", 12, time.Now(), true)

				mock.ExpectQuery(`SELECT (.+) FROM lib_card WHERE (.+)`).
					WithArgs(args.readerID).
					WillReturnRows(rows)
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
			name: "Error: no rows found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM lib_card WHERE (.+)`).
					WithArgs(args.readerID).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Nil(t, libCard)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: query execution",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM lib_card WHERE (.+)`).
					WithArgs(args.readerID).
					WillReturnError(errors.New("query error"))
			},
			args: args{
				readerID: uuid.New(),
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Error(t, err)
				expectedError := errors.New("error retrieving lib card by reader ID: query error")
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

			libCard, err := lcr.GetByReaderID(context.Background(), testCase.args.readerID)

			testCase.expected(t, libCard, err)
		})
	}
}

func TestLibCardRepo_GetByNum(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	lcr := postgres.NewLibCardRepo(db)

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
			name: "Success: get lib card by num",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "readerid", "libcardnum", "validity", "issuedate", "actionstatus"}).
					AddRow(uuid.New(), uuid.New(), args.libCardNum, 12, time.Now(), true)

				mock.ExpectQuery(`SELECT (.+) FROM lib_card WHERE (.+)`).
					WithArgs(args.libCardNum).
					WillReturnRows(rows)
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
			name: "Error: no rows found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM lib_card WHERE (.+)`).
					WithArgs(args.libCardNum).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				libCardNum: "1234567890123",
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Nil(t, libCard)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: query execution",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM lib_card WHERE (.+)`).
					WithArgs(args.libCardNum).
					WillReturnError(errors.New("query error"))
			},
			args: args{
				libCardNum: "1234567890123",
			},
			expected: func(t *testing.T, libCard *models.LibCardModel, err error) {
				assert.Error(t, err)
				expectedError := errors.New("error retrieving lib card by num: query error")
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

			libCard, err := lcr.GetByNum(context.Background(), testCase.args.libCardNum)

			testCase.expected(t, libCard, err)
		})
	}
}

func TestLibCardRepo_Update(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	lcr := postgres.NewLibCardRepo(db)

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
			name: "Success: update lib card",
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE lib_card SET (.+) WHERE (.+)`).
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
			name: "Error: executing query",
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE lib_card SET (.+) WHERE (.+)`).
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
				expectedError := "[!] ERROR! Error updating lib card: update error"
				assert.Equal(t, expectedError, err.Error())
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

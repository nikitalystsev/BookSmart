package repositoryTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/implRepo/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestReaderRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	rr := postgres.NewReaderRepo(db, nil)

	type args struct {
		reader *models.ReaderModel
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
			name: "Success: create reader",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO reader VALUES`).
					WithArgs(args.reader.ID, args.reader.Fio, args.reader.PhoneNumber, args.reader.Age, args.reader.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "Test Reader",
					PhoneNumber: "1234567890",
					Age:         25,
					Password:    "password",
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
				mock.ExpectExec(`INSERT INTO reader VALUES`).
					WithArgs(args.reader.ID, args.reader.Fio, args.reader.PhoneNumber, args.reader.Age, args.reader.Password).
					WillReturnError(errors.New("insert error"))
			},
			args: args{
				reader: &models.ReaderModel{
					ID:          uuid.New(),
					Fio:         "Test Reader",
					PhoneNumber: "1234567890",
					Age:         25,
					Password:    "password",
				},
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("error inserting reader: insert error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = rr.Create(context.Background(), testCase.args.reader)

			testCase.expected(t, err)
		})
	}
}

func TestReaderRepo_GetByPhoneNumber(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	rr := postgres.NewReaderRepo(db, nil)

	type args struct {
		phoneNumber string
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, reader *models.ReaderModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success: get reader by phone number",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "fio", "phonenumber", "age", "password"}).
					AddRow(uuid.New(), "Test Reader", args.phoneNumber, 25, "password")

				mock.ExpectQuery(`SELECT (.+) FROM reader WHERE (.+)`).
					WithArgs(args.phoneNumber).WillReturnRows(rows)
			},
			args: args{
				phoneNumber: "1234567890",
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				assert.Equal(t, "Test Reader", reader.Fio)
				assert.Equal(t, "1234567890", reader.PhoneNumber)
				assert.Equal(t, uint(25), reader.Age)
				assert.Equal(t, "password", reader.Password)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: reader not found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM reader WHERE (.+)`).
					WithArgs(args.phoneNumber).WillReturnError(errors.New("sql: no rows in result set"))
			},
			args: args{
				phoneNumber: "1234567890",
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reader)
				expectedError := fmt.Errorf("error fetching reader by phone number: sql: no rows in result set")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			reader, err := rr.GetByPhoneNumber(context.Background(), testCase.args.phoneNumber)

			testCase.expected(t, reader, err)
		})
	}
}

func TestReaderRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	rr := postgres.NewReaderRepo(db, nil)

	type args struct {
		id uuid.UUID
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, reader *models.ReaderModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success: get reader by phone number",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "fio", "phonenumber", "age", "password"}).
					AddRow(uuid.New(), "Test Reader", "1234567890", 25, "password")

				mock.ExpectQuery(`SELECT (.+) FROM reader WHERE (.+)`).
					WithArgs(args.id).WillReturnRows(rows)
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				assert.Equal(t, "Test Reader", reader.Fio)
				assert.Equal(t, "1234567890", reader.PhoneNumber)
				assert.Equal(t, uint(25), reader.Age)
				assert.Equal(t, "password", reader.Password)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: reader not found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM reader WHERE (.+)`).
					WithArgs(args.id).WillReturnError(errors.New("sql: no rows in result set"))
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, reader)
				expectedError := fmt.Errorf("error fetching reader by phone number: sql: no rows in result set")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			reader, err := rr.GetByID(context.Background(), testCase.args.id)

			testCase.expected(t, reader, err)
		})
	}
}

func TestReaderRepo_IsFavorite(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	rr := postgres.NewReaderRepo(db, nil)

	type args struct {
		readerID uuid.UUID
		bookID   uuid.UUID
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, isFavorite bool, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success: book is favorite",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

				mock.ExpectQuery(`SELECT (.+) FROM favorite_books WHERE (.+)`).
					WithArgs(args.readerID, args.bookID).WillReturnRows(rows)
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			expected: func(t *testing.T, isFavorite bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isFavorite)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Success: book is not favorite",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

				mock.ExpectQuery(`SELECT (.+) FROM favorite_books WHERE (.+)`).
					WithArgs(args.readerID, args.bookID).WillReturnRows(rows)
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			expected: func(t *testing.T, isFavorite bool, err error) {
				assert.NoError(t, err)
				assert.False(t, isFavorite)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: query execution",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM favorite_books WHERE (.+)`).
					WithArgs(args.readerID, args.bookID).WillReturnError(errors.New("query error"))
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			expected: func(t *testing.T, isFavorite bool, err error) {
				assert.Error(t, err)
				assert.False(t, isFavorite)
				expectedError := fmt.Errorf("error checking if book is favorite: query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			isFavorite, err := rr.IsFavorite(context.Background(), testCase.args.readerID, testCase.args.bookID)

			testCase.expected(t, isFavorite, err)
		})
	}
}

func TestReaderRepo_AddToFavorites(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	rr := postgres.NewReaderRepo(db, nil)

	type args struct {
		readerID uuid.UUID
		bookID   uuid.UUID
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
			name: "Success: add book to favorites",
			mockBehavior: func(args args) {
				mock.ExpectExec(`INSERT INTO favorite_books`).
					WithArgs(args.readerID, args.bookID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
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
				mock.ExpectExec(`INSERT INTO favorite_books`).
					WithArgs(args.readerID, args.bookID).WillReturnError(errors.New("insert error"))
			},
			args: args{
				readerID: uuid.New(),
				bookID:   uuid.New(),
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := fmt.Errorf("error adding book to favorites: insert error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = rr.AddToFavorites(context.Background(), testCase.args.readerID, testCase.args.bookID)

			testCase.expected(t, err)
		})
	}
}

func TestReaderRepo_SaveRefreshToken(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when running miniredis", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	rr := postgres.NewReaderRepo(nil, rdb)

	type args struct {
		id    uuid.UUID
		token string
		ttl   time.Duration
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, args args, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success: save refresh token",
			mockBehavior: func(args args) {
				// No special behavior for success case
			},
			args: args{
				id:    uuid.New(),
				token: "test_refresh_token",
				ttl:   time.Hour,
			},
			expected: func(t *testing.T, args args, err error) {

				curErr := rdb.Set(context.Background(), args.token, args.id.String(), args.ttl).Err()

				if got, err := mr.Get(args.token); err != nil || got != args.id.String() {
					t.Errorf("'%s' has the wrong value", args.token)
				}

				val, _ := rdb.Get(context.Background(), args.token).Result()

				assert.Equal(t, curErr, err)
				assert.Equal(t, args.id.String(), val)
			},
		},
		{
			name: "Error: set token fails",
			mockBehavior: func(args args) {
				mr.Close() // force an error by closing the miniredis instance
			},
			args: args{
				id:    uuid.New(),
				token: "test_refresh_token",
				ttl:   time.Hour,
			},
			expected: func(t *testing.T, args args, err error) {
				curErr := rdb.Set(context.Background(), args.token, args.id.String(), 0).Err()

				expectedError := fmt.Errorf("error saving refresh token: %v", curErr)
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err := rr.SaveRefreshToken(context.Background(), testCase.args.id, testCase.args.token, testCase.args.ttl)
			testCase.expected(t, testCase.args, err)
		})
	}
}

func TestReaderRepo_GetByRefreshToken(t *testing.T) {
	// Initialize miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when running miniredis", err)
	}
	defer mr.Close()

	// Initialize go-redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Initialize sqlmock
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing db connection %v", err)
		}
	}(db)

	rr := postgres.NewReaderRepo(db, rdb)

	type args struct {
		token string
	}

	type expectedFunc func(t *testing.T, reader *models.ReaderModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior func(args args, readerID uuid.UUID)
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success: get reader by refresh token",
			mockBehavior: func(args args, readerID uuid.UUID) {
				// Set token in miniredis with readerID
				mr.Set(args.token, readerID.String())

				// Mock database query
				rows := sqlmock.NewRows([]string{"id", "fio", "phonenumber", "age", "password"}).
					AddRow(readerID, "Test Reader", "1234567890", 30, "password")
				mock.ExpectQuery(`SELECT (.+) FROM reader WHERE (.+)`).
					WithArgs(readerID).WillReturnRows(rows)
			},
			args: args{
				token: "test_refresh_token",
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				assert.Equal(t, "Test Reader", reader.Fio)
				assert.Equal(t, "1234567890", reader.PhoneNumber)
				assert.Equal(t, uint(30), reader.Age)
				assert.Equal(t, "password", reader.Password)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error: refresh token not found",
			mockBehavior: func(args args, readerID uuid.UUID) {
				// No need to set anything in miniredis
			},
			args: args{
				token: "invalid_refresh_token",
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.Error(t, err)
				expectedError := fmt.Errorf("refresh token not found")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			name: "Error: invalid reader ID",
			mockBehavior: func(args args, readerID uuid.UUID) {
				// Set invalid readerID in miniredis
				mr.Set(args.token, "invalid_uuid")
			},
			args: args{
				token: "test_refresh_token",
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.Error(t, err)
				expectedError := fmt.Errorf("invalid reader ID: %s", "invalid UUID length: 12")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
		{
			name: "Error: retrieving reader from database",
			mockBehavior: func(args args, readerID uuid.UUID) {
				// Set token in miniredis with readerID
				mr.Set(args.token, readerID.String())

				// Mock database query error
				mock.ExpectQuery(`SELECT (.+) FROM reader WHERE (.+)`).
					WithArgs(readerID).WillReturnError(errors.New("database error"))
			},
			args: args{
				token: "test_refresh_token",
			},
			expected: func(t *testing.T, reader *models.ReaderModel, err error) {
				assert.Error(t, err)
				expectedError := fmt.Errorf("error retrieving reader: %w", errors.New("database error"))
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			readerID := uuid.New()
			testCase.mockBehavior(testCase.args, readerID)

			reader, err := rr.GetByRefreshToken(context.Background(), testCase.args.token)

			testCase.expected(t, reader, err)
		})
	}
}

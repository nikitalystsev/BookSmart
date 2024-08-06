package repositoryTests

import (
	"BookSmart-postgres/impl"
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"Booksmart/pkg/logging"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestBookRepo_Create(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	br := impl.NewBookRepo(db, logging.GetLoggerForTests())

	type args struct {
		book *models.BookModel
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
			name: "Success insert book",
			mockBehavior: func(args args) {

				mock.ExpectExec(`INSERT INTO bs.book VALUES`).
					WithArgs(args.book.ID, args.book.Title, args.book.Author, args.book.Publisher,
						args.book.CopiesNumber, args.book.Rarity, args.book.Genre,
						args.book.PublishingYear, args.book.Language, args.book.AgeLimit).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				book: &models.BookModel{
					ID:             uuid.New(),
					Title:          "Test Book",
					Author:         "Test Author",
					Publisher:      "Test Publisher",
					CopiesNumber:   10,
					Rarity:         "Common",
					Genre:          "Fiction",
					PublishingYear: 2021,
					Language:       "English",
					AgeLimit:       12,
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

				mock.ExpectExec(`INSERT INTO bs.book VALUES`).
					WithArgs(args.book.ID, args.book.Title, args.book.Author, args.book.Publisher,
						args.book.CopiesNumber, args.book.Rarity, args.book.Genre,
						args.book.PublishingYear, args.book.Language, args.book.AgeLimit).
					WillReturnError(errors.New("insert error"))
			},
			args: args{
				book: &models.BookModel{
					ID:             uuid.New(),
					Title:          "Test Book",
					Author:         "Test Author",
					Publisher:      "Test Publisher",
					CopiesNumber:   10,
					Rarity:         "Common",
					Genre:          "Fiction",
					PublishingYear: 2021,
					Language:       "English",
					AgeLimit:       12,
				},
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("insert error")
				assert.Equal(t, expectedError, err)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = br.Create(context.Background(), testCase.args.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookRepo_GetByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	br := impl.NewBookRepo(db, logging.GetLoggerForTests())

	type args struct {
		id uuid.UUID
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, book *models.BookModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success get book by ID",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"}).
					AddRow(args.id, "Test Book", "Test Author", "Test Publisher", "10", "Common", "Fiction", 2021, "English", 12)

				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.id).WillReturnRows(rows)
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, book *models.BookModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, book)
				assert.Equal(t, "Test Book", book.Title)
				assert.Equal(t, "Test Author", book.Author)
				assert.Equal(t, "Test Publisher", book.Publisher)
				assert.Equal(t, uint(10), book.CopiesNumber)
				assert.Equal(t, "Common", book.Rarity)
				assert.Equal(t, "Fiction", book.Genre)
				assert.Equal(t, uint(2021), book.PublishingYear)
				assert.Equal(t, "English", book.Language)
				assert.Equal(t, uint(12), book.AgeLimit)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error book not found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.id).WillReturnError(sql.ErrNoRows)
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, book *models.BookModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, book)
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: executing query",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.id).
					WillReturnError(errors.New("query error"))
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, book *models.BookModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, book)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args)

			var book *models.BookModel
			book, err = br.GetByID(context.Background(), testCase.args.id)

			testCase.expected(t, book, err)
		})
	}
}

func TestBookRepo_GetByTitle(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	br := impl.NewBookRepo(db, logging.GetLoggerForTests())

	type args struct {
		title string
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, book *models.BookModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success get book by title",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"}).
					AddRow(uuid.New(), args.title, "Test Author", "Test Publisher", "10", "Common", "Fiction", 2021, "English", 12)

				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.title).WillReturnRows(rows)
			},
			args: args{
				title: "Test Book",
			},
			expected: func(t *testing.T, book *models.BookModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, book)
				assert.Equal(t, "Test Book", book.Title)
				assert.Equal(t, "Test Author", book.Author)
				assert.Equal(t, "Test Publisher", book.Publisher)
				assert.Equal(t, uint(10), book.CopiesNumber)
				assert.Equal(t, "Common", book.Rarity)
				assert.Equal(t, "Fiction", book.Genre)
				assert.Equal(t, uint(2021), book.PublishingYear)
				assert.Equal(t, "English", book.Language)
				assert.Equal(t, uint(12), book.AgeLimit)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error book not found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.title).WillReturnError(sql.ErrNoRows)
			},
			args: args{
				title: "Test Book",
			},
			expected: func(t *testing.T, book *models.BookModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, book)
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error executing query",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.title).
					WillReturnError(errors.New("query error"))
			},
			args: args{
				title: "Test Book",
			},
			expected: func(t *testing.T, book *models.BookModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, book)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {

			testCase.mockBehavior(testCase.args)

			var book *models.BookModel
			book, err = br.GetByTitle(context.Background(), testCase.args.title)

			testCase.expected(t, book, err)
		})
	}
}

func TestBookRepo_Delete(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	br := impl.NewBookRepo(db, logging.GetLoggerForTests())

	type args struct {
		id uuid.UUID
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
			name: "Success delete book by ID",
			mockBehavior: func(args args) {
				mock.ExpectExec(`DELETE FROM bs.book WHERE (.+)`).
					WithArgs(args.id).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				id: uuid.New(),
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
				mock.ExpectExec(`DELETE FROM bs.book WHERE (.+)`).
					WithArgs(args.id).
					WillReturnError(errors.New("delete error"))
			},
			args: args{
				id: uuid.New(),
			},
			expected: func(t *testing.T, err error) {
				assert.Error(t, err)
				expectedError := errors.New("delete error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err = br.Delete(context.Background(), testCase.args.id)

			testCase.expected(t, err)
		})
	}
}

func TestBookRepo_Update(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	br := impl.NewBookRepo(db, logging.GetLoggerForTests())

	type args struct {
		book *models.BookModel
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
			name: "Success update book copies",
			mockBehavior: func(args args) {
				mock.ExpectExec(`UPDATE bs.book SET (.+) WHERE (.+)`).
					WithArgs(args.book.CopiesNumber, args.book.ID).
					WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					CopiesNumber: 5,
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
				mock.ExpectExec(`UPDATE bs.book SET (.+) WHERE (.+)`).
					WithArgs(args.book.CopiesNumber, args.book.ID).
					WillReturnError(errors.New("update error"))
			},
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					CopiesNumber: 5,
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

			err = br.Update(context.Background(), testCase.args.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookRepo_GetByParams(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	br := impl.NewBookRepo(db, logging.GetLoggerForTests())

	type args struct {
		params *dto.BookParamsDTO
	}

	type mockBehavior func(args args)
	type expectedFunc func(t *testing.T, books []*models.BookModel, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success get books by params",
			mockBehavior: func(args args) {
				rows := sqlxmock.NewRows([]string{"id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"}).
					AddRow(uuid.New(), "Test Book", "Test Author", "Test Publisher", 10, "Common", "Fiction", 2021, "English", 12)

				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.params.Title, args.params.Author, args.params.Publisher, args.params.CopiesNumber, args.params.Rarity, args.params.Genre, args.params.PublishingYear, args.params.Language, args.params.AgeLimit, args.params.Limit, args.params.Offset).
					WillReturnRows(rows)
			},
			args: args{
				params: &dto.BookParamsDTO{
					Title:          "Test Book",
					Author:         "Test Author",
					Publisher:      "Test Publisher",
					CopiesNumber:   10,
					Rarity:         "Common",
					Genre:          "Fiction",
					PublishingYear: 2021,
					Language:       "English",
					AgeLimit:       12,
					Limit:          10,
					Offset:         0,
				},
			},
			expected: func(t *testing.T, books []*models.BookModel, err error) {
				assert.NoError(t, err)
				assert.Len(t, books, 1)
				assert.Equal(t, "Test Book", books[0].Title)
				assert.Equal(t, "Test Author", books[0].Author)
				assert.Equal(t, "Test Publisher", books[0].Publisher)
				assert.Equal(t, uint(10), books[0].CopiesNumber)
				assert.Equal(t, "Common", books[0].Rarity)
				assert.Equal(t, "Fiction", books[0].Genre)
				assert.Equal(t, uint(2021), books[0].PublishingYear)
				assert.Equal(t, "English", books[0].Language)
				assert.Equal(t, uint(12), books[0].AgeLimit)
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error executing query",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.params.Title, args.params.Author, args.params.Publisher, args.params.CopiesNumber, args.params.Rarity, args.params.Genre, args.params.PublishingYear, args.params.Language, args.params.AgeLimit, args.params.Limit, args.params.Offset).
					WillReturnError(errors.New("query error"))
			},
			args: args{
				params: &dto.BookParamsDTO{
					Title:          "Test Book",
					Author:         "Test Author",
					Publisher:      "Test Publisher",
					CopiesNumber:   10,
					Rarity:         "Common",
					Genre:          "Fiction",
					PublishingYear: 2021,
					Language:       "English",
					AgeLimit:       12,
					Limit:          10,
					Offset:         0,
				},
			},
			expected: func(t *testing.T, books []*models.BookModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, books)
				expectedError := errors.New("query error")
				assert.Equal(t, expectedError.Error(), err.Error())
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
		{
			name: "Error books not found",
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM bs.book WHERE (.+)`).
					WithArgs(args.params.Title, args.params.Author, args.params.Publisher, args.params.CopiesNumber,
						args.params.Rarity, args.params.Genre, args.params.PublishingYear,
						args.params.Language, args.params.AgeLimit, args.params.Limit, args.params.Offset).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				params: &dto.BookParamsDTO{
					Title:          "Test Book",
					Author:         "Test Author",
					Publisher:      "Test Publisher",
					CopiesNumber:   10,
					Rarity:         "Common",
					Genre:          "Fiction",
					PublishingYear: 2021,
					Language:       "English",
					AgeLimit:       12,
					Limit:          10,
					Offset:         0,
				},
			},
			expected: func(t *testing.T, books []*models.BookModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, books)
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError.Error(), err.Error())
				err = mock.ExpectationsWereMet()
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			var books []*models.BookModel
			books, err = br.GetByParams(context.Background(), testCase.args.params)

			testCase.expected(t, books, err)
		})
	}
}

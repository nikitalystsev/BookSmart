package repositoryTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/implRepo/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestBookRepo_Create(t *testing.T) {
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

	br := postgres.NewBookRepo(db)

	type args struct {
		book *models.BookModel
	}

	type mockBehavior func(args args, id uuid.UUID)
	type expectedFunc func(t *testing.T, err error)

	testsTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expected     expectedFunc
	}{
		{
			name: "Success: create book",
			mockBehavior: func(args args, id uuid.UUID) {

				mock.ExpectExec(`INSERT INTO book VALUES`).
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
			name: "Error: executing query",
			mockBehavior: func(args args, id uuid.UUID) {

				mock.ExpectExec(`INSERT INTO book VALUES`).
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
				expectedError := errors.New("error inserting book: insert error")
				assert.Equal(t, expectedError.Error(), err.Error())
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.book.ID)

			err = br.Create(context.Background(), testCase.args.book)

			testCase.expected(t, err)
		})
	}
}

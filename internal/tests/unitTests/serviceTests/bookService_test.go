package serviceTests

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBookService_Create(t *testing.T) {
	type args struct {
		book *models.BookModel
	}
	type mockBehaviour func(m *mockrepo.MockIBookRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testsTable := []struct {
		name         string
		args         args
		mockBehavior mockBehaviour
		expected     expectedFunc
	}{
		{
			name: "Success successful creation",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
				m.EXPECT().Create(gomock.Any(), args.book).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error checking book existence",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "Common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book already exists",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "Common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(args.book, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookAlreadyExist
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error empty book title",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrEmptyBookTitle
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error empty book author",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrEmptyBookAuthor
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error empty book rarity",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrEmptyBookRarity
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error invalid book copies number",
			args: args{
				book: &models.BookModel{
					ID:     uuid.New(),
					Title:  "The Great Gatsby",
					Author: "F. Scott Fitzgerald",
					Rarity: "common",
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrInvalidBookCopiesNum
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error creating book",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
				m.EXPECT().Create(gomock.Any(), args.book).Return(errors.New("create error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("create error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args)

			err := bookService.Create(context.Background(), testCase.args.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_Delete(t *testing.T) {
	type args struct {
		book *models.BookModel
	}
	type mockBehaviour func(m *mockrepo.MockIBookRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehaviour
		expected     expectedFunc
	}{
		{
			name: "Success successfully delete book",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(args.book, nil)
				m.EXPECT().Delete(gomock.Any(), args.book.ID).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error getting book",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book not found",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error error deleting book",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(args.book, nil)
				m.EXPECT().Delete(gomock.Any(), args.book.ID).Return(fmt.Errorf("delete error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("delete error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args)

			err := bookService.Delete(context.Background(), testCase.args.book.ID)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_GetByID(t *testing.T) {
	type args struct {
		book *models.BookModel
	}
	type mockBehaviour func(m *mockrepo.MockIBookRepo, args args)
	type expectedFunc func(t *testing.T, err error)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehaviour
		expected     expectedFunc
	}{
		{
			name: "Success successful get book",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(args.book, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error getting book",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errors.New("database error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error book not found",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, args args) {
				m.EXPECT().GetByID(gomock.Any(), args.book.ID).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args)

			_, err := bookService.GetByID(context.Background(), testCase.args.book.ID)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_GetByParams(t *testing.T) {
	type mockBehavior func(m *mockrepo.MockIBookRepo, params *dto.BookParamsDTO)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		bookDTO *dto.BookParamsDTO
	}

	testBooks := []*models.BookModel{
		{
			ID:           uuid.New(),
			Title:        "The Great Gatsby",
			Author:       "F. Scott Fitzgerald",
			Rarity:       "common",
			CopiesNumber: 10,
		},
		{
			ID:           uuid.New(),
			Title:        "1984",
			Author:       "George Orwell",
			Rarity:       "common",
			CopiesNumber: 8,
		},
		{
			ID:           uuid.New(),
			Title:        "To Kill a Mockingbird",
			Author:       "Harper Lee",
			Rarity:       "common",
			CopiesNumber: 5,
		},
	}

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success successful get books",
			args: args{
				bookDTO: &dto.BookParamsDTO{
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       "common",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, params *dto.BookParamsDTO) {
				m.EXPECT().GetByParams(gomock.Any(), params).Return([]*models.BookModel{testBooks[0]}, nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error error on getting books",
			args: args{
				bookDTO: &dto.BookParamsDTO{
					Title:        "Non-existent Book",
					Author:       "Unknown Author",
					Rarity:       "common",
					CopiesNumber: 0,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, params *dto.BookParamsDTO) {
				m.EXPECT().GetByParams(gomock.Any(), params).Return(nil, fmt.Errorf("error getting books"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := fmt.Errorf("error getting books")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error books not found",
			args: args{
				bookDTO: &dto.BookParamsDTO{
					Title:        "Non-existent Book",
					Author:       "Unknown Author",
					Rarity:       "common",
					CopiesNumber: 0,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, params *dto.BookParamsDTO) {
				m.EXPECT().GetByParams(gomock.Any(), params).Return(nil, errs.ErrBookDoesNotExists)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errs.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args.bookDTO)

			_, err := bookService.GetByParams(context.Background(), testCase.args.bookDTO)
			testCase.expected(t, err)
		})
	}
}

package serviceTests

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/services/errsService"
	"BookSmart/internal/services/implServices"
	mockrepo "BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"BookSmart/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBookService_Create(t *testing.T) {
	type mockBehaviour func(m *mockrepo.MockIBookRepo, book *models.BookModel)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		book *models.BookModel
	}

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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), book).Return(nil)
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrBookAlreadyExist
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error empty book title",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Author:       "F. Scott Fitzgerald",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrEmptyBookTitle
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error empty book author",
			args: args{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrEmptyBookAuthor
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
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrEmptyBookRarity
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
					Rarity: implServices.BookRarityCommon,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrInvalidBookCopiesNum
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), book).Return(errors.New("create error"))
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
			bookService := implServices.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args.book)

			err := bookService.Create(context.Background(), testCase.args.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_Delete(t *testing.T) {
	type mockBehaviour func(m *mockrepo.MockIBookRepo, book *models.BookModel)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		book *models.BookModel
	}

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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
				m.EXPECT().Delete(gomock.Any(), book.ID).Return(nil)
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrBookDoesNotExists
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
				m.EXPECT().Delete(gomock.Any(), book.ID).Return(fmt.Errorf("delete error"))
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
			bookService := implServices.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args.book)

			err := bookService.Delete(context.Background(), testCase.args.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_GetByID(t *testing.T) {
	type mockBehaviour func(m *mockrepo.MockIBookRepo, book *models.BookModel)
	type expectedFunc func(t *testing.T, err error)
	type args struct {
		book *models.BookModel
	}

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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsService.ErrBookDoesNotExists
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := implServices.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args.book)

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
			Rarity:       implServices.BookRarityCommon,
			CopiesNumber: 10,
		},
		{
			ID:           uuid.New(),
			Title:        "1984",
			Author:       "George Orwell",
			Rarity:       implServices.BookRarityCommon,
			CopiesNumber: 8,
		},
		{
			ID:           uuid.New(),
			Title:        "To Kill a Mockingbird",
			Author:       "Harper Lee",
			Rarity:       implServices.BookRarityCommon,
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
					Rarity:       implServices.BookRarityCommon,
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
					Rarity:       implServices.BookRarityCommon,
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
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 0,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, params *dto.BookParamsDTO) {
				m.EXPECT().GetByParams(gomock.Any(), params).Return(nil, errsRepo.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errsRepo.ErrNotFound
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := implServices.NewBookService(mockBookRepo, logging.GetLoggerForTests())

			testCase.mockBehavior(mockBookRepo, testCase.args.bookDTO)

			_, err := bookService.GetByParams(context.Background(), testCase.args.bookDTO)
			testCase.expected(t, err)
		})
	}
}

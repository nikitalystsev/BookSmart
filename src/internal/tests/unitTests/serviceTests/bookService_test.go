package serviceTests

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/services/implServices"
	mockrepo "BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO тесты типа написаны, но их бы как то нужно глянуть, а то они пустые

func TestBookService_Create(t *testing.T) {
	type mockBehaviour func(m *mockrepo.MockIBookRepo, book *models.BookModel)
	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		book *models.BookModel
	}

	testsTable := []struct {
		name         string
		input        inputStruct
		mockBehavior mockBehaviour
		expected     expectedFunc
	}{
		{
			name: "Success: successful creation",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), book).Return(nil)
			},
			expected: func(t *testing.T, err error) {
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "Error: error checking book existence",
			input: inputStruct{
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
				expectedError := errors.New("[!] ERROR! Error checking book existence: database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: book already exists",
			input: inputStruct{
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
				expectedError := errors.New("[!] ERROR! Book with this title already exists")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: empty book title",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Author:       "F. Scott Fitzgerald",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Empty book title")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: empty book author",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Empty book author")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: empty book rarity",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Empty book rarity")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: empty book copies number",
			input: inputStruct{
				book: &models.BookModel{
					ID:     uuid.New(),
					Title:  "The Great Gatsby",
					Author: "F. Scott Fitzgerald",
					Rarity: implServices.BookRarityCommon,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Invalid book copies number")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error creating book",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
				m.EXPECT().Create(gomock.Any(), book).Return(errors.New("create error"))
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Error creating book: create error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testsTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := implServices.NewBookService(mockBookRepo)

			testCase.mockBehavior(mockBookRepo, testCase.input.book)

			err := bookService.Create(context.Background(), testCase.input.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_Delete(t *testing.T) {
	type mockBehaviour func(m *mockrepo.MockIBookRepo, book *models.BookModel)
	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		book *models.BookModel
	}

	testTable := []struct {
		name         string
		input        inputStruct
		mockBehavior mockBehaviour
		expected     expectedFunc
	}{
		{
			name: "Success: successfully delete book",
			input: inputStruct{
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
			name: "Error: error getting book",
			input: inputStruct{
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
				expectedError := errors.New("[!] ERROR! Error checking book existence: database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: book not found",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Book with this title does not exist")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: error deleting book",
			input: inputStruct{
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
				expectedError := errors.New("[!] ERROR! Error deleting book: delete error")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := implServices.NewBookService(mockBookRepo)

			testCase.mockBehavior(mockBookRepo, testCase.input.book)

			err := bookService.Delete(context.Background(), testCase.input.book)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_GetByID(t *testing.T) {
	type mockBehaviour func(m *mockrepo.MockIBookRepo, book *models.BookModel)
	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
		book *models.BookModel
	}

	testTable := []struct {
		name         string
		input        inputStruct
		mockBehavior mockBehaviour
		expected     expectedFunc
	}{
		{
			name: "Success: successful get book",
			input: inputStruct{
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
			name: "Error: error getting book",
			input: inputStruct{
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
				expectedError := errors.New("[!] ERROR! Error retrieving book information: database error")
				assert.Equal(t, expectedError, err)
			},
		},
		{
			name: "Error: book not found",
			input: inputStruct{
				book: &models.BookModel{
					ID:           uuid.New(),
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Rarity:       implServices.BookRarityCommon,
					CopiesNumber: 10,
				},
			},
			mockBehavior: func(m *mockrepo.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrNotFound)
			},
			expected: func(t *testing.T, err error) {
				expectedError := errors.New("[!] ERROR! Book with this ID does not exist")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := implServices.NewBookService(mockBookRepo)

			testCase.mockBehavior(mockBookRepo, testCase.input.book)

			_, err := bookService.GetByID(context.Background(), testCase.input.book.ID)

			testCase.expected(t, err)
		})
	}
}

func TestBookService_GetByParams(t *testing.T) {
	type mockBehavior func(m *mockrepo.MockIBookRepo, params *dto.BookParamsDTO)
	type expectedFunc func(t *testing.T, err error)
	type inputStruct struct {
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
		input        inputStruct
		mockBehavior mockBehavior
		expected     expectedFunc
	}{
		{
			name: "Success: successful get books",
			input: inputStruct{
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
			name: "Error: error on getting books",
			input: inputStruct{
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
				expectedError := fmt.Errorf("[!] ERROR! Error searching for books: error getting books")
				assert.Equal(t, expectedError, err)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
			bookService := implServices.NewBookService(mockBookRepo)

			testCase.mockBehavior(mockBookRepo, testCase.input.bookDTO)

			_, err := bookService.GetByParams(context.Background(), testCase.input.bookDTO)
			testCase.expected(t, err)
		})
	}
}

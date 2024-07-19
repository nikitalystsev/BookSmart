package unitTests

import (
	"BookSmart/internal/models"
	"BookSmart/internal/services/impl"
	mockrepositories "BookSmart/internal/tests/unitTests/mocks"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBookService_Create(t *testing.T) {

	type mockBehaviour func(m *mockrepositories.MockIBookRepo, book *models.BookModel)

	testTable := []struct {
		name          string
		mockBehavior  mockBehaviour
		expectedError error
	}{
		{
			name: "successful creation",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(nil, errors.New("[!] ERROR! Object not found"))
				m.EXPECT().Create(gomock.Any(), book).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "book already exists",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(book, nil)
			},
			expectedError: errors.New("[!] ERROR! Book with this title already exists"),
		},
		{
			name: "error checking book existence",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(nil, errors.New("some error"))
			},
			expectedError: errors.New("[!] ERROR! Error checking book existence: some error"),
		},
		{
			name: "error creating book",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(nil, errors.New("[!] ERROR! Object not found"))
				m.EXPECT().Create(gomock.Any(), book).Return(errors.New("create error"))
			},
			expectedError: errors.New("[!] ERROR! Error creating book: create error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepositories.NewMockIBookRepo(ctrl)
			bookService := impl.NewBookService(mockBookRepo)
			newBook := &models.BookModel{
				ID:     uuid.New(),
				Title:  "The Great Gatsby",
				Author: "F. Scott Fitzgerald",
			}

			testCase.mockBehavior(mockBookRepo, newBook)

			err := bookService.Create(context.Background(), newBook)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestBookService_DeleteByTitle(t *testing.T) {

	type mockBehaviour func(m *mockrepositories.MockIBookRepo, book *models.BookModel)

	testTable := []struct {
		name          string
		mockBehavior  mockBehaviour
		expectedError error
	}{
		{
			name: "successfully delete book",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(book, nil)
				m.EXPECT().DeleteByTitle(gomock.Any(), book.Title).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "book not found",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(nil, errors.New("[!] ERROR! Object not found"))
			},
			expectedError: errors.New("[!] ERROR! Book with this title does not exist"),
		},
		{
			name: "error getting book",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(nil, errors.New("some error"))
			},
			expectedError: errors.New("[!] ERROR! Error checking book existence: some error"),
		},
		{
			name: "error deleting book",
			mockBehavior: func(m *mockrepositories.MockIBookRepo, book *models.BookModel) {
				m.EXPECT().GetByTitle(gomock.Any(), book.Title).Return(book, nil)
				m.EXPECT().DeleteByTitle(gomock.Any(), book.Title).Return(fmt.Errorf("delete error"))
			},
			expectedError: errors.New("[!] ERROR! Error deleting book: delete error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockBookRepo := mockrepositories.NewMockIBookRepo(ctrl)
			bookService := impl.NewBookService(mockBookRepo)
			newBook := &models.BookModel{
				ID:     uuid.New(),
				Title:  "The Great Gatsby",
				Author: "F. Scott Fitzgerald",
			}

			testCase.mockBehavior(mockBookRepo, newBook)

			err := bookService.DeleteByTitle(context.Background(), newBook)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

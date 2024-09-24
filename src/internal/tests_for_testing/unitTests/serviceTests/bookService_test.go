package serviceTests

import (
	mockrepo "Booksmart/internal/tests/unitTests/serviceTests/mocks"
	dto2 "Booksmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/dto"
	models2 "Booksmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/models"
	"Booksmart/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Лондонский вариант
*/

func TestBookService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	var err error

	runner.Run(t, "success create book", func(t provider.T) {
		// Arrange
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book := models2.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrBookDoesNotExists)
		mockBookRepo.EXPECT().Create(gomock.Any(), book).Return(nil)

		// Act
		err = bookService.Create(context.Background(), book)

		// Assert
		assert.Nil(t, err)
	})
}

func TestBookService_Create_ErrorCheckBookExistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	var err error

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))

	runner.Run(t, "error check book existence", func(t provider.T) {

		// Act
		err = bookService.Create(context.Background(), book)

	})

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("database error"), err)
}

func TestBookService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	var err error

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
	mockBookRepo.EXPECT().Delete(gomock.Any(), book.ID).Return(nil)

	runner.Run(t, "success delete book", func(t provider.T) {
		// Act
		err = bookService.Delete(context.Background(), book.ID)
	})

	// Assert
	assert.Nil(t, err)
}

func TestBookService_Delete_ErrorGetBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	var err error

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))

	runner.Run(t, "error get book", func(t provider.T) {
		// Act
		err = bookService.Delete(context.Background(), book.ID)
	})

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("database error"), err)
}

func TestBookService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	var (
		err      error
		findBook *models.BookModel
	)

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)

	runner.Run(t, "success get book by id", func(t provider.T) {
		// Act
		findBook, err = bookService.GetByID(context.Background(), book.ID)
	})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, book, findBook)
}

func TestBookService_GetByID_ErrorGetBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	var (
		err      error
		findBook *models.BookModel
	)

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrBookDoesNotExists)

	runner.Run(t, "error get book by id", func(t provider.T) {
		// Act
		findBook, err = bookService.GetByID(context.Background(), book.ID)
	})

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errs.ErrBookDoesNotExists, err)
	assert.Nil(t, findBook)
}

func TestBookService_GetByParams_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	var (
		err      error
		findBook []*models.BookModel
	)

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	params := dto2.NewBookParamsDTOObjectMother().DefaultBookParams()

	mockBookRepo.EXPECT().GetByParams(gomock.Any(), params).Return([]*models.BookModel{book}, nil)

	runner.Run(t, "success get books by params", func(t provider.T) {
		// Act
		findBook, err = bookService.GetByParams(context.Background(), params)
	})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, []*models.BookModel{book}, findBook)
}

func TestBookService_GetByParams_ErrorGetBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	var (
		err      error
		findBook []*models.BookModel
	)

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	params := dto2.NewBookParamsDTOObjectMother().DefaultBookParams()

	mockBookRepo.EXPECT().GetByParams(gomock.Any(), params).Return(nil, fmt.Errorf("database error"))

	runner.Run(t, "error get books by params", func(t provider.T) {
		// Act
		findBook, err = bookService.GetByParams(context.Background(), params)
	})

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("database error"), err)
	assert.Nil(t, findBook)
}

/*
	Классический вариант
*/

func TestBookService_Create_Success_Classic(t *testing.T) {
	var err error
	ctx := context.Background()

	// Arrange
	container, err := getContainerForClassicUnitTests()
	if err != nil {
		t.Fatal(err)
	}
	db, err := applyMigrations(container)
	if err != nil {
		t.Fatal(err)
	}
	bookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
	bookService := impl.NewBookService(bookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()
	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
	}(db)

	defer func() {
		if err = container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()

	runner.Run(t, "success create book", func(t provider.T) {

		// Act
		err = bookService.Create(context.Background(), book)

	})

	// Assert
	assert.Nil(t, err)
}

func TestBookService_Create_ErrorCheckBookExistence_Classic(t *testing.T) {
	var err error
	ctx := context.Background()

	// Arrange
	container, err := getContainerForClassicUnitTests()
	if err != nil {
		t.Fatal(err)
	}
	db, err := applyMigrations(container)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatalf("error closing db: %v", err)
	}
	bookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
	bookService := impl.NewBookService(bookRepo, logging.GetLoggerForTests())
	book := models2.NewBookModelObjectMother().DefaultBook()

	defer func() {
		if err = container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()

	runner.Run(t, "success create book", func(t provider.T) {

		// Act
		err = bookService.Create(context.Background(), book)

	})

	// Assert
	assert.Error(t, err)
}

package serviceTests

import (
	mockrepo "Booksmart/internal/tests/unitTests/serviceTests/mocks"
	"Booksmart/internal/tests_for_testing/unitTests/serviceTests/objectMother"
	"Booksmart/pkg/logging"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Лондонский вариант

func TestBookService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	var err error

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := objectMother.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrBookDoesNotExists)
	mockBookRepo.EXPECT().Create(gomock.Any(), book).Return(nil)

	runner.Run(t, "success create book", func(t provider.T) {

		// Act
		err = bookService.Create(context.Background(), book)

	})

	// Assert
	assert.Nil(t, err)
}

func TestBookService_Create_ErrorCheckBookExistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	var err error

	// Arrange
	mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
	bookService := impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
	book := objectMother.NewBookModelObjectMother().DefaultBook()
	mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))

	runner.Run(t, "error check book existence", func(t provider.T) {

		// Act
		err = bookService.Create(context.Background(), book)

	})

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("database error"), err)
}

// Классический вариант

func TestBookService_Create_Success_Classic(t *testing.T) {
	var err error
	ctx := context.Background()

	// Arrange
	container, db, err := getConnectionsForClassicUnitTests()
	if err != nil {
		t.Fatal(err)
	}
	bookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
	bookService := impl.NewBookService(bookRepo, logging.GetLoggerForTests())
	book := objectMother.NewBookModelObjectMother().DefaultBook()
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
	container, db, err := getConnectionsForClassicUnitTests()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatalf("error closing db: %v", err)
	}
	bookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
	bookService := impl.NewBookService(bookRepo, logging.GetLoggerForTests())
	book := objectMother.NewBookModelObjectMother().DefaultBook()

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

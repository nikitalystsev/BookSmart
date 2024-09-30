package serviceTests

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	omdto "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/dto"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

type BookServiceTestsSuite struct {
	suite.Suite
}

func (bsts *BookServiceTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "book service", "suite", "steps")
}

/*
	Лондонский вариант
*/

func (bsts *BookServiceTestsSuite) Test_Create_Success(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		err         error
	)

	t.Title("Test Create Book Success")
	t.Description("The new book was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrBookDoesNotExists)
		mockBookRepo.EXPECT().Create(gomock.Any(), book).Return(nil)
		sCtx.WithNewParameters("book", book)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (bsts *BookServiceTestsSuite) Test_Create_ErrorCheckBookExistence(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		err         error
	)

	t.Title("Test Create Book Error: check book existence")
	t.Description("The new book was not created successfully due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))
		sCtx.WithNewParameters("book", book)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

func (bsts *BookServiceTestsSuite) Test_Create_ErrorBookAlreadyExists(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		err         error
	)

	t.Title("Test Create Book Error: book already exists")
	t.Description("The new book was not successfully created because a book with the same data already exists")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrBookAlreadyExist, err)
	})
}

func (bsts *BookServiceTestsSuite) Test_Delete_Success(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		err         error
	)

	t.Title("Test Delete Book Success")
	t.Description("The book was successfully deleted")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
		mockBookRepo.EXPECT().Delete(gomock.Any(), book.ID).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Delete(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (bsts *BookServiceTestsSuite) Test_Delete_ErrorGettingBook(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		err         error
	)

	t.Title("Test Delete Book Error: getting book")
	t.Description("The new book was not deleted successfully due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Delete(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

func (bsts *BookServiceTestsSuite) Test_GetByID_Success(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBook    *models.BookModel
		err         error
	)

	t.Title("Test Get Book By ID Success")
	t.Description("The book was successfully retrieved by ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookService.GetByID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(findBook)
		t.Assert().Equal(book, findBook)
	})
}

func (bsts *BookServiceTestsSuite) Test_GetByID_ErrorGettingBook(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBook    *models.BookModel
		err         error
	)

	t.Title("Test Get Book By ID Error: getting book")
	t.Description("The new book was not retrieved successfully due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookService.GetByID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("database error"), err)
		t.Assert().Nil(findBook)
	})
}

func (bsts *BookServiceTestsSuite) Test_GetByParams_Success(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBooks   []*models.BookModel
		params      *dto.BookParamsDTO
		err         error
	)

	t.Title("Test Get Book By Params Success")
	t.Description("Books were successfully received according to the parameters")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		params = omdto.NewBookParamsDTOObjectMother().DefaultBookParams()
		mockBookRepo.EXPECT().GetByParams(gomock.Any(), params).Return([]*models.BookModel{book}, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBooks, err = bookService.GetByParams(context.Background(), params)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(findBooks)
		t.Assert().Equal([]*models.BookModel{book}, findBooks)
	})
}

func (bsts *BookServiceTestsSuite) Test_GetByParams_ErrorGettingBooks(t provider.T) {
	var (
		bookService intf.IBookService
		findBooks   []*models.BookModel
		params      *dto.BookParamsDTO
		err         error
	)

	t.Title("Test Get Book By Params Error: getting book")
	t.Description("Books were not successfully retrieved due to error checking error retrieving from repositories")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		params = omdto.NewBookParamsDTOObjectMother().DefaultBookParams()
		mockBookRepo.EXPECT().GetByParams(gomock.Any(), params).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBooks, err = bookService.GetByParams(context.Background(), params)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Nil(findBooks)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

/*
	Классический вариант
*/

func (bsts *BookServiceTestsSuite) Test_Create_Success_Classic(t provider.T) {
	var (
		bookService intf.IBookService
		container   *postgres.PostgresContainer
		db          *sqlx.DB
		book        *models.BookModel
		err         error
	)

	t.Title("Test Create Book Success Classic")
	t.Description("The new book was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		container, err = getPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = applyMigrations(container)
		if err != nil {
			t.Fatal(err)
		}
		mockBookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		sCtx.WithNewParameters("book", book)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})

	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
	}(db)

	defer func() {
		if err = container.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()
}

func (bsts *BookServiceTestsSuite) Test_Create_ErrorCheckBookExistence_Classic(t provider.T) {
	var (
		bookService intf.IBookService
		container   *postgres.PostgresContainer
		db          *sqlx.DB
		book        *models.BookModel
		err         error
	)

	t.Title("Test Create Book Error: check book existence")
	t.Description("The new book was not created successfully due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		container, err = getPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = applyMigrations(container)
		if err != nil {
			t.Fatal(err)
		}
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
		mockBookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(mockBookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		sCtx.WithNewParameters("book", book)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Error(err)
	})

	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
	}(db)

	defer func() {
		if err = container.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()
}

func TestBookServiceTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(BookServiceTestsSuite))
}

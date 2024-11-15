package integrationTests

import (
	"context"
	repomodels "github.com/nikitalystsev/BookSmart-repo-postgres/core/models"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbdto "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/dto"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (its *IntegrationTestSuite) TestBook_Create_Success(t provider.T) {
	t.Parallel()
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBooks   []*repomodels.BookModel
		err         error
	)
	t.Title("Integration Test Create Book Success")
	t.Description("The new book was successfully created")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(bookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.book where id = $1`
		err = its.db.SelectContext(context.Background(), &findBooks, query, book.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findBooks, 1)
		t.Assert().Equal(book.ID, findBooks[0].ID)
	})
}

func (its *IntegrationTestSuite) TestBook_Delete_Success(t provider.T) {
	t.Parallel()
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBooks   []*repomodels.BookModel
		err         error
	)

	t.Title("Integration Test Delete Book Success")
	t.Description("The new book was successfully deleted")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(bookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.book values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			book.ID,
			book.Title,
			book.Author,
			book.Publisher,
			book.CopiesNumber,
			book.Rarity,
			book.Genre,
			book.PublishingYear,
			book.Language,
			book.AgeLimit,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Delete(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.book where id = $1`
		err = its.db.SelectContext(context.Background(), &findBooks, query, book.ID)
		t.Assert().Len(findBooks, 0)
	})
}

func (its *IntegrationTestSuite) TestBook_GetByID_Success(t provider.T) {
	t.Parallel()
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBook    *models.BookModel
		err         error
	)

	t.Title("Integration Test Get Book By ID Success")
	t.Description("Book was successfully getting by ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(bookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.book values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			book.ID,
			book.Title,
			book.Author,
			book.Publisher,
			book.CopiesNumber,
			book.Rarity,
			book.Genre,
			book.PublishingYear,
			book.Language,
			book.AgeLimit,
		)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookService.GetByID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(book, findBook)
	})
}

func (its *IntegrationTestSuite) TestBook_GetByParams_Success(t provider.T) {
	t.Parallel()
	var (
		bookService intf.IBookService
		book        *models.BookModel
		findBooks   []*models.BookModel
		params      *dto.BookParamsDTO
		err         error
	)

	t.Title("Integration Test Get Book By Params Success")
	t.Description("Book was successfully getting by params")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(bookRepo, logging.GetLoggerForTests())
		params = tdbdto.NewBookParamsDTOBuilder().WithTitle("new test title").Build()
		book = tdbmodels.NewBookModelBuilder().WithTitle("new test title").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.book values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			book.ID,
			book.Title,
			book.Author,
			book.Publisher,
			book.CopiesNumber,
			book.Rarity,
			book.Genre,
			book.PublishingYear,
			book.Language,
			book.AgeLimit,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBooks, err = bookService.GetByParams(context.Background(), params)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Len(findBooks, 1)
		t.Assert().Equal([]*models.BookModel{book}, findBooks)
	})
}

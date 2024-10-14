package integrationTest

import (
	"context"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (its *IntegrationTestSuite) TestBook_Create_Success(t provider.T) {
	var (
		bookService intf.IBookService
		book        *models.BookModel
		err         error
	)

	t.Title("Integration Test Create Book Success")
	t.Description("The new book was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		bookService = impl.NewBookService(bookRepo, logging.GetLoggerForTests())
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		sCtx.WithNewParameters("book", book)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookService.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

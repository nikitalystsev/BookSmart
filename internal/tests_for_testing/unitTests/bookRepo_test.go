package unitTests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	omdto "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/dto"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type BookRepoTestsSuite struct {
	suite.Suite
}

func (brts *BookRepoTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "book repo", "suite", "steps")
}

func (brts *BookRepoTestsSuite) Test_Create_Success(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		err      error
	)

	t.Title("Test Create Book Success")
	t.Description("The new book was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		db, mock, err := sqlxmock.Newx()
		if err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.book values`).
			WithArgs(book.ID, book.Title, book.Author, book.Publisher,
				book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear,
				book.Language, book.AgeLimit).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookRepo.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (brts *BookRepoTestsSuite) Test_Create_ErrorExecutionQuery(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Create Book Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.book values`).
			WithArgs(book.ID, book.Title, book.Author, book.Publisher,
				book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear,
				book.Language, book.AgeLimit).
			WillReturnError(errors.New("insert error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookRepo.Create(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("insert error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *BookRepoTestsSuite) Test_GetByID_Success(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		findBook *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Get Book By ID Success")
	t.Description("The book was successfully retrieved by ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"}).
			AddRow(book.ID, book.Title, book.Author, book.Publisher,
				book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear,
				book.Language, book.AgeLimit)

		mock.ExpectQuery(`select (.+) from bs.book where (.+)`).
			WithArgs(book.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookRepo.GetByID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(book, findBook)
	})
}

func (brts *BookRepoTestsSuite) Test_GetByID_ErrorBookNotFound(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		findBook *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Get Book By ID Error: book not found")
	t.Description("Book not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.book where (.+)`).
			WithArgs(book.ID).WillReturnError(sql.ErrNoRows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookRepo.GetByID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrBookDoesNotExists, err)
		t.Assert().Nil(findBook)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *BookRepoTestsSuite) Test_GetByTitle_Success(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		findBook *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Get Book By Title Success")
	t.Description("The book was successfully retrieved by title")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"}).
			AddRow(book.ID, book.Title, book.Author, book.Publisher,
				book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear,
				book.Language, book.AgeLimit)

		mock.ExpectQuery(`select (.+) from bs.book where (.+)`).
			WithArgs(book.Title).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookRepo.GetByTitle(context.Background(), book.Title)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(book, findBook)
	})
}

func (brts *BookRepoTestsSuite) Test_GetByTitle_ErrorBookNotFound(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		findBook *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Get Book By Title Error: book not found")
	t.Description("Book not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.book where (.+)`).
			WithArgs(book.Title).WillReturnError(sql.ErrNoRows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBook, err = bookRepo.GetByTitle(context.Background(), book.Title)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrBookDoesNotExists, err)
		t.Assert().Nil(findBook)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *BookRepoTestsSuite) Test_Delete_Success(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Delete Book Success")
	t.Description("The book was successfully deleted")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`delete from bs.book where (.+)`).
			WithArgs(book.ID).WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookRepo.Delete(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (brts *BookRepoTestsSuite) Test_Delete_ErrorExecutionQuery(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Delete Book Error: execution query")
	t.Description("Error execution query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`delete from bs.book where (.+)`).
			WithArgs(book.ID).WillReturnError(errors.New("delete error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookRepo.Delete(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("delete error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *BookRepoTestsSuite) Test_Update_Success(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Update Book Success")
	t.Description("The book was successfully updated")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`update bs.book set (.+) where (.+)`).
			WithArgs(book.Title, book.Author, book.Publisher, book.CopiesNumber,
				book.Rarity, book.Genre, book.PublishingYear, book.Language,
				book.AgeLimit, book.ID).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookRepo.Update(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (brts *BookRepoTestsSuite) Test_Update_ErrorExecutionQuery(t provider.T) {
	var (
		bookRepo intfRepo.IBookRepo
		book     *models.BookModel
		mock     sqlxmock.Sqlmock
		db       *sqlx.DB
		err      error
	)

	t.Title("Test Update Book Error: execution query")
	t.Description("Error execution query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`update bs.book set (.+) where (.+)`).
			WithArgs(book.Title, book.Author, book.Publisher, book.CopiesNumber,
				book.Rarity, book.Genre, book.PublishingYear, book.Language,
				book.AgeLimit, book.ID).
			WillReturnError(errors.New("update error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = bookRepo.Update(context.Background(), book)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("update error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *BookRepoTestsSuite) Test_GetByParams_Success(t provider.T) {
	var (
		bookRepo  intfRepo.IBookRepo
		book      *models.BookModel
		params    *dto.BookParamsDTO
		findBooks []*models.BookModel
		mock      sqlxmock.Sqlmock
		db        *sqlx.DB
		err       error
	)

	t.Title("Test Get Book By Params Success")
	t.Description("The book was successfully retrieved by params")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		params = omdto.NewBookParamsDTOObjectMother().DefaultBookParams()
		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"}).
			AddRow(book.ID, book.Title, book.Author, book.Publisher,
				book.CopiesNumber, book.Rarity, book.Genre, book.PublishingYear,
				book.Language, book.AgeLimit)

		mock.ExpectQuery(`select (.+) from bs.book where (.+)`).
			WithArgs(params.Title, params.Author, params.Publisher, params.CopiesNumber,
				params.Rarity, params.Genre, params.PublishingYear, params.Language,
				params.AgeLimit, params.Limit, params.Offset).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBooks, err = bookRepo.GetByParams(context.Background(), params)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.BookModel{book}, findBooks)
	})
}

func (brts *BookRepoTestsSuite) Test_GetByParams_ErrorBookNotFound(t provider.T) {
	var (
		bookRepo  intfRepo.IBookRepo
		findBooks []*models.BookModel
		params    *dto.BookParamsDTO
		mock      sqlxmock.Sqlmock
		db        *sqlx.DB
		err       error
	)

	t.Title("Test Get Book By Params Error: book not found")
	t.Description("Book not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		bookRepo = implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		params = omdto.NewBookParamsDTOObjectMother().DefaultBookParams()

		mock.ExpectQuery(`select (.+) from bs.book where (.+)`).
			WithArgs(params.Title, params.Author, params.Publisher, params.CopiesNumber,
				params.Rarity, params.Genre, params.PublishingYear, params.Language,
				params.AgeLimit, params.Limit, params.Offset).WillReturnError(sql.ErrNoRows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findBooks, err = bookRepo.GetByParams(context.Background(), params)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrBookDoesNotExists, err)
		t.Assert().Nil(findBooks)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func TestBookRepoTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(BookRepoTestsSuite))
}

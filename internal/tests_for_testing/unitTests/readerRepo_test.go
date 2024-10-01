package unitTests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

type ReaderRepoTestsSuite struct {
	suite.Suite
}

func (rrts *ReaderRepoTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "reader repo", "suite", "steps")
}

func (rrts *ReaderRepoTestsSuite) Test_Create_Success(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Create Reader Success")
	t.Description("The new reader was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.reader values`).
			WithArgs(reader.ID, reader.Fio, reader.PhoneNumber,
				reader.Age, reader.Password, reader.Role).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerRepo.Create(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_Create_ErrorExecutionQuery(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Create Reader Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.reader values`).
			WithArgs(reader.ID, reader.Fio, reader.PhoneNumber,
				reader.Age, reader.Password, reader.Role).
			WillReturnError(errors.New("insert error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerRepo.Create(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("insert error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_GetByPhoneNumber_Success(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		findReader *models.ReaderModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Get Reader By PhoneNumber Success")
	t.Description("Reader was successfully retrieved by PhoneNumber")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "fio", "phone_number", "age", "password", "role"}).
			AddRow(reader.ID, reader.Fio, reader.PhoneNumber, reader.Age, reader.Password, reader.Role)

		mock.ExpectQuery(`select (.+) from bs.reader where (.+)`).
			WithArgs(reader.PhoneNumber).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerRepo.GetByPhoneNumber(context.Background(), reader.PhoneNumber)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reader, findReader)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_GetByPhoneNumber_ErrorReaderNotFound(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		findReader *models.ReaderModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Get Reader By PhoneNumber Error: reader not found")
	t.Description("Reader not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reader where (.+)`).
			WithArgs(reader.PhoneNumber).WillReturnError(sql.ErrNoRows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerRepo.GetByPhoneNumber(context.Background(), reader.PhoneNumber)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrReaderDoesNotExists, err)
		t.Assert().Nil(findReader)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_GetByID_Success(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		findReader *models.ReaderModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Get Reader By ID Success")
	t.Description("Reader was successfully retrieved by ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "fio", "phone_number", "age", "password", "role"}).
			AddRow(reader.ID, reader.Fio, reader.PhoneNumber, reader.Age, reader.Password, reader.Role)

		mock.ExpectQuery(`select (.+) from bs.reader where (.+)`).
			WithArgs(reader.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerRepo.GetByID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reader, findReader)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_GetByID_ErrorReaderNotFound(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		findReader *models.ReaderModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Get Reader By PhoneNumber Error: reader not found")
	t.Description("Reader not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reader where (.+)`).
			WithArgs(reader.ID).WillReturnError(sql.ErrNoRows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerRepo.GetByID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrReaderDoesNotExists, err)
		t.Assert().Nil(findReader)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_IsFavorite_Success(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		book       *models.BookModel
		isFavorite bool
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Book Is Favorite Success")
	t.Description("Book is favorite")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(`select (.+) from bs.favorite_books where (.+)`).
			WithArgs(reader.ID, book.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		isFavorite, err = readerRepo.IsFavorite(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().True(isFavorite)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_IsFavorite_ErrorBookIsNotFavorite(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		book       *models.BookModel
		isFavorite bool
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Book Is Favorite Error")
	t.Description("Book is not favorite")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery(`select (.+) from bs.favorite_books where (.+)`).
			WithArgs(reader.ID, book.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		isFavorite, err = readerRepo.IsFavorite(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().False(isFavorite)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_AddToFavorites_Success(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		book       *models.BookModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Add Book To Favorites Success")
	t.Description("Book was successfully added to favorites")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.favorite_books`).
			WithArgs(reader.ID, book.ID).WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerRepo.AddToFavorites(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_AddToFavorites_ErrorExecutingQuery(t provider.T) {
	var (
		readerRepo intfRepo.IReaderRepo
		reader     *models.ReaderModel
		book       *models.BookModel
		mock       sqlxmock.Sqlmock
		db         *sqlx.DB
		err        error
	)

	t.Title("Test Add Book To Favorites Error: executing query")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		readerRepo = implRepo.NewReaderRepo(db, nil, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.favorite_books`).
			WithArgs(reader.ID, book.ID).WillReturnError(errors.New("insert error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerRepo.AddToFavorites(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("insert error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_SaveRefreshToken_Success(t provider.T) {
	var (
		readerRepo   intfRepo.IReaderRepo
		reader       *models.ReaderModel
		mock         redismock.ClientMock
		redisClient  *redis.Client
		refreshToken string
		err          error
	)

	t.Title("Test Save Refresh Token Success")
	t.Description("Refresh token was successfully saved")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		redisClient, mock = redismock.NewClientMock()

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(nil, redisClient, logging.GetLoggerForTests())

		refreshToken = "test_refresh_token"
		mock.ExpectSet(refreshToken, reader.ID.String(), time.Hour).
			SetVal(reader.ID.String())
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerRepo.SaveRefreshToken(context.Background(), reader.ID, refreshToken, time.Hour)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_SaveRefreshToken_ErrorFailedSetToken(t provider.T) {
	var (
		readerRepo   intfRepo.IReaderRepo
		reader       *models.ReaderModel
		mock         redismock.ClientMock
		redisClient  *redis.Client
		refreshToken string
		err          error
	)

	t.Title("Test Save Refresh Token Error: failed set token")
	t.Description("Refresh was not successfully saved")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		redisClient, mock = redismock.NewClientMock()

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(nil, redisClient, logging.GetLoggerForTests())

		refreshToken = "test_refresh_token"
		mock.ExpectSet(refreshToken, reader.ID.String(), time.Hour).
			SetErr(errors.New("set token error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerRepo.SaveRefreshToken(context.Background(), reader.ID, refreshToken, time.Hour)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("set token error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_GetByRefreshToken_Success(t provider.T) {
	var (
		readerRepo   intfRepo.IReaderRepo
		reader       *models.ReaderModel
		findReader   *models.ReaderModel
		mock         redismock.ClientMock
		dbmock       sqlxmock.Sqlmock
		db           *sqlx.DB
		redisClient  *redis.Client
		refreshToken string
		err          error
	)

	t.Title("Test Get By Refresh Token Success")
	t.Description("Reader was successfully retrieved by refresh token")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, dbmock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}
		redisClient, mock = redismock.NewClientMock()

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		readerRepo = implRepo.NewReaderRepo(db, redisClient, logging.GetLoggerForTests())

		refreshToken = "test_refresh_token"
		mock.ExpectGet(refreshToken).SetVal(reader.ID.String())

		rows := sqlxmock.NewRows([]string{"id", "fio", "phone_number", "age", "password", "role"}).
			AddRow(reader.ID, reader.Fio, reader.PhoneNumber, reader.Age, reader.Password, reader.Role)
		dbmock.ExpectQuery(`select (.+) from bs.reader where (.+)`).
			WithArgs(reader.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerRepo.GetByRefreshToken(context.Background(), refreshToken)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reader, findReader)
	})
}

func (rrts *ReaderRepoTestsSuite) Test_GetByRefreshToken_ErrorRefreshTokenNotFound(t provider.T) {
	var (
		readerRepo   intfRepo.IReaderRepo
		findReader   *models.ReaderModel
		mock         redismock.ClientMock
		redisClient  *redis.Client
		refreshToken string
		err          error
	)

	t.Title("Test Get By Refresh Token Error: refresh token not found")
	t.Description("Refresh token was not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		redisClient, mock = redismock.NewClientMock()

		readerRepo = implRepo.NewReaderRepo(nil, redisClient, logging.GetLoggerForTests())

		refreshToken = "test_refresh_token"
		mock.ExpectGet(refreshToken).RedisNil()
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerRepo.GetByRefreshToken(context.Background(), refreshToken)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrReaderDoesNotExists, err)
		t.Assert().Nil(findReader)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func TestReaderRepoTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ReaderRepoTestsSuite))
}

package unitTests

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

type ReservationRepoTestsSuite struct {
	suite.Suite
}

func (brts *ReservationRepoTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "reservation repo", "suite", "steps")
}

func (brts *ReservationRepoTestsSuite) Test_Create_Success(t provider.T) {
	var (
		reservationRepo intfRepo.IReservationRepo
		mock            sqlxmock.Sqlmock
		db              *sqlx.DB
		reservation     *models.ReservationModel
		err             error
	)

	t.Title("Test Create Reservation Success")
	t.Description("The new reservation was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.reservation`).
			WithArgs(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationRepo.Create(context.Background(), reservation)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_Create_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo intfRepo.IReservationRepo
		mock            sqlxmock.Sqlmock
		db              *sqlx.DB
		reservation     *models.ReservationModel
		err             error
	)

	t.Title("Test Create Reservation Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.reservation`).
			WithArgs(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State).
			WillReturnError(errors.New("insert error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationRepo.Create(context.Background(), reservation)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("insert error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetByReaderAndBook_Success(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		reservation      *models.ReservationModel
		book             *models.BookModel
		reader           *models.ReaderModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Reservations By Reader And Book Success")
	t.Description("Reservations was successfully retrieved by reader and book")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
			AddRow(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State)

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reader.ID, book.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetByReaderAndBook(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation}, findReservations)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetByReaderAndBook_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		book             *models.BookModel
		reader           *models.ReaderModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Reservations By Reader And Book Error: execution query error")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reader.ID, book.ID).WillReturnError(errors.New("query error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetByReaderAndBook(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("query error"), err)
		t.Assert().Nil(findReservations)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetByID_Success(t provider.T) {
	var (
		reservationRepo intfRepo.IReservationRepo
		mock            sqlxmock.Sqlmock
		db              *sqlx.DB
		reservation     *models.ReservationModel
		findReservation *models.ReservationModel
		err             error
	)

	t.Title("Test Get Reservation By ID Success")
	t.Description("Reservation was successfully retrieved by ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
			AddRow(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State)

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reservation.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservation, err = reservationRepo.GetByID(context.Background(), reservation.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reservation, findReservation)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetByID_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo intfRepo.IReservationRepo
		mock            sqlxmock.Sqlmock
		db              *sqlx.DB
		reservation     *models.ReservationModel
		findReservation *models.ReservationModel
		err             error
	)

	t.Title("Test Get Reservations By ID Error: execution query error")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reservation.ID).WillReturnError(errors.New("query error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservation, err = reservationRepo.GetByID(context.Background(), reservation.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("query error"), err)
		t.Assert().Nil(findReservation)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetByBookID_Success(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		reservation      *models.ReservationModel
		book             *models.BookModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Reservations By Book ID Success")
	t.Description("Reservations was successfully retrieved by book id")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
			AddRow(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State)

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(book.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetByBookID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation}, findReservations)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetByBookID_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		book             *models.BookModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Reservations By Book ID Error: execution query error")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		book = ommodels.NewBookModelObjectMother().DefaultBook()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(book.ID).WillReturnError(errors.New("query error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetByBookID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("query error"), err)
		t.Assert().Nil(findReservations)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_Update_Success(t provider.T) {
	var (
		reservationRepo intfRepo.IReservationRepo
		mock            sqlxmock.Sqlmock
		db              *sqlx.DB
		reservation     *models.ReservationModel
		err             error
	)

	t.Title("Test Update Reservations Success")
	t.Description("Reservations was successfully updated")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`update bs.reservation set (.+) where (.+)`).
			WithArgs(reservation.ReaderID, reservation.BookID, reservation.IssueDate,
				reservation.ReturnDate, reservation.State, reservation.ID).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationRepo.Update(context.Background(), reservation)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_Update_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo intfRepo.IReservationRepo
		mock            sqlxmock.Sqlmock
		db              *sqlx.DB
		reservation     *models.ReservationModel
		err             error
	)

	t.Title("Test Get Reservations By Book ID Error: execution query error")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reservation = ommodels.NewReservationModelObjectMother().DefaultReservation()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`update bs.reservation set (.+) where (.+)`).
			WithArgs(reservation.ReaderID, reservation.BookID, reservation.IssueDate,
				reservation.ReturnDate, reservation.State, reservation.ID).
			WillReturnError(errors.New("update error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationRepo.Update(context.Background(), reservation)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("update error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetExpiredByReaderID_Success(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		reservation      *models.ReservationModel
		reader           *models.ReaderModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Expired Reservations By Reader ID Success")
	t.Description("Expired reservations was successfully retrieved by reader ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservation = tdbmodels.NewReservationModelBuilder().
			WithReaderID(reader.ID).
			WithIssueDate(time.Now().AddDate(0, 0, -10)).
			WithReturnDate(time.Now().AddDate(0, 0, -5)).
			WithState("Expired").Build()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
			AddRow(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State)

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reader.ID, time.Now()).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetExpiredByReaderID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation}, findReservations)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetExpiredByReaderID_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		reader           *models.ReaderModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Expired Reservations By Reader ID Error: execution query error")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reader.ID, time.Now()).WillReturnError(errors.New("select error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetExpiredByReaderID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("select error"), err)
		t.Assert().Nil(findReservations)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetActiveByReaderID_Success(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		reservation      *models.ReservationModel
		reader           *models.ReaderModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Active Reservations By Reader ID Success")
	t.Description("Active reservations was successfully retrieved by reader ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservation = tdbmodels.NewReservationModelBuilder().
			WithReaderID(reader.ID).Build()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "book_id", "issue_date", "return_date", "state"}).
			AddRow(reservation.ID, reservation.ReaderID, reservation.BookID,
				reservation.IssueDate, reservation.ReturnDate, reservation.State)

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reader.ID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetActiveByReaderID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation}, findReservations)
	})
}

func (brts *ReservationRepoTestsSuite) Test_GetActiveByReaderID_ErrorExecutionQuery(t provider.T) {
	var (
		reservationRepo  intfRepo.IReservationRepo
		mock             sqlxmock.Sqlmock
		db               *sqlx.DB
		reader           *models.ReaderModel
		findReservations []*models.ReservationModel
		err              error
	)

	t.Title("Test Get Active Reservations By Reader ID Error: execution query error")
	t.Description("Error executing query")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservationRepo = implRepo.NewReservationRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.reservation_view where (.+)`).
			WithArgs(reader.ID).WillReturnError(errors.New("select error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationRepo.GetActiveByReaderID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("select error"), err)
		t.Assert().Nil(findReservations)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func TestReservationRepoTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ReservationRepoTestsSuite))
}

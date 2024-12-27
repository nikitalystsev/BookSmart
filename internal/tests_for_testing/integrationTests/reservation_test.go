package integrationTests

import (
	"context"
	repomodels "github.com/nikitalystsev/BookSmart-repo-postgres/core/models"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/pkg/transact"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (its *IntegrationTestSuite) TestReservation_Create_Success(t provider.T) {
	t.Parallel()
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		findReservations   []*repomodels.ReservationModel
		err                error
	)

	t.Title("Integration Test Create Reservation Success")
	t.Description("The new reservation was successfully created")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		var transactionManager transact.ITransactionManager
		if transactionManager, err = getTransactionManagerForIntegrationTests(its.db); err != nil {
			t.Fatal(err)
		}
		reservationService = impl.NewReservationService(
			reservationRepo,
			bookRepo,
			readerRepo,
			libCardRepo,
			transactionManager,
			logging.GetLoggerForTests(),
		)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("87654321098").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithLibCardNum("4676874323421").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationService.Create(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.reservation where reader_id = $1 and book_id = $2`
		err = its.db.SelectContext(context.Background(), &findReservations, query, reader.ID, book.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findReservations, 1)
	})
}

func (its *IntegrationTestSuite) TestReservation_Update_Success(t provider.T) {
	t.Parallel()
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		reservation        *models.ReservationModel
		findReservations   []*repomodels.ReservationModel
		err                error
	)

	t.Title("Integration Test Update Reservation Success")
	t.Description("The new reservation was successfully updated")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		var transactionManager transact.ITransactionManager
		if transactionManager, err = getTransactionManagerForIntegrationTests(its.db); err != nil {
			t.Fatal(err)
		}
		reservationService = impl.NewReservationService(
			reservationRepo,
			bookRepo,
			readerRepo,
			libCardRepo,
			transactionManager,
			logging.GetLoggerForTests(),
		)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("76543210987").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithLibCardNum("4676874376989").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
		reservation = tdbmodels.NewReservationModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reservation values ($1, $2, $3, $4, $5, $6)`,
			reservation.ID,
			reservation.ReaderID,
			reservation.BookID,
			reservation.IssueDate,
			reservation.ReturnDate,
			reservation.State,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationService.Update(context.Background(), reservation, 6)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.reservation where reader_id = $1 and book_id = $2`
		err = its.db.SelectContext(context.Background(), &findReservations, query, reader.ID, book.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findReservations, 1)
		t.Assert().Equal(impl.ReservationExtended, findReservations[0].State)
	})
}

func (its *IntegrationTestSuite) TestReservation_GetByBookID_Success(t provider.T) {
	t.Parallel()
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		reservation        *models.ReservationModel
		findReservations   []*models.ReservationModel
		err                error
	)

	t.Title("Integration Test Get Reservations By Book ID Success")
	t.Description("Reservations was successfully getting by book ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		var transactionManager transact.ITransactionManager
		if transactionManager, err = getTransactionManagerForIntegrationTests(its.db); err != nil {
			t.Fatal(err)
		}
		reservationService = impl.NewReservationService(
			reservationRepo,
			bookRepo,
			readerRepo,
			libCardRepo,
			transactionManager,
			logging.GetLoggerForTests(),
		)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("65432109876").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithLibCardNum("4676874346545").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
		reservation = tdbmodels.NewReservationModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reservation values ($1, $2, $3, $4, $5, $6)`,
			reservation.ID,
			reservation.ReaderID,
			reservation.BookID,
			reservation.IssueDate,
			reservation.ReturnDate,
			reservation.State,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationService.GetByBookID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation}[0].ID, findReservations[0].ID)
	})
}

func (its *IntegrationTestSuite) TestReservation_GetAllReservationsByReaderID_Success(t provider.T) {
	t.Parallel()
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		reservation        *models.ReservationModel
		findReservations   []*models.ReservationModel
		err                error
	)

	t.Title("Integration Test Get Reservations By Reader ID Success")
	t.Description("Reservations was successfully getting by reader ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		var transactionManager transact.ITransactionManager
		if transactionManager, err = getTransactionManagerForIntegrationTests(its.db); err != nil {
			t.Fatal(err)
		}
		reservationService = impl.NewReservationService(
			reservationRepo,
			bookRepo,
			readerRepo,
			libCardRepo,
			transactionManager,
			logging.GetLoggerForTests(),
		)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("54321098765").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithLibCardNum("3543654757893").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
		reservation = tdbmodels.NewReservationModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reservation values ($1, $2, $3, $4, $5, $6)`,
			reservation.ID,
			reservation.ReaderID,
			reservation.BookID,
			reservation.IssueDate,
			reservation.ReturnDate,
			reservation.State,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationService.GetByReaderID(context.Background(), reader.ID, impl.ReviewsPageLimit, 0)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation}[0].ID, findReservations[0].ID)
	})
}

func (its *IntegrationTestSuite) TestReservation_GetByID_Success(t provider.T) {
	t.Parallel()
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		reservation        *models.ReservationModel
		findReservation    *models.ReservationModel
		err                error
	)

	t.Title("Integration Test Get Reservations By ID Success")
	t.Description("Reservations was successfully getting by ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		var transactionManager transact.ITransactionManager
		if transactionManager, err = getTransactionManagerForIntegrationTests(its.db); err != nil {
			t.Fatal(err)
		}
		reservationService = impl.NewReservationService(
			reservationRepo,
			bookRepo,
			readerRepo,
			libCardRepo,
			transactionManager,
			logging.GetLoggerForTests(),
		)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("43210987654").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithLibCardNum("4676874346547").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
		reservation = tdbmodels.NewReservationModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reservation values ($1, $2, $3, $4, $5, $6)`,
			reservation.ID,
			reservation.ReaderID,
			reservation.BookID,
			reservation.IssueDate,
			reservation.ReturnDate,
			reservation.State,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservation, err = reservationService.GetByID(context.Background(), reservation.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reservation.ID, findReservation.ID)
	})
}

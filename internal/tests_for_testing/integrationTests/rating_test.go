package integrationTests

import (
	"context"
	repomodels "github.com/nikitalystsev/BookSmart-repo-postgres/core/models"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (its *IntegrationTestSuite) TestRating_Create_Success(t provider.T) {
	t.Parallel()
	var (
		ratingService intf.IRatingService
		reader        *models.ReaderModel
		book          *models.BookModel
		rating        *models.RatingModel
		reservation   *models.ReservationModel
		findRatings   []*repomodels.RatingModel
		err           error
	)

	t.Title("Integration Test Rating Create Success")
	t.Description("The new rating was successfully created")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ratingRepo := implRepo.NewRatingRepo(its.db, logging.GetLoggerForTests())
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		ratingService = impl.NewRatingService(ratingRepo, reservationRepo, logging.GetLoggerForTests())
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("34567890123").Build()
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

		rating = tdbmodels.NewRatingModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).Build()
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = ratingService.Create(context.Background(), rating)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.rating where reader_id = $1 and book_id = $2`
		err = its.db.SelectContext(context.Background(), &findRatings, query, reader.ID, book.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findRatings, 1)
		t.Assert().Equal(reader.ID, findRatings[0].ReaderID)
		t.Assert().Equal(book.ID, findRatings[0].BookID)
	})
}

func (its *IntegrationTestSuite) TestRating_GetByBookID_Success(t provider.T) {
	t.Parallel()
	var (
		ratingService intf.IRatingService
		reader        *models.ReaderModel
		book          *models.BookModel
		rating        *models.RatingModel
		findRatings   []*models.RatingModel
		err           error
	)

	t.Title("Integration Test Get Rating By Book ID Success")
	t.Description("Rating was successfully getting by book ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ratingRepo := implRepo.NewRatingRepo(its.db, logging.GetLoggerForTests())
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		ratingService = impl.NewRatingService(ratingRepo, reservationRepo, logging.GetLoggerForTests())
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("45678901234").Build()
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

		rating = tdbmodels.NewRatingModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).Build()

		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.rating values ($1, $2, $3, $4, $5)`,
			rating.ID,
			rating.ReaderID,
			rating.BookID,
			rating.Review,
			rating.Rating,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findRatings, err = ratingService.GetByBookID(context.Background(), book.ID, impl.ReviewsPageLimit, 0)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Len(findRatings, 1)
		t.Assert().Equal(reader.ID, findRatings[0].ReaderID)
		t.Assert().Equal(book.ID, findRatings[0].BookID)
	})
}

func (its *IntegrationTestSuite) TestRating_GetAvgRatingByBookID_Success(t provider.T) {
	t.Parallel()
	var (
		ratingService intf.IRatingService
		reader        *models.ReaderModel
		book          *models.BookModel
		rating        *models.RatingModel
		avgRating     float32
		err           error
	)

	t.Title("Integration Test Get Avg Rating For Book With ID Success")
	t.Description("Avg rating was successfully getting for book with ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ratingRepo := implRepo.NewRatingRepo(its.db, logging.GetLoggerForTests())
		reservationRepo := implRepo.NewReservationRepo(its.db, logging.GetLoggerForTests())
		ratingService = impl.NewRatingService(ratingRepo, reservationRepo, logging.GetLoggerForTests())
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("56789012345").Build()
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

		rating = tdbmodels.NewRatingModelBuilder().
			WithReaderID(reader.ID).
			WithBookID(book.ID).
			WithRating(5).Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.rating values ($1, $2, $3, $4, $5)`,
			rating.ID,
			rating.ReaderID,
			rating.BookID,
			rating.Review,
			rating.Rating,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		avgRating, err = ratingService.GetAvgRatingByBookID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(float32(5), avgRating)
	})
}

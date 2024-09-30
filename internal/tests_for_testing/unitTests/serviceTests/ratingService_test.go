package serviceTests

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

type RatingServiceTestsSuite struct {
	suite.Suite
}

func (bsts *RatingServiceTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "rating service", "suite", "steps")
}

/*
Лондонский вариант
*/
func (bsts *RatingServiceTestsSuite) Test_Create_Success(t provider.T) {
	var (
		ratingService intf.IRatingService
		reader        *models.ReaderModel
		reservation   *models.ReservationModel
		rating        *models.RatingModel
		bookID        uuid.UUID
		err           error
	)

	t.Title("Test Create Rating Success")
	t.Description("The new rating was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockRatingRepo := mockrepo.NewMockIRatingRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		bookID = uuid.New()
		rating = tdbmodels.NewRatingModelBuilder().WithReaderID(reader.ID).WithBookID(bookID).Build()
		reservation = tdbmodels.NewReservationModelBuilder().WithBookID(bookID).WithReaderID(reader.ID).Build()
		mockRatingRepo.EXPECT().GetByReaderAndBook(gomock.Any(), reader.ID, bookID).Return(nil, errs.ErrRatingDoesNotExists)
		mockReservationRepo.EXPECT().GetByReaderAndBook(gomock.Any(), reader.ID, bookID).Return([]*models.ReservationModel{reservation}, nil)
		mockRatingRepo.EXPECT().Create(gomock.Any(), rating).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = ratingService.Create(context.Background(), rating)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (bsts *RatingServiceTestsSuite) Test_Create_ErrorRatingAlreadyExists(t provider.T) {
	var (
		ratingService intf.IRatingService
		reader        *models.ReaderModel
		rating        *models.RatingModel
		bookID        uuid.UUID
		err           error
	)

	t.Title("Test Create Rating Error: rating already exists")
	t.Description("Rating already exists")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockRatingRepo := mockrepo.NewMockIRatingRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		bookID = uuid.New()
		rating = tdbmodels.NewRatingModelBuilder().WithReaderID(reader.ID).WithBookID(bookID).Build()

		mockRatingRepo.EXPECT().GetByReaderAndBook(gomock.Any(), reader.ID, bookID).Return(rating, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = ratingService.Create(context.Background(), rating)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrRatingAlreadyExist, err)
	})
}

func (bsts *RatingServiceTestsSuite) Test_GetByBookID_Success(t provider.T) {
	var (
		ratingService intf.IRatingService
		rating        *models.RatingModel
		findRatings   []*models.RatingModel
		bookID        uuid.UUID
		err           error
	)

	t.Title("Test Get Rating By Book ID Success")
	t.Description("Rating was successfully retrieved by book ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockRatingRepo := mockrepo.NewMockIRatingRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		bookID = uuid.New()
		rating = tdbmodels.NewRatingModelBuilder().WithBookID(bookID).Build()

		mockRatingRepo.EXPECT().GetByBookID(gomock.Any(), bookID).Return([]*models.RatingModel{rating}, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findRatings, err = ratingService.GetByBookID(context.Background(), bookID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.RatingModel{rating}, findRatings)
	})
}

func (bsts *RatingServiceTestsSuite) Test_GetByBookID_ErrorGettingRatings(t provider.T) {
	var (
		ratingService intf.IRatingService
		findRatings   []*models.RatingModel
		bookID        uuid.UUID
		err           error
	)

	t.Title("Test Get Rating By Book ID Error: getting ratings")
	t.Description("Rating was not retrieved because of error getting ratings")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockRatingRepo := mockrepo.NewMockIRatingRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		bookID = uuid.New()
		mockRatingRepo.EXPECT().GetByBookID(gomock.Any(), bookID).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findRatings, err = ratingService.GetByBookID(context.Background(), bookID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Nil(findRatings)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

func (bsts *RatingServiceTestsSuite) Test_GetAvgRatingByBookID_Success(t provider.T) {
	var (
		ratingService intf.IRatingService
		rating        *models.RatingModel
		bookID        uuid.UUID
		avgRating     float32
		err           error
	)

	t.Title("Test Get Avg Rating By Book ID Success")
	t.Description("Avg rating was successfully retrieved for book with ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockRatingRepo := mockrepo.NewMockIRatingRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		bookID = uuid.New()
		rating = tdbmodels.NewRatingModelBuilder().WithBookID(bookID).Build()
		mockRatingRepo.EXPECT().GetByBookID(gomock.Any(), bookID).Return([]*models.RatingModel{rating}, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		avgRating, err = ratingService.GetAvgRatingByBookID(context.Background(), bookID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(float32(5), avgRating)
	})
}

func (bsts *RatingServiceTestsSuite) Test_GetAvgRatingByBookID_ErrorRatingsNotFound(t provider.T) {
	var (
		ratingService intf.IRatingService
		bookID        uuid.UUID
		avgRating     float32
		err           error
	)

	t.Title("Test Get Avg Rating By Book ID Error: ratings not found")
	t.Description("Avg Rating was not retrieved because of error ratings not found")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockRatingRepo := mockrepo.NewMockIRatingRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		bookID = uuid.New()
		mockRatingRepo.EXPECT().GetByBookID(gomock.Any(), bookID).Return(nil, errs.ErrRatingDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		avgRating, err = ratingService.GetAvgRatingByBookID(context.Background(), bookID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrRatingDoesNotExists, err)
		t.Assert().Equal(float32(-1), avgRating)
	})
}

/*
	Классический вариант
*/

func (bsts *RatingServiceTestsSuite) Test_Create_Success_Classic(t provider.T) {
	var (
		postgresContainer *postgres.PostgresContainer
		db                *sqlx.DB
		ratingService     intf.IRatingService
		reader            *models.ReaderModel
		rating            *models.RatingModel
		reservation       *models.ReservationModel
		book              *models.BookModel
		err               error
	)

	t.Title("Test Create Rating Success Classic")
	t.Description("The new rating was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		postgresContainer, err = unitTests.GetPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = unitTests.ApplyMigrations(postgresContainer)
		if err != nil {
			t.Fatal(err)
		}
		mockRatingRepo := implRepo.NewRatingRepo(db, logging.GetLoggerForTests())
		mockReservationRepo := implRepo.NewReservationRepo(db, logging.GetLoggerForTests())
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		_, err = db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		_, err = db.ExecContext(
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
		reservation = tdbmodels.NewReservationModelBuilder().WithBookID(book.ID).WithReaderID(reader.ID).Build()
		_, err = db.ExecContext(
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
		rating = tdbmodels.NewRatingModelBuilder().WithReaderID(reader.ID).WithBookID(book.ID).Build()
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = ratingService.Create(context.Background(), rating)
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
		if err = postgresContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()
}

func (bsts *RatingServiceTestsSuite) Test_Create_ErrorRatingAlreadyExists_Classic(t provider.T) {
	var (
		postgresContainer *postgres.PostgresContainer
		db                *sqlx.DB
		ratingService     intf.IRatingService
		reader            *models.ReaderModel
		rating            *models.RatingModel
		book              *models.BookModel
		err               error
	)

	t.Title("Test Create Rating Error: rating already exists")
	t.Description("Rating already exists")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		postgresContainer, err = unitTests.GetPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = unitTests.ApplyMigrations(postgresContainer)
		if err != nil {
			t.Fatal(err)
		}
		mockRatingRepo := implRepo.NewRatingRepo(db, logging.GetLoggerForTests())
		mockReservationRepo := implRepo.NewReservationRepo(db, logging.GetLoggerForTests())
		ratingService = impl.NewRatingService(
			mockRatingRepo,
			mockReservationRepo,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		_, err = db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		_, err = db.ExecContext(
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
		rating = tdbmodels.NewRatingModelBuilder().WithReaderID(reader.ID).WithBookID(book.ID).Build()

		_, err = db.ExecContext(
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
		err = ratingService.Create(context.Background(), rating)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrRatingAlreadyExist, err)
	})

	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
	}(db)

	defer func() {
		if err = postgresContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()
}

func TestRatingServiceTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(RatingServiceTestsSuite))
}

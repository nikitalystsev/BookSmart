package unitTests

import (
	"context"
	"errors"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/pkg/transact"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"log"
	"testing"
	"time"
)

type ReservationServiceTestsSuite struct {
	suite.Suite
}

func (rsts *ReservationServiceTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "reservation service", "suite", "steps")
}

/*
Лондонский вариант
*/

func (rsts *ReservationServiceTestsSuite) Test_Create_Success(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		err                error
	)

	t.Title("Test Reservation Create Success")
	t.Description("The new reservation was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()

		mockReaderRepo.EXPECT().GetByID(gomock.Any(), reader.ID).Return(reader, nil)
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), reader.ID).Return(libCard, nil)
		mockReservationRepo.EXPECT().GetExpiredByReaderID(gomock.Any(), reader.ID).Return(nil, nil)
		mockReservationRepo.EXPECT().GetActiveByReaderID(gomock.Any(), reader.ID).Return(nil, nil)
		mockReservationRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockBookRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
		mockTransactionManager.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationService.Create(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_Create_ErrorExpiredLibCard(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		err                error
	)

	t.Title("Test Reservation Create Error: expired libCard")
	t.Description("Reservation cannot be made due to expired libCard")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		libCard = ommodels.NewLibCardModelObjectMother().ExpiredLibCard()

		mockReaderRepo.EXPECT().GetByID(gomock.Any(), reader.ID).Return(reader, nil)
		mockReservationRepo.EXPECT().GetExpiredByReaderID(gomock.Any(), reader.ID).Return(nil, nil)
		mockReservationRepo.EXPECT().GetActiveByReaderID(gomock.Any(), reader.ID).Return(nil, nil)
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), reader.ID).Return(libCard, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationService.Create(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrLibCardIsInvalid, err)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_Update_Success(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		reservation        *models.ReservationModel
		libCard            *models.LibCardModel
		err                error
	)

	t.Title("Test Reservation Update Success")
	t.Description("Reservation was successfully updated")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)

		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		reservation = tdbmodels.NewReservationModelBuilder().WithReaderID(reader.ID).WithBookID(book.ID).Build()

		mockReservationRepo.EXPECT().Update(gomock.Any(), reservation).Return(nil)
		mockReservationRepo.EXPECT().GetExpiredByReaderID(gomock.Any(), reader.ID).Return(nil, nil)
		mockBookRepo.EXPECT().GetByID(gomock.Any(), reservation.BookID).Return(book, nil)
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), reader.ID).Return(libCard, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationService.Update(context.Background(), reservation, 5)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_Update_ExpiredLibCard(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		reservation        *models.ReservationModel
		libCard            *models.LibCardModel
		err                error
	)

	t.Title("Test Reservation Update Error: expired libCard")
	t.Description("Reservation was not updated due to invalid library card")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		libCard = ommodels.NewLibCardModelObjectMother().ExpiredLibCard()
		reservation = tdbmodels.NewReservationModelBuilder().WithReaderID(reader.ID).WithBookID(book.ID).Build()

		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), reader.ID).Return(libCard, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = reservationService.Update(context.Background(), reservation, 5)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrLibCardIsInvalid, err)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_GetByBookID_Success(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		reservation        *models.ReservationModel
		findReservations   []*models.ReservationModel
		err                error
	)

	t.Title("Test Get Reservation By Book ID Success")
	t.Description("Reservation was successfully retrieved by book ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		reservation = tdbmodels.NewReservationModelBuilder().WithReaderID(reader.ID).WithBookID(book.ID).Build()

		mockReservationRepo.EXPECT().GetByBookID(gomock.Any(), book.ID).Return([]*models.ReservationModel{reservation}, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationService.GetByBookID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(findReservations, []*models.ReservationModel{reservation})
	})
}

func (rsts *ReservationServiceTestsSuite) Test_GetByBookID_ErrorReservationDoesNotExists(t provider.T) {
	var (
		reservationService intf.IReservationService
		book               *models.BookModel
		findReservations   []*models.ReservationModel
		err                error
	)

	t.Title("Test Get Reservation By Book ID Error: reservation not found")
	t.Description("The reservation was not received because it does not exist")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		book = ommodels.NewBookModelObjectMother().DefaultBook()

		mockReservationRepo.EXPECT().GetByBookID(gomock.Any(), book.ID).Return(nil, errs.ErrReservationDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationService.GetByBookID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrReservationDoesNotExists, err)
		t.Assert().Nil(findReservations)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_GetByID_Success(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		book               *models.BookModel
		reservation        *models.ReservationModel
		findReservation    *models.ReservationModel
		err                error
	)

	t.Title("Test Get Reservation By Book ID Success")
	t.Description("Reservation was successfully retrieved by book ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		reservation = tdbmodels.NewReservationModelBuilder().WithReaderID(reader.ID).WithBookID(book.ID).Build()

		mockReservationRepo.EXPECT().GetByID(gomock.Any(), reservation.ID).Return(reservation, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservation, err = reservationService.GetByID(context.Background(), reservation.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(findReservation, reservation)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_GetByID_ErrorReservationDoesNotExists(t provider.T) {
	var (
		reservationService intf.IReservationService
		book               *models.BookModel
		findReservation    *models.ReservationModel
		err                error
	)

	t.Title("Test Get Reservation By Book ID Error: reservation not found")
	t.Description("The reservation was not received because it does not exist")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		book = ommodels.NewBookModelObjectMother().DefaultBook()

		mockReservationRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(nil, errs.ErrReservationDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservation, err = reservationService.GetByID(context.Background(), book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrReservationDoesNotExists, err)
		t.Assert().Nil(findReservation)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_GetByReaderID_Success(t provider.T) {
	var (
		reservationService intf.IReservationService
		reservation        *models.ReservationModel
		expiredReservation *models.ReservationModel
		reader             *models.ReaderModel
		findReservations   []*models.ReservationModel
		err                error
	)

	t.Title("Test Get All Reservation By Reader ID Success")
	t.Description("The all reservation by reader ID was received")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		reservation = tdbmodels.NewReservationModelBuilder().WithReaderID(reader.ID).Build()
		expiredReservation = tdbmodels.NewReservationModelBuilder().WithReaderID(reader.ID).
			WithReturnDate(time.Now().AddDate(0, -1, 0)).Build()
		mockReservationRepo.EXPECT().GetByReaderID(gomock.Any(), reader.ID, impl.ReservationsPageLimit, 0).Return([]*models.ReservationModel{reservation, expiredReservation}, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationService.GetByReaderID(context.Background(), reader.ID, impl.ReservationsPageLimit, 0)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal([]*models.ReservationModel{reservation, expiredReservation}, findReservations)
	})
}

func (rsts *ReservationServiceTestsSuite) Test_GetByReaderID_ErrorCheckExpiredReservations(t provider.T) {
	var (
		reservationService intf.IReservationService
		reader             *models.ReaderModel
		findReservations   []*models.ReservationModel
		err                error
	)

	t.Title("Test Get All Reservation By Reader ID Error: check expired reservations")
	t.Description("All reader's reservations were not received due to error receiving valid reservations")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		mockReservationRepo := mockrepo.NewMockIReservationRepo(ctrl)
		mockTransactionManager := mockrepo.NewMockITransactionManager(ctrl)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
			logging.GetLoggerForTests(),
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReservationRepo.EXPECT().GetByReaderID(gomock.Any(), reader.ID, impl.ReservationsPageLimit, 0).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReservations, err = reservationService.GetByReaderID(context.Background(), reader.ID, impl.ReservationsPageLimit, 0)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Nil(findReservations)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

/*
Классический вариант
*/

func (rsts *ReservationServiceTestsSuite) Test_Create_Success_Classic(t provider.T) {
	var (
		postgresContainer  *postgres.PostgresContainer
		db                 *sqlx.DB
		redisContainer     *testredis.RedisContainer
		redisClient        *redis.Client
		reservationService intf.IReservationService
		trm                *manager.Manager
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		err                error
	)

	t.Title("Test Reservation Create Success Classic")
	t.Description("The new reservation was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		postgresContainer, err = GetPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = ApplyMigrations(postgresContainer)
		if err != nil {
			t.Fatal(err)
		}
		redisContainer, err = GetRedisForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		redisClient, err = GetRedisClientForClassicUnitTests(redisContainer)
		if err != nil {
			t.Fatal(err)
		}
		if trm, err = manager.New(trmsqlx.NewDefaultFactory(db)); err != nil {
			t.Fatal(err)
		}
		mockReaderRepo := implRepo.NewReaderRepo(db, redisClient, logging.GetLoggerForTests())
		mockBookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		mockLibCardRepo := implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())
		mockReservationRepo := implRepo.NewReservationRepo(db, logging.GetLoggerForTests())
		mockTransactionManager := transact.NewTransactionManager(trm)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
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
		if err != nil {
			t.Fatal(err)
		}
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
		if err != nil {
			t.Fatal(err)
		}
		libCard = tdbmodels.NewLibCardModelBuilder().WithReaderID(reader.ID).Build()
		_, err = db.ExecContext(
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

	defer func() {
		if err = redisContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
}

func (rsts *ReservationServiceTestsSuite) Test_Create_ErrorExpiredLibCard_Classic(t provider.T) {
	var (
		postgresContainer  *postgres.PostgresContainer
		db                 *sqlx.DB
		redisContainer     *testredis.RedisContainer
		redisClient        *redis.Client
		reservationService intf.IReservationService
		trm                *manager.Manager
		reader             *models.ReaderModel
		book               *models.BookModel
		libCard            *models.LibCardModel
		err                error
	)

	t.Title("Test Reservation Create Classic Error: expired libCard")
	t.Description("Reservation cannot be made due to expired libCard")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		postgresContainer, err = GetPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = ApplyMigrations(postgresContainer)
		if err != nil {
			t.Fatal(err)
		}
		redisContainer, err = GetRedisForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		redisClient, err = GetRedisClientForClassicUnitTests(redisContainer)
		if err != nil {
			t.Fatal(err)
		}
		if trm, err = manager.New(trmsqlx.NewDefaultFactory(db)); err != nil {
			t.Fatal(err)
		}
		mockReaderRepo := implRepo.NewReaderRepo(db, redisClient, logging.GetLoggerForTests())
		mockBookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		mockLibCardRepo := implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())
		mockReservationRepo := implRepo.NewReservationRepo(db, logging.GetLoggerForTests())
		mockTransactionManager := transact.NewTransactionManager(trm)
		reservationService = impl.NewReservationService(
			mockReservationRepo,
			mockBookRepo,
			mockReaderRepo,
			mockLibCardRepo,
			mockTransactionManager,
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
		if err != nil {
			t.Fatal(err)
		}
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
		if err != nil {
			t.Fatal(err)
		}
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithActionStatus(false).Build()
		_, err = db.ExecContext(
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
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrLibCardIsInvalid, err)
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

	defer func() {
		if err = redisContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
}

func TestReservationServiceTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ReservationServiceTestsSuite))
}

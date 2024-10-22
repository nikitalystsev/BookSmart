package e2eTest

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/pkg/auth"
	"github.com/nikitalystsev/BookSmart-services/pkg/hash"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbdto "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/dto"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"testing"
)

type E2ETestSuite struct {
	suite.Suite

	logger           *logrus.Entry
	postgreContainer *testpostgres.PostgresContainer
	redisContainer   *testredis.RedisContainer
	db               *sqlx.DB
	client           *redis.Client
}

func (ets *E2ETestSuite) BeforeAll(t provider.T) {
	var err error

	if ets.postgreContainer, err = GetPostgresForE2ETests(); err != nil {
		t.Fatal(err)
	}

	if ets.redisContainer, err = GetRedisForE2ETests(); err != nil {
		t.Fatal(err)
	}
}

func (ets *E2ETestSuite) BeforeEach(t provider.T) {
	connectionString, err := getGenericPostgresConnectionString(context.Background(), ets.postgreContainer)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf(
		"%sbooksmart?sslmode=disable&search_path=bs",
		connectionString,
	)

	if ets.db, err = GetPostgresClientForE2ETests(url); err != nil {
		t.Fatal(err)
	}

	if ets.client, err = GetRedisClientForE2ETests(ets.redisContainer); err != nil {
		t.Fatal(err)
	}
}

func (ets *E2ETestSuite) AfterEach(t provider.T) {
	if err := ets.db.Close(); err != nil {
		t.Fatalf("failed to close connection to postgres container: %s", err)
	}

	if err := ets.client.Close(); err != nil {
		t.Fatalf("failed to close redis connection to postgres container: %s", err)
	}
}

func (ets *E2ETestSuite) AfterAll(t provider.T) {
	if err := ets.postgreContainer.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}

	if err := ets.redisContainer.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func (ets *E2ETestSuite) Test_ReserveBook(t provider.T) {
	t.Title("E2E Test For Reserve Book")
	t.Description("Reader SignUp, SignIn and Reserve Book")
	if isIntegrationTestsFailed() {
		t.Skip()
	}

	ctx := context.Background()
	reservationRepo := implRepo.NewReservationRepo(ets.db, logging.GetLoggerForTests())
	bookRepo := implRepo.NewBookRepo(ets.db, logging.GetLoggerForTests())
	readerRepo := implRepo.NewReaderRepo(ets.db, ets.client, logging.GetLoggerForTests())
	libCardRepo := implRepo.NewLibCardRepo(ets.db, logging.GetLoggerForTests())
	tokenManager, err := auth.NewTokenManager("signing_key")
	if err != nil {
		t.Fatal(err)
	}
	hasher := hash.NewPasswordHasher("salt_string")
	transactionManager, err := getTransactionManagerForE2ETests(ets.db)
	if err != nil {
		t.Fatal(err)
	}
	readerService := impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
		logging.GetLoggerForTests(), 1000, 2000)
	bookService := impl.NewBookService(bookRepo, logging.GetLoggerForTests())
	libCardService := impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
	reservationService := impl.NewReservationService(reservationRepo, bookRepo,
		readerRepo, libCardRepo,
		transactionManager, logging.GetLoggerForTests(),
	)

	reader := ommodels.NewReaderModelObjectMother().DefaultReader()
	password := reader.Password
	err = readerService.SignUp(ctx, reader)
	t.Assert().Nil(err)

	tokens, err := readerService.SignIn(ctx, reader.PhoneNumber, password)
	t.Assert().Nil(err)
	t.Assert().NotNil(tokens)

	params := tdbdto.NewBookParamsDTOBuilder().
		WithTitle("The Hunger Games").
		WithAuthor("Suzanne Collins").
		WithPublisher("Scholastic Press").
		WithCopiesNumber(14).
		WithRarity("Common").
		WithGenre("Young Adult,Fiction,Dystopia,Fantasy,Science Fiction,Romance,Adventure,Teen,Post Apocalyptic,Action").
		WithPublishingYear(2008).
		WithLanguage("English").
		WithAgeLimit(6).Build()
	findBooks, err := bookService.GetByParams(ctx, params)
	t.Assert().Nil(err)
	t.Assert().Len(findBooks, 1)

	needBook := findBooks[0]

	err = libCardService.Create(ctx, reader.ID)
	t.Assert().Nil(err)

	err = reservationService.Create(ctx, reader.ID, needBook.ID)
	t.Assert().Nil(err)
}

func TestE2ETestSuite(t *testing.T) {
	suite.RunSuite(t, new(E2ETestSuite))
}

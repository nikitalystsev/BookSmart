package integrationTests

import (
	"BookSmart/internal/repositories/implRepo/postgres"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/internal/services/implServices"
	"BookSmart/internal/services/intfServices"
	"BookSmart/pkg/auth"
	"BookSmart/pkg/hash"
	"BookSmart/pkg/logging"
	"BookSmart/pkg/transact"
	"errors"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	migrations "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	// database
	db     *sqlx.DB
	client *redis.Client

	// services
	bookService        intfServices.IBookService
	libCardService     intfServices.ILibCardService
	readerService      intfServices.IReaderService
	reservationService intfServices.IReservationService

	// repositories
	bookRepo        intfRepo.IBookRepo
	libCardRepo     intfRepo.ILibCardRepo
	readerRepo      intfRepo.IReaderRepo
	reservationRepo intfRepo.IReservationRepo

	hasher             hash.IPasswordHasher
	tokenManager       auth.ITokenManager
	transactionManager transact.ITransactionManager
}

func (s *IntegrationTestSuite) SetupSuite() {
	dsn := "postgres://postgres:postgres@localhost:5437/testdb?sslmode=disable"

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		s.FailNow("Failed to ping database: " + err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6381",
		Password: "",
		DB:       0,
	})

	s.db = db
	s.client = client
	s.initDeps()

	if err = s.migrateUpDB(); err != nil {
		s.FailNowf("Failed to migrate up DB", err.Error())
	}
}

func (s *IntegrationTestSuite) initDeps() {
	logger := logging.GetLoggerForTests()

	tokenManager, err := auth.NewTokenManager("signing_key")
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	hasher := hash.NewPasswordHasher("salt_string")

	_manager, err := manager.New(trmsqlx.NewDefaultFactory(s.db))
	if err != nil {
		panic(err)
	}

	transactionManager := transact.NewTransactionManager(_manager)
	s.bookRepo = postgres.NewBookRepo(s.db, logger)
	s.libCardRepo = postgres.NewLibCardRepo(s.db, logger)
	s.readerRepo = postgres.NewReaderRepo(s.db, s.client, logger)
	s.reservationRepo = postgres.NewReservationRepo(s.db, logger)

	s.bookService = implServices.NewBookService(s.bookRepo, logger)
	s.libCardService = implServices.NewLibCardService(s.libCardRepo, logger)
	s.readerService = implServices.NewReaderService(s.readerRepo, s.bookRepo, tokenManager, hasher, logger)
	s.reservationService = implServices.NewReservationService(s.reservationRepo, s.bookRepo, s.readerRepo, s.libCardRepo, transactionManager, logger)

	s.hasher = hasher
	s.tokenManager = tokenManager
	s.transactionManager = transactionManager
}

func (s *IntegrationTestSuite) TearDownSuite() {
	err := s.migrateDownDB()
	if err != nil {
		fmt.Println("Failed to migrate down DB")
	}

	err = s.db.Close()
	if err != nil {
		fmt.Println("Failed to close DB", err)
	}
}

func (s *IntegrationTestSuite) migrateUpDB() error {
	driver, err := migrations.WithInstance(s.db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://C:/Users/nikitalystsev/Documents/bmstu/ppo/BookSmart/src/internal/tests/integrationTests/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) migrateDownDB() error {
	driver, err := migrations.WithInstance(s.db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://C:/Users/nikitalystsev/Documents/bmstu/ppo/BookSmart/src/internal/tests/integrationTests/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(IntegrationTestSuite))
}

package integrationTests

import (
	implRepo "BookSmart-postgres/impl"
	"BookSmart-services/impl"
	"BookSmart-services/intf"
	"BookSmart-services/intfRepo"
	"BookSmart-services/pkg/auth"
	"BookSmart-services/pkg/hash"
	"BookSmart-services/pkg/transact"
	"Booksmart/pkg/logging"
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
	"time"
)

type IntegrationTestSuite struct {
	suite.Suite

	// database
	db     *sqlx.DB
	client *redis.Client

	// services
	bookService        intf.IBookService
	libCardService     intf.ILibCardService
	readerService      intf.IReaderService
	reservationService intf.IReservationService

	// repositories
	bookRepo        intfRepo.IBookRepo
	libCardRepo     intfRepo.ILibCardRepo
	readerRepo      intfRepo.IReaderRepo
	reservationRepo intfRepo.IReservationRepo

	hasher             hash.IPasswordHasher
	tokenManager       auth.ITokenManager
	transactionManager transact.ITransactionManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (s *IntegrationTestSuite) SetupSuite() {
	dsn := os.Getenv("DB_DSN_TEST")

	fmt.Printf("%s", dsn)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		s.FailNow("Failed to ping database: " + err.Error())
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6380",
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_USER_PASSWORD"),
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

	s.accessTokenTTL = time.Hour
	s.refreshTokenTTL = time.Hour * 24

	transactionManager := transact.NewTransactionManager(_manager)
	s.bookRepo = implRepo.NewBookRepo(s.db, logger)
	s.libCardRepo = implRepo.NewLibCardRepo(s.db, logger)
	s.readerRepo = implRepo.NewReaderRepo(s.db, s.client, logger)
	s.reservationRepo = implRepo.NewReservationRepo(s.db, logger)

	s.bookService = impl.NewBookService(s.bookRepo, logger)
	s.libCardService = impl.NewLibCardService(s.libCardRepo, logger)
	s.readerService = impl.NewReaderService(s.readerRepo, s.bookRepo, tokenManager, hasher, logger, s.accessTokenTTL, s.refreshTokenTTL)
	s.reservationService = impl.NewReservationService(s.reservationRepo, s.bookRepo, s.readerRepo, s.libCardRepo, transactionManager, logger)

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
		os.Getenv("DB_MIGRATION_PATH_TEST"),
		"postgres", driver,
	)
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
		os.Getenv("DB_MIGRATION_PATH_TEST"),
		"postgres", driver,
	)
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

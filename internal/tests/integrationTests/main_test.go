package integrationTests

import (
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
	"github.com/joho/godotenv"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	"github.com/nikitalystsev/BookSmart-services/pkg/auth"
	"github.com/nikitalystsev/BookSmart-services/pkg/hash"
	"github.com/nikitalystsev/BookSmart-services/pkg/transact"
	"github.com/stretchr/testify/suite"
	"log"
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
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := s.createDBMigration(); err != nil {
		log.Fatal(err)
	}
	if err := s.createSchemaMigration(); err != nil {
		log.Fatal(err)
	}
	db, err := s.fillDBMigration()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST_TEST") + ":" + os.Getenv("REDIS_PORT_TEST"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_USER_PASSWORD"),
		DB:       0,
	})

	s.db = db
	s.client = client
	s.initDeps()
}

func (s *IntegrationTestSuite) createDBMigration() error {
	dsn := os.Getenv("POSTGRES_CREATE_TEST_DB_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m1, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_CREATE_TEST_DB_MIGRATION_PATH"),
		"postgres", driver,
	)
	if err != nil {
		fmt.Println(os.Getenv("POSTGRES_CREATE_TEST_DB_MIGRATION_PATH"))
		fmt.Println("error")
		return err
	}

	err = m1.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) createSchemaMigration() error {
	dsn := os.Getenv("POSTGRES_CREATE_TEST_SCHEMA_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m1, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_CREATE_TEST_SCHEMA_MIGRATION_PATH"),
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	err = m1.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) fillDBMigration() (*sqlx.DB, error) {
	dsn := os.Getenv("POSTGRES_FILL_TEST_DB_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
		return nil, err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return nil, err
	}

	m1, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_FILL_TEST_DB_MIGRATION_PATH"),
		"postgres", driver,
	)
	if err != nil {
		return nil, err
	}

	err = m1.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return db, err
}

func (s *IntegrationTestSuite) initDeps() {
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
	s.bookRepo = implRepo.NewBookRepo(s.db, logging.GetLoggerForTests())
	s.libCardRepo = implRepo.NewLibCardRepo(s.db, logging.GetLoggerForTests())
	s.readerRepo = implRepo.NewReaderRepo(s.db, s.client, logging.GetLoggerForTests())
	s.reservationRepo = implRepo.NewReservationRepo(s.db, logging.GetLoggerForTests())

	s.bookService = impl.NewBookService(s.bookRepo, logging.GetLoggerForTests())
	s.libCardService = impl.NewLibCardService(s.libCardRepo, logging.GetLoggerForTests())
	s.readerService = impl.NewReaderService(s.readerRepo, s.bookRepo, tokenManager, hasher, logging.GetLoggerForTests(), s.accessTokenTTL, s.refreshTokenTTL)
	s.reservationService = impl.NewReservationService(s.reservationRepo, s.bookRepo, s.readerRepo, s.libCardRepo, transactionManager, logging.GetLoggerForTests())

	s.hasher = hasher
	s.tokenManager = tokenManager
	s.transactionManager = transactionManager
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if err := s.clearDBMigration(); err != nil {
		fmt.Println("Failed clearDBMigration")
	}
	if err := s.dropSchemaMigration(); err != nil {
		fmt.Println("Failed dropSchemaMigration")
	}
	if err := s.dropDBMigration(); err != nil {
		fmt.Println("Failed dropDBMigration")
	}

	if err := s.db.Close(); err != nil {
		fmt.Println("Failed to close DB", err)
	}
}

func (s *IntegrationTestSuite) dropDBMigration() error {
	dsn := os.Getenv("POSTGRES_CREATE_TEST_DB_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m1, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_CREATE_TEST_DB_MIGRATION_PATH"),
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	err = m1.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) dropSchemaMigration() error {
	dsn := os.Getenv("POSTGRES_CREATE_TEST_SCHEMA_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m1, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_CREATE_TEST_SCHEMA_MIGRATION_PATH"),
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	err = m1.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) clearDBMigration() error {
	dsn := os.Getenv("POSTGRES_FILL_TEST_DB_URL")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		s.FailNow("Failed to connect to database: " + err.Error())
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m1, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_FILL_TEST_DB_MIGRATION_PATH"),
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	err = m1.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return err
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

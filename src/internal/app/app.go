package app

import (
	"BookSmart/internal/repositories/implRepo/postgres"
	"BookSmart/internal/services/implServices"
	"BookSmart/internal/ui/cli"
	"BookSmart/internal/ui/cli/handlers"
	"BookSmart/pkg/auth"
	"BookSmart/pkg/hash"
	"BookSmart/pkg/logging"
	"BookSmart/pkg/transact"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func Run() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DB_DSN")

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_USER_PASSWORD"),
		DB:       0,
	})

	logger := logging.GetLoggerForTests()

	tokenManager, err := auth.NewTokenManager("signing_key")
	if err != nil {
		panic(err)
	}

	hasher := hash.NewPasswordHasher("salt_string")

	_manager, err := manager.New(trmsqlx.NewDefaultFactory(db))
	if err != nil {
		panic(err)
	}
	transactionManager := transact.NewTransactionManager(_manager)

	bookRepo := postgres.NewBookRepo(db, logger)
	libCardRepo := postgres.NewLibCardRepo(db, logger)
	readerRepo := postgres.NewReaderRepo(db, client, logger)
	reservationRepo := postgres.NewReservationRepo(db, logger)

	bookService := implServices.NewBookService(bookRepo, logger)
	libCardService := implServices.NewLibCardService(libCardRepo, logger)
	readerService := implServices.NewReaderService(readerRepo, bookRepo, tokenManager, hasher, logger)
	reservationService := implServices.NewReservationService(reservationRepo, bookRepo, readerRepo, libCardRepo, transactionManager, logger)

	bookHandler := handlers.NewBookHandler(bookService, logger)
	libCardHandler := handlers.NewLibCardHandler(libCardService, logger)
	readerHandler := handlers.NewReaderHandler(readerService, bookService, libCardService, reservationService, logger)
	reservationHandler := handlers.NewReservationHandler(reservationService, logger)

	server := cli.NewServer(bookHandler, libCardHandler, readerHandler, reservationHandler)

	server.Run()
}

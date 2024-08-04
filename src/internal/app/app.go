package app

import (
	"BookSmart-repositories/impl/postgres"
	"BookSmart-services/impl"
	"BookSmart-services/pkg/auth"
	"BookSmart-services/pkg/hash"
	"BookSmart-services/pkg/transact"
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/requesters"
	"Booksmart/internal/config"
	"Booksmart/pkg/logging"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Run(configDir string) {
	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Init(configDir)
	if err != nil {
		logger.Errorf("error initializing config: %v", err)
		return
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.DBName,
		cfg.Postgres.Password, cfg.Postgres.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		logger.Errorf("error connecting to database: %v", err)
		return
	}

	err = db.Ping()
	if err != nil {
		logger.Errorf("error pinging database: %v", err)
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	tokenManager, err := auth.NewTokenManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Errorf("error initializing token manager: %v", err)
		return
	}

	hasher := hash.NewPasswordHasher(cfg.Auth.PasswordSalt)

	_manager, err := manager.New(trmsqlx.NewDefaultFactory(db))
	if err != nil {
		logger.Errorf("error initializing manager: %v", err)
		return
	}

	transactionManager := transact.NewTransactionManager(_manager)

	bookRepo := postgres.NewBookRepo(db, logger)
	libCardRepo := postgres.NewLibCardRepo(db, logger)
	readerRepo := postgres.NewReaderRepo(db, client, logger)
	reservationRepo := postgres.NewReservationRepo(db, logger)

	bookService := impl.NewBookService(bookRepo, logger)
	libCardService := impl.NewLibCardService(libCardRepo, logger)
	readerService := impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher, logger, cfg.Auth.JWT.AccessTokenTTL, cfg.Auth.JWT.RefreshTokenTTL)
	reservationService := impl.NewReservationService(reservationRepo, bookRepo, readerRepo, libCardRepo, transactionManager, logger)

	handler := handlers.NewHandler(bookService, libCardService, readerService, reservationService, tokenManager, logger)

	router := handler.InitRoutes()

	go func() {
		err = router.Run(":" + cfg.Port)
		if err != nil {
			logger.Errorf("error running server: %v", err)
			return
		}
	}()

	requester := requesters.NewRequester(logger, cfg.Auth.JWT.AccessTokenTTL, cfg.Auth.JWT.RefreshTokenTTL)
	requester.Run()
}

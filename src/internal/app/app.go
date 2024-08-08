package app

import (
	repoPostgres "BookSmart-postgres"
	implPostgres "BookSmart-postgres/impl"
	"BookSmart-services/impl"
	"BookSmart-services/intfRepo"
	"BookSmart-services/pkg/auth"
	"BookSmart-services/pkg/hash"
	"BookSmart-services/pkg/transact"
	"BookSmart-techUI/handlers"
	"BookSmart-techUI/requesters"
	repoMongo "Booksmart-mongo"
	implMongo "Booksmart-mongo/impl"
	"Booksmart/internal/config"
	"Booksmart/pkg/logging"
	"fmt"
	trmmongo "github.com/avito-tech/go-transaction-manager/drivers/mongo/v2"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/go-redis/redis/v8"
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

	var (
		bookRepo        intfRepo.IBookRepo
		libCardRepo     intfRepo.ILibCardRepo
		readerRepo      intfRepo.IReaderRepo
		reservationRepo intfRepo.IReservationRepo

		_manager *manager.Manager
	)

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	switch cfg.DBType {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.DBName,
			cfg.Postgres.Password, cfg.Postgres.SSLMode)

		db, err := repoPostgres.NewClient(dsn)
		if err != nil {
			logger.Errorf("error connect to postgres: %v", err)
			return
		}

		bookRepo = implPostgres.NewBookRepo(db, logger)
		libCardRepo = implPostgres.NewLibCardRepo(db, logger)
		readerRepo = implPostgres.NewReaderRepo(db, client, logger)
		reservationRepo = implPostgres.NewReservationRepo(db, logger)

		_manager, err = manager.New(trmsqlx.NewDefaultFactory(db))
		if err != nil {
			logger.Errorf("error initializing manager: %v", err)
			return
		}

	default:
		fmt.Println("choice branch with mongodb")
		mongoClient, err := repoMongo.NewClient(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.DBName)
		if err != nil {
			logger.Error(err)
			return
		}
		db := mongoClient.Database(cfg.Mongo.DBName)

		bookRepo = implMongo.NewBookRepo(db, logger)
		libCardRepo = implMongo.NewLibCardRepo(db, logger)
		readerRepo = implMongo.NewReaderRepo(db, client, logger)
		reservationRepo = implMongo.NewReservationRepo(db, logger)

		_manager, err = manager.New(trmmongo.NewDefaultFactory(mongoClient))
		if err != nil {
			logger.Errorf("error initializing manager: %v", err)
			return
		}
	}

	tokenManager, err := auth.NewTokenManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Errorf("error initializing token manager: %v", err)
		return
	}

	hasher := hash.NewPasswordHasher(cfg.Auth.PasswordSalt)

	transactionManager := transact.NewTransactionManager(_manager)

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

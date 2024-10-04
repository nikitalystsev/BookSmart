package app

import (
	"fmt"
	trmmongo "github.com/avito-tech/go-transaction-manager/drivers/mongo/v2"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	_ "github.com/lib/pq"
	repoMongo "github.com/nikitalystsev/BookSmart-repo-mongo"
	implMongo "github.com/nikitalystsev/BookSmart-repo-mongo/impl"
	repoPostgres "github.com/nikitalystsev/BookSmart-repo-postgres"
	implPostgres "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	"github.com/nikitalystsev/BookSmart-services/pkg/auth"
	"github.com/nikitalystsev/BookSmart-services/pkg/hash"
	"github.com/nikitalystsev/BookSmart-services/pkg/transact"
	"github.com/nikitalystsev/BookSmart-web-api/handlers"
	"github.com/nikitalystsev/BookSmart/internal/config"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/redis/go-redis/v9"
)

func Run(configDir string) {
	cfg, err := config.Init(configDir)
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}

	var (
		bookRepo        intfRepo.IBookRepo
		libCardRepo     intfRepo.ILibCardRepo
		readerRepo      intfRepo.IReaderRepo
		reservationRepo intfRepo.IReservationRepo
		ratingRepo      intfRepo.IRatingRepo

		trm *manager.Manager
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
		ratingRepo = implPostgres.NewRatingRepo(db, logger)

		trm, err = manager.New(trmsqlx.NewDefaultFactory(db))
		if err != nil {
			logger.Errorf("error initializing manager: %v", err)
			return
		}
	default:
		mongoClient, err := repoMongo.NewClient(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.DBName)
		if err != nil {
			logger.Errorf("error connect to mongo: %v", err)
			return
		}
		db := mongoClient.Database(cfg.Mongo.DBName)

		bookRepo = implMongo.NewBookRepo(db, logger)
		libCardRepo = implMongo.NewLibCardRepo(db, logger)
		readerRepo = implMongo.NewReaderRepo(db, client, logger)
		reservationRepo = implMongo.NewReservationRepo(db, logger)
		ratingRepo = implMongo.NewRatingRepo(db, logger)

		trm, err = manager.New(trmmongo.NewDefaultFactory(mongoClient))
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

	transactionManager := transact.NewTransactionManager(trm)

	bookService := impl.NewBookService(bookRepo, logger)
	libCardService := impl.NewLibCardService(libCardRepo, logger)
	readerService := impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher, logger, cfg.Auth.JWT.AccessTokenTTL, cfg.Auth.JWT.RefreshTokenTTL)
	reservationService := impl.NewReservationService(reservationRepo, bookRepo, readerRepo, libCardRepo, transactionManager, logger)
	ratingService := impl.NewRatingService(ratingRepo, reservationRepo, logger)

	handler := handlers.NewHandler(
		bookService,
		libCardService,
		readerService,
		reservationService,
		ratingService,
		tokenManager,
		cfg.Auth.JWT.AccessTokenTTL,
		cfg.Auth.JWT.RefreshTokenTTL,
	)

	router := handler.InitRoutes()

	err = router.Run(":" + cfg.Port)
	if err != nil {
		logger.Errorf("error running server: %v", err)
		return
	}
}

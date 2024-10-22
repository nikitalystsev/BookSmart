package e2eTest

import (
	"context"
	"errors"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/docker/docker/api/types/container"
	"github.com/golang-migrate/migrate/v4"
	migrations "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/nikitalystsev/BookSmart-services/pkg/transact"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func GetRedisForE2ETests() (*testredis.RedisContainer, error) {
	ctx := context.Background()

	redisContainer, err := testredis.Run(
		ctx,
		"redis:latest",
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
	)

	if err != nil {
		fmt.Printf("Failed to start postgres container: %v\n", err)
		return nil, err
	}

	return redisContainer, nil
}

func GetRedisClientForE2ETests(container *testredis.RedisContainer) (*redis.Client, error) {
	ctx := context.Background()
	uri, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: uri[8:],
	})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetPostgresForE2ETests() (*testpostgres.PostgresContainer, error) {
	ctx := context.Background()

	if err := godotenv.Load("../../../.env"); err != nil {
		log.Print("No .env file found")
	}

	postgresContainer, err := testpostgres.Run(
		ctx,
		"postgres:latest",
		testpostgres.WithDatabase("postgres"),
		testpostgres.WithUsername("postgres"),
		testpostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.Binds = []string{
				os.Getenv("DB_DATASETS_PATH_FOR_TESTS") + ":/data",
			}
		}),
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
	)

	if err != nil {
		fmt.Printf("Failed to start postgres container: %v\n", err)
		return nil, err
	}

	db, err := ApplyMigrations(postgresContainer)
	if err != nil {
		fmt.Printf("Failed to apply migrations: %v\n", err)
		return nil, err
	}

	if err = db.Close(); err != nil {
		fmt.Printf("Failed to close postgres container: %v\n", err)
		return nil, err
	}

	return postgresContainer, err
}

func ApplyMigrations(container *testpostgres.PostgresContainer) (*sqlx.DB, error) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Print("No .env file found")
	}

	ctx := context.Background()
	genericConnStr, err := getGenericPostgresConnectionString(ctx, container)
	if err != nil {
		return nil, err
	}

	if err = createDBMigration(genericConnStr); err != nil {
		return nil, err
	}
	if err = createSchemaMigration(genericConnStr); err != nil {
		return nil, err
	}
	db, err := fillDBMigration(genericConnStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDBMigration(genericConnStr string) error {
	dsn := fmt.Sprintf("%s?sslmode=disable", genericConnStr)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	genericMigrationPath, err := getGenericPostgresMigrationPath()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+genericMigrationPath+"/create_db",
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

func createSchemaMigration(genericConnStr string) error {
	dsn := fmt.Sprintf("%sbooksmart?sslmode=disable", genericConnStr)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	genericMigrationPath, err := getGenericPostgresMigrationPath()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+genericMigrationPath+"/create_schema",
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

func fillDBMigration(genericConnStr string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%sbooksmart?sslmode=disable&search_path=bs", genericConnStr)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return nil, err
	}

	genericMigrationPath, err := getGenericPostgresMigrationPath()
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+genericMigrationPath+"/fill_db",
		"postgres", driver,
	)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return db, err
}

func clearDBMigration(genericConnStr string) error {
	dsn := fmt.Sprintf("%sbooksmart?sslmode=disable&search_path=bs", genericConnStr)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	genericMigrationPath, err := getGenericPostgresMigrationPath()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+genericMigrationPath+"/fill_db",
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return err
}

func GetPostgresClientForE2ETests(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func parsePostgresURL(postgresURL string) (string, error) {
	parsedURL, err := url.Parse(postgresURL)
	if err != nil {
		return "", err
	}

	password, isSet := parsedURL.User.Password()
	if !isSet {
		return "", errors.New("postgres URL does not contain a password")
	}

	result := fmt.Sprintf(
		"%s://%s:%s@%s:%s/",
		parsedURL.Scheme,
		parsedURL.User.Username(),
		password,
		parsedURL.Hostname(),
		parsedURL.Port(),
	)

	return result, nil
}

func getGenericPostgresConnectionString(ctx context.Context, container *testpostgres.PostgresContainer) (string, error) {
	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		return "", err
	}

	genericConnectionString, err := parsePostgresURL(connectionString)
	if err != nil {
		return "", err
	}

	return genericConnectionString, nil
}

func getGenericPostgresMigrationPath() (string, error) {
	relativePath := "../../../components/component-repo-postgres/impl/migrations"
	absolutePath, err := filepath.Abs(relativePath)

	if err != nil {
		return "", err
	}

	return filepath.ToSlash(absolutePath), nil
}

func isIntegrationTestsFailed() bool {
	if os.Getenv("INTEGRATION_TESTS_IS_SUCCESS") == "1" {
		return false
	}

	return true
}

func getTransactionManagerForE2ETests(db *sqlx.DB) (transact.ITransactionManager, error) {
	trm, err := manager.New(trmsqlx.NewDefaultFactory(db))
	if err != nil {
		return nil, err
	}

	transactionManager := transact.NewTransactionManager(trm)

	return transactionManager, err
}

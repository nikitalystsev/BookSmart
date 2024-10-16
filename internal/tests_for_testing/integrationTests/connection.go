package integrationTest

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/golang-migrate/migrate/v4"
	migrations "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"os"
	"time"
)

func GetRedisForIntegrationTests() (*testredis.RedisContainer, error) {
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

func GetRedisClientForIntegrationTests(container *testredis.RedisContainer) (*redis.Client, error) {
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

func GetPostgresForIntegrationTests() (*testpostgres.PostgresContainer, error) {
	ctx := context.Background()

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
				"/c/Users/nikitalystsev/Documents/bmstu/ppo/BookSmart/data/mydatasets:/data",
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

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	if err = createDBMigration(port.Port()); err != nil {
		return nil, err
	}
	if err = createSchemaMigration(port.Port()); err != nil {
		return nil, err
	}
	db, err := fillDBMigration(port.Port())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDBMigration(port string) error {
	dsn := fmt.Sprintf("postgres://postgres:postgres@0.0.0.0:%s/?sslmode=disable", port)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_CREATE_TEST_DB_MIGRATION_PATH"),
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

func createSchemaMigration(port string) error {
	dsn := fmt.Sprintf("postgres://postgres:postgres@0.0.0.0:%s/booksmart?sslmode=disable", port)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_CREATE_TEST_SCHEMA_MIGRATION_PATH"),
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

func fillDBMigration(port string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://postgres:postgres@0.0.0.0:%s/booksmart?sslmode=disable&search_path=bs", port)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := migrations.WithInstance(db.DB, &migrations.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("POSTGRES_FILL_TEST_DB_MIGRATION_PATH"),
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

func GetPostgresClientForIntegrationTests(url string) (*sqlx.DB, error) {
	fmt.Printf("url: %s\n", url)
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

package serviceTests

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/golang-migrate/migrate/v4"
	migrations "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"os"
	"time"
)

func getContainerForClassicUnitTests() (*postgres.PostgresContainer, error) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.Binds = []string{
				"/c/Users/nikitalystsev/Documents/bmstu/ppo/BookSmart/docs/data/mydatasets:/data",
			}
		}),
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
	)

	if err != nil {
		fmt.Printf("Failed to start postgres container: %v\n", err)
		return nil, err
	}

	return postgresContainer, err
}

func applyMigrations(container *postgres.PostgresContainer) (*sqlx.DB, error) {
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

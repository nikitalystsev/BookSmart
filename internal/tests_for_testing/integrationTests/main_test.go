package integrationTest

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/sirupsen/logrus"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	logger           *logrus.Entry
	postgreContainer *testpostgres.PostgresContainer
	db               *sqlx.DB
}

func (its *IntegrationTestSuite) BeforeAll(t provider.T) {
	var err error

	if its.postgreContainer, err = GetPostgresForIntegrationTests(); err != nil {
		t.Fatal(err)
	}
}

func (its *IntegrationTestSuite) BeforeEach(t provider.T) {
	port, err := its.postgreContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf(
		"postgres://postgres:postgres@0.0.0.0:%s/booksmart?sslmode=disable&search_path=bs",
		port.Port(),
	)

	if its.db, err = GetPostgresClientForIntegrationTests(url); err != nil {
		t.Fatal(err)
	}
}

func (its *IntegrationTestSuite) AfterEach(t provider.T) {
	if err := its.db.Close(); err != nil {
		t.Fatalf("failed to close connection to postgres container: %s", err)
	}
}

func (its *IntegrationTestSuite) AfterAll(t provider.T) {
	if err := its.postgreContainer.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(IntegrationTestSuite))
}

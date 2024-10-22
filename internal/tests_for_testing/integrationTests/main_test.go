package integrationTests

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	logger           *logrus.Entry
	postgreContainer *testpostgres.PostgresContainer
	redisContainer   *testredis.RedisContainer
	db               *sqlx.DB
	client           *redis.Client
}

func (its *IntegrationTestSuite) BeforeAll(t provider.T) {
	var err error

	if its.postgreContainer, err = GetPostgresForIntegrationTests(); err != nil {
		t.Fatal(err)
	}

	if its.redisContainer, err = GetRedisForIntegrationTests(); err != nil {
		t.Fatal(err)
	}
}

func (its *IntegrationTestSuite) BeforeEach(t provider.T) {
	connectionString, err := getGenericPostgresConnectionString(context.Background(), its.postgreContainer)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf(
		"%sbooksmart?sslmode=disable&search_path=bs",
		connectionString,
	)

	if its.db, err = GetPostgresClientForIntegrationTests(url); err != nil {
		t.Fatal(err)
	}

	if its.client, err = GetRedisClientForIntegrationTests(its.redisContainer); err != nil {
		t.Fatal(err)
	}
}

func (its *IntegrationTestSuite) AfterAll(t provider.T) {
	if err := its.postgreContainer.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}

	if err := its.redisContainer.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.RunSuite(t, new(IntegrationTestSuite))
}

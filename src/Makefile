
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

migration_path := $(DB_MIGRATION_PATH)
database_url := '$(DB_DSN)'

.PHONY: run build utest-srv utest-repo itest migrate-up migrate-down clean

all:
	docker compose up bs-app bs-postgres bs-mongo bs-redis

build:
	docker build -t booksmart:local .

utest-srv:
	go test -v ./internal/tests/unitTests/serviceTests/

utest-repo:
	go test -v ./internal/tests/unitTests/repositoryTests/

itest:
	docker compose up -d test-postgres-db test-redis
	go test -v ./internal/tests/integrationTests
	docker stop test_postgres_db test_redis

migrate-up:
	migrate -database $(database_url) -path $(migration_path) up

migrate-down:
	migrate -database $(database_url) -path $(migration_path) down

mmigrate-up:
	cd internal/repoMongo/impl/migrations && migrate-mongo up
	cd docs/data/mydatasets/ && python to_mongodb.py

mmigrate-down:
	cd internal/repoMongo/impl/migrations && migrate-mongo down

clean:
	rm *.exe

ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

.PHONY: run build utest-srv utest-repo itest migrate-up migrate-down clean

all:build
	./app

build:
	docker compose up -d bs-postgres bs-mongo bs-redis
	go build -o app cmd/app/main.go

utest-srv:
	go test -v ./internal/tests/unitTests/serviceTests/

utest-repo:
	go test -v ./internal/tests/unitTests/repositoryTests/

itest:
	docker compose up -d bs-postgres-test bs-redis-test
	go test -v ./internal/tests/integrationTests

migrate-up:
	migrate -database '$(POSTGRES_CREATE_DB_URL)' -path $(POSTGRES_CREATE_DB_MIGRATION_PATH) up
	migrate -database '$(POSTGRES_CREATE_SCHEMA_URL)' -path $(POSTGRES_CREATE_SCHEMA_MIGRATION_PATH) up
	migrate -database '$(POSTGRES_FILL_DB_URL)' -path $(POSTGRES_FILL_DB_MIGRATION_PATH) up

migrate-down:
	migrate -database '$(POSTGRES_FILL_DB_URL)' -path $(POSTGRES_FILL_DB_MIGRATION_PATH) down
	migrate -database '$(POSTGRES_CREATE_SCHEMA_URL)' -path $(POSTGRES_CREATE_SCHEMA_MIGRATION_PATH) down
	migrate -database '$(POSTGRES_CREATE_DB_URL)' -path $(POSTGRES_CREATE_DB_MIGRATION_PATH) down


mmigrate-up:
	cd internal/repoMongo/impl/migrations && migrate-mongo up
	cd docs/data/mydatasets/ && python to_mongodb.py

mmigrate-down:
	cd internal/repoMongo/impl/migrations && migrate-mongo down

clean:
	rm *.exe
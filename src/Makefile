
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

migration_path := $(DB_MIGRATION_PATH)
database_url := '$(DB_DSN)'

.PHONY: run build migrate_up migrate_down

run: build
	./main.exe

build:
	docker compose up -d postgres-db redis
	go build cmd/app/main.go

utest_srv:
	go test -v ./internal/tests/unitTests/serviceTests/

utest_repo:
	go test -v ./internal/tests/unitTests/repositoryTests/

itest:
	docker compose up -d test-postgres-db test-redis
	go test -v ./internal/tests/integrationTests

migrate_up:
	migrate -database $(database_url) -path $(migration_path) up

migrate_down:
	migrate -database $(database_url) -path $(migration_path) down

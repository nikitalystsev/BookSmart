
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

migration_path := $(DB_MIGRATION_PATH)
database_url := '$(DB_DSN)'

.PHONY: run build utest-srv utest-repo itest migrate-up migrate-down clean

run: build
	./main.exe

build:
	docker compose up -d postgres-db redis mongo-db
	go build cmd/app/main.go

utest-srv:
	go test -v ./internal/tests/unitTests/serviceTests/

utest-repo:
	go test -v ./internal/tests/unitTests/repositoryTests/

itest:
	docker compose up -d test-postgres-db test-redis
	go test -v ./internal/tests/integrationTests

migrate-up:
	migrate -database $(database_url) -path $(migration_path) up

migrate-down:
	migrate -database $(database_url) -path $(migration_path) down

mmigrate-up:
	cd internal/repoMongo/impl/migrations && migrate-mongo up
	eval $(MIGRATION_COMMAND)

mmigrate-down:
	cd internal/repoMongo/impl/migrations && migrate-mongo down

clean:
	rm *.exe
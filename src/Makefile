
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

migration_path := $(DB_MIGRATION_PATH)
database_url := '$(DB_DSN)'

build:
	docker compose up -d postgres-db redis
	go build cmd/app/main.go

migration_up:
	migrate -database $(database_url) -path $(migration_path) up


migration_down:
	migrate -database $(database_url) -path $(migration_path) down

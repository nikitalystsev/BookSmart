
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

.PHONY: run build utest-srv utest-repo itest migrate-up migrate-down clean

run: build-ui run-app
	./techUI

build-ui:
	go build -o techUI cmd/techUI/main.go

run-app: build-app
	docker compose up -d bs-ppo-app bs-ppo-postgres bs-ppo-mongo bs-ppo-redis bs-ppo-nginx bs-ppo-pgadmin

build-app:
	docker build -t booksmart:local .

stop-app:
	docker stop bs-ppo-app bs-ppo-postgres bs-ppo-mongo bs-ppo-redis bs-ppo-nginx bs-ppo-pgadmin

rerun-app:
	make stop-app && make run-app

get-swagger:
	swag init -g cmd/app/main.go -o ./docs_swagger

# тесты ППО (исправить)
#utest-srv:
#	go test -v ./internal/tests/unitTests/serviceTests/
#
#utest-repo:
#	go test -v ./internal/tests/unitTests/repositoryTests/
#
#itest:
#	docker compose up -d bs-ppo-postgres-test bs-ppo-redis-test
#	go test -v ./internal/tests/integrationTests
#

# тесты тестирования
utest-srv:
	go test -v  -shuffle on ./internal/tests_for_testing/unitTests/
	cp ./internal/tests_for_testing/unitTests/environment.properties ./internal/tests_for_testing/unitTests/allure-results
	cd ./internal/tests_for_testing/unitTests/ && allure serve

migrate-up:
	migrate -database '$(POSTGRES_CREATE_DB_URL)' -path $(POSTGRES_CREATE_DB_MIGRATION_PATH) up
	migrate -database '$(POSTGRES_CREATE_SCHEMA_URL)' -path $(POSTGRES_CREATE_SCHEMA_MIGRATION_PATH) up
	migrate -database '$(POSTGRES_FILL_DB_URL)' -path $(POSTGRES_FILL_DB_MIGRATION_PATH) up

migrate-down:
	migrate -database '$(POSTGRES_FILL_DB_URL)' -path $(POSTGRES_FILL_DB_MIGRATION_PATH) down
	migrate -database '$(POSTGRES_CREATE_SCHEMA_URL)' -path $(POSTGRES_CREATE_SCHEMA_MIGRATION_PATH) down
	migrate -database '$(POSTGRES_CREATE_DB_URL)' -path $(POSTGRES_CREATE_DB_MIGRATION_PATH) down

mmigrate-up:
	cd ./components/component-repo-mongo/impl/migrations && migrate-mongo up
	cd data/mydatasets/ && python to_mongodb.py

mmigrate-down:
	cd ./components/component-repo-mongo/impl/migrations && migrate-mongo down

clean:
	rm *.exe ./app ./techUI
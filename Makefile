
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

.PHONY: run build utest-srv utest-repo itest migrate-up migrate-down clean

run: build-ui run-app
	./techUI

build-ui:
	go build -o techUI cmd/techUI/main.go

run-app: build-all
	docker compose up -d bs-app-gin bs-app-echo \
 		bs-postgres-master bs-redis bs-prometheus bs-grafana

build-all:
	docker compose build

stop-app:
	docker stop bs-app-gin bs-app-echo \
		bs-postgres-master bs-redis bs-prometheus bs-grafana

rerun-app:
	make stop-app && make run-app

# тесты ППО (исправить)
#utest-srv:
#	go tests_for_testing -v ./internal/tests/unitTests/serviceTests/
#
#utest-repo:
#	go tests_for_testing -v ./internal/tests/unitTests/repositoryTests/
#
#itest:
#	docker compose up -d bs-postgres-tests_for_testing bs-redis-tests_for_testing
#	go tests_for_testing -v ./internal/tests/integrationTests
#

# тесты тестирования
tests:
	./run_tests.sh

migrate-up:
	./migrate.sh up

migrate-down:
	./migrate.sh down

mmigrate-up:
	cd ./components/component-repo-mongo/impl/migrations && migrate-mongo up
	cd data/mydatasets/ && python to_mongodb.py

mmigrate-down:
	cd ./components/component-repo-mongo/impl/migrations && migrate-mongo down

clean:
	rm *.exe ./app ./techUI

# docker inspect --format='{{range .NetworkSettings.Networks}}{{.MacAddress}}{{end}}' $INSTANCE_ID -- посмотреть ip адрес сервера
.PHONY: run build utest-srv utest-repo itest migrate-up migrate-down clean

#run: build-ui run-app // не работает с текущим хендлером
#	./techUI
#
#build-ui:
#	go build -o techUI cmd/techUI/main.go


run-app: build-all
	docker compose up -d bs-app-main bs-app-inst1 bs-app-inst2 bs-app-mirror1 \
 		bs-postgres-master bs-postgres-slave bs-redis bs-nginx bs-pgadmin bs-react

build-all:
	docker compose build

stop-app:
	docker stop bs-app-main bs-app-inst1 bs-app-inst2 bs-app-mirror1 \
		bs-postgres-master bs-postgres-slave bs-redis bs-nginx bs-pgadmin bs-react

rerun-app:
	make stop-app && docker rm bs-nginx && make run-app

get-swagger:
	swag init -g cmd/app/main.go -o ./docs_swagger

# тесты ППО // частично не работают с текущей бизнес логикой
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

# тесты тестирования (они тоже не работают в местах, где я добавил пагинацию)
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

# docker inspect --format='{{range .NetworkSettings.Networks}}{{.MacAddress}}{{end}}' $INSTANCE_ID -- посмотреть ip адрес сервера
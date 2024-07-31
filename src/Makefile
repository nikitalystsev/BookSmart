
migration_path := internal/repositories/implRepo/postgres/migrations/
database_url := 'postgres://nikitalystsev:zhpiix69@localhost:5436/booksmart?sslmode=disable&search_path=bs'

migration_up:
	migrate -database $(database_url) -path $(migration_path) up


migration_down:
	migrate -database $(database_url) -path $(migration_path) down

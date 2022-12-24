DB_URL="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"

create_postgres:
	docker run --name postgres15 -p 5432:5432\
		-e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.1-alpine

start_postgres:
	docker start postgres15

stop_postgres:
	docker stop postgres15

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate --path db/migration -database\
		"$(DB_URL)" -verbose up

migratedown:
	migrate --path db/migration -database\
		"$(DB_URL)" -verbose down

migrateup1:
	migrate --path db/migration -database\
		"$(DB_URL)" -verbose up 1

migratedown1:
	migrate --path db/migration -database\
		"$(DB_URL)" -verbose down 1


sqlc:
	sqlc generate

test:
	# Test a specific function
	# go test -timeout 30s github.com/khiemledev/simple_bank_golang/db/sqlc -run ^TestCreateAccount$

	# Test cover all
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/khiemledev/simple_bank_golang/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup start_postgres stop_postgres sqlc server dbdocs

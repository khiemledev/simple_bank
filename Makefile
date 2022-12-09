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
		"postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate --path db/migration -database\
		"postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	# Test a specific function
	# go test -timeout 30s github.com/khiemledev/simple_bank_golang/db/sqlc -run ^TestCreateAccount$

	# Test cover all
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup start_postgres stop_postgres sqlc

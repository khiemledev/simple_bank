package main

import (
	"database/sql"
	"log"

	"github.com/khiemledev/simple_bank_golang/api"
	db "github.com/khiemledev/simple_bank_golang/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	serverAddress = "0.0.0.0:8080"
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", conn)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}

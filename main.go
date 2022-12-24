package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/khiemledev/simple_bank_golang/api"
	db "github.com/khiemledev/simple_bank_golang/db/sqlc"
	"github.com/khiemledev/simple_bank_golang/gapi"
	"github.com/khiemledev/simple_bank_golang/pb"
	"github.com/khiemledev/simple_bank_golang/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", conn)
	}

	store := db.NewStore(conn)
	startGrpcServer(config, store)

}

func startGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("grpc server listening on %s\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func startGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: %w", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}

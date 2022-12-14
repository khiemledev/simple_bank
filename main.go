package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/khiemledev/simple_bank_golang/api"
	db "github.com/khiemledev/simple_bank_golang/db/sqlc"
	"github.com/khiemledev/simple_bank_golang/gapi"
	"github.com/khiemledev/simple_bank_golang/pb"
	"github.com/khiemledev/simple_bank_golang/util"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/khiemledev/simple_bank_golang/doc/statik"
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

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	go startGatewayServer(config, store)
	startGrpcServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create migration:", err)
	}

	if err := migration.Up(); err != migrate.ErrNoChange {
		log.Fatal("cannot run migration up:", err)
	}

	log.Print("DBMigration successful!")
}

func startGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: %w", err)
	}

	muxOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(muxOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statikFS:", err)
	}
	fileHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", fileHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("HTTP gateway server listening on %s\n", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}

func startGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("grpc server listening on %s\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server:", err)
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

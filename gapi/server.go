package gapi

import (
	db "github.com/khiemledev/simple_bank_golang/db/sqlc"
	"github.com/khiemledev/simple_bank_golang/pb"
	"github.com/khiemledev/simple_bank_golang/token"
	"github.com/khiemledev/simple_bank_golang/util"
)

// Server servers all gRPC requests for banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates new gRPC Server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.CreatePasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}

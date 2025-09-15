package main

import (
	"net"

	"github.com/cprakhar/uber-clone/services/trip-service/repo"
	"github.com/cprakhar/uber-clone/services/trip-service/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewgRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) run() error {

	tripRepository := repo.NewInMemoRepository()
	_ = service.NewgRPCService(tripRepository)

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		// handle error
	}
	srv := grpc.NewServer()
	// srv.RegisterService()
	return srv.Serve(lis)
}
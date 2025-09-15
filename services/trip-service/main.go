package main

import (
	"github.com/cprakhar/uber-clone/shared/env"
)

var (
	grpcAddr = env.GetString("GRPC_ADDR", ":9000")
)

func main() {
	gRPCServer := NewgRPCServer(grpcAddr)
	gRPCServer.run()
}

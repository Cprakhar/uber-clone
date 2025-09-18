package handler

import (
	"context"

	"github.com/cprakhar/uber-clone/services/driver-service/service"
	pb "github.com/cprakhar/uber-clone/shared/proto/driver"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	svc service.GrpcDriverService
}

func NewgRPCHandler(srv *grpc.Server, svc service.GrpcDriverService) {
	handler := &gRPCHandler{svc: svc}
	pb.RegisterDriverServiceServer(srv, handler)
}

func (h *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driverID := req.GetDriverID()
	packageSlug := req.GetPackageSlug()
	// Implement the logic to register a driver
	h.svc.RegisterDriver(ctx, driverID, packageSlug)
	return nil, nil
}

func (h *gRPCHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	// Implement the logic to unregister a driver
	return nil, nil
}
package service

import (
	"context"

	"github.com/cprakhar/uber-clone/services/driver-service/repo"
	"github.com/cprakhar/uber-clone/services/driver-service/types"
)

type grpcDriverService struct {
	repo repo.DriverRepo
}

type GrpcDriverService interface {
	RegisterDriver(ctx context.Context, driverID, packageSlug string) (*types.DriverModel, error)
	UnregisterDriver(ctx context.Context, driverID string) error
}

func NewgRPCService(repo repo.DriverRepo) *grpcDriverService {
	return &grpcDriverService{repo: repo}
}

func (s *grpcDriverService) RegisterDriver(ctx context.Context, driverID, packageSlug string) (*types.DriverModel, error) {
	// profilePic := sharedUtil.GetRandomProfilePic()
	// carPlate := util.GenerateRandomPlate()
	// driver := &types.DriverModel{
	// 	ID:          driverID,
	// 	PackageSlug: packageSlug,
	// }
	// Implement the logic to register a driver
	return nil, nil
}

func (s *grpcDriverService) UnregisterDriver(ctx context.Context, driverID string) error {
	// Implement the logic to unregister a driver
	return nil
}

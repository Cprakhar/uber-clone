package service

import (
	"context"

	"github.com/cprakhar/uber-clone/services/trip-service/repo"
	"github.com/cprakhar/uber-clone/services/trip-service/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type grpcTripService struct {
	repo repo.TripRepo
}

type GRPCTripService interface {
	CreateTrip(ctx context.Context, fare types.RideFareModel) (*types.TripModel, error)
}

func NewgRPCService(repo repo.TripRepo) *grpcTripService {
	return &grpcTripService{repo: repo}
}

func (s *grpcTripService) CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error) {
	trip := &types.TripModel{
		ID:       primitive.NewObjectID(),
		RiderID:  fare.RiderID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.Create(ctx, trip)
}

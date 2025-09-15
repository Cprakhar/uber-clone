package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cprakhar/uber-clone/services/trip-service/repo"
	"github.com/cprakhar/uber-clone/services/trip-service/types"
	sharedtypes "github.com/cprakhar/uber-clone/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type grpcTripService struct {
	repo repo.TripRepo
}

type GRPCTripService interface {
	CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *sharedtypes.Coordinate) (*types.OSRMApiResponse, error)
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

func (s *grpcTripService) GetRoute(ctx context.Context, pickup, destination *sharedtypes.Coordinate) (*types.OSRMApiResponse, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM api: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var routeResponse types.OSRMApiResponse
	if err := json.Unmarshal(body, &routeResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &routeResponse, nil
}

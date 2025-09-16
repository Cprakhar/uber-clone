package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cprakhar/uber-clone/services/trip-service/repo"
	"github.com/cprakhar/uber-clone/services/trip-service/types"
	"github.com/cprakhar/uber-clone/shared/proto/trip"
	sharedtypes "github.com/cprakhar/uber-clone/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type grpcTripService struct {
	repo repo.TripRepo
}

type GRPCTripService interface {
	CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *sharedtypes.Coordinate) (*types.OSRMApiResponse, error)
	EstimatePackagesPriceWithRoute(route *types.OSRMApiResponse) []*types.RideFareModel
	GenerateTripFares(ctx context.Context, fares []*types.RideFareModel, riderID string, route *types.OSRMApiResponse) ([]*types.RideFareModel, error)
	GetAndValidateRideFare(ctx context.Context, fareID, riderID string) (*types.RideFareModel, error)
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
		Driver:   &trip.TripDriver{},
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

func (s *grpcTripService) EstimatePackagesPriceWithRoute(route *types.OSRMApiResponse) []*types.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*types.RideFareModel, len(baseFares))

	for i, fare := range baseFares {
		estimatedFares[i] = estimateFareRoute(route, fare)
	}
	return estimatedFares
}

func (s *grpcTripService) GenerateTripFares(ctx context.Context, rideFares []*types.RideFareModel, riderID string, route *types.OSRMApiResponse) ([]*types.RideFareModel, error) {
	fares := make([]*types.RideFareModel, len(rideFares))
	for i, fare := range rideFares {
		f := &types.RideFareModel{
			RiderID:           riderID,
			ID:                primitive.NewObjectID(),
			PackageSlug:       fare.PackageSlug,
			TotalFareInPaise: fare.TotalFareInPaise,
			Route:             route,
		}

		if err := s.repo.SaveRideFare(ctx, f); err != nil {
			return nil, fmt.Errorf("failed to save ride fare: %v", err)
		}
		fares[i] = f
	}

	return fares, nil
}

func (s *grpcTripService) GetAndValidateRideFare(ctx context.Context, fareID, riderID string) (*types.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ride fare: %v", err)
	}

	if fare == nil {
		return nil, fmt.Errorf("ride fare not found")
	}

	if fare.RiderID != riderID {
		return nil, fmt.Errorf("ride fare does not belong to the rider")
	}

	return fare, nil
}

func estimateFareRoute(route *types.OSRMApiResponse, fare *types.RideFareModel) *types.RideFareModel {
	pricingCfg := types.DefaultPricingConfig()
	carPackagePrice := fare.TotalFareInPaise

	distanceInKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceInKm * pricingCfg.PricePerUnitDistance
	durationFare := durationInMinutes * pricingCfg.PricePerMinute

	totalFare := carPackagePrice + distanceFare + durationFare

	return &types.RideFareModel{
		PackageSlug:       fare.PackageSlug,
		TotalFareInPaise: totalFare,
	}
}

func getBaseFares() []*types.RideFareModel {
	return []*types.RideFareModel{
		{
			PackageSlug:       "bike",
			TotalFareInPaise: 50.0,
		},
		{
			PackageSlug:       "auto",
			TotalFareInPaise: 70.0,
		},
		{
			PackageSlug:       "sedan",
			TotalFareInPaise: 100.0,
		},
		{
			PackageSlug:       "suv",
			TotalFareInPaise: 150.0,
		},
	}
}

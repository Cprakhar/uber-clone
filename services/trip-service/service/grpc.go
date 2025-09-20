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

type tripService struct {
	repo repo.TripRepo
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *sharedtypes.Coordinate) (*types.OSRMApiResponse, error)
	EstimatePackagesPriceWithRoute(route *types.OSRMApiResponse) []*types.RideFareModel
	GenerateTripFares(ctx context.Context, fares []*types.RideFareModel, riderID string, route *types.OSRMApiResponse) ([]*types.RideFareModel, error)
	GetAndValidateRideFare(ctx context.Context, fareID, riderID string) (*types.RideFareModel, error)
	AcceptRide(ctx context.Context, tripID string, driver *trip.TripDriver) (*types.TripModel, error)
	GetTripByID(ctx context.Context, tripID string) (*types.TripModel, error)
}

// NewService creates a new instance of GrpcTripService
func NewService(repo repo.TripRepo) *tripService {
	return &tripService{repo: repo}
}

func (s *tripService) GetTripByID(ctx context.Context, tripID string) (*types.TripModel, error) {
	return s.repo.GetByID(ctx, tripID)
}

// CreateTrip creates a new trip based on the provided fare
func (s *tripService) CreateTrip(ctx context.Context, fare *types.RideFareModel) (*types.TripModel, error) {
	trip := &types.TripModel{
		ID:       primitive.NewObjectID(),
		RiderID:  fare.RiderID,
		Status:   "pending",
		RideFare: fare,
		Driver:   &trip.TripDriver{},
	}
	return s.repo.Create(ctx, trip)
}

// GetRoute fetches the route from OSRM API between pickup and destination coordinates
func (s *tripService) GetRoute(ctx context.Context, pickup, destination *sharedtypes.Coordinate) (*types.OSRMApiResponse, error) {
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

// EstimatePackagesPriceWithRoute estimates prices for different car packages based on the provided route
func (s *tripService) EstimatePackagesPriceWithRoute(route *types.OSRMApiResponse) []*types.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*types.RideFareModel, len(baseFares))

	for i, fare := range baseFares {
		estimatedFares[i] = estimateFareRoute(route, fare)
	}
	return estimatedFares
}

// GenerateTripFares generates and saves ride fares for a rider based on the provided estimated fares and route
func (s *tripService) GenerateTripFares(ctx context.Context, rideFares []*types.RideFareModel, riderID string, route *types.OSRMApiResponse) ([]*types.RideFareModel, error) {
	fares := make([]*types.RideFareModel, len(rideFares))
	for i, fare := range rideFares {
		f := &types.RideFareModel{
			RiderID:          riderID,
			ID:               primitive.NewObjectID(),
			PackageSlug:      fare.PackageSlug,
			TotalFareInPaise: fare.TotalFareInPaise,
			Route:            route,
		}

		if err := s.repo.SaveRideFare(ctx, f); err != nil {
			return nil, fmt.Errorf("failed to save ride fare: %v", err)
		}
		fares[i] = f
	}

	return fares, nil
}

// GetAndValidateRideFare retrieves a ride fare by ID and validates that it belongs to the specified rider
func (s *tripService) GetAndValidateRideFare(ctx context.Context, fareID, riderID string) (*types.RideFareModel, error) {
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

// AcceptRide allows a driver to accept a trip, updating the trip with the driver's details
func (s *tripService) AcceptRide(ctx context.Context, tripID string, driver *trip.TripDriver) (*types.TripModel, error) {
	return s.repo.UpdateWithDriver(ctx, tripID, driver)
}

// estimateFareRoute estimates the total fare for a given route and base fare
func estimateFareRoute(route *types.OSRMApiResponse, fare *types.RideFareModel) *types.RideFareModel {
	pricingCfg := types.DefaultPricingConfig()
	carPackagePrice := fare.TotalFareInPaise

	distanceInKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceInKm * pricingCfg.PricePerUnitDistance
	durationFare := durationInMinutes * pricingCfg.PricePerMinute

	totalFare := carPackagePrice + distanceFare + durationFare

	return &types.RideFareModel{
		PackageSlug:      fare.PackageSlug,
		TotalFareInPaise: totalFare,
	}
}

// getBaseFares returns a list of base fares for different car packages
func getBaseFares() []*types.RideFareModel {
	return []*types.RideFareModel{
		{
			PackageSlug:      "bike",
			TotalFareInPaise: 50.0,
		},
		{
			PackageSlug:      "auto",
			TotalFareInPaise: 70.0,
		},
		{
			PackageSlug:      "sedan",
			TotalFareInPaise: 100.0,
		},
		{
			PackageSlug:      "suv",
			TotalFareInPaise: 150.0,
		},
	}
}

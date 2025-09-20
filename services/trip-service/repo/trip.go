package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/cprakhar/uber-clone/services/trip-service/types"
	pbd "github.com/cprakhar/uber-clone/shared/proto/trip"
)

var (
	ErrNotFound = fmt.Errorf("resource not found")
)

type inMemoRepo struct {
	sync.RWMutex
	trips     map[string]*types.TripModel
	rideFares map[string]*types.RideFareModel
}

type TripRepo interface {
	Create(ctx context.Context, trip *types.TripModel) (*types.TripModel, error)
	SaveRideFare(ctx context.Context, fare *types.RideFareModel) error
	GetRideFareByID(ctx context.Context, fareID string) (*types.RideFareModel, error)
	UpdateWithDriver(ctx context.Context, tripID string, driver *pbd.TripDriver) (*types.TripModel, error)
	GetByID(ctx context.Context, tripID string) (*types.TripModel, error)
}

// NewInMemoRepository creates a new instance of in-memory TripRepo
func NewInMemoRepository() *inMemoRepo {
	return &inMemoRepo{
		trips:     make(map[string]*types.TripModel),
		rideFares: make(map[string]*types.RideFareModel),
	}
}

func (r *inMemoRepo) GetByID(ctx context.Context, tripID string) (*types.TripModel, error) {
	r.RLock()
	trip, exists := r.trips[tripID]
	r.RUnlock()
	if !exists {
		return nil, ErrNotFound
	}
	return trip, nil
}

// Create adds a new trip to the in-memory store
func (r *inMemoRepo) Create(ctx context.Context, trip *types.TripModel) (*types.TripModel, error) {
	r.Lock()
	r.trips[trip.ID.Hex()] = trip
	r.Unlock()
	return trip, nil
}

// SaveRideFare saves a ride fare to the in-memory store
func (r *inMemoRepo) SaveRideFare(ctx context.Context, fare *types.RideFareModel) error {
	r.Lock()
	r.rideFares[fare.ID.Hex()] = fare
	r.Unlock()
	return nil
}

// GetRideFareByID retrieves a ride fare by its ID
func (r *inMemoRepo) GetRideFareByID(ctx context.Context, fareID string) (*types.RideFareModel, error) {
	r.RLock()
	fare, exists := r.rideFares[fareID]
	r.RUnlock()
	if !exists {
		return nil, ErrNotFound
	}
	return fare, nil
}

// UpdateWithDriver updates a trip with the given driver details and changes its status to "accepted"
func (r *inMemoRepo) UpdateWithDriver(ctx context.Context, tripID string, driver *pbd.TripDriver) (*types.TripModel, error) {
	r.Lock()
	r.trips[tripID].Driver = driver
	r.trips[tripID].Status = "accepted"
	r.Unlock()

	r.RLock()
	defer r.RUnlock()
	return r.trips[tripID], nil
}

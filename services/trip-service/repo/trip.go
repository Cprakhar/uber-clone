package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/cprakhar/uber-clone/services/trip-service/types"
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
}

func NewInMemoRepository() *inMemoRepo {
	return &inMemoRepo{
		trips:     make(map[string]*types.TripModel),
		rideFares: make(map[string]*types.RideFareModel),
	}
}

func (r *inMemoRepo) Create(ctx context.Context, trip *types.TripModel) (*types.TripModel, error) {
	r.Lock()
	r.trips[trip.ID.Hex()] = trip
	r.Unlock()
	return trip, nil
}

func (r *inMemoRepo) SaveRideFare(ctx context.Context, fare *types.RideFareModel) error {
	r.Lock()
	r.rideFares[fare.ID.Hex()] = fare
	r.Unlock()
	return nil
}

func (r *inMemoRepo) GetRideFareByID(ctx context.Context, fareID string) (*types.RideFareModel, error) {
	r.RLock()
	fare, exists := r.rideFares[fareID]
	r.RUnlock()
	if !exists {
		return nil, ErrNotFound
	}
	return fare, nil
}

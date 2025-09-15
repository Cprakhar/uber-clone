package repo

import (
	"context"

	"github.com/cprakhar/uber-clone/services/trip-service/types"
)

type inMemoRepo struct {
	trips     map[string]*types.TripModel
	rideFares map[string]*types.RideFareModel
}

type TripRepo interface {
	Create(ctx context.Context, trip *types.TripModel) (*types.TripModel, error)
}

func NewInMemoRepository() *inMemoRepo {
	return &inMemoRepo{
		trips:     make(map[string]*types.TripModel),
		rideFares: make(map[string]*types.RideFareModel),
	}
}

func (r *inMemoRepo) Create(ctx context.Context, trip *types.TripModel) (*types.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

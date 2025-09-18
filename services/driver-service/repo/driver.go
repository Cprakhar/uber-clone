package repo

import (
	"sync"

	"github.com/cprakhar/uber-clone/services/driver-service/types"
)

type inMemoRepo struct {
	sync.RWMutex
	drivers map[string]*types.DriverModel
}

type DriverRepo interface {
	Create(driver *types.DriverModel) (*types.DriverModel, error)
	Delete(driverID string) error
}

func NewDriverRepository() *inMemoRepo {
	return &inMemoRepo{
		drivers: make(map[string]*types.DriverModel),
	}
}

func (r *inMemoRepo) Create(driver *types.DriverModel) (*types.DriverModel, error) {
	r.Lock()
	r.drivers[driver.ID] = driver
	r.Unlock()
	return driver, nil
}

func (r *inMemoRepo) Delete(driverID string) error {
	r.Lock()
	delete(r.drivers, driverID)
	r.Unlock()
	return nil
}

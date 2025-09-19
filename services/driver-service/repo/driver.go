package repo

import (
	"sync"

	pb "github.com/cprakhar/uber-clone/shared/proto/driver"
)

type inMemoRepo struct {
	sync.RWMutex
	drivers []*pb.Driver
}

type DriverRepo interface {
	Create(driver *pb.Driver) (*pb.Driver, error)
	Delete(driverID string) error
	GetAll() []*pb.Driver
}

func NewDriverRepository() *inMemoRepo {
	return &inMemoRepo{
		drivers: []*pb.Driver{},
	}
}

func (r *inMemoRepo) Create(driver *pb.Driver) (*pb.Driver, error) {
	r.Lock()
	r.drivers = append(r.drivers, driver)
	r.Unlock()
	return driver, nil
}

func (r *inMemoRepo) Delete(driverID string) error {
	r.Lock()
	for i, d := range r.drivers {
		if d.Id == driverID {
			r.drivers = append(r.drivers[:i], r.drivers[i+1:]...)
			break
		}
	}
	r.Unlock()
	return nil
}

func (r *inMemoRepo) GetAll() []*pb.Driver {
	return r.drivers
}

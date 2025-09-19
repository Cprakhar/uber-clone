package service

import (
	"context"
	"log"
	"math/rand/v2"

	"github.com/cprakhar/uber-clone/services/driver-service/repo"
	"github.com/cprakhar/uber-clone/services/driver-service/util"
	pb "github.com/cprakhar/uber-clone/shared/proto/driver"
	sharedUtil "github.com/cprakhar/uber-clone/shared/util"
	"github.com/mmcloughlin/geohash"
)

type driverService struct {
	repo repo.DriverRepo
}

type DriverService interface {
	RegisterDriver(ctx context.Context, driverID, packageSlug string) (*pb.Driver, error)
	UnregisterDriver(ctx context.Context, driverID string) error
	FindAvailableDrivers(ctx context.Context, packageSlug string) []string
}

func NewDriverService(repo repo.DriverRepo) *driverService {
	return &driverService{repo: repo}
}

func (s *driverService) RegisterDriver(ctx context.Context, driverID, packageSlug string) (*pb.Driver, error) {
	randomIdx := rand.IntN(len(util.PredefinedRoutes))
	profilePic := sharedUtil.GetRandomProfilePic(randomIdx)
	randomRoute := util.PredefinedRoutes[randomIdx]
	carPlate := util.GenerateRandomPlate()

	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Id:          driverID,
		Name:        "Prakhar Chhalotre",
		ProfilePic:  profilePic,
		CarPlate:    carPlate,
		PackageSlug: packageSlug,
		Geohash:     geohash,
		Location: &pb.Location{
			Latitude:  randomRoute[0][0],
			Longitude: randomRoute[0][1],
		},
	}

	driver, err := s.repo.Create(driver)
	if err != nil {
		return nil, err
	}
	log.Printf("All drivers: %v", s.repo.GetAll())
	return driver, nil
}

func (s *driverService) UnregisterDriver(ctx context.Context, driverID string) error {
	return s.repo.Delete(driverID)
}

func (s *driverService) FindAvailableDrivers(ctx context.Context, packageSlug string) []string {
	matchingDrivers := []string{}
	for _, d := range s.repo.GetAll() {
		if d.PackageSlug == packageSlug {
			matchingDrivers = append(matchingDrivers, d.Id)
		}
	}

	return matchingDrivers
}

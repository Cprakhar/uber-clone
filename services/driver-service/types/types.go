package types

import (
	sharedtypes "github.com/cprakhar/uber-clone/shared/types"
)

type DriverModel struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	ProfilePic  string                 `json:"profilePic"`
	CarPlate    string                 `json:"carPlate"`
	PackageSlug string                 `json:"packageSlug"`
	Geohash     string                 `json:"geohash"`
	Location    sharedtypes.Coordinate `json:"location"`
}

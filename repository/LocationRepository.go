package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type LocationRepository interface {
	GetLocations() []model.Location
	GetLocationByIDs(locationIDs []string) []model.Location
}

type LocationRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewLocationRepository(db *connection.DBConnection) LocationRepository {
	return &LocationRepositoryImpl{
		dbConn: db,
	}
}

func (repository *LocationRepositoryImpl) GetLocations() []model.Location {
	locations := []model.Location{}
	err := repository.dbConn.DB.Model(&locations).Find(&locations).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return locations
}

func (repository *LocationRepositoryImpl) GetLocationByIDs(locationIDs []string) []model.Location {
	locations := []model.Location{}

	err := repository.dbConn.DB.Model(&locations).Where("id IN ?", locationIDs).Find(&locations).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return locations
}

package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type LocationRepository interface {
	GetLocations(*[]model.Location)
}

type LocationRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewLocationRepository(db *connection.DBConnection) LocationRepository {
	return &LocationRepositoryImpl{
		dbConn: db,
	}
}

func (repository *LocationRepositoryImpl) GetLocations(locations *[]model.Location) {

	err := repository.dbConn.DB.Model(&locations).Find(&locations).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}
}

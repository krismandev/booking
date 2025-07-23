package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type RoomRepository interface {
	GetRooms(filter model.ListRoomQueryFilter) []model.Room
	GetRoomByIDs(roomIDs []string) []model.Room
}

type RoomRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewRoomRepository(db *connection.DBConnection) RoomRepository {
	return &RoomRepositoryImpl{
		dbConn: db,
	}
}

func (repository *RoomRepositoryImpl) GetRooms(filter model.ListRoomQueryFilter) []model.Room {
	Rooms := []model.Room{}
	qry := repository.dbConn.DB.Model(&Rooms)
	if len(filter.LocationID) > 0 {
		qry = qry.Where("locationId = ?", filter.LocationID)
	}

	err := qry.Find(&Rooms)

	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return Rooms
}

func (repository *RoomRepositoryImpl) GetRoomByIDs(roomIDs []string) []model.Room {
	rooms := []model.Room{}

	err := repository.dbConn.DB.Model(&rooms).Where("id IN ?", roomIDs).Find(&rooms).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return rooms
}

package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type RoomRepository interface {
	GetRooms() []model.Room
}

type RoomRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewRoomRepository(db *connection.DBConnection) RoomRepository {
	return &RoomRepositoryImpl{
		dbConn: db,
	}
}

func (repository *RoomRepositoryImpl) GetRooms() []model.Room {
	Rooms := []model.Room{}
	err := repository.dbConn.DB.Model(&Rooms).Find(&Rooms).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return Rooms
}

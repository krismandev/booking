package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type RoomRepository interface {
	GetRooms() []model.Room
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

func (repository *RoomRepositoryImpl) GetRooms() []model.Room {
	Rooms := []model.Room{}
	err := repository.dbConn.DB.Model(&Rooms).Find(&Rooms).Error
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

package service

import (
	"booking/model"
	"booking/model/response"
	"booking/repository"
	"context"
)

type RoomService interface {
	GetRooms(ctx context.Context) []response.RoomResponse
}

type RoomServiceImpl struct {
	repository         repository.RoomRepository
	locationRepository repository.LocationRepository
}

func NewRoomService(repository repository.RoomRepository, locationRepository repository.LocationRepository) RoomService {
	return &RoomServiceImpl{
		repository:         repository,
		locationRepository: locationRepository,
	}
}

func (service *RoomServiceImpl) GetRooms(ctx context.Context) []response.RoomResponse {
	var output []response.RoomResponse

	rooms := service.repository.GetRooms()

	locations := service.locationRepository.GetLocations()

	if len(rooms) > 0 {
		for _, each := range rooms {
			var location model.Location
			for _, l := range locations {
				if each.LocationID == l.ID {
					location = l
					break
				}
			}
			output = append(output, response.ToRoomResponse(each, location))
		}
	}

	return output
}

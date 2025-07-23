package service

import (
	"booking/model"
	"booking/model/request"
	"booking/model/response"
	"booking/repository"
	"booking/utils"
	"context"
	"net/url"

	"github.com/sirupsen/logrus"
)

type RoomService interface {
	GetRooms(ctx context.Context, request request.BookingListRequest) ([]response.RoomResponse, error)
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

func (service *RoomServiceImpl) GetRooms(ctx context.Context, request request.BookingListRequest) ([]response.RoomResponse, error) {
	var output []response.RoomResponse

	var filter model.ListRoomQueryFilter

	decodedFilter, err := url.QueryUnescape(request.Filter)
	if err != nil {
		return output, err
	}

	if len(decodedFilter) > 0 {
		err = utils.Decode(decodedFilter, &filter)
		if err != nil {
			logrus.Errorf("Error parsing filter: %v", err)
			return output, &utils.BadRequestError{Message: "Invalid filter format. must be encoded string of json"}
		}
	}

	rooms := service.repository.GetRooms(filter)

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

	return output, err
}

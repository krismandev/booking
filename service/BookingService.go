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

type BookingService interface {
	GetBookings(ctx context.Context, request request.BookingListRequest) (response.BookingListResponse, error)
	CreateBooking(ctx context.Context, request request.CreateBookingRequest) (response.BookingResponse, error)
}

type BookingServiceImpl struct {
	repository         repository.BookingRepository
	locationRepository repository.LocationRepository
	roomRepository     repository.RoomRepository
}

func NewBookingService(repository repository.BookingRepository, locationRepository repository.LocationRepository, roomRepository repository.RoomRepository) BookingService {
	return &BookingServiceImpl{
		repository:         repository,
		locationRepository: locationRepository,
		roomRepository:     roomRepository,
	}
}

func (service *BookingServiceImpl) GetBookings(ctx context.Context, request request.BookingListRequest) (response.BookingListResponse, error) {
	var resp response.BookingListResponse

	var filter model.BookingListQueryFilter

	decodedFilter, err := url.QueryUnescape(request.Filter)
	if err != nil {
		return resp, err
	}

	if len(decodedFilter) > 0 {
		err = utils.Decode(decodedFilter, &filter)
		if err != nil {
			logrus.Errorf("Error parsing filter: %v", err)
			return resp, &utils.BadRequestError{Message: "Invalid filter format. must be encoded string of json"}
		}
	}

	bookings := service.repository.GetBookings(filter)

	var roomIDs []string
	for _, each := range bookings {
		roomIDs = append(roomIDs, each.RoomID)
	}

	rooms := service.roomRepository.GetRoomByIDs(roomIDs)

	var locationIDs []string
	for _, each := range rooms {
		locationIDs = append(locationIDs, each.LocationID)
	}

	locations := service.locationRepository.GetLocationByIDs(locationIDs)

	if len(bookings) > 0 {
		for _, each := range bookings {
			var room model.Room
			var location model.Location
			for _, r := range rooms {
				if r.ID == each.RoomID {
					room = r
					break
				}
			}

			for _, l := range locations {
				if l.ID == room.LocationID {
					location = l
					break
				}
			}
			single := response.ToBookingResponse(each, &room, &location)
			resp.Data = append(resp.Data, single)
		}
	}

	return resp, err
}

func (service *BookingServiceImpl) CreateBooking(ctx context.Context, request request.CreateBookingRequest) (response.BookingResponse, error) {
	var resp response.BookingResponse
	var err error

	var roomIDs []string = []string{request.RoomID}
	rooms := service.roomRepository.GetRoomByIDs(roomIDs)
	if len(rooms) == 0 {
		logrus.Errorf("Error in service. Room is not valid : %v", err)
		return resp, &utils.UnprocessableContentError{Message: "Invalid Room"}
	}

	var data model.Booking
	data.Category = request.Category
	data.UserID = request.UserID
	data.Description = request.Description
	data.StartDate = request.StartDate
	data.EndDate = request.EndDate
	data.RoomID = request.RoomID
	data.Title = request.Title

	err = service.repository.CreateBooking(&data)
	if err != nil {
		logrus.Errorf("Failed to create booking : %v", err)
		return resp, err
	}

	resp = response.ToBookingResponse(data, nil, nil)

	return resp, err
}

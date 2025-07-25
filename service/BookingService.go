package service

import (
	"booking/model"
	"booking/model/request"
	"booking/model/response"
	"booking/repository"
	"booking/utils"
	"context"
	"encoding/json"
	"net/url"

	"github.com/sirupsen/logrus"
)

type BookingService interface {
	GetBookings(ctx context.Context, request request.BookingListRequest) (response.BookingListResponse, error)
	CreateBooking(ctx context.Context, request request.CreateBookingRequest) (response.BookingResponse, error)
	CancelBooking(ctx context.Context, request request.CancelBookingRequest) error
	ApproveBooking(ctx context.Context, request request.ApproveBookingRequest) error
}

type BookingServiceImpl struct {
	repository         repository.BookingRepository
	locationRepository repository.LocationRepository
	userRepository     repository.UserRepository
	roomRepository     repository.RoomRepository
}

func NewBookingService(repository repository.BookingRepository, locationRepository repository.LocationRepository, roomRepository repository.RoomRepository, userRepository repository.UserRepository) BookingService {
	return &BookingServiceImpl{
		repository:         repository,
		locationRepository: locationRepository,
		roomRepository:     roomRepository,
		userRepository:     userRepository,
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

	dttbFilter := request.GlobalListDataRequest.ParseToJson()
	if len(dttbFilter) > 0 {
		json.Unmarshal([]byte(dttbFilter), &filter.GlobalQueryFilter)
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

	var userIDs []string
	for _, each := range bookings {
		userIDs = append(userIDs, each.UserID)
	}
	users, err := service.userRepository.FindUserByIDs(userIDs)
	if err != nil {
		logrus.Errorf("Error in service. Failed to retrieve user data : %v", err)
		return resp, err
	}

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

			var user model.User
			for _, usr := range users {
				if each.UserID == usr.ID {
					user = usr
				}
			}
			single := response.ToBookingResponse(each, &room, &location, &user)
			resp.Data = append(resp.Data, single)
		}
	}

	count, _ := service.repository.CountBooking(filter)

	perPage, currentPage, totalPage := request.CollectMetadata(int(count))
	resp.Count = int(count)
	resp.Limit = perPage
	resp.TotalPage = totalPage
	resp.Page = currentPage

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

	alreadyBooked := service.repository.CheckRoomAlreadyBooked(request.RoomID, request.StartDate)
	if alreadyBooked {
		logrus.Errorf("Error in service. Room Already booked : %v", err)
		return resp, &utils.UnprocessableContentError{Message: "Room already booked at given time"}
	}

	user, err := service.userRepository.FindUserById(request.UserID)
	if !user.IsActive {
		logrus.Errorf("User Inactive! You are not allowed to booking a room")
		return resp, &utils.ForbiddenError{Message: "User Inactive! You are not allowed to booking a room"}
	}

	var data model.Booking
	data.Category = request.Category
	data.UserID = request.UserID
	data.Description = &request.Description
	data.StartDate = request.StartDate
	data.EndDate = request.EndDate
	data.RoomID = request.RoomID
	data.Title = request.Title
	data.Status = model.REQUESTED

	err = service.repository.CreateBooking(&data)
	if err != nil {
		logrus.Errorf("Failed to create booking : %v", err)
		return resp, err
	}

	resp = response.ToBookingResponse(data, nil, nil, nil)

	return resp, err
}

func (service *BookingServiceImpl) CancelBooking(ctx context.Context, request request.CancelBookingRequest) error {
	var err error

	booking := service.repository.FindBookingByID(request.BookingID)

	if len(booking.ID) == 0 {
		logrus.Errorf("Error in service. Booking not found : %v", err)
		return &utils.NotFoundError{Message: "Data not found"}
	}

	booking.Status = model.CANCELED

	err = service.repository.CancelBooking(booking)
	if err != nil {
		logrus.Errorf("Error in service. Failed to cancel booking : %v", err)
		return &utils.InternalServerError{}
	}

	return err
}

func (service *BookingServiceImpl) ApproveBooking(ctx context.Context, request request.ApproveBookingRequest) error {
	var err error

	booking := service.repository.FindBookingByID(request.BookingID)

	if len(booking.ID) == 0 {
		logrus.Errorf("Error in service. Booking not found : %v", err)
		return &utils.NotFoundError{Message: "Data not found"}
	}

	booking.Status = request.Status

	err = service.repository.ApproveBooking(booking)
	if err != nil {
		logrus.Errorf("Error in service. Failed to cancel booking : %v", err)
		return &utils.InternalServerError{}
	}

	return err
}

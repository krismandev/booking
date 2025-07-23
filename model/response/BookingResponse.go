package response

import (
	"booking/model"
)

type BookingResponse struct {
	ID          string        `json:"id"`
	RoomID      string        `json:"roomId"`
	Title       string        `json:"title"`
	Category    string        `json:"category"`
	Description string        `json:"description"`
	StartDate   string        `json:"startDate"`
	EndDate     string        `json:"endDate"`
	UserID      string        `json:"userId"`
	Status      string        `json:"status"`
	Room        *RoomResponse `json:"room,omitempty"`
	User        *UserResponse `json:"user,omitempty"`
}

type BookingListResponse struct {
	MetadataResponse
	Data []BookingResponse
}

func ToBookingResponse(dt model.Booking, dtRoom *model.Room, dtLocation *model.Location, dtUser *model.User) BookingResponse {
	var resp BookingResponse

	resp.ID = dt.ID
	resp.RoomID = dt.RoomID
	resp.Category = dt.Category
	resp.Description = *dt.Description
	resp.StartDate = dt.StartDate
	resp.EndDate = dt.EndDate
	resp.Title = dt.Title
	resp.UserID = dt.UserID
	resp.Status = dt.Status

	if dtRoom != nil {
		var room RoomResponse
		room = ToRoomResponse(*dtRoom, *dtLocation)

		resp.Room = &room
	}

	if dtUser != nil {
		var userResp UserResponse
		userResp = ToUserResponse(*dtUser)

		resp.User = &userResp
	}
	return resp
}

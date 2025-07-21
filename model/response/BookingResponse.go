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
	Room        *RoomResponse `json:"room,omitempty"`
}

type BookingListResponse struct {
	MetadataResponse
	Data []BookingResponse
}

func ToBookingResponse(dt model.Booking, dtRoom *model.Room, dtLocation *model.Location) BookingResponse {
	var resp BookingResponse

	resp.ID = dt.ID
	resp.RoomID = dt.RoomID
	resp.Category = dt.Category
	resp.Description = *dt.Description
	resp.StartDate = dt.StartDate
	resp.EndDate = dt.EndDate
	resp.Title = dt.Title
	resp.UserID = dt.UserID

	if dtRoom != nil {
		var room RoomResponse
		room = ToRoomResponse(*dtRoom, *dtLocation)

		resp.Room = &room
	}
	return resp
}

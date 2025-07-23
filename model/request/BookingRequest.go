package request

type BookingListRequest struct {
	GlobalListDataRequest
	Filter string
}

type CreateBookingRequest struct {
	RoomID      string `json:"roomId" validate:"required"`
	UserID      string `json:"userId" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	StartDate   string `json:"startDate" validate:"required"`
	EndDate     string `json:"endDate" validate:"required"`
	Category    string `json:"category"`
}

type CancelBookingRequest struct {
	BookingID string `json:"bookingId" validate:"required"`
	UserID    string
}

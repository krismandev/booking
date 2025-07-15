package response

type BookingResponse struct {
	ID          string `json:"id"`
	RoomID      string `json:"roomId"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	UserID      string `json:"userId"`
}

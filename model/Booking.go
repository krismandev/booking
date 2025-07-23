package model

const (
	REQUESTED = "REQUESTED"
	APPROVED  = "APPROVED"
	CANCELED  = "CANCELED"
	REJECTED  = "REJECTED"
)

type Booking struct {
	ID          string  `gorm:"column:id;default:uuid_generate_v4();primaryKey"`
	RoomID      string  `gorm:"column:roomid"`
	Title       string  `gorm:"column:title"`
	Category    string  `gorm:"column:category"`
	Description *string `gorm:"column:description"`
	StartDate   string  `gorm:"column:startdate"`
	EndDate     string  `gorm:"column:enddate"`
	UserID      string  `gorm:"column:userid"`
	Status      string  `gorm:"column:status"`
}

type BookingListQueryFilter struct {
	Booking
	RoomIDs    []string `json:"roomIds"`
	LocationID []string `json:"locationIds"`
}

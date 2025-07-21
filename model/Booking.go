package model

type Booking struct {
	ID          string `gorm:"column:id;default:uuid_generate_v4();primaryKey"`
	RoomID      string `json:"roomId" gorm:"column:roomid"`
	Title       string `json:"title" gorm:"column:title"`
	Category    string `json:"category" gorm:"column:category"`
	Description string `json:"description" gorm:"column:description"`
	StartDate   string `json:"startDate" gorm:"column:startdate"`
	EndDate     string `json:"endDate" gorm:"column:enddate"`
	UserID      string `json:"userId" gorm:"column:userid"`
}

type BookingListQueryFilter struct {
	Booking
	RoomID     []string `json:"roomIds"`
	LocationID []string `json:"locationIds"`
}

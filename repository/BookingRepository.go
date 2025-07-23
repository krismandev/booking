package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type BookingRepository interface {
	GetBookings(model.BookingListQueryFilter) []model.Booking
	CreateBooking(*model.Booking) error
	CheckRoomAlreadyBooked(roomId string, startDateRequested string) bool
	CancelBooking(model.Booking) error
	ApproveBooking(model.Booking) error
	FindBookingByID(bookingID, userID string) model.Booking
}

type BookingRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewBookingRepository(db *connection.DBConnection) BookingRepository {
	return &BookingRepositoryImpl{
		dbConn: db,
	}
}

func (repository *BookingRepositoryImpl) GetBookings(filter model.BookingListQueryFilter) []model.Booking {
	bookings := []model.Booking{}

	qry := repository.dbConn.DB.Model(&bookings)

	if len(filter.StartDate) > 0 {
		qry = qry.Where("startDate >= ?", filter.StartDate)
	}
	if len(filter.Title) > 0 {
		qry = qry.Where("title like ?", "%"+filter.Title+"%")
	}
	if len(filter.EndDate) > 0 {
		qry = qry.Where("endDate <= ?", filter.EndDate)
	}
	if len(filter.Category) > 0 {
		qry = qry.Where("category = ?", filter.Category)
	}

	if len(filter.LocationID) > 0 {
		qry = qry.Where("locationid IN ?", filter.LocationID)
	}

	if len(filter.RoomID) > 0 {
		qry = qry.Where("roomid = ?", filter.RoomID)
	}

	if len(filter.Status) > 0 {
		qry = qry.Where("status = ?", filter.Status)
	}

	err := qry.Find(&bookings).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return bookings
}

func (repository *BookingRepositoryImpl) CreateBooking(dt *model.Booking) error {
	// var booking model.Booking

	err := repository.dbConn.DB.Create(&dt).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return err
}

func (repository *BookingRepositoryImpl) CheckRoomAlreadyBooked(roomId string, startDateRequested string) bool {
	var available bool

	qry := "SELECT EXISTS (SELECT 1 from bookings WHERE ? >= startDate AND ? < endDate) "
	repository.dbConn.Raw(qry, startDateRequested, startDateRequested).Scan(&available)
	return available
}
func (repository *BookingRepositoryImpl) CancelBooking(dt model.Booking) error {
	var err error

	err = repository.dbConn.DB.Model(&dt).Select("status").Updates(model.Booking{Status: model.CANCELED}).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return err
}

func (repository *BookingRepositoryImpl) FindBookingByID(bookingID, userID string) model.Booking {
	var dt model.Booking

	err := repository.dbConn.DB.Where("id = ?", bookingID).Where("userid = ?", userID).First(&dt).Error

	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return dt
}

func (repo *BookingRepositoryImpl) ApproveBooking(dt model.Booking) error {
	var err error
	err = repo.dbConn.DB.Model(&dt).Select("status").Updates(model.Booking{Status: dt.Status}).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return err
}

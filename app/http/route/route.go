package route

import (
	"booking/app/http/controller"
	"booking/app/http/middleware"
	"booking/utils"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	Echo               *echo.Echo
	LocationController controller.LocationController
	RoomController     controller.RoomController
	BookingController  controller.BookingController
	AuthController     controller.AuthController
	UserController     controller.UserController
}

func (r *RouteConfig) InitRoute() {
	r.Echo.HTTPErrorHandler = utils.CustomHTTPErrorHandler
	r.InitPublicRoute()
}

func (r *RouteConfig) InitPrivateRoute() {
	route := r.Echo.Group("/private", middleware.JWTAuth())
	route.POST("/bookings", r.BookingController.CreateBooking)
}

func (r *RouteConfig) InitPublicRoute() {
	route := r.Echo.Group("/api")
	route.GET("/locations", r.LocationController.GetLocations)
	route.GET("/bookings", r.BookingController.GetBookings)
	route.POST("/booking", r.BookingController.CreateBooking)
	route.GET("/rooms", r.RoomController.GetRooms)

	route.POST("/register", r.UserController.CreateUser)
}

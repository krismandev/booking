package route

import (
	"booking/app/http/controller"
	"booking/app/http/middleware"
	connection "booking/connection/database"
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
	DB                 connection.DBConnection
}

func (r *RouteConfig) InitRoute() {
	r.Echo.HTTPErrorHandler = utils.CustomHTTPErrorHandler
	r.InitPublicRoute()
	r.InitPrivateRoute()
}

func (r *RouteConfig) InitPrivateRoute() {

	authorizationMiddleware := middleware.NewAuthorizationMiddleware(r.DB)

	r.Echo.GET("/auth/detail", r.AuthController.AuthUserDetail, middleware.JWTAuth(), authorizationMiddleware.Authorize("users.read"))

	route := r.Echo.Group("/api", middleware.JWTAuth(), authorizationMiddleware.Authorize("user.create"))
	route.POST("/booking", r.BookingController.CreateBooking)
	// route.GET("/auth/detail", r.AuthController.AuthUserDetail)

}

func (r *RouteConfig) InitPublicRoute() {
	r.Echo.POST("/login", r.AuthController.Login)

	route := r.Echo.Group("/api")
	route.GET("/locations", r.LocationController.GetLocations)
	route.GET("/bookings", r.BookingController.GetBookings)
	route.GET("/rooms", r.RoomController.GetRooms)

	route.POST("/register", r.UserController.CreateUser)
}

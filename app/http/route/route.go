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
	RoleController     controller.RoleController
	DB                 connection.DBConnection
}

func (r *RouteConfig) InitRoute() {
	r.Echo.HTTPErrorHandler = utils.CustomHTTPErrorHandler
	r.InitPublicRoute()
	r.InitPrivateRoute()
}

func (r *RouteConfig) InitPrivateRoute() {

	authorizationMiddleware := middleware.NewAuthorizationMiddleware(r.DB)

	r.Echo.GET("/auth/detail", r.AuthController.AuthUserDetail, middleware.JWTAuth(), authorizationMiddleware.Authorize("users.detail"))

	route := r.Echo.Group("/api", middleware.JWTAuth())

	route.POST("/bookings", r.BookingController.CreateBooking, authorizationMiddleware.Authorize("bookings.create"))
	route.DELETE("/bookings/cancel", r.BookingController.CancelBooking, authorizationMiddleware.Authorize("bookings.cancel"))
	route.POST("/bookings/approval", r.BookingController.ApproveBooking, authorizationMiddleware.Authorize("bookings.approval"))

	route.POST("/users", r.UserController.CreateUser, authorizationMiddleware.Authorize("users.create"))
	route.GET("/users", r.UserController.ListUser, authorizationMiddleware.Authorize("users.read"))
	route.PATCH("/users", r.UserController.UpdateUser, authorizationMiddleware.Authorize("users.update"))
	route.DELETE("/users", r.UserController.DeactivateUser, authorizationMiddleware.Authorize("users.delete"))

}

func (r *RouteConfig) InitPublicRoute() {
	r.Echo.POST("/login", r.AuthController.Login)

	route := r.Echo.Group("/api")
	route.GET("/locations", r.LocationController.GetLocations)
	route.GET("/rooms", r.RoomController.GetRooms)

	route.POST("/register", r.UserController.CreateUser)

	route.GET("/roles", r.RoleController.GetRoles)

	route.GET("/bookings", r.BookingController.GetBookings)
	// route.GET("/roles",)
}

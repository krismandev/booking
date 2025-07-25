package app

import (
	"booking/app/http/controller"
	"booking/app/http/route"
	connection "booking/connection/database"
	"booking/repository"
	"booking/service"

	"github.com/labstack/echo/v4"
)

// "aggregator-service/connection"
// cache "aggregator-service/connection/cache"
// rabbitmq "aggregator-service/connection/queue/rabbitmq"

type Application struct {
	DB *connection.DBConnection
	// MailQueueProducer *rabbitmq.Consumer
	Echo *echo.Echo
	// Validate          *validator.Validate
	// Redis             *cache.RedisConnection
	// RoleMap           map[string]AccessPermission
}

type AccessPermission struct {
	Resource    string
	Scopes      []string
	Description string
}

func InitApp(app *Application) {

	locationRepository := repository.NewLocationRepository(app.DB)
	locationService := service.NewLocationService(locationRepository)
	locationController := controller.NewLocationController(locationService)

	roomRepository := repository.NewRoomRepository(app.DB)
	roomService := service.NewRoomService(roomRepository, locationRepository)
	roomController := controller.NewRoomController(roomService)

	bookingRepository := repository.NewBookingRepository(app.DB)

	userRepository := repository.NewUserRepository(app.DB)
	bookingService := service.NewBookingService(bookingRepository, locationRepository, roomRepository, userRepository)
	bookingController := controller.NewBookingController(bookingService)

	roleRepository := repository.NewRoleRepository(app.DB)

	authService := service.NewAuthService(userRepository, app.DB, roleRepository)
	authController := controller.NewAuthController(authService)

	roleService := service.NewRoleService(roleRepository)
	roleController := controller.NewRoleController(roleService)

	departmentRepository := repository.NewDepartmentRepository(app.DB)
	departmentService := service.NewDepartmentService(departmentRepository)
	departmentController := controller.NewDepartmentController(departmentService)

	userService := service.NewUserService(userRepository, app.DB, roleRepository, departmentRepository)
	userController := controller.NewUserController(userService)

	routeConfig := route.RouteConfig{
		Echo:                 app.Echo,
		LocationController:   locationController,
		RoomController:       roomController,
		BookingController:    bookingController,
		UserController:       userController,
		AuthController:       authController,
		RoleController:       roleController,
		DepartmentController: departmentController,
		DB:                   *app.DB,
	}

	routeConfig.InitRoute()
}

func (app *Application) LoadRoles() {

}

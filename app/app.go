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

	// cfg, err := envcfg.New(envcfg.ConfigFile("./config"))
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	exchange = "mpg"
	// 	// exchangeType         = string(rabbitmq.ExchangeType_DIRECT)
	// 	mailRequestQueueName = "mpg.mail.request"
	// )

	// validator := validator.New()

	// mailProducer, err := rabbitmq.NewProducer(cfg.GetString("amqp.url"), mailRequestQueueName, rabbitmq.WithExchange(exchange))
	// if err != nil {
	// 	logrus.Errorf("failed start worker %v", err)
	// 	panic(err)
	// }

	// commonhttpConfig := api.Config{}
	// httpClient, err := api.New(commonhttpConfig)
	// if err != nil {
	// 	panic(err)
	// }

	// // proxyUrl := cfg.GetString("proxyUrl")
	// // httpConfigExternalClient := api.Config{UseProxy: true, ProxyURL: proxyUrl, LogReqResBodyEnable: true}
	// // externalHttpClient, err := api.New(httpConfigExternalClient)
	// // if err != nil {
	// // 	panic(err)
	// // }

	// cacheRepository := repository.NewCacheRepository(app.Redis)

	// storagePath := cfg.GetString("documentFilepath")
	// fileStorageRepository := repository.NewFileStorageRepository(storagePath)
	// if fileStorageRepository == nil {
	// 	logrus.Fatal("Failed to initialize fileStorageRepository")
	// }

	locationRepository := repository.NewLocationRepository(app.DB)
	locationService := service.NewLocationService(locationRepository)
	locationController := controller.NewLocationController(locationService)

	roomRepository := repository.NewRoomRepository(app.DB)
	roomService := service.NewRoomService(roomRepository, locationRepository)
	roomController := controller.NewRoomController(roomService)

	routeConfig := route.RouteConfig{
		Echo:               app.Echo,
		LocationController: locationController,
		RoomController:     roomController,
	}

	routeConfig.InitRoute()
}

func (app *Application) LoadRoles() {

}

package main

import (
	"fmt"
	"os"

	"booking/app"
	connection "booking/connection/database"
	"booking/utils"

	"github.com/kenshaw/envcfg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, latency=${latency_human}, bytes_in=${bytes_in}, bytes_out=${bytes_out}` + "\n",
	}))

	e.Validator = utils.NewCustomValidator()

	cfg, err := envcfg.New(envcfg.ConfigFile("./config"))
	if err != nil {
		panic(errors.Wrap(err, "loading configuration"))
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			func() string {
				if cfg.GetString("cors.origin") != "" && cfg.GetString("cors.origin") != "*" {
					return cfg.GetString("cors.origin")
				}
				return "*"
			}(),
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	port := cfg.GetString("port")
	logLevel := cfg.GetString("logLevel")

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&utils.CustomJSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if logLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	dbUrl := cfg.GetString("dbUrl")
	if dbUrl == "" {
		panic(fmt.Sprintf("%s dbUrl is not set", dbUrl))
	}

	dbHost := cfg.GetString("database.host")
	dbUser := cfg.GetString("database.user")
	dbPassword := cfg.GetString("database.password")
	dbPort := cfg.GetString("database.port")
	dbName := cfg.GetString("database.database")

	dbConfig := make(map[string]string)
	dbConfig["host"] = dbHost
	dbConfig["port"] = dbPort
	dbConfig["user"] = dbUser
	dbConfig["dbname"] = dbName
	dbConfig["password"] = dbPassword
	dbConfig["sslmode"] = "disable"

	db, err := connection.NewConnection(dbConfig)
	if err != nil {
		logrus.Error("error initiate db connection")
		panic(err)
	}

	// redisHost := cfg.GetString("cache.redisHost")
	// redisPort := cfg.GetString("cache.redisPort")
	// redisPassword := cfg.GetString("cache.redisPassword")
	// redisDB := cfg.GetInt("cache.redisDB")
	// redisUrl := redisHost + ":" + redisPort

	// redisClient, err := cache.NewRedisClient(redisUrl, redisPassword, redisDB)
	// if err != nil {
	// 	logrus.Errorf("error initiate redis connection")
	// 	panic(err)
	// }

	app.InitApp(&app.Application{DB: db, Echo: e})
	logrus.Info("onboarding service started on port " + port)
	err = e.Start(":" + port)
	if err != nil {
		panic(err)
	}
}

package controller

import (
	"booking/model/response"
	"booking/service"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type LocationController interface {
	GetLocations(c echo.Context) error
}

type LocationControllerImpl struct {
	LocationService service.LocationService
}

func NewLocationController(locationService service.LocationService) LocationController {
	return &LocationControllerImpl{
		LocationService: locationService,
	}
}

func (controller *LocationControllerImpl) GetLocations(c echo.Context) error {
	var err error
	var dataResp response.GlobalListDataResponse
	ctx := c.Request().Context()
	var responseData interface{}

	// request := request.LocationRequest{}

	resp := controller.LocationService.GetLocations(ctx)
	logrus.Info(resp)
	if err != nil {
		response.WriteResponseSingleJSON(c, responseData, err)
	}

	for _, each := range resp {
		dataResp.List = append(dataResp.List, each)
	}
	response.WriteResponseListJSON(c, dataResp, nil)
	return err
}

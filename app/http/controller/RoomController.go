package controller

import (
	"booking/model/request"
	"booking/model/response"
	"booking/service"

	"github.com/labstack/echo/v4"
)

type RoomController interface {
	GetRooms(c echo.Context) error
}

type RoomControllerImpl struct {
	RoomService service.RoomService
}

func NewRoomController(locationService service.RoomService) RoomController {
	return &RoomControllerImpl{
		RoomService: locationService,
	}
}

func (controller *RoomControllerImpl) GetRooms(c echo.Context) error {
	var err error
	var dataResp response.GlobalListDataResponse
	ctx := c.Request().Context()
	var responseData interface{}

	var req request.RoomListRequest
	limit := c.QueryParams().Get("limit")

	req.Limit = limit
	page := c.QueryParams().Get("page")
	req.Page = page
	req.Filter = c.QueryParams().Get("filter")

	resp, err := controller.RoomService.GetRooms(ctx, req)
	if err != nil {
		response.WriteResponseSingleJSON(c, responseData, err)
	}

	for _, each := range resp {
		dataResp.List = append(dataResp.List, each)
	}
	response.WriteResponseListJSON(c, dataResp, nil)
	return err
}

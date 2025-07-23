package controller

import (
	"booking/model/request"
	"booking/model/response"
	"booking/service"
	"strconv"

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

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	req.Limit = limitInt
	page, err := strconv.Atoi(c.QueryParams().Get("page"))
	if err != nil {
		page = 0
	}
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

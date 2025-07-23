package controller

import (
	"booking/app/http/middleware"
	"booking/model/request"
	"booking/model/response"
	"booking/service"
	"booking/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type BookingController interface {
	GetBookings(c echo.Context) error
	CreateBooking(c echo.Context) error
	CancelBooking(c echo.Context) error
	ApproveBooking(c echo.Context) error
}

type BookingControllerImpl struct {
	BookingService service.BookingService
}

func NewBookingController(locationService service.BookingService) BookingController {
	return &BookingControllerImpl{
		BookingService: locationService,
	}
}

func (controller *BookingControllerImpl) GetBookings(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	var responseData response.GlobalListDataResponse

	var req request.BookingListRequest
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

	resp, err := controller.BookingService.GetBookings(ctx, req)
	if err != nil {
		response.WriteResponseSingleJSON(c, responseData, err)
	}

	for _, each := range resp.Data {
		responseData.List = append(responseData.List, each)
	}
	responseData.MetadataResponse = resp.MetadataResponse

	response.WriteResponseListJSON(c, responseData, nil)
	return err
}

func (controller *BookingControllerImpl) CreateBooking(c echo.Context) error {
	var err error

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	var request request.CreateBookingRequest
	var resp response.GlobalSingleResponse
	err = utils.ParseRequestBody(c, &request)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{})
		logrus.Errorf("Error in controller. Parsing error : %v", err)
		return err
	}

	request.UserID = claims.UserID

	if err = c.Validate(&request); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		logrus.Errorf("Error in controller. Validation error : %v", err)
		return err
	}

	ctx := c.Request().Context()

	createBookingResponse, err := controller.BookingService.CreateBooking(ctx, request)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, err)
		logrus.Errorf("Failed to create booking : %v", err)
		return err
	}
	resp.Data = createBookingResponse
	response.WriteResponseSingleJSON(c, resp.Data, err)

	return err
}

func (controller *BookingControllerImpl) CancelBooking(c echo.Context) error {
	var err error

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	var request request.CancelBookingRequest
	var resp response.GlobalSingleResponse
	err = utils.ParseRequestBody(c, &request)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{})
		logrus.Errorf("Error in controller. Parsing error : %v", err)
		return err
	}

	request.UserID = claims.UserID

	if err = c.Validate(&request); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		logrus.Errorf("Error in controller. Validation error : %v", err)
		return err
	}

	ctx := c.Request().Context()

	err = controller.BookingService.CancelBooking(ctx, request)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, err)
		logrus.Errorf("Failed to create booking : %v", err)
		return err
	}
	response.WriteResponseSingleJSON(c, resp.Data, err)

	return err
}

func (controller *BookingControllerImpl) ApproveBooking(c echo.Context) error {
	var err error

	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*middleware.JWTCustomClaims)

	var request request.ApproveBookingRequest
	var resp response.GlobalSingleResponse
	err = utils.ParseRequestBody(c, &request)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{})
		logrus.Errorf("Error in controller. Parsing error : %v", err)
		return err
	}

	// request.UserID = claims.UserID

	if err = c.Validate(&request); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		logrus.Errorf("Error in controller. Validation error : %v", err)
		return err
	}

	ctx := c.Request().Context()

	err = controller.BookingService.ApproveBooking(ctx, request)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, err)
		logrus.Errorf("Failed to create booking : %v", err)
		return err
	}
	response.WriteResponseSingleJSON(c, resp.Data, err)

	return err
}

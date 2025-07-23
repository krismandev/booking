package controller

import (
	"booking/model/request"
	"booking/model/response"
	"booking/service"
	"booking/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	CreateUser(c echo.Context) error
	// GetUser(c echo.Context) error
	ListUser(c echo.Context) error
	// UpdateUser(c echo.Context) error
	// DeleteUser(c echo.Context) error
	// SetPassword(c echo.Context) error
}

type UserControllerImpl struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &UserControllerImpl{
		service: service,
	}
}

func (controller *UserControllerImpl) CreateUser(c echo.Context) error {
	var err error

	ctx := c.Request().Context()

	createUserRequest := request.CreateUserRequest{}
	if err := c.Bind(&createUserRequest); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		return err
	}

	if err := c.Validate(&createUserRequest); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		return err
	}

	data, err := controller.service.CreateUser(ctx, createUserRequest)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, err)
		return err
	}

	response.WriteResponseSingleJSON(c, data, nil)

	return err
}

func (controller *UserControllerImpl) ListUser(c echo.Context) error {

	var responseList response.GlobalListDataResponse
	var err error

	ctx := c.Request().Context()

	var request request.UserListRequest

	limit, err := strconv.Atoi(c.QueryParams().Get("limit"))
	if limit >= 0 {
		limit = 10
		request.Limit = limit
	}
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: "Limit must be number",
		})
		return err
	}

	page, err := strconv.Atoi(c.QueryParams().Get("page"))
	if page >= 0 {
		page = 1
		request.Page = page
	}
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: "Page must be number",
		})
		return err
	}

	resp, err := controller.service.GetUsers(ctx, request)
	if err != nil {
		response.WriteResponseListJSON(c, responseList, err)
		return err
	}

	responseList.Page = resp.Page
	responseList.Limit = resp.Limit
	responseList.TotalPage = resp.TotalPage
	responseList.Count = resp.Count
	responseList.List = make([]any, len(resp.Data))
	for i, each := range resp.Data {
		responseList.List[i] = each
	}

	response.WriteResponseListJSON(c, responseList, err)

	return err
}

package controller

import (
	"booking/model/request"
	"booking/model/response"
	"booking/service"
	"booking/utils"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	CreateUser(c echo.Context) error
	// GetUser(c echo.Context) error
	// ListUser(c echo.Context) error
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

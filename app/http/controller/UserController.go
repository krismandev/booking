package controller

import (
	"booking/model/request"
	"booking/model/response"
	"booking/service"
	"booking/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type UserController interface {
	CreateUser(c echo.Context) error
	// GetUser(c echo.Context) error
	ListUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeactivateUser(c echo.Context) error
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
	request.Limit = c.QueryParams().Get("limit")
	request.Page = c.QueryParams().Get("page")
	request.Filter = c.QueryParams().Get("filter")

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

func (controller *UserControllerImpl) UpdateUser(c echo.Context) error {

	var err error

	ctx := c.Request().Context()

	updateUserRequest := request.UpdateUserRequest{}
	err = utils.ParseRequestBody(c, &updateUserRequest)
	if err != nil {
		logrus.Errorf("Error in controller. Parsing error : %v", err)
		return &utils.BadRequestError{
			Message: "Invalid format",
		}
	}

	if err := c.Validate(&updateUserRequest); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		return err
	}

	if !utils.IsValidUUID(updateUserRequest.UserId) {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: "Invalid ID Fromat",
		})
		return err
	}

	data, err := controller.service.UpdateUser(ctx, updateUserRequest)
	if err != nil {
		logrus.Info("Error updating user : ", err)
		response.WriteResponseSingleJSON(c, nil, err)
		return err
	}

	response.WriteResponseSingleJSON(c, data, err)

	return err
}

func (controller *UserControllerImpl) DeactivateUser(c echo.Context) error {
	var err error

	ctx := c.Request().Context()

	deactivateUserRequest := request.DeactivateUserRequest{}
	err = utils.ParseRequestBody(c, &deactivateUserRequest)
	if err != nil {
		logrus.Errorf("Error in controller. Parsing error : %v", err)
		return &utils.BadRequestError{
			Message: "Invalid format",
		}
	}

	if err := c.Validate(&deactivateUserRequest); err != nil {
		response.WriteResponseSingleJSON(c, nil, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		return err
	}

	err = controller.service.DeactivateUser(ctx, deactivateUserRequest)
	if err != nil {
		logrus.Info("Error deactivating user : ", err)
		response.WriteResponseSingleJSON(c, nil, err)
		return err
	}

	response.WriteResponseSingleJSON(c, nil, err)

	return err
}

package controller

import (
	"booking/model/response"
	"booking/service"

	"github.com/labstack/echo/v4"
)

type RoleController interface {
	GetRoles(c echo.Context) error
}

type RoleControllerImpl struct {
	roleService service.RoleService
}

func NewRoleController(roleService service.RoleService) RoleController {
	return &RoleControllerImpl{
		roleService: roleService,
	}
}

func (controller *RoleControllerImpl) GetRoles(c echo.Context) error {
	var err error
	var dataResp response.GlobalListDataResponse
	ctx := c.Request().Context()

	// request := request.RoleRequest{}

	resp, err := controller.roleService.GetRoles(ctx)
	if err != nil {
		response.WriteResponseListJSON(c, dataResp, err)
		return err
	}

	for _, each := range resp {
		dataResp.List = append(dataResp.List, each)
	}
	response.WriteResponseListJSON(c, dataResp, nil)
	return err
}

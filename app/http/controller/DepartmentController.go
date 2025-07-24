package controller

import (
	"booking/model/response"
	"booking/service"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DepartmentController interface {
	GetDepartments(c echo.Context) error
}

type DepartmentControllerImpl struct {
	DepartmentService service.DepartmentService
}

func NewDepartmentController(departmentService service.DepartmentService) DepartmentController {
	return &DepartmentControllerImpl{
		DepartmentService: departmentService,
	}
}

func (controller *DepartmentControllerImpl) GetDepartments(c echo.Context) error {
	var err error
	var dataResp response.GlobalListDataResponse
	ctx := c.Request().Context()
	var responseData interface{}

	resp := controller.DepartmentService.GetDepartments(ctx)
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

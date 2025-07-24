package service

import (
	"booking/model/response"
	"booking/repository"
	"context"
)

type DepartmentService interface {
	GetDepartments(ctx context.Context) []response.DepartmentResponse
}

type DepartmentServiceImpl struct {
	repository repository.DepartmentRepository
}

func NewDepartmentService(repository repository.DepartmentRepository) DepartmentService {
	return &DepartmentServiceImpl{
		repository: repository,
	}
}

func (service *DepartmentServiceImpl) GetDepartments(ctx context.Context) []response.DepartmentResponse {
	var output []response.DepartmentResponse

	departments := service.repository.GetDepartments()

	if len(departments) > 0 {
		for _, each := range departments {
			output = append(output, response.ToDepartmentResponse(each))
		}
	}

	return output
}

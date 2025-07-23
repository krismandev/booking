package service

import (
	"booking/model/response"
	"booking/repository"
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]response.RoleResponse, error)
}

type RoleServiceImpl struct {
	repository repository.RoleRepository
}

func NewRoleService(repository repository.RoleRepository) RoleService {
	return &RoleServiceImpl{
		repository: repository,
	}
}

func (service *RoleServiceImpl) GetRoles(ctx context.Context) ([]response.RoleResponse, error) {
	var output []response.RoleResponse

	var err error

	roles := service.repository.GetRoles()

	if len(roles) > 0 {
		for i, _ := range roles {
			privilegesEncoded, err := json.Marshal(roles[i].Privileges)
			if err != nil {
				logrus.Errorf("Error in service. failed to marshalling : %v", err)
				return output, err
			}

			roles[i].Privileges = string(privilegesEncoded)
			output = append(output, response.ToRoleResponse(roles[i]))
		}
	}

	return output, err
}

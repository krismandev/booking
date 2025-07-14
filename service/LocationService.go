package service

import (
	"booking/model"
	"booking/model/response"
	"booking/repository"
	"context"
)

type LocationService interface {
	GetLocations(ctx context.Context) []response.LocationResponse
}

type LocationServiceImpl struct {
	repository repository.LocationRepository
}

func NewLocationService(repository repository.LocationRepository) LocationService {
	return &LocationServiceImpl{
		repository: repository,
	}
}

func (service *LocationServiceImpl) GetLocations(ctx context.Context) []response.LocationResponse {
	var output []response.LocationResponse

	var locations *[]model.Location = new([]model.Location)

	service.repository.GetLocations(locations)

	if locations != nil {
		for _, each := range *locations {
			output = append(output, response.ToLocationResponse(each))
		}
	}

	return output
}

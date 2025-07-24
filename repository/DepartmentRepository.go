package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type DepartmentRepository interface {
	GetDepartments() []model.Department
	GetDepartmentByIDs(departmentIDs []string) []model.Department
}

type DepartmentRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewDepartmentRepository(db *connection.DBConnection) DepartmentRepository {
	return &DepartmentRepositoryImpl{
		dbConn: db,
	}
}

func (repository *DepartmentRepositoryImpl) GetDepartments() []model.Department {
	departments := []model.Department{}
	err := repository.dbConn.DB.Model(&departments).Find(&departments).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return departments
}

func (repository *DepartmentRepositoryImpl) GetDepartmentByIDs(departmentIDs []string) []model.Department {
	departments := []model.Department{}

	err := repository.dbConn.DB.Model(&departments).Where("id IN ?", departmentIDs).Find(&departments).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return departments
}

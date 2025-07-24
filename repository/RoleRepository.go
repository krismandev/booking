package repository

import (
	connection "booking/connection/database"
	"booking/model"

	"github.com/sirupsen/logrus"
)

type RoleRepository interface {
	GetRoles() []model.Role
	GetUserRole(userId string) model.UserRole
	GetRoleByIDs(roleIDs []string) []model.Role
	GetRoleByID(roleID string) model.Role
	CreateUserRole(userID, roleID string) error
	GetListUserRole(userIDs []string) []model.UserRole
}

type RoleRepositoryImpl struct {
	dbConn *connection.DBConnection
}

func NewRoleRepository(db *connection.DBConnection) RoleRepository {
	return &RoleRepositoryImpl{
		dbConn: db,
	}
}

func (repository *RoleRepositoryImpl) GetRoles() []model.Role {
	roles := []model.Role{}
	err := repository.dbConn.DB.Model(&roles).Find(&roles).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return roles
}

func (repository *RoleRepositoryImpl) GetRoleByIDs(roleIDs []string) []model.Role {
	roles := []model.Role{}

	err := repository.dbConn.DB.Model(&roles).Where("id IN ?", roleIDs).Find(&roles).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return roles
}

func (repository *RoleRepositoryImpl) GetUserRole(userId string) model.UserRole {
	dt := model.UserRole{}

	err := repository.dbConn.DB.Model(&dt).Where("userid = ?", userId).First(&dt).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}
	return dt
}

func (repository *RoleRepositoryImpl) GetRoleByID(roleID string) model.Role {

	dt := model.Role{}

	err := repository.dbConn.DB.Model(&dt).Where("id = ?", roleID).First(&dt).Error
	if err != nil {
		logrus.Errorf("SQL Error : %v", err)
	}

	return dt
}

func (repository *RoleRepositoryImpl) CreateUserRole(userID, roleID string) error {
	var err error

	err = repository.dbConn.DB.Create(&model.UserRole{UserID: userID, RoleID: roleID}).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return err
}

func (repository *RoleRepositoryImpl) GetListUserRole(userIDs []string) []model.UserRole {
	var output []model.UserRole

	err := repository.dbConn.DB.Preload("Role").Where("userid IN ?", userIDs).Find(&output).Error

	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}

	return output
}

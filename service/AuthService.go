package service

import (
	"context"
	"encoding/json"
	"errors"

	"booking/app/http/middleware"
	connection "booking/connection/database"
	"booking/model"
	"booking/model/request"
	"booking/model/response"
	"booking/repository"
	"booking/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, request request.LoginRequest) (response.LoginResponse, error)
	AuthUserDetail(ctx context.Context, userID string) (response.AuthUserDetailResponse, error)
	// RefreshToken(ctx context.Context, userID string) (map[string]any, error)
}

type AuthServiceImpl struct {
	repository     repository.UserRepository
	dbConn         *connection.DBConnection
	roleRepository repository.RoleRepository
	// validate   *validator.Validate
}

func NewAuthService(repository repository.UserRepository, dbConn *connection.DBConnection, roleRepository repository.RoleRepository) AuthService {
	return &AuthServiceImpl{
		repository:     repository,
		dbConn:         dbConn,
		roleRepository: roleRepository,
		// validate:   validate,
	}
}

func (service *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (response.LoginResponse, error) {
	var err error
	var response response.LoginResponse

	user, err := service.repository.FindOneUserByEmail(req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response, &utils.BadRequestError{
			Code:    400,
			Message: "Invalid email address or password",
		}
	}

	role := service.roleRepository.GetUserRole(user.ID)
	if len(role.ID) == 0 {
		return response, &utils.InternalServerError{Message: "Something went wrong"}
	}

	checkPassword := utils.ComparePass([]byte(*user.Password), []byte(req.Password))
	if !checkPassword {
		return response, &utils.BadRequestError{
			Code:    400,
			Message: "Invalid email address or password",
		}
	}

	accessToken, refreshToken, expiredAtStr, err := middleware.GenerateJWT(user.ID, role.RoleID)
	if err != nil {
		logrus.Errorf("Error Generating JWT: %v", err)
		return response, &utils.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}
	}

	response.AccessToken = accessToken
	response.RefreshToken = refreshToken
	response.ExpiryTime = expiredAtStr

	return response, err
}

func (service *AuthServiceImpl) AuthUserDetail(ctx context.Context, userID string) (response.AuthUserDetailResponse, error) {
	var resp response.AuthUserDetailResponse

	user, err := service.repository.FindUserById(userID)

	if err != nil {
		logrus.Errorf("Failed get user data %v", err)
		return resp, err
	}

	userRole := service.roleRepository.GetUserRole(user.ID)

	var role model.Role
	role = service.roleRepository.GetRoleByID(userRole.RoleID)

	if len(userRole.ID) == 0 {
		logrus.Errorf("Error in service. Role Not Found : %v", err)
		return resp, err
	}

	resp.User.ID = user.ID
	resp.User.Name = user.Name
	resp.User.Email = user.Email
	resp.User.CreatedAt = user.CreatedAt
	resp.User.Role = &role.Name

	resp.Role.ID = userRole.RoleID
	resp.Role.Name = role.Name

	privilegesStr, _ := json.Marshal(role.Privileges)
	resp.Role.Privileges = string(privilegesStr)

	return resp, err
}

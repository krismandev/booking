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
	RefreshToken(ctx context.Context, userID string) (map[string]any, error)
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

	checkPassword := utils.ComparePass([]byte(*user.Password), []byte(req.Password))
	if !checkPassword {
		return response, &utils.BadRequestError{
			Code:    400,
			Message: "Invalid email address or password",
		}
	}

	accessToken, refreshToken, expiredAtStr, err := middleware.GenerateJWT(user.ID)
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

func (service *AuthServiceImpl) RefreshToken(ctx context.Context, userID string) (map[string]any, error) {

	newAccessToken, _, expiryTimeStr, err := middleware.GenerateJWT(userID)
	if err != nil {
		return map[string]any{"error": "could not generate token"}, err
	}

	return map[string]any{
		"accessToken": newAccessToken,
		"expiryTime":  expiryTimeStr,
	}, err
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

	var privileges []string

	err = json.Unmarshal([]byte(role.Privileges), &privileges)
	if err != nil {
		logrus.Errorf("Error when unmarshalling %v", err)
		return resp, err
	}

	resp.User.ID = user.ID
	resp.User.Name = user.Name
	resp.User.Email = user.Email
	resp.User.CreatedAt = user.CreatedAt
	resp.User.Role = role.Name

	resp.Role.ID = userRole.RoleID
	resp.Role.Name = role.Name
	resp.Role.Privileges = privileges

	return resp, err
}

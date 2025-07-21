package service

import (
	"context"
	"errors"

	"booking/app/http/middleware"
	connection "booking/connection/database"
	"booking/model/request"
	"booking/model/response"
	"booking/repository"
	"booking/utils"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, request request.LoginRequest) (response.LoginResponse, error)
	RefreshToken(ctx context.Context, userID string) (map[string]any, error)
}

type AuthServiceImpl struct {
	repository repository.UserRepository
	dbConn     *connection.DBConnection
	validate   *validator.Validate
}

func NewAuthService(repository repository.UserRepository, dbConn *connection.DBConnection, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		repository: repository,
		dbConn:     dbConn,
		validate:   validate,
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

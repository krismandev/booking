package controller

import (
	"booking/app/http/middleware"
	"booking/model/request"
	"booking/model/response"
	"booking/service"
	"booking/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthController interface {
	Login(c echo.Context) error
	// RefreshToken(c echo.Context) error
	AuthUserDetail(c echo.Context) error
}

type AuthControllerImpl struct {
	service service.AuthService
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func NewAuthController(service service.AuthService) AuthController {
	return &AuthControllerImpl{
		service: service,
	}
}

func (controller *AuthControllerImpl) Login(c echo.Context) error {
	var responseData any
	var err error

	ctx := c.Request().Context()

	loginRequest := request.LoginRequest{}

	if err := c.Bind(&loginRequest); err != nil {
		response.WriteResponseSingleJSON(c, responseData, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		return err
	}

	if err := c.Validate(&loginRequest); err != nil {
		response.WriteResponseSingleJSON(c, responseData, &utils.BadRequestError{
			Code:    400,
			Message: utils.FormatValidationErrors(err),
		})
		return err
	}

	responseData, err = controller.service.Login(ctx, loginRequest)
	if err != nil {
		response.WriteResponseSingleJSON(c, nil, err)
		return err
	}

	response.WriteResponseSingleJSON(c, responseData, nil)

	return err
}

// func (controller *AuthControllerImpl) RefreshToken(c echo.Context) error {
// 	var request RefreshRequest
// 	if err := c.Bind(&request); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
// 	}

// 	refreshToken := request.RefreshToken

// 	ctx := c.Request().Context()

// 	token, err := middleware.ParseJWT(refreshToken)

// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]string{
// 			"error": "invalid refresh token",
// 		})
// 	}

// 	userID := token.UserID

// 	data, err := controller.service.RefreshToken(ctx, userID)
// 	if err != nil {
// 		response.WriteResponseSingleJSON(c, nil, err)
// 		return err
// 	}

// 	response.WriteResponseSingleJSON(c, data, nil)

// 	return err
// }

func (controller *AuthControllerImpl) AuthUserDetail(c echo.Context) error {
	var err error

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	ctx := c.Request().Context()
	data, err := controller.service.AuthUserDetail(ctx, claims.UserID)
	if err != nil {
		logrus.Errorf("Failed AuthUserDetail : %v", err)
		response.WriteResponseSingleJSON(c, nil, err)
		return err
	}

	response.WriteResponseSingleJSON(c, data, err)

	return err
}

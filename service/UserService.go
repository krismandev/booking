package service

import (
	"context"
	"time"

	connection "booking/connection/database"
	"booking/model"
	"booking/model/request"
	"booking/model/response"
	"booking/repository"
	"booking/utils"

	"github.com/google/uuid"
	"github.com/kenshaw/envcfg"
)

type UserService interface {
	CreateUser(ctx context.Context, request request.CreateUserRequest) (response.CreateUserResponse, error)
	// GetUser(ctx context.Context, id string) (response.CreateUserResponse, error)
	// ListUser(ctx context.Context, request request.UserListRequest) ([]response.UserResponse, int, int, int, int, error)
	// UpdateUser(ctx context.Context, request request.UpdateUserRequest) (response.UpdateUserResponse, error)
	// DeleteUser(ctx context.Context, id string) (response.GlobalJSONResponse, error)
	// SetPassword(ctx context.Context, token string, req request.SetPasswordRequest) error
}

type userServiceImpl struct {
	repository repository.UserRepository
	dbConn     *connection.DBConnection
	config     *envcfg.Envcfg
}

func NewUserService(repository repository.UserRepository, dbConn *connection.DBConnection) UserService {
	return &userServiceImpl{
		repository: repository,
		dbConn:     dbConn,
	}
}

type UserStatus string

const (
	Active   UserStatus = "ACTIVE"
	Review   UserStatus = "REVIEW"
	Disabled UserStatus = "DISABLED"
)

func (service *userServiceImpl) CreateUser(ctx context.Context, req request.CreateUserRequest) (response.CreateUserResponse, error) {
	var err error
	var resp response.CreateUserResponse

	nowStr := time.Now().Format("2006-01-02 15:04:05")

	var user model.User
	user.ID = uuid.New().String()
	user.Email = req.Email
	user.Name = req.Name
	pwd := utils.HashPassword(req.Password)
	user.Password = &pwd
	user.CreatedAt = &nowStr

	tx := service.dbConn.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// if err := tx.Error; err != nil {
	// 	logrus.Errorf("CreateUser Database Connection Error: %v", err)
	// 	return resp, &utils.InternalServerError{
	// 		Code:    500,
	// 		Message: "Internal Server Error",
	// 	}
	// }

	// checkUser := service.repository.FindUserByEmail(tx, req.Email)
	// if len(checkUser) > 0 {
	// 	return resp, &utils.UnprocessableContentError{
	// 		Code:    422,
	// 		Message: "User with email " + req.Email + " already exists",
	// 	}
	// }

	// if err := tx.Create(&user).Error; err != nil {
	// 	logrus.Errorf("CreateUser Database Connection Error: %v", err)
	// 	tx.Rollback()
	// 	return resp, &utils.InternalServerError{
	// 		Code:    500,
	// 		Message: "Internal Server Error",
	// 	}
	// }

	err = service.repository.InsertUser(tx, user)
	tx.Commit()

	resp.Id = user.ID
	resp.Name = user.Name
	resp.Email = user.Email
	resp.CreatedTime = *user.CreatedAt

	// mailReq := request.MailRequest{}
	// mailReq.MailType = "userverification"
	// mailReq.Subject = "Verification Mail"
	// mailReq.To = req.Email
	// mailReq.Url = service.config.GetString("appUrl") + "/verify/" + "verifyToken"

	// encoded, err := json.Marshal(mailReq)
	// jsonString := string(encoded)
	// service.mailProducer.PublishMessage(jsonString)
	return resp, err
}

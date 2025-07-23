package service

import (
	"context"
	"errors"
	"math"
	"time"

	connection "booking/connection/database"
	"booking/model"
	"booking/model/request"
	"booking/model/response"
	"booking/repository"
	"booking/utils"

	"github.com/google/uuid"
	"github.com/kenshaw/envcfg"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, request request.CreateUserRequest) (response.CreateUserResponse, error)
	// GetUser(ctx context.Context, id string) (response.CreateUserResponse, error)
	GetUsers(ctx context.Context, request request.UserListRequest) (response.UserListResponse, error)
	// UpdateUser(ctx context.Context, request request.UpdateUserRequest) (response.UpdateUserResponse, error)
	// DeleteUser(ctx context.Context, id string) (response.GlobalJSONResponse, error)
	// SetPassword(ctx context.Context, token string, req request.SetPasswordRequest) error
}

type userServiceImpl struct {
	repository     repository.UserRepository
	roleRepository repository.RoleRepository
	dbConn         *connection.DBConnection
	config         *envcfg.Envcfg
}

func NewUserService(repository repository.UserRepository, dbConn *connection.DBConnection, roleRepository repository.RoleRepository) UserService {
	return &userServiceImpl{
		repository:     repository,
		dbConn:         dbConn,
		roleRepository: roleRepository,
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

	id, err := service.repository.InsertUser(user)
	if err != nil {
		logrus.Errorf("Failed to create user : %v", err)
		return resp, err
	}

	// userRole := model.UserRole{UserID: id, RoleID: req.RoleID}

	err = service.roleRepository.CreateUserRole(id, req.RoleID)
	if err != nil {
		logrus.Errorf("Failed to create user role : %v", err)
		return resp, err
	}

	resp.Id = user.ID
	resp.Name = user.Name
	resp.Email = user.Email
	resp.CreatedTime = *user.CreatedAt

	return resp, err
}

func (service *userServiceImpl) GetUser(ctx context.Context, id string) (response.UserResponse, error) {
	var resp response.UserResponse

	user, err := service.repository.FindUserById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return resp, &utils.NotFoundError{
			Code:    400,
			Message: "User not found",
		}
	}

	resp = response.ToUserResponse(user)

	return resp, nil
}

func (service *userServiceImpl) ListUser(ctx context.Context, request request.UserListRequest) ([]response.UserResponse, int, int, int, int, error) {
	var data []response.UserResponse

	limit := request.Limit
	page := request.Page
	offset := (page - 1) * limit

	users, count := service.repository.GetUserList(limit, offset)

	for _, each := range users {
		dt := response.ToUserResponse(each)
		data = append(data, dt)
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return data, page, limit, totalPages, int(count), nil
}

func (service *userServiceImpl) UpdateUser(ctx context.Context, request request.UpdateUserRequest) (response.UpdateUserResponse, error) {
	var resp response.UpdateUserResponse

	existingUser, err := service.repository.FindUserById(request.UserId)
	if err != nil {
		logrus.Errorf("Error : %v", err)
		return resp, &utils.InternalServerError{}
	}

	if len(existingUser.ID) == 0 {
		logrus.Errorf("Error in service. User not found id=%v", request.UserId)
		return resp, &utils.NotFoundError{Message: "User not found"}
	}

	if len(request.Name) > 0 {
		existingUser.Name = request.Name
	}
	if len(request.Email) > 0 {
		existingUser.Email = request.Email
	}

	if len(request.Password) > 0 {
		pwd := utils.HashPassword(request.Password)
		existingUser.Password = &pwd
	}

	timeNow := utils.TimeNowString()
	existingUser.UpdatedAt = &timeNow

	err = service.repository.UpdateUser(existingUser)
	if err != nil {
		logrus.Errorf("Failed to update user : %v", err)
		return resp, err
	}
	// 	tx.Commit()

	resp.ID = existingUser.ID
	resp.Name = existingUser.Name
	resp.Email = existingUser.Email
	resp.CreatedAt = *existingUser.CreatedAt
	resp.UpdatedAt = *existingUser.UpdatedAt

	return resp, nil
}

func (service *userServiceImpl) DeleteUser(ctx context.Context, id string) (response.GlobalJSONResponse, error) {
	var resp response.GlobalJSONResponse

	tx := service.dbConn.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		logrus.Errorf("Database Connection Error: %v", err)
		return resp, &utils.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}
	}

	var userModel model.User
	checkUser := tx.Where(&model.User{ID: id}).First(&userModel)

	if errors.Is(checkUser.Error, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return resp, &utils.UnprocessableContentError{
			Code:    422,
			Message: "User with ID " + id + " not found or already deleted",
		}
	}

	if err := tx.Delete(&userModel).Error; err != nil {
		logrus.Errorf("Database Connection Error: %v", err)
		tx.Rollback()
		return resp, &utils.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}
	}

	tx.Commit()

	resp.Message = "success"

	return resp, nil
}

func (service *userServiceImpl) GetUsers(ctx context.Context, request request.UserListRequest) (response.UserListResponse, error) {
	var resp response.UserListResponse
	var err error

	return resp, err
}

// set password for the first time
// func (service *userServiceImpl) SetPassword(ctx context.Context, token string, req request.SetPasswordRequest) error {
// 	var err error

// 	merchant := service.merchantRepository.FindMerchantByVerifyToken(token)
// 	if len(merchant.ID) == 0 {
// 		logrus.Errorf("Error in Service : Merchant not found")
// 		return &utils.BadRequestError{Message: "Merchant not found"}
// 	}

// 	user := service.repository.FindUserByMerchantID(merchant.ID)
// 	if len(user.ID) == 0 {
// 		logrus.Errorf("Error in Service : Merchant not found")
// 		return &utils.BadRequestError{Message: "Merchant not found"}
// 	}

// 	if user.Password != nil {
// 		logrus.Errorf("Password has already been set")
// 		return &utils.ForbiddenError{Message: "Password has already been set"}
// 	}

// 	password := utils.HashPassword(req.Password)
// 	err = service.repository.SetPassword(user.ID, password)
// 	if err != nil {
// 		logrus.Errorf("Error in Service : Failed to set password %v", err)
// 		return err
// 	}
// 	return err
// }

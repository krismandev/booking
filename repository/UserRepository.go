package repository

import (
	connection "booking/connection/database"
	"booking/model"
	"booking/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(dt model.User) (string, error)
	GetUser(userID string) map[string]string
	DeleteUser(userID string) error
	FindUserByEmail(tx *gorm.DB, email string) []model.User
	FindOneUserByEmail(email string) (model.User, error)
	FindUserById(userId string) (model.User, error)
	FindUserByIDs(userIDs []string) ([]model.User, error)
	GetUserList(filter model.UserListQueryFilter) ([]model.User, int64)
	SetPassword(userId string, password string) error
	FindUserByMerchantID(merchantID string) model.User
	UpdateUser(dt model.User) error
	DeactivateUser(userID string) error
}

type UserRepositoryImpl struct {
	db *connection.DBConnection
}

func NewUserRepository(db *connection.DBConnection) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (repository *UserRepositoryImpl) InsertUser(dt model.User) (string, error) {
	err := repository.db.DB.Create(&dt).Error
	if err != nil {
		logrus.Errorf("SQL Error : %v", err)
	}

	return dt.ID, err
}

func (repository *UserRepositoryImpl) GetUser(userID string) map[string]string {
	var resultQuery map[string]any
	var output map[string]string

	qry := `SELECT id, createdtime, email, name
	FROM public.users 
	WHERE id = ?::UUID;`

	err := repository.db.Raw(qry, userID).Scan(&resultQuery).Error
	if err != nil {
		logrus.Errorf("SQL Error : %v", err)
	}

	output = utils.ConvertToMapOfString(resultQuery)
	return output
}

func (repository *UserRepositoryImpl) DeleteUser(userID string) error {
	var resultQuery any

	qry := `DELETE FROM users 
	WHERE id = ?::UUID;`

	err := repository.db.Raw(qry, userID).Scan(&resultQuery).Error
	if err != nil {
		logrus.Errorf("SQL Error : %v", err)
	}

	return err
}

func (repo *UserRepositoryImpl) FindUserByEmail(tx *gorm.DB, email string) []model.User {
	var output []model.User
	if err := tx.Where(model.User{Email: email}).Find(&output).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("FindUserByEmail SQL Error : %v", err)
		return output
	}

	return output
}

func (repo *UserRepositoryImpl) FindOneUserByEmail(email string) (model.User, error) {
	var output model.User
	var err error

	if err := repo.db.DB.Where(model.User{Email: email}).First(&output).Error; err != nil {
		logrus.Errorf("FindOneUserByEmail SQL Error : %v", err)
		return output, err
	}

	return output, err
}

func (repo *UserRepositoryImpl) FindUserById(userId string) (model.User, error) {
	var output model.User
	var err error

	if err := repo.db.DB.Where(model.User{ID: userId}).First(&output).Error; err != nil {
		logrus.Errorf("FindUserById SQL Error : %v", err)
		return output, err
	}

	return output, err
}

func (repo *UserRepositoryImpl) GetUserList(filter model.UserListQueryFilter) ([]model.User, int64) {
	var output []model.User
	var count int64

	qry := repo.db.DB.Model(&output).Scopes(repo.db.Order(filter.GlobalQueryFilter)).Scopes(repo.db.Paginate(filter.GlobalQueryFilter))

	if len(filter.Name) > 0 {
		qry = qry.Where("name like ?", "%"+filter.Name+"%")
	}

	err := qry.Find(&output).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
		return output, count
	}

	// result := repo.db.DB.Where(&model.User{}).Limit(limit).Offset(offset).Find(&output)
	// if result.Error != nil {
	// 	logrus.Errorf("GetUserList SQL Error : %v", result.Error)
	// 	return output, count
	// }

	countResult := repo.db.DB.Where(&model.User{}).Count(&count)
	if countResult.Error != nil {
		logrus.Errorf("Error in repository : %v", countResult.Error)
		return output, count
	}

	return output, count
}

func (repo *UserRepositoryImpl) SetPassword(userId string, password string) error {
	var err error

	qry := "UPDATE users SET password = ? WHERE id = ?"
	err = repo.db.Exec(qry, password, userId).Error
	if err != nil {
		logrus.Errorf("SQL Error : %v", err)
	}
	return err
}

func (repo *UserRepositoryImpl) FindUserByMerchantID(merchantID string) model.User {
	var user model.User
	err := repo.db.DB.Where("merchantId = ?", merchantID).First(&user).Error
	if err != nil {
		logrus.Errorf("SQL Error : %v", err)
	}

	return user
}

func (repo *UserRepositoryImpl) UpdateUser(dt model.User) error {
	var err error

	err = repo.db.DB.Save(&model.User{ID: dt.ID, Name: dt.Name, Email: dt.Email, Password: dt.Password, UpdatedAt: dt.UpdatedAt}).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
		return err
	}

	return err
}

func (repository *UserRepositoryImpl) FindUserByIDs(userIDs []string) ([]model.User, error) {
	var output []model.User
	var err error

	err = repository.db.DB.Where("id IN ?", userIDs).Find(&output).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
		return output, err
	}

	return output, err
}

func (repository *UserRepositoryImpl) DeactivateUser(userID string) error {
	var err error

	err = repository.db.DB.Model(&model.User{}).Where("id = ?", userID).Update("isactive", false).Error
	if err != nil {
		logrus.Errorf("Error in repository : %v", err)
	}
	return err
}

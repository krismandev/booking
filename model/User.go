package model

const (
	REVIEW   = "REVIEW"
	ACTIVE   = "ACTIVE"
	DISABLED = "DISABLED"
)

type User struct {
	ID        string  `gorm:"column:id; primaryKey;"`
	CreatedAt *string `gorm:"column:createdat"`
	UpdatedAt *string `gorm:"column:updatedat"`
	Password  *string `gorm:"column:password"`
	Email     string  `gorm:"column:email"`
	Name      string  `gorm:"column:name"`
}

type UserListQueryFilter struct {
	GlobalQueryFilter
	User
}

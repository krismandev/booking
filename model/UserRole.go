package model

type UserRole struct {
	ID     string `gorm:"column:id"`
	UserID string `gorm:"column:userid"`
	RoleID string `gorm:"column:roleid"`
}

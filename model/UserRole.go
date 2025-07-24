package model

type UserRole struct {
	ID     string `gorm:"column:id;default:uuid_generate_v4();primaryKey"`
	UserID string `gorm:"column:userid"`
	RoleID string `gorm:"column:roleid"`

	Role Role
}

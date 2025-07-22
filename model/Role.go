package model

type Role struct {
	ID         string `gorm:"column:id;primaryKey"`
	Name       string `gorm:"column:name"`
	Privileges string `gorm:"column:privileges"`
}

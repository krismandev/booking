package model

type Location struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

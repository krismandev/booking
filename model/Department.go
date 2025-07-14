package model

type Department struct {
	ID   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

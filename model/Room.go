package model

const (
	SMALL  = "SMALL"
	MEDIUM = "MEDIUM"
	LARGE  = "LARGE"
)

type Room struct {
	ID         string `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	LocationID string `gorm:"column:locationid"`
	Floor      string `gorm:"column:floor"`
	Capacity   string `gorm:"column:capacity"`
	CreatedAt  string `gorm:"column:createdat"`
	UpdatedAt  string `gorm:"column:updatedat"`
}

package model

type PropertyImage struct {
	ID         uint   `gorm:"primaryKey"`
	PropertyID int    `gorm:"not null;index" validate:"required"`
	URL        string `gorm:"type:text" validate:"required,url"`
}

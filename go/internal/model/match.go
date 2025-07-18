package model

import "time"

type Match struct {
    ID          uint      `gorm:"primaryKey"`
    PropertyID  uint      `gorm:"not null" validate:"required"`
    ContactID   uint      `gorm:"not null" validate:"required"`
    ClientID    uint      `gorm:"not null" validate:"required"`
    Score       float64   `gorm:"not null" validate:"required,gte=0,lte=100"`
    Explanation string    `gorm:"type:text"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
}


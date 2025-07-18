package model

import "time"

type Agent struct {
    ID                   int       `gorm:"primaryKey" json:"id"`
    UserID               int       `gorm:"not null" validate:"required" json:"user_id"`
    NationalID           string    `gorm:"not null;unique" validate:"required,len=10" json:"national_id"`
    CommercialRegister   string    `gorm:"not null" validate:"required" json:"commercial_register"`
    Verified             bool      `json:"verified"`
    CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
}


package model

import "time"

type Client struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    UserID    int       `gorm:"not null;unique" validate:"required" json:"user_id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

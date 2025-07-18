package model

import "time"

type User struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    FirstName string    `gorm:"not null" validate:"required,alpha" json:"first_name"`
    LastName  string    `gorm:"not null" validate:"required,alpha" json:"last_name"`
    Phone     string    `gorm:"not null;unique" validate:"required,e164" json:"phone_number"`
    Email     string    `gorm:"not null;unique" validate:"required,email" json:"email"`
    Password  string    `gorm:"not null" validate:"required,min=6" json:"-"`
    Role      string    `gorm:"not null" validate:"required,oneof=agent client" json:"role"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}


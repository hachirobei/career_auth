package models

import (
    "gorm.io/gorm"
    "time"
)

type UserToken struct {
    gorm.Model
    UserID    uint
    Token     string    `gorm:"size:500"`
    ExpiresAt time.Time
}
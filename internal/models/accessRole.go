package models

import (
    "gorm.io/gorm"
)

type AccessRole struct {
    gorm.Model
	ID          int    `json:"id" gorm:"primaryKey"`
    UserID    uint
	RoleID    uint
    Status      int         `json:"status" gorm:"type:int;not null;default:1"`
	Base
}
package models

import (
    "gorm.io/gorm"
)

type Role struct {
    ID          int         `json:"id" gorm:"primaryKey"`
    Description string      `json:"description" gorm:"type:text;not null;unique"`
    Status      int         `json:"status" gorm:"type:int;not null;default:1"`
	Base
}
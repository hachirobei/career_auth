package models

import (
    "gorm.io/gorm"
    "time"
)

type Users struct {
    ID          int         `json:"id" gorm:"primaryKey"`
    FullName    string      `json:"full_name" gorm:"type:text;not null"`
    Email       string      `json:"email" gorm:"type:text;not null;unique"`
    Phone       string      `json:"phone" gorm:"type:text;not null"`
    Username    string      `json:"username" gorm:"type:text;not null;unique"`
    Password    string      `json:"password" gorm:"type:text;not null"`
    Status      int         `json:"status" gorm:"type:int;not null;default:1"`
    Base
}

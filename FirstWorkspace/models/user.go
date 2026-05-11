package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `json:"email" gorm:"unique;not null"`
	Password     string `json:"password" gorm:"not null"`
	RefreshToken string `json:"refresh_token" gorm:"default:null"`
}

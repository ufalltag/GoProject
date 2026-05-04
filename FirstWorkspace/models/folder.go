package models

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name      string     `json:"name" gorm:"unique"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	User      User       `json:"-" gorm:"foreignKey:UserID"`
	Bookmarks []Bookmark `json:"-" gorm:"foreignKey:FolderID"`
}

type FolderResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

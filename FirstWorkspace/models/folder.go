package models

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name      string     `json:"name" gorm:"uniqueIndex:idx_user_folder_name"`
	UserID    uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_folder_name"`
	User      User       `json:"-" gorm:"foreignKey:UserID"`
	Bookmarks []Bookmark `json:"-" gorm:"foreignKey:FolderID"`
}

type FolderResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

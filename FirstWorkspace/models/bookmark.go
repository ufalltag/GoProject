package models

import "gorm.io/gorm"

type Bookmark struct {
	gorm.Model
	URL      string `json:"url" gorm:"not null"`
	Title    string `json:"title"`
	UserID   uint   `json:"user_id" gorm:"not null"`
	FolderID uint   `json:"folder_id" gorm:"not null"`
	User     User   `json:"-" gorm:"foreignKey:UserID"`
	Folder   Folder `json:"-" gorm:"foreignKey:FolderID"`
}
type BookmarkResponse struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	FolderID uint   `json:"folder_id"`
}

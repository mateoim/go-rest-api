package models

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	Name  string `json:"name" binding:"required"`
	Users []User `json:"users"`
}

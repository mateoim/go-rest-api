package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName      string `json:"first-name" binding:"required"`
	LastName       string `json:"last-name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	OrganizationID uint   `json:"organization" binding:"required"`
}

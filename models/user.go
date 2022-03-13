package models

type User struct {
	ID             uint         `json:"id" gorm:"primarykey"`
	FirstName      string       `json:"first-name" binding:"required"`
	LastName       string       `json:"last-name" binding:"required"`
	Email          string       `json:"email" binding:"required,email"`
	OrganizationID uint         `json:"organization" binding:"required"`
	Events         []Event      `gorm:"many2many:event_users;"`
	Meetings       []Meeting    `json:"-" gorm:"many2many:user_meetings"`
	Invitations    []Invitation `json:"-"`
}

package models

type Organization struct {
	ID    uint   `json:"id" gorm:"primarykey"`
	Name  string `json:"name" binding:"required"`
	Users []User `json:"users"`
}

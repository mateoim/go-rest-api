package models

type Event struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	StartDate   CustomTime `json:"start-date" binding:"required"`
	EndDate     CustomTime `json:"end-date" binding:"required"`
	Users       []User     `json:"-" gorm:"many2many:event_users;"`
	Meetings    []Meeting  `json:"-"`
}

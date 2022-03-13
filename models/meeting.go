package models

type Meeting struct {
	ID          uint         `json:"id" gorm:"primarykey"`
	Title       string       `json:"title"`
	StartDate   CustomTime   `json:"start-date" binding:"required"`
	EndDate     CustomTime   `json:"end-date" binding:"required"`
	EventID     uint         `json:"event"`
	Invitations []Invitation `json:"-"`
}

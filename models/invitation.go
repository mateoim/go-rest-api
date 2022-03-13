package models

import (
	"encoding/json"
)

type Invitation struct {
	ID        uint             `json:"id" gorm:"primarykey"`
	MeetingID uint             `json:"meeting"`
	UserID    uint             `json:"user" binding:"required"`
	Status    InvitationStatus `json:"status"`
}

type InvitationStatus int

const (
	Pending InvitationStatus = iota + 1
	Accepted
	Rejected
)

func (status InvitationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

func (status InvitationStatus) String() string {
	return [...]string{"Pending", "Accepted", "Rejected"}[status-1]
}

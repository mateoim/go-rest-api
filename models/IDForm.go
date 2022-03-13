package models

type IDForm struct {
	ID uint `json:"id" binding:"required"`
}

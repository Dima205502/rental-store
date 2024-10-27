package models

type User struct {
	Nickname string `json:"nickname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

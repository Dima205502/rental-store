package models

type Thing struct {
	Owner       string `json:"owner" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" validate:"required"`
	Available   bool   `json:"available" validate:"required"`
}

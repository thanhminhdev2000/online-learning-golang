package models

type Contact struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

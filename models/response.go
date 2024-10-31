package models

type Message struct {
	Message string `json:"message" validate:"required"`
}

type Error struct {
	Error string `json:"error" validate:"required"`
}

type Paging struct {
	Total int `json:"total" validate:"required"`
	Page  int `json:"page" validate:"required"`
	Limit int `json:"limit" validate:"required"`
}

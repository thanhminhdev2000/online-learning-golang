package models

type Message struct {
	Message string `json:"message" validate:"required"`
}

type Error struct {
	Error string `json:"error" validate:"required"`
}

package models

type Message struct {
	Message string `json:"message"`
}

type Error struct {
	Error string `json:"error"`
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta Pagination  `json:"meta"`
}

type Pagination struct {
	Total       int  `json:"total"`
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalPages  int  `json:"totalPages"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
}

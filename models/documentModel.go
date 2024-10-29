package models

type Class struct {
	ID       int       `json:"id" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Count    int       `json:"count" validate:"required"`
	Subjects []Subject `json:"subjects" validate:"required"`
}

type Subject struct {
	ID    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count" validate:"required"`
}

type Document struct {
	ID        int    `json:"id" validate:"required"`
	ClassID   int    `json:"classId" validate:"required"`
	SubjectID int    `json:"subjectId" validate:"required"`
	Category  string `json:"category" validate:"required"`
	Title     string `json:"title" validate:"required"`
	FileURL   string `json:"fileUrl" validate:"required"`
	Views     int    `json:"views" validate:"required"`
	Downloads int    `json:"downloads" validate:"required"`
	Author    string `json:"author" validate:"required"`
}

type CreateDocument struct {
	SubjectID int    `json:"subjectId" validate:"required" `
	Title     string `json:"title" validate:"required"`
	Author    string `json:"author" validate:"required"`
}

package models

import "time"

type Class struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Subject struct {
	ID        int       `json:"id"`
	ClassID   int       `json:"classId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Document struct {
	ID           int       `json:"id"`
	SubjectID    int       `json:"subjectId"`
	Title        string    `json:"title"`
	FileURL      string    `json:"fileUrl"`
	DocumentType string    `json:"documentType"`
	CreatedAt    time.Time `json:"createdAt"`
}

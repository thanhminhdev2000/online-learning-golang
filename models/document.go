package models

import "mime/multipart"

type Class struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Subject struct {
	ID      int    `json:"id"`
	ClassID int    `json:"classId"`
	Name    string `json:"name"`
}

type Document struct {
	ID           int    `json:"id" validate:"required"`
	SubjectId    int    `json:"subjectId" validate:"required"`
	Title        string `json:"title" validate:"required"`
	FileUrl      string `json:"fileUrl" validate:"required"`
	DocumentType string `json:"documentType" validate:"required"`
	Views        int    `json:"views" validate:"required"`
	Downloads    int    `json:"downloads" validate:"required"`
	Author       string `json:"author" validate:"required"`
}

type ClassWithSubjects struct {
	ClassId   int         `json:"classId" validate:"required"`
	ClassName string      `json:"className" validate:"required"`
	Count     int         `json:"count" validate:"required"`
	Subjects  []SubjectId `json:"subjects" validate:"required"`
}

type SubjectId struct {
	SubjectId   int    `json:"subjectId" validate:"required"`
	SubjectName string `json:"subjectName" validate:"required"`
	Count       int    `json:"count" validate:"required"`
}

type UploadRequest struct {
	SubjectId int                   `form:"subjectId" validate:"required" json:"subjectId"`
	Title     string                `form:"title" validate:"required" json:"title"`
	Author    string                `json:"author"`
	File      *multipart.FileHeader `form:"file" swaggerignore:"true" validate:"required" `
}

type DocumentsResponse struct {
	ID           int    `json:"id" validate:"required"`
	Category     string `json:"category" validate:"required"`
	Title        string `json:"title" validate:"required"`
	FileUrl      string `json:"fileUrl" validate:"required"`
	DocumentType string `json:"documentType" validate:"required"`
	Views        int    `json:"views" validate:"required"`
	Downloads    int    `json:"downloads" validate:"required"`
	Author       string `json:"author" validate:"required"`
}

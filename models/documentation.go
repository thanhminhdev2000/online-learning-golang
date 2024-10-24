package models

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
	ID           int    `json:"id"`
	SubjectId    int    `json:"subjectId"`
	Title        string `json:"title"`
	FileUrl      string `json:"fileUrl"`
	DocumentType string `json:"documentType"`
}

type ClassWithSubjects struct {
	ClassId   int         `json:"classId"`
	ClassName string      `json:"className"`
	Count     int         `json:"count"`
	Subjects  []SubjectId `json:"subjects"`
}

type SubjectId struct {
	SubjectId   int    `json:"subjectId"`
	SubjectName string `json:"subjectName"`
	Count       int    `json:"count"`
}

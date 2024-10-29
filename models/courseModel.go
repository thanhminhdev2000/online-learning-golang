package models

type Course struct {
	ID           int     `json:"id"`
	SubjectID    int     `json:"subjectId" binding:"required" validate:"required"`
	Title        string  `json:"title" binding:"required" validate:"required"`
	ThumbnailURL string  `json:"thumbnailUrl" binding:"required" validate:"required"`
	Description  string  `json:"description" binding:"required" validate:"required"`
	Price        float64 `json:"price" binding:"required" validate:"required"`
	Instructor   string  `json:"instructor" binding:"required" validate:"required"`
}

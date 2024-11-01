package models

type Lesson struct {
	ID       int    `json:"id" validate:"required"`
	CourseID int    `json:"courseId" validate:"required"`
	Title    string `json:"title" validate:"required"`
	VideoURL string `json:"videoUrl" validate:"required"`
	Duration int    `json:"duration" validate:"required"`
}

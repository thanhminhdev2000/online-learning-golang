package models

type Course struct {
	ID          int     `json:"id" db:"id"`
	SubjectID   int     `json:"subjectId" db:"subjectId" binding:"required"`
	Title       string  `json:"title" db:"title" binding:"required"`
	Description string  `json:"description,omitempty" db:"description"`
	Price       float64 `json:"price" db:"price" binding:"required"`
	Instructor  string  `json:"instructor,omitempty" db:"instructor"`
}

type CreateCourseRequest struct {
	SubjectID   int     `form:"subjectId" binding:"required"`
	Title       string  `form:"title" binding:"required"`
	Description string  `form:"description"`
	Price       float64 `form:"price" binding:"required"`
	Instructor  string  `form:"instructor"`
}

type UpdateCourseRequest struct {
	Title       *string  `form:"title"`
	Description *string  `form:"description"`
	Price       *float64 `form:"price"`
	Instructor  *string  `form:"instructor"`
}

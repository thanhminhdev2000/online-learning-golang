package models

type Lesson struct {
	ID       int    `json:"id" db:"id"`
	CourseID int    `json:"courseId" db:"courseId"`
	Title    string `json:"title" db:"title"`
	VideoURL string `json:"videoUrl" db:"videoUrl"`
}

type CreateLessonRequest struct {
	CourseID int    `form:"courseId" binding:"required"`
	Title    string `form:"title" binding:"required"`
}

type UpdateLessonRequest struct {
	CourseID *int    `form:"courseId"`
	Title    *string `form:"title"`
}

type VideoMetadata struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

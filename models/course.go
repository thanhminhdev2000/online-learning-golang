package models

import (
	"fmt"
)

type Course struct {
	ID           int     `json:"id" validate:"required"`
	SubjectID    int     `json:"subjectId" validate:"required"`
	Title        string  `json:"title" validate:"required"`
	ThumbnailURL string  `json:"thumbnailUrl" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
	Instructor   string  `json:"instructor" validate:"required"`
}

// Validation constants
const (
	MinTitleLength       = 3
	MaxTitleLength       = 100
	MinDescriptionLength = 10
	MaxDescriptionLength = 1000
	MinPrice             = 0
)

func (c *Course) Validate() error {
	if len(c.Title) < MinTitleLength || len(c.Title) > MaxTitleLength {
		return fmt.Errorf("title must be between %d and %d characters", MinTitleLength, MaxTitleLength)
	}

	if len(c.Description) < MinDescriptionLength || len(c.Description) > MaxDescriptionLength {
		return fmt.Errorf("description must be between %d and %d characters", MinDescriptionLength, MaxDescriptionLength)
	}

	if c.Price <= MinPrice {
		return fmt.Errorf("price must be greater than %d", MinPrice)
	}

	return nil
}

type CourseListResponse struct {
	Data   []Course `json:"data" validate:"required"`
	Paging Paging   `json:"paging" validate:"required"`
}

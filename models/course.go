package models

import (
	"fmt"
)

type Course struct {
	ID           int     `json:"id"`
	SubjectID    int     `json:"subjectId"`
	Title        string  `json:"title"`
	ThumbnailURL string  `json:"thumbnailUrl"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Instructor   string  `json:"instructor"`
}

// Validation constants
const (
	MinTitleLength       = 3
	MaxTitleLength       = 100
	MinDescriptionLength = 10
	MaxDescriptionLength = 1000
	MinPrice            = 0
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
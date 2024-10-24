package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"online-learning-golang/models"

	"github.com/gin-gonic/gin"
)

func GetClassesWithSubjects(db *sql.DB) ([]models.ClassWithSubjects, error) {
	classQuery := `
        SELECT c.id, c.name, COUNT(d.id) as documentCount
        FROM classes c
        LEFT JOIN subjects s ON c.id = s.classId
        LEFT JOIN documents d ON s.id = d.subjectId
        GROUP BY c.id;
    `
	rows, err := db.Query(classQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}
	defer rows.Close()

	var classList []models.ClassWithSubjects

	for rows.Next() {
		var class models.ClassWithSubjects
		if err := rows.Scan(&class.ClassId, &class.ClassName, &class.Count); err != nil {
			return nil, err
		}

		subjectQuery := `
            SELECT s.id, s.name, COUNT(d.id) as documentCount
            FROM subjects s
            LEFT JOIN documents d ON s.id = d.subjectId
            WHERE s.classId = ?
            GROUP BY s.id;
        `
		subjectRows, err := db.Query(subjectQuery, class.ClassId)
		if err != nil {
			return nil, fmt.Errorf("failed to get subjects for classId %d: %w", class.ClassId, err)
		}
		defer subjectRows.Close()

		var subjects []models.SubjectId
		for subjectRows.Next() {
			var subject models.SubjectId
			if err := subjectRows.Scan(&subject.SubjectId, &subject.SubjectName, &subject.Count); err != nil {
				return nil, err
			}
			subjects = append(subjects, subject)
		}

		class.Subjects = subjects
		classList = append(classList, class)
	}

	return classList, nil
}

// GetListClassesWithSubjects godoc
// @Summary List of classes with their subjects and document counts
// @Description List of classes with their subjects and document counts
// @Tags Documentation
// @Security BearerAuth
// @Success 200 {array} models.ClassWithSubjects
// @Failure 500 {object} models.Error
// @Router /documentations/ [get]
func GetListClassesWithSubjects(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		classesWithSubjects, err := GetClassesWithSubjects(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classesWithSubjects)
	}
}

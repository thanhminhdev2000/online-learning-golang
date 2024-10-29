package controllers

import (
	"database/sql"
	"net/http"
	"online-learning-golang/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCourses godoc
// @Summary Get a list of courses with pagination
// @Description Retrieve details of all available courses with pagination support
// @Tags Course
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of courses per page" default(10)
// @Success 200 {array} models.Course
// @Failure 400 {object} models.Error "Invalid pagination parameters"
// @Failure 500 {object} models.Error "Failed to retrieve courses"
// @Router /courses/ [get]
func GetCourses(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid page number"})
			return
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid limit"})
			return
		}

		offset := (page - 1) * limit

		query := "SELECT id, subjectId, title, description, price, instructor FROM courses LIMIT ? OFFSET ?"
		rows, err := db.Query(query, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve courses"})
			return
		}
		defer rows.Close()

		var courses []models.Course
		for rows.Next() {
			var course models.Course
			err := rows.Scan(&course.ID, &course.SubjectID, &course.Title, &course.Description, &course.Price, &course.Instructor)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve courses"})
				return
			}
			courses = append(courses, course)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve courses"})
			return
		}

		c.JSON(http.StatusOK, courses)
	}
}

// GetCourse godoc
// @Summary Get a course by ID
// @Description Retrieve details of a specific course by its ID
// @Tags Course
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 404 {object} models.Error "Course not found"
// @Failure 500 {object} models.Error "Failed to retrieve course"
// @Router /courses/{id} [get]
func GetCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var course models.Course

		query := "SELECT id, subjectId, title, description, price, instructor FROM courses WHERE id = ?"
		err := db.QueryRow(query, courseID).Scan(
			&course.ID, &course.SubjectID, &course.Title, &course.Description,
			&course.Price, &course.Instructor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
			} else {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve course"})
			}
			return
		}

		c.JSON(http.StatusOK, course)
	}
}

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a new course under a specific subject
// @Tags Course
// @Security BearerAuth
// @Param course body models.CreateCourseRequest true "Course data"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.Error
// @Router /courses/ [post]
func CreateCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course models.CreateCourseRequest
		if err := c.ShouldBind(&course); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request data"})
			return
		}

		query := `INSERT INTO courses (subjectId, title, description, price, instructor) VALUES (?, ?, ?, ?, ?)`
		_, err := db.Exec(query, course.SubjectID, course.Title, course.Description, course.Price, course.Instructor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create course"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Course created successfully"})
	}
}

// UpdateCourse godoc
// @Summary Update a course
// @Description Update course details like title, description, price, or instructor
// @Tags Course
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param course body models.UpdateCourseRequest true "Course data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error "Invalid request data"
// @Failure 500 {object} models.Error "Failed to update course"
// @Router /courses/{id} [put]
func UpdateCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRequest models.UpdateCourseRequest
		if err := c.ShouldBind(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request data"})
			return
		}

		courseID := c.Param("id")

		query := "UPDATE courses SET "
		args := []interface{}{}
		if updateRequest.Title != nil {
			query += "title = ?, "
			args = append(args, *updateRequest.Title)
		}
		if updateRequest.Description != nil {
			query += "description = ?, "
			args = append(args, *updateRequest.Description)
		}
		if updateRequest.Price != nil {
			query += "price = ?, "
			args = append(args, *updateRequest.Price)
		}
		if updateRequest.Instructor != nil {
			query += "instructor = ?, "
			args = append(args, *updateRequest.Instructor)
		}
		query = query[:len(query)-2] // Remove the last comma and space
		query += " WHERE id = ?"
		args = append(args, courseID)

		_, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update course"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Course updated successfully"})
	}
}

// DeleteCourse godoc
// @Summary Delete a course
// @Description Delete a course by its ID
// @Tags Course
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} models.Message
// @Failure 500 {object} models.Error "Failed to delete course"
// @Router /courses/{id} [delete]
func DeleteCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")

		query := "DELETE FROM courses WHERE id = ?"
		_, err := db.Exec(query, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete course"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Course deleted successfully"})
	}
}

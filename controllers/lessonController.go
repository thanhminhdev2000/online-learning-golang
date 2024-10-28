package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	cloudinarySetup "online-learning-golang/cloudinary"
	"online-learning-golang/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetLessons godoc
// @Summary Get a list of lessons with pagination
// @Description Retrieve a paginated list of lessons
// @Tags Lesson
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of lessons per page" default(10)
// @Success 200 {array} models.Lesson
// @Failure 400 {object} models.Error "Invalid pagination parameters"
// @Failure 500 {object} models.Error "Failed to retrieve lessons"
// @Router /lessons [get]
func GetLessons(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve and validate page and limit query parameters
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

		query := "SELECT id, courseId, title, videoUrl FROM lessons LIMIT ? OFFSET ?"
		rows, err := db.Query(query, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve lessons"})
			return
		}
		defer rows.Close()

		var lessons []models.Lesson
		for rows.Next() {
			var lesson models.Lesson
			err := rows.Scan(&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.VideoURL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve lessons"})
				return
			}
			lessons = append(lessons, lesson)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve lessons"})
			return
		}

		c.JSON(http.StatusOK, lessons)
	}
}

// GetLesson godoc
// @Summary Get a lesson by ID
// @Description Retrieve details of a specific lesson by its ID
// @Tags Lesson
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Success 200 {object} models.Lesson
// @Failure 404 {object} models.Error "Lesson not found"
// @Failure 500 {object} models.Error "Failed to retrieve lesson"
// @Router /lessons/{lessonId} [get]
func GetLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonID := c.Param("lessonId")
		var lesson models.Lesson

		query := "SELECT id, courseId, title, videoUrl FROM lessons WHERE id = ?"
		err := db.QueryRow(query, lessonID).Scan(
			&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.VideoURL,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Lesson not found"})
			} else {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve lesson"})
			}
			return
		}

		c.JSON(http.StatusOK, lesson)
	}
}

// CreateLesson godoc
// @Summary Upload a lesson video
// @Description Upload a video for a lesson and save lesson details in the database
// @Tags Lesson
// @Security BearerAuth
// @Param courseId formData int true "Course ID"
// @Param title formData string true "Lesson title"
// @Param file formData file true "Video file to upload"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error "Invalid request data"
// @Failure 500 {object} models.Error "Failed to create lesson"
// @Router /lessons/ [post]
func CreateLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lesson models.CreateLessonRequest
		if err := c.ShouldBind(&lesson); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request data"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No file provided"})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid file upload"})
			return
		}
		defer fileContent.Close()

		fileContent.Seek(0, 0)

		cld, err := cloudinarySetup.SetupCloudinary()
		if err != nil {
			log.Fatalf("Error setting up Cloudinary: %v", err)
		}

		url, err := cloudinarySetup.UploadVideo(cld, fileContent)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload video"})
			return
		}

		query := "INSERT INTO lessons (courseId, title, videoUrl) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(query, lesson.CourseID, lesson.Title, url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create lesson"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Create lesson successful"})
	}
}

// UpdateLesson godoc
// @Summary Update a lesson's details
// @Description Update a lesson's details and optionally upload a new video
// @Tags Lesson
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Param courseId formData int false "Course ID"
// @Param title formData string false "Lesson title"
// @Param file formData file false "New video file to upload"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error "Invalid request data"
// @Failure 500 {object} models.Error "Failed to update lesson"
// @Router /lessons/{lessonId} [put]
func UpdateLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonID := c.Param("lessonId")
		var updateRequest models.UpdateLessonRequest
		if err := c.ShouldBind(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request data"})
			return
		}

		var videoURL string
		file, err := c.FormFile("file")
		if err == nil {
			fileContent, err := file.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid file upload"})
				return
			}
			defer fileContent.Close()

			cld, err := cloudinarySetup.SetupCloudinary()
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
				return
			}

			videoURL, err = cloudinarySetup.UploadVideo(cld, fileContent)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload new video"})
				return
			}
		}

		query := "UPDATE lessons SET "
		args := []interface{}{}

		if updateRequest.CourseID != nil {
			query += "courseId = ?, "
			args = append(args, *updateRequest.CourseID)
		}
		if updateRequest.Title != nil {
			query += "title = ?, "
			args = append(args, *updateRequest.Title)
		}
		if videoURL != "" {
			query += "videoUrl = ?, "
			args = append(args, videoURL)
		}
		query = query[:len(query)-2]
		query += " WHERE id = ?"
		args = append(args, lessonID)

		_, err = db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update lesson"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Lesson updated successfully"})
	}
}

// DeleteLesson godoc
// @Summary Delete a lesson
// @Description Delete a lesson by its ID
// @Tags Lesson
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Success 200 {object} models.Message
// @Failure 404 {object} models.Error "Lesson not found"
// @Failure 500 {object} models.Error "Failed to delete lesson"
// @Router /lessons/{lessonId} [delete]
func DeleteLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonID := c.Param("lessonId")

		query := "DELETE FROM lessons WHERE id = ?"
		result, err := db.Exec(query, lessonID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete lesson"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to confirm deletion"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, models.Error{Error: "Lesson not found"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Lesson deleted successfully"})
	}
}

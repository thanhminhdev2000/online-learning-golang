package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateLesson godoc
// @Summary Create a new lesson
// @Description Create a new lesson with video upload
// @Tags Lesson
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param courseId formData int true "Course ID"
// @Param title formData string true "Lesson Title"
// @Param video formData file true "Video File"
// @Success 200 {object} models.Lesson
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /lessons/ [post]
func CreateLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		courseIdStr := c.PostForm("courseId")
		title := c.PostForm("title")

		courseId, err := strconv.Atoi(courseIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid courseId"})
			return
		}

		file, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "No file provided or invalid file",
			})
			return
		}

		if file.Size > 200*1024*1024 {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "File size exceeds maximum limit of 200MB",
			})
			return
		}

		if !utils.IsValidVideoType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid file type. Only MP4, AVI, and MOV are allowed",
			})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to process uploaded file",
			})
			return
		}
		defer fileContent.Close()

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		cld, err := utils.SetupCloudinary()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to initialize upload service",
			})
			return
		}

		videoUrl, duration, err := utils.UploadVideo(cld, fileContent)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to upload video",
			})
			return
		}

		var lesson models.Lesson
		lesson.CourseID = courseId
		lesson.Title = title
		lesson.VideoURL = videoUrl
		lesson.Duration = duration

		result, err := tx.Exec("INSERT INTO lessons (courseId, title, videoUrl, duration) VALUES (?, ?, ?, ?)", lesson.CourseID, lesson.Title, lesson.VideoURL, 100)
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to insert lesson into database",
			})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit transaction",
			})
			return
		}

		lessonID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve lesson ID"})
			return
		}
		lesson.ID = int(lessonID)

		c.JSON(http.StatusOK, lesson)
	}
}

// UpdateLesson godoc
// @Summary Update an existing lesson
// @Description Update the title and video of a lesson
// @Tags Lesson
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Param title formData string false "Lesson Title"
// @Param video formData file false "Video File"
// @Success 200 {object} models.Lesson
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /lessons/{id} [put]
func UpdateLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonIdStr := c.Param("id")
		lessonId, err := strconv.Atoi(lessonIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid lesson ID"})
			return
		}

		var existingLesson models.Lesson
		query := `SELECT videoUrl WHERE id = ?`
		err = db.QueryRow(query, lessonId).Scan(&existingLesson.VideoURL)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve course"})
				return
			}
		}

		var lesson models.Lesson
		if title := c.PostForm("title"); title != "" {
			lesson.Title = title
		}

		if file, err := c.FormFile("video"); err == nil {
			if file.Size > 200*1024*1024 {
				c.JSON(http.StatusBadRequest, models.Error{
					Error: "File size exceeds maximum limit of 200MB",
				})
				return
			}

			if !utils.IsValidVideoType(file.Header.Get("Content-Type")) {
				c.JSON(http.StatusBadRequest, models.Error{
					Error: "Invalid file type. Only MP4, AVI, and MOV are allowed",
				})
				return
			}

			fileContent, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{
					Error: "Failed to process uploaded file",
				})
				return
			}
			defer fileContent.Close()

			cld, err := utils.SetupCloudinary()
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{
					Error: "Failed to initialize upload service",
				})
				return
			}

			if err := utils.DeleteVideo(cld, existingLesson.VideoURL); err != nil {
				log.Printf("Failed to delete old thumbnail: %v", err)
			}

			videoUrl, duration, err := utils.UploadVideo(cld, fileContent)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{
					Error: "Failed to upload video",
				})
				return
			}
			lesson.VideoURL = videoUrl
			lesson.Duration = duration
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec("UPDATE lessons SET title = ?, videoUrl = ?, duration = ? WHERE id = ?", lesson.Title, lesson.VideoURL, lesson.Duration, lessonId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to update lesson in database",
			})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit transaction",
			})
			return
		}

		c.JSON(http.StatusOK, lesson)
	}
}

// DeleteLesson godoc
// @Summary Delete an existing lesson
// @Description Delete a lesson by ID
// @Tags Lesson
// @Security BearerAuth
// @Param id path int true "Lesson ID"
// @Success 200 {object} models.Message
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /lessons/{id} [delete]
func DeleteLesson(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonIdStr := c.Param("id")
		lessonId, err := strconv.Atoi(lessonIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid lesson ID"})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to begin transaction",
			})
			return
		}
		defer tx.Rollback()

		var existingLesson models.Lesson
		query := `SELECT videoUrl FROM lessons WHERE id = ?`
		err = db.QueryRow(query, lessonId).Scan(&existingLesson.VideoURL)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Lesson not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve lesson"})
				return
			}
		}

		cld, err := utils.SetupCloudinary()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to initialize upload service",
			})
			return
		}

		if err := utils.DeleteVideo(cld, existingLesson.VideoURL); err != nil {
			log.Printf("Failed to delete video from cloud: %v", err)
		}

		_, err = tx.Exec("DELETE FROM lessons WHERE id = ?", lessonId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to delete lesson from database",
			})
			return
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to commit transaction",
			})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Lesson deleted successfully"})
	}
}

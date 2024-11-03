package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateCourse handles the creation of a new course
// @Summary      Create a new course
// @Description  Create a new course with the provided details
// @Tags         Course
// @Accept       multipart/form-data
// @Produce      json
// @Param        subjectId   formData  int     true   "Subject ID"
// @Param        title       formData  string  true   "Course Title"
// @Param        description formData  string  true   "Course Description"
// @Param        price       formData  number  true   "Course Price"
// @Param        instructor  formData  string  true   "Instructor Name"
// @Param        thumbnail   formData  file    true   "Thumbnail Image"
// @Success      200         {object}  models.Course
// @Failure      400         {object}  models.Error
// @Failure      500         {object}  models.Error
// @Router       /courses/ [post]
func CreateCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var course models.Course

		// Parse form data
		subjectIdStr := c.PostForm("subjectId")
		title := c.PostForm("title")
		description := c.PostForm("description")
		priceStr := c.PostForm("price")
		instructor := c.PostForm("instructor")

		// Validate required fields
		if subjectIdStr == "" || title == "" || description == "" || priceStr == "" || instructor == "" {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		// Convert subjectId and price to appropriate types
		subjectId, err := strconv.Atoi(subjectIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid subjectId"})
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid price"})
			return
		}

		if price <= 0 {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Price must be greater than 0"})
			return
		}

		if len(title) < 3 || len(title) > 100 {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Title must be between 3 and 100 characters"})
			return
		}

		if len(description) < 10 || len(description) > 1000 {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Description must be between 10 and 1000 characters"})
			return
		}

		file, err := c.FormFile("thumbnail")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No file provided"})
			return
		}

		if file.Header.Get("Content-Type") != "image/jpeg" &&
			file.Header.Get("Content-Type") != "image/png" &&
			file.Header.Get("Content-Type") != "image/gif" {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Only JPEG, PNG and GIF images are allowed"})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to open thumbnail file"})
			return
		}
		defer fileContent.Close()

		cld, err := utils.SetupCloudinary()
		if err != nil {
			log.Printf("Error setting up Cloudinary: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
			return
		}

		thumbnailUrl, err := utils.UploadImage(cld, fileContent)
		if err != nil {
			log.Printf("Error uploading thumbnail: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload thumbnail"})
			return
		}

		course = models.Course{
			SubjectID:    subjectId,
			Title:        title,
			ThumbnailURL: thumbnailUrl,
			Description:  description,
			Price:        price,
			Instructor:   instructor,
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to start transaction"})
			return
		}

		query := `INSERT INTO courses (subjectId, title, thumbnailUrl, description, price, instructor)
				  VALUES (?, ?, ?, ?, ?, ?)`

		result, err := tx.Exec(query, course.SubjectID, course.Title, course.ThumbnailURL,
			course.Description, course.Price, course.Instructor)
		if err != nil {
			tx.Rollback()
			if deleteErr := utils.DeleteImage(cld, course.ThumbnailURL); deleteErr != nil {
				log.Printf("Failed to delete image from Cloudinary: %v", deleteErr)
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create course"})
			return
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			if deleteErr := utils.DeleteImage(cld, course.ThumbnailURL); deleteErr != nil {
				log.Printf("Failed to delete image from Cloudinary: %v", deleteErr)
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to commit transaction"})
			return
		}

		// Get the last inserted ID
		courseID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve course ID"})
			return
		}
		course.ID = int(courseID)

		c.JSON(http.StatusOK, course)
	}
}

// UpdateCourse handles updating an existing course
// @Summary      Update an existing course
// @Description  Update course details by ID
// @Tags         Course
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          path      int     true   "Course ID"
// @Param        subjectId   formData  int     false  "Subject ID"
// @Param        title       formData  string  false  "Course Title"
// @Param        description formData  string  false  "Course Description"
// @Param        price       formData  number  false  "Course Price"
// @Param        instructor  formData  string  false  "Instructor Name"
// @Param        thumbnail   formData  file    false  "Thumbnail Image"
// @Success      200         {object}  models.Course
// @Failure      400         {object}  models.Error
// @Failure      404         {object}  models.Error
// @Failure      500         {object}  models.Error
// @Router       /courses/{id} [put]
func UpdateCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid course ID"})
			return
		}

		var existingCourse models.Course
		query := `SELECT id, subjectId, title, thumbnailUrl, description, price, instructor FROM courses WHERE id = ?`
		err = db.QueryRow(query, id).Scan(&existingCourse.ID, &existingCourse.SubjectID, &existingCourse.Title, &existingCourse.ThumbnailURL, &existingCourse.Description, &existingCourse.Price, &existingCourse.Instructor)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve course"})
				return
			}
		}

		subjectIdStr := c.PostForm("subjectId")
		title := c.PostForm("title")
		description := c.PostForm("description")
		priceStr := c.PostForm("price")
		instructor := c.PostForm("instructor")

		// Update fields only if they are provided
		if subjectIdStr != "" {
			subjectId, err := strconv.Atoi(subjectIdStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid subjectId"})
				return
			}
			existingCourse.SubjectID = subjectId
		}

		if title != "" {
			if len(title) < 3 || len(title) > 100 {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Title must be between 3 and 100 characters"})
				return
			}
			existingCourse.Title = title
		}

		if description != "" {
			if len(description) < 10 || len(description) > 1000 {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Description must be between 10 and 1000 characters"})
				return
			}
			existingCourse.Description = description
		}

		if priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid price"})
				return
			}
			if price <= 0 {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Price must be greater than 0"})
				return
			}
			existingCourse.Price = price
		}

		if instructor != "" {
			existingCourse.Instructor = instructor
		}

		// Handle optional thumbnail update
		file, err := c.FormFile("thumbnail")
		if err == nil {
			if file.Header.Get("Content-Type") != "image/jpeg" &&
				file.Header.Get("Content-Type") != "image/png" &&
				file.Header.Get("Content-Type") != "image/gif" {
				c.JSON(http.StatusBadRequest, models.Error{Error: "Only JPEG, PNG and GIF images are allowed"})
				return
			}

			fileContent, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to open thumbnail file"})
				return
			}
			defer fileContent.Close()

			cld, err := utils.SetupCloudinary()
			if err != nil {
				log.Printf("Error setting up Cloudinary: %v", err)
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
				return
			}

			thumbnailURL, err := utils.UploadImage(cld, fileContent)
			if err != nil {
				log.Printf("Error uploading thumbnail: %v", err)
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload thumbnail"})
				return
			}

			if err := utils.DeleteImage(cld, existingCourse.ThumbnailURL); err != nil {
				log.Printf("Failed to delete old thumbnail: %v", err)
			}

			existingCourse.ThumbnailURL = thumbnailURL
		}

		// Update the course in the database
		query = `UPDATE courses SET subjectId = ?, title = ?, thumbnailUrl = ?, description = ?, price = ?, instructor = ? WHERE id = ?`
		_, err = db.Exec(query, existingCourse.SubjectID, existingCourse.Title, existingCourse.ThumbnailURL,
			existingCourse.Description, existingCourse.Price, existingCourse.Instructor, existingCourse.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update course"})
			return
		}

		c.JSON(http.StatusOK, existingCourse)
	}
}

// DeleteCourse handles deleting a course
// @Summary      Delete a course
// @Description  Delete a course by ID
// @Tags         Course
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  models.Message
// @Failure      404  {object}  models.Error
// @Failure      500  {object}  models.Error
// @Router       /courses/{id} [delete]
func DeleteCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the course ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid course ID"})
			return
		}

		// Get the existing course to retrieve thumbnail URL
		var thumbnailURL string
		err = db.QueryRow("SELECT thumbnailUrl FROM courses WHERE id = ?", id).Scan(&thumbnailURL)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve course"})
			return
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to start transaction"})
			return
		}

		// Delete the course from the database
		query := `DELETE FROM courses WHERE id = ?`
		result, err := tx.Exec(query, id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete course"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to determine affected rows"})
			return
		}
		if rowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
			return
		}

		// Commit the transaction
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to commit transaction"})
			return
		}

		// Delete thumbnail from Cloudinary after successful database deletion
		if thumbnailURL != "" {
			cld, err := utils.SetupCloudinary()
			if err != nil {
				log.Printf("Error setting up Cloudinary while deleting thumbnail: %v", err)
				// Continue even if Cloudinary setup fails
			} else {
				if err := utils.DeleteImage(cld, thumbnailURL); err != nil {
					log.Printf("Failed to delete thumbnail from Cloudinary: %v", err)
				}
			}
		}

		c.JSON(http.StatusOK, models.Message{Message: "Course deleted successfully"})
	}
}

// GetCourse handles fetching a single course by ID
// @Summary      Get a course by ID
// @Description  Retrieve a single course using its ID
// @Tags         Course
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  models.Course
// @Failure      400  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Failure      500  {object}  models.Error
// @Router       /courses/{id} [get]
func GetCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: fmt.Sprintf("Invalid course ID format: %s", idStr),
			})
			return
		}

		var course models.Course

		query := `
			SELECT 
				c.id, 
				cls.id,
				s.id, 
				CONCAT(s.name, ' - ', cls.name) AS category,
				c.title, 
				c.thumbnailUrl, 
				c.description, 
				c.price, 
				c.instructor
			FROM courses c
			LEFT JOIN subjects s ON c.subjectId = s.id
			LEFT JOIN classes cls ON s.classId = cls.id
			WHERE c.id = ?`

		err = db.QueryRow(query, id).Scan(
			&course.ID,
			&course.ClassID,
			&course.SubjectID,
			&course.Category,
			&course.Title,
			&course.ThumbnailURL,
			&course.Description,
			&course.Price,
			&course.Instructor,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{
					Error: fmt.Sprintf("Course with ID %d not found", id),
				})
				return
			}
			log.Printf("Error retrieving course %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to retrieve course. Please try again later.",
			})
			return
		}

		userID, _ := c.Get("userId")

		var isActive bool
		err = db.QueryRow("SELECT COUNT(*) > 0 FROM user_courses WHERE userId = ? AND courseId = ?", userID, id).Scan(&isActive)
		if err != nil {
			log.Printf("Error checking user course activation: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to check course activation status",
			})
			return
		}

		currentUserRole, _ := c.Get("role")
		if currentUserRole == "admin" {
			isActive = true
		}

		var lessons []models.Lesson

		lessonsQuery := `SELECT id, courseId, title, videoUrl, duration, position FROM lessons WHERE courseId = ? ORDER BY position ASC`
		rows, err := db.Query(lessonsQuery, id)
		if err != nil {
			log.Printf("Error retrieving lessons for course %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to retrieve lessons. Please try again later.",
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var lesson models.Lesson
			if err := rows.Scan(&lesson.ID, &lesson.CourseID, &lesson.Title, &lesson.VideoURL, &lesson.Duration, &lesson.Position); err != nil {
				log.Printf("Error scanning lesson: %v", err)
				c.JSON(http.StatusInternalServerError, models.Error{
					Error: "Failed to process lesson data",
				})
				return
			}

			if !isActive {
				lesson.VideoURL = ""
			}
			lessons = append(lessons, lesson)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating lessons: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to process lessons",
			})
			return
		}

		course.IsActive = isActive
		course.Lessons = lessons

		c.JSON(http.StatusOK, course)
	}
}

// GetCourses handles fetching all courses
// @Summary      Get all courses
// @Description  Retrieve a list of all courses with optional filtering and pagination
// @Tags         Course
// @Produce      json
// @Param        page     query    int     false  "Page number (default: 1)"
// @Param        limit    query    int     false  "Items per page (default: 10)"
// @Param        subject  query    int     false  "Filter by subject ID"
// @Param        search   query    string  false  "Search in title and description"
// @Param        sort     query    string  false  "Sort field (title, price) (default: id)"
// @Param        order    query    string  false  "Sort order (asc, desc) (default: asc)"
// @Success      200      {object} models.CourseListResponse
// @Failure      400      {object} models.Error
// @Failure      500      {object} models.Error
// @Router       /courses/ [get]
func GetCourses(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		subjectID := c.Query("subject")
		search := c.Query("search")
		sortField := c.DefaultQuery("sort", "id")
		sortOrder := strings.ToLower(c.DefaultQuery("order", "asc"))

		// Validate pagination parameters
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		// Validate sort parameters
		validSortFields := map[string]bool{"id": true, "title": true, "price": true}
		if !validSortFields[sortField] {
			sortField = "id"
		}
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc"
		}

		// Build query with JOIN and WHERE clauses
		query := `
			SELECT 
				c.id, 
				cls.id,
				s.id, 
				CONCAT(s.name, ' - ', cls.name) AS category,
				c.title, 
				c.thumbnailUrl, 
				c.description, 
				c.price, 
				c.instructor
			FROM courses c
			LEFT JOIN subjects s ON c.subjectId = s.id
			LEFT JOIN classes cls ON s.classId = cls.id
			WHERE 1=1`

		countQuery := "SELECT COUNT(*) FROM courses c WHERE 1=1"
		params := []interface{}{}

		// Add filters
		if subjectID != "" {
			query += " AND c.subjectId = ?"
			countQuery += " AND c.subjectId = ?"
			params = append(params, subjectID)
		}

		if search != "" {
			query += " AND (c.title LIKE ? OR c.description LIKE ?)"
			countQuery += " AND (c.title LIKE ? OR c.description LIKE ?)"
			searchParam := "%" + search + "%"
			params = append(params, searchParam, searchParam)
		}

		query += fmt.Sprintf(" ORDER BY c.%s %s", sortField, sortOrder)

		query += " LIMIT ? OFFSET ?"
		params = append(params, limit, offset)

		var total int
		err := db.QueryRow(countQuery, params[:len(params)-2]...).Scan(&total)
		if err != nil {
			log.Printf("Error counting courses: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to retrieve course count",
			})
			return
		}

		// Execute main query
		rows, err := db.Query(query, params...)
		if err != nil {
			log.Printf("Error retrieving courses: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to retrieve courses",
			})
			return
		}
		defer rows.Close()

		courses := make([]models.Course, 0)
		for rows.Next() {
			var course models.Course
			err := rows.Scan(
				&course.ID,
				&course.ClassID,
				&course.SubjectID,
				&course.Category,
				&course.Title,
				&course.ThumbnailURL,
				&course.Description,
				&course.Price,
				&course.Instructor,
			)
			if err != nil {
				log.Printf("Error scanning course: %v", err)
				c.JSON(http.StatusInternalServerError, models.Error{
					Error: "Failed to process course data",
				})
				return
			}
			courses = append(courses, course)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error iterating courses: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to process courses",
			})
			return
		}

		response := models.CourseListResponse{
			Data: courses,
			Paging: models.Paging{
				Page:  page,
				Limit: limit,
				Total: total,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// ActivateCourseForUser handles activating a course for a user by admin
// @Summary      Activate a course for a user
// @Description  Admin activates a course for a specific user using email
// @Tags         Course
// @Accept       json
// @Produce      json
// @Param        email      formData      string  true   "User Email"
// @Param        courseId   formData      int     true   "Course ID"
// @Success      200        {object}  models.Message
// @Failure      400        {object}  models.Error
// @Failure      404        {object}  models.Error
// @Failure      500        {object}  models.Error
// @Router       /courses/activate [post]
func ActivateCourseForUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var activationRequest struct {
			Email    string `json:"email" binding:"required"`
			CourseID int    `json:"courseId" binding:"required"`
		}

		if err := c.ShouldBindJSON(&activationRequest); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var userID int
		err := db.QueryRow("SELECT id FROM users WHERE email = ?", activationRequest.Email).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "User with the given email not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to check user existence"})
			return
		}

		_, err = db.Exec("INSERT INTO user_courses (userId, courseId) VALUES (?, ?)", userID, activationRequest.CourseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to activate course for user"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Course activated successfully for user"})
	}
}

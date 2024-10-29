package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	cloudinaryUtils "online-learning-golang/cloudinary"
	"online-learning-golang/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateCourse handles the creation of a new course
// @Summary      Create a new course
// @Description  Create a new course with the provided details
// @Tags         courses
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
// @Router       /courses [post]
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
			c.JSON(http.StatusBadRequest, models.Error{Error: "All fields are required"})
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

		// Additional validations
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
		// Get the thumbnail file
		file, err := c.FormFile("thumbnail")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No file provided"})
			return
		}

		// Validate file type
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

		// Setup Cloudinary
		cld, err := cloudinaryUtils.SetupCloudinary()
		if err != nil {
			log.Printf("Error setting up Cloudinary: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
			return
		}

		// Upload the file to Cloudinary
		thumbnailUrl, err := cloudinaryUtils.UploadImage(cld, fileContent)
		if err != nil {
			log.Printf("Error uploading thumbnail: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload thumbnail"})
			return
		}

		// Create the course struct
		course = models.Course{
			SubjectID:    subjectId,
			Title:        title,
			ThumbnailURL: thumbnailUrl,
			Description:  description,
			Price:        price,
			Instructor:   instructor,
		}

		// Insert into database with transaction
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
			if deleteErr := cloudinaryUtils.DeleteImage(cld, course.ThumbnailURL); deleteErr != nil {
				log.Printf("Failed to delete image from Cloudinary: %v", deleteErr)
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create course"})
			return
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			if deleteErr := cloudinaryUtils.DeleteImage(cld, course.ThumbnailURL); deleteErr != nil {
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
// @Tags         courses
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          path      int     true   "Course ID"
// @Param        subjectId   formData  int     true   "Subject ID"
// @Param        title       formData  string  true   "Course Title"
// @Param        description formData  string  true   "Course Description"
// @Param        price       formData  number  true   "Course Price"
// @Param        instructor  formData  string  true   "Instructor Name"
// @Param        thumbnail   formData  file    true   "Thumbnail Image"
// @Success      200         {object}  models.Course
// @Failure      400         {object}  models.Error
// @Failure      404         {object}  models.Error
// @Failure      500         {object}  models.Error
// @Router       /courses/{id} [put]
func UpdateCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the course ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid course ID"})
			return
		}

		// Get the existing course from the database
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

		// Parse form data
		subjectIdStr := c.PostForm("subjectId")
		title := c.PostForm("title")
		description := c.PostForm("description")
		priceStr := c.PostForm("price")
		instructor := c.PostForm("instructor")

		// Validate required fields
		if subjectIdStr == "" || title == "" || description == "" || priceStr == "" || instructor == "" {
			c.JSON(http.StatusBadRequest, models.Error{Error: "All fields are required"})
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

		// Additional validations
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

		// Update basic course information
		existingCourse.SubjectID = subjectId
		existingCourse.Title = title
		existingCourse.Description = description
		existingCourse.Price = price
		existingCourse.Instructor = instructor

		// Handle optional thumbnail update
		file, err := c.FormFile("thumbnail")
		if err == nil {
			// Validate file type
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

			// Setup Cloudinary
			cld, err := cloudinaryUtils.SetupCloudinary()
			if err != nil {
				log.Printf("Error setting up Cloudinary: %v", err)
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
				return
			}

			// Upload new thumbnail
			thumbnailURL, err := cloudinaryUtils.UploadImage(cld, fileContent)
			if err != nil {
				log.Printf("Error uploading thumbnail: %v", err)
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload thumbnail"})
				return
			}

			// Delete old thumbnail
			if err := cloudinaryUtils.DeleteImage(cld, existingCourse.ThumbnailURL); err != nil {
				log.Printf("Failed to delete old thumbnail: %v", err)
				// Continue with update even if deletion fails
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
// @Tags         courses
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
			cld, err := cloudinaryUtils.SetupCloudinary()
			if err != nil {
				log.Printf("Error setting up Cloudinary while deleting thumbnail: %v", err)
				// Continue even if Cloudinary setup fails
			} else {
				if err := cloudinaryUtils.DeleteImage(cld, thumbnailURL); err != nil {
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
// @Tags         courses
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  models.Course
// @Failure      400  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Failure      500  {object}  models.Error
// @Router       /courses/{id} [get]
func GetCourse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the course ID from the URL parameter
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: fmt.Sprintf("Invalid course ID format: %s", idStr),
			})
			return
		}

		// Retrieve the course with subject information using JOIN
		var course struct {
			models.Course
			SubjectName string `json:"subjectName"`
		}

		query := `
			SELECT 
				c.id, 
				c.subjectId, 
				c.title, 
				c.thumbnailUrl, 
				c.description, 
				c.price, 
				c.instructor,
				s.name as subject_name
			FROM courses c
			LEFT JOIN subjects s ON c.subjectId = s.id
			WHERE c.id = ?`

		err = db.QueryRow(query, id).Scan(
			&course.ID,
			&course.SubjectID,
			&course.Title,
			&course.ThumbnailURL,
			&course.Description,
			&course.Price,
			&course.Instructor,
			&course.SubjectName,
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

		// Validate the retrieved data
		if course.Title == "" || course.Description == "" {
			log.Printf("Warning: Course %d has invalid data: title=%s, desc=%s",
				id, course.Title, course.Description)
		}

		c.JSON(http.StatusOK, course)
	}
}

// GetCourses handles fetching all courses
// @Summary      Get all courses
// @Description  Retrieve a list of all courses with optional filtering and pagination
// @Tags         courses
// @Produce      json
// @Param        page     query    int     false  "Page number (default: 1)"
// @Param        limit    query    int     false  "Items per page (default: 10)"
// @Param        subject  query    int     false  "Filter by subject ID"
// @Param        search   query    string  false  "Search in title and description"
// @Param        sort     query    string  false  "Sort field (title, price) (default: id)"
// @Param        order    query    string  false  "Sort order (asc, desc) (default: asc)"
// @Success      200      {object} models.PaginatedResponse
// @Failure      400      {object} models.Error
// @Failure      500      {object} models.Error
// @Router       /courses [get]
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
				c.subjectId, 
				c.title, 
				c.thumbnailUrl, 
				c.description, 
				c.price, 
				c.instructor,
				s.name as subject_name
			FROM courses c
			LEFT JOIN subjects s ON c.subjectId = s.id
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

		// Add sorting
		query += fmt.Sprintf(" ORDER BY c.%s %s", sortField, sortOrder)

		// Add pagination
		query += " LIMIT ? OFFSET ?"
		params = append(params, limit, offset)

		// Get total count
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

		type CourseWithSubject struct {
			models.Course
			SubjectName string `json:"subjectName"`
		}

		courses := make([]CourseWithSubject, 0)
		for rows.Next() {
			var course CourseWithSubject
			err := rows.Scan(
				&course.ID,
				&course.SubjectID,
				&course.Title,
				&course.ThumbnailURL,
				&course.Description,
				&course.Price,
				&course.Instructor,
				&course.SubjectName,
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

		// Calculate pagination metadata
		totalPages := (total + limit - 1) / limit
		hasNext := page < totalPages
		hasPrev := page > 1

		response := gin.H{
			"data": courses,
			"meta": gin.H{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"totalPages":  totalPages,
				"hasNext":     hasNext,
				"hasPrevious": hasPrev,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

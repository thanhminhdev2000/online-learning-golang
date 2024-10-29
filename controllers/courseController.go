package controllers

import (
	"database/sql"
	"log"
	"net/http"
	cloudinarySetup "online-learning-golang/cloudinary"
	"online-learning-golang/models"
	"strconv"

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

		// Get the thumbnail file
		file, err := c.FormFile("thumbnail")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No file provided"})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to open thumbnail file"})
			return
		}
		defer fileContent.Close()

		// Setup Cloudinary
		cld, err := cloudinarySetup.SetupCloudinary()
		if err != nil {
			log.Printf("Error setting up Cloudinary: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
			return
		}

		// Upload the file to Cloudinary
		thumbnailUrl, err := cloudinarySetup.UploadImage(cld, fileContent)
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

		// Insert the course into the database
		query := `INSERT INTO courses (subjectId, title, thumbnailUrl, description, price, instructor)
              VALUES (?, ?, ?, ?, ?, ?)`

		result, err := db.Exec(query, course.SubjectID, course.Title, course.ThumbnailURL, course.Description, course.Price, course.Instructor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to create course"})
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

		existingCourse.SubjectID = subjectId
		existingCourse.Title = title
		existingCourse.Description = description
		existingCourse.Price = price
		existingCourse.Instructor = instructor

		// Handle thumbnail file
		file, err := c.FormFile("thumbnail")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No file provided"})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to open thumbnail file"})
			return
		}
		defer fileContent.Close()

		// Setup Cloudinary
		cld, err := cloudinarySetup.SetupCloudinary()
		if err != nil {
			log.Printf("Error setting up Cloudinary: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to setup Cloudinary"})
			return
		}

		// Upload the file to Cloudinary
		thumbnailURL, err := cloudinarySetup.UploadImage(cld, fileContent)
		if err != nil {
			log.Printf("Error uploading thumbnail: %v", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to upload thumbnail"})
			return
		}

		existingCourse.ThumbnailURL = thumbnailURL

		// Update the course in the database
		query = `UPDATE courses SET subjectId = ?, title = ?, thumbnailUrl = ?, description = ?, price = ?, instructor = ? WHERE id = ?`
		_, err = db.Exec(query, existingCourse.SubjectID, existingCourse.Title, existingCourse.ThumbnailURL, existingCourse.Description, existingCourse.Price, existingCourse.Instructor, existingCourse.ID)
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

		// Delete the course from the database
		query := `DELETE FROM courses WHERE id = ?`
		result, err := db.Exec(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete course"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to determine affected rows"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Course deleted"})
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
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid course ID"})
			return
		}

		// Retrieve the course from the database
		var course models.Course
		query := `SELECT id, subjectId, title, thumbnailUrl, description, price, instructor FROM courses WHERE id = ?`
		err = db.QueryRow(query, id).Scan(
			&course.ID,
			&course.SubjectID,
			&course.Title,
			&course.ThumbnailURL,
			&course.Description,
			&course.Price,
			&course.Instructor,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Course not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve course"})
			return
		}

		c.JSON(http.StatusOK, course)
	}
}

// GetCourses handles fetching all courses
// @Summary      Get all courses
// @Description  Retrieve a list of all courses
// @Tags         courses
// @Produce      json
// @Success      200  {array}   models.Course
// @Failure      500  {object}  models.Error
// @Router       /courses [get]
func GetCourses(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve all courses from the database
		query := `SELECT id, subjectId, title, thumbnailUrl, description, price, instructor FROM courses`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to retrieve courses"})
			return
		}
		defer rows.Close()

		var courses []models.Course

		for rows.Next() {
			var course models.Course
			err := rows.Scan(
				&course.ID,
				&course.SubjectID,
				&course.Title,
				&course.ThumbnailURL,
				&course.Description,
				&course.Price,
				&course.Instructor,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to scan course"})
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

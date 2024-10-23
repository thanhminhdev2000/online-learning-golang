package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"online-learning-golang/cloudinary"
	"online-learning-golang/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 409 {object} models.Error
// @Router /users/ [post]
func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.CreateUserRequest
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var exists bool
		query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
		err := db.QueryRow(query, user.Email, user.Username).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to check for existing user"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, models.Error{Error: "Email or username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to hash password"})
			return
		}

		query = "INSERT INTO users (email, username, fullName, password, gender, avatar, dateOfBirth) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err = db.Exec(query, user.Email, user.Username, user.FullName, hashedPassword, user.Gender, user.Avatar, user.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to register user"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "User registered successfully"})
	}
}

// CreateUser godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 409 {object} models.Error
// @Router /users/admin [post]
func CreateUserAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.CreateUserRequest
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var exists bool
		query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
		err := db.QueryRow(query, user.Email, user.Username).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to check for existing user"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, models.Error{Error: "Email or username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to hash password"})
			return
		}

		currentUserRole, _ := c.Get("role")

		if currentUserRole != "admin" && user.Role == "admin" {
			c.JSON(http.StatusForbidden, models.Error{Error: "Only admins are allowed to create admin users."})
			return
		}
		query = "INSERT INTO users (email, username, fullName, password, gender, avatar, dateOfBirth, role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		_, err = db.Exec(query, user.Email, user.Username, user.FullName, hashedPassword, user.Gender, user.Avatar, user.DateOfBirth, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to register user"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "User registered successfully"})
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users, with optional filters for email, username, full name, date of birth, role, and pagination.
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param email query string false "Filter by email"
// @Param username query string false "Filter by username"
// @Param fullName query string false "Filter by full name"
// @Param dateOfBirth query string false "Filter by date of birth"
// @Param role query string false "Filter by role"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit number of items per page (max 100)"
// @Success 200 {object} models.UserResponse
// @Failure 500 {object} models.Error
// @Router /users/ [get]
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var queryParams models.UserQueryParams

		if err := c.ShouldBindQuery(&queryParams); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid query parameters"})
			return
		}

		query := "SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role FROM users WHERE deleted_at IS NULL"
		countQuery := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"

		var params []interface{}
		var countParams []interface{}

		if queryParams.Email != "" {
			query += " AND email LIKE ?"
			countQuery += " AND email LIKE ?"
			params = append(params, "%"+queryParams.Email+"%")
			countParams = append(countParams, "%"+queryParams.Email+"%")
		}
		if queryParams.Username != "" {
			query += " AND username LIKE ?"
			countQuery += " AND username LIKE ?"
			params = append(params, "%"+queryParams.Username+"%")
			countParams = append(countParams, "%"+queryParams.Username+"%")
		}
		if queryParams.FullName != "" {
			query += " AND fullName LIKE ?"
			countQuery += " AND fullName LIKE ?"
			params = append(params, "%"+queryParams.FullName+"%")
			countParams = append(countParams, "%"+queryParams.FullName+"%")
		}
		if queryParams.DateOfBirth != "" {
			query += " AND dateOfBirth LIKE ?"
			countQuery += " AND dateOfBirth LIKE ?"
			params = append(params, "%"+queryParams.DateOfBirth+"%")
			countParams = append(countParams, "%"+queryParams.DateOfBirth+"%")
		}
		if queryParams.Role != "" {
			query += " AND role = ?"
			countQuery += " AND role = ?"
			params = append(params, queryParams.Role)
			countParams = append(countParams, queryParams.Role)
		}

		page := queryParams.Page
		limit := queryParams.Limit

		if page == 0 {
			page = 1
		}
		if limit == 0 {
			limit = 10
		}

		var totalCount int
		err := db.QueryRow(countQuery, countParams...).Scan(&totalCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to count users"})
			return
		}

		offset := (page - 1) * limit
		query += " LIMIT ? OFFSET ?"
		params = append(params, limit, offset)

		rows, err := db.Query(query, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.UserDetail

		for rows.Next() {
			var user models.UserDetail
			if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &user.Gender, &user.Avatar, &user.DateOfBirth, &user.Role); err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to scan user"})
				return
			}

			users = append(users, user)
		}

		if len(users) == 0 {
			users = []models.UserDetail{}
		}

		response := models.UserResponse{
			Data: users,
			Paging: models.PagingInfo{
				Page:       page,
				Limit:      limit,
				TotalCount: totalCount,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetUserDetail(db *sql.DB, userId string) models.UserDetail {
	row := db.QueryRow("SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role FROM users WHERE id = ?", userId)

	var user models.UserDetail
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &user.Gender, &user.Avatar, &user.DateOfBirth, &user.Role); err != nil {
		log.Fatal("Failed to fetch user")
	}
	return user
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} models.UserDetail
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{userId} [get]
func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		user := GetUserDetail(db, userId)
		c.JSON(http.StatusOK, user)
	}
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information by user ID
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body models.UserDetail true "User data"
// @Success 200 {object} models.UpdateUserResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{userId} [put]
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var updateUser models.UserDetail

		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		query := "UPDATE users SET email = ?, username = ?, fullName = ?, avatar = ? WHERE id = ?"
		_, err := db.Exec(query, updateUser.Email, updateUser.Username, updateUser.FullName, updateUser.Avatar, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update user"})
			return
		}

		user := GetUserDetail(db, userId)
		c.JSON(http.StatusOK, models.UpdateUserResponse{Message: "User updated successfully", User: user})
	}
}

// Password godoc
// @Summary Change user password
// @Description Change the user's password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param password body models.PasswordUpdateRequest true "Password data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /users/{userId}/password [put]
func UpdateUserPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		var req struct {
			CurrentPassword string `json:"currentPassword"`
			NewPassword     string `json:"newPassword"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE id = ?", userId).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to fetch user"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.CurrentPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "Current password is incorrect"})
			return
		}

		hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to hash new password"})
			return
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", hashedNewPassword, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update password"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Password changed successfully"})
	}
}

// UpdateUserAvatar godoc
// @Summary Update user avatar
// @Description Update the avatar for a specific user
// @Tags User
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param userId path int true "User ID"
// @Param avatar formData file true "User Avatar"
// @Success 200 {object} models.UpdateUserResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{userId}/avatar [put]
func UpdateUserAvatar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "No file provided"})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to open avatar file"})
			return
		}
		defer fileContent.Close()

		cld, err := cloudinary.SetupCloudinary()
		if err != nil {
			log.Fatalf("Error setting up Cloudinary: %v", err)
		}

		avatarURL, err := cloudinary.UploadAvatar(cld, fileContent)
		if err != nil {
			log.Fatalf("Error uploading avatar: %v", err)
		}

		query := "UPDATE users SET avatar = ? WHERE id = ?"
		_, err = db.Exec(query, avatarURL, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update avatar in database"})
			return
		}

		user := GetUserDetail(db, userId)
		c.JSON(http.StatusOK, models.UpdateUserResponse{Message: "Avatar updated successfully", User: user})
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by user ID
// @Tags User
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} models.Message
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{userId} [delete]
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			return
		}

		var currentUserIDStr string

		switch v := currentUserID.(type) {
		case int:
			currentUserIDStr = fmt.Sprintf("%d", v)
		case string:
			currentUserIDStr = v
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}

		userIDToDelete := c.Param("userId")
		if userIDToDelete == currentUserIDStr {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete your own account"})
			return
		}

		query := "UPDATE users SET deleted_at = NOW() WHERE id = ?"
		result, err := db.Exec(query, userIDToDelete)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to delete user"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to get rows affected"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "User deleted successfully"})
	}
}

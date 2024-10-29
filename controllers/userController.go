package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	cloudinaryUtils "online-learning-golang/cloudinary"
	"online-learning-golang/models"
	"strconv"

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
// @Param data query models.UserQueryParams false "Filter"
// @Success 200 {object} models.UserResponse
// @Failure 500 {object} models.Error
// @Router /users/ [get]
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role")
		if currentUserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to get list user"})
			return
		}

		email := c.Query("email")
		username := c.Query("username")
		fullName := c.Query("fullName")
		dateOfBirth := c.Query("dateOfBirth")
		role := c.Query("role")
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		query := "SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role FROM users WHERE deletedAt IS NULL"
		countQuery := "SELECT COUNT(*) FROM users WHERE deletedAt IS NULL"

		var params []interface{}
		var countParams []interface{}

		if email != "" {
			query += " AND email LIKE ?"
			countQuery += " AND email LIKE ?"
			params = append(params, "%"+email+"%")
			countParams = append(countParams, "%"+email+"%")
		}
		if username != "" {
			query += " AND username LIKE ?"
			countQuery += " AND username LIKE ?"
			params = append(params, "%"+username+"%")
			countParams = append(countParams, "%"+username+"%")
		}
		if fullName != "" {
			query += " AND fullName LIKE ?"
			countQuery += " AND fullName LIKE ?"
			params = append(params, "%"+fullName+"%")
			countParams = append(countParams, "%"+fullName+"%")
		}
		if dateOfBirth != "" {
			query += " AND dateOfBirth LIKE ?"
			countQuery += " AND dateOfBirth LIKE ?"
			params = append(params, "%"+dateOfBirth+"%")
			countParams = append(countParams, "%"+dateOfBirth+"%")
		}
		if role != "" {
			query += " AND role = ?"
			countQuery += " AND role = ?"
			params = append(params, role)
			countParams = append(countParams, role)
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}

		var totalCount int
		err = db.QueryRow(countQuery, countParams...).Scan(&totalCount)
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

func GetUserDetail(db *sql.DB, userId string) (models.UserDetail, error) {
	row := db.QueryRow("SELECT id, email, username, fullName, gender, avatar, dateOfBirth, role FROM users WHERE id = ? AND deletedAt IS NULL", userId)

	var user models.UserDetail
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &user.Gender, &user.Avatar, &user.DateOfBirth, &user.Role); err != nil {
		return user, err
	}

	return user, nil
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
// @Router /users/{id} [get]
func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Param("id")
		currentUserID, _ := c.Get("userId")

		currentUserRole, _ := c.Get("role")
		if currentUserRole == "user" && currentUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to get this user"})
			return
		}

		user, err := GetUserDetail(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
			return
		}

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
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{id} [put]
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role")
		if currentUserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this user"})
			return
		}

		userID := c.Param("id")
		var updateUser models.UserDetail

		if err := c.ShouldBindJSON(&updateUser); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		query := "UPDATE users SET email = ?, username = ?, fullName = ?, gender = ?, dateOfBirth = ? WHERE id = ? AND deletedAt IS NULL"
		_, err := db.Exec(query, updateUser.Email, updateUser.Username, updateUser.FullName, updateUser.Gender, updateUser.DateOfBirth, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update user"})
			return
		}

		user, err := GetUserDetail(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
			return
		}

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
// @Router /users/{id}/password [put]
func UpdateUserPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role")
		if currentUserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this user"})
			return
		}

		userID := c.Param("id")

		var req struct {
			CurrentPassword string `json:"currentPassword"`
			NewPassword     string `json:"newPassword"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE id = ? AND deletedAt IS NULL", userID).Scan(&storedPassword)
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

		_, err = db.Exec("UPDATE users SET password = ? WHERE id = ? AND deletedAt IS NULL", hashedNewPassword, userID)
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
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{id}/avatar [put]
func UpdateUserAvatar(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserRole, _ := c.Get("role")
		if currentUserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this user"})
			return
		}

		userID := c.Param("id")

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

		cld, err := cloudinaryUtils.SetupCloudinary()
		if err != nil {
			log.Fatalf("Error setting up Cloudinary: %v", err)
		}

		avatarURL, err := cloudinaryUtils.UploadImage(cld, fileContent)
		if err != nil {
			log.Fatalf("Error uploading avatar: %v", err)
		}

		query := "UPDATE users SET avatar = ? WHERE id = ? AND deletedAt IS NULL"
		_, err = db.Exec(query, avatarURL, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update avatar in database"})
			return
		}

		user, err := GetUserDetail(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
			return
		}
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
// @Router /users/{id} [delete]
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserId, _ := c.Get("userId")

		currentUserRole, _ := c.Get("role")
		if currentUserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this user"})
			return
		}

		userIdToDelete := c.Param("id")
		if userIdToDelete == currentUserId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete your own account"})
			return
		}

		query := "UPDATE users SET deletedAt = NOW() WHERE id = ?"
		result, err := db.Exec(query, userIdToDelete)
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

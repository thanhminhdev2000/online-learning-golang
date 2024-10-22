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

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.UserDetail
// @Failure 500 {object} models.Error
// @Router /users/ [get]
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, email, username, fullName, gender, avatar, dateOfBirth FROM users WHERE deleted_at IS NULL")
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.UserDetail

		for rows.Next() {
			var user models.UserDetail
			if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &user.Gender, &user.Avatar, &user.DateOfBirth); err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to scan user"})
				return
			}

			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUserDetail(db *sql.DB, userId string) models.UserDetail {
	row := db.QueryRow("SELECT id, email, username, fullName, gender, avatar, dateOfBirth FROM users WHERE id = ?", userId)

	var user models.UserDetail
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &user.Gender, &user.Avatar, &user.DateOfBirth); err != nil {
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
		fmt.Println("fileContent", fileContent)

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
		userId := c.Param("userId")

		query := "UPDATE users SET deleted_at = NOW() WHERE id = ?"
		result, err := db.Exec(query, userId)
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

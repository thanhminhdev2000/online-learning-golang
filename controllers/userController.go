package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.SignUpRequest true "User data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 409 {object} models.Error
// @Router /users/signup [post]
func SignUp(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.SignUpRequest
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var exists bool
		query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
		err := db.QueryRow(query, newUser.Email, newUser.Username).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to check for existing user"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, models.Error{Error: "Email or username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to hash password"})
			return
		}

		query = "INSERT INTO users (email, username, fullName, password) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(query, newUser.Email, newUser.Username, newUser.FullName, hashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to register user"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "User registered successfully"})
	}
}

// Login godoc
// @Summary Log in an existing user
// @Description Log in a user using email or username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /users/login [post]
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData models.LoginRequest
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		var storedPassword string
		var user models.UserDetail
		isEmail, _ := regexp.MatchString(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, loginData.Identifier)
		var query string
		if isEmail {
			query = "SELECT id, email, username, fullName, password FROM users WHERE email = ?"
		} else {
			query = "SELECT id, email, username, fullName, password FROM users WHERE username = ?"
		}

		err := db.QueryRow(query, loginData.Identifier).
			Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.Error{Error: "Invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Database query error"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginData.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "Invalid email or password"})
			return
		}

		accessToken, err := utils.CreateAccessToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to generate access token"})
			return
		}

		refreshToken, err := utils.CreateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to generate refresh token"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   7 * 24 * 60 * 60,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		response := models.LoginResponse{
			Message:     "Login successful",
			User:        user,
			AccessToken: accessToken,
		}

		c.JSON(http.StatusOK, response)
	}
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags users
// @Produce json
// @Success 200 {object} models.AccessTokenReponse
// @Failure 401 {object} models.Error
// @Router /users/refresh [post]
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtKey = []byte(os.Getenv("JWT_KEY"))
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "No refresh token found in cookies"})
			return
		}

		token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "Invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "Expired or invalid refresh token"})
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Invalid userID"})
			return
		}

		accessToken, err := utils.CreateAccessToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to generate access token"})
			return
		}

		c.JSON(http.StatusOK, models.AccessTokenReponse{
			AccessToken: accessToken,
		})
	}
}

// Logout godoc
// @Summary Log out the user
// @Description Log out the user by clearing the refresh token
// @Tags users
// @Produce json
// @Success 200 {object} models.Message
// @Router /users/logout [post]
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Message{Message: "Logout successful"})
	}
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset link to the user's email
// @Tags users
// @Accept json
// @Produce json
// @Param email body models.ForgotPasswordRequest true "User email"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /users/forgot-password [post]
func ForgotPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Email is required"})
			return
		}

		var user models.UserDetail
		err := db.QueryRow("SELECT id, email FROM users WHERE email = ?", req.Email).Scan(&user.ID, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "Email not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to query user"})
			return
		}

		token, err := utils.GenerateResetToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to generate reset token"})
			return
		}

		expiry := time.Now().Add(1 * time.Hour)
		_, err = db.Exec("INSERT INTO password_reset_tokens (user_id, token, expiry) VALUES (?, ?, ?)", user.ID, token, expiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to store reset token"})
			return
		}

		err = utils.SendResetEmail(req.Email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to send reset email"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Password reset link sent to your email"})
	}
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset the user's password using a valid token
// @Tags users
// @Accept json
// @Produce json
// @Param token query string true "Reset token"
// @Param password body models.ResetPasswordRequest true "New password"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /users/reset-password [post]
func ResetPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Token is required"})
			return
		}

		var req models.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Password == "" {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Password is required"})
			return
		}

		var userID int
		var tokenExpiryRaw []uint8
		query := "SELECT user_id, expiry FROM password_reset_tokens WHERE token = ?"

		err := db.QueryRow(query, token).Scan(&userID, &tokenExpiryRaw)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.Error{Error: "Invalid or expired token"})
				return
			}

			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to verify token"})
			return
		}

		tokenExpiryString := string(tokenExpiryRaw)
		tokenExpiry, err := time.Parse("2006-01-02 15:04:05", tokenExpiryString)
		if err != nil {
			fmt.Printf("Error parsing expiry time: %v\n", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to parse token expiry time"})
			return
		}

		if time.Now().After(tokenExpiry) {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "Token has expired"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to hash password"})
			return
		}

		updateQuery := "UPDATE users SET password = ? WHERE id = ?"
		_, err = db.Exec(updateQuery, hashedPassword, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to reset password"})
			return
		}

		_, err = db.Exec("DELETE FROM password_reset_tokens WHERE token = ?", token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to invalidate token"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Password reset successful"})
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.UserDetail
// @Failure 500 {object} models.Error
// @Router /users/ [get]
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, email, username, fullName FROM users FROM users WHERE deleted_at IS NULL")
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.UserDetail

		for rows.Next() {
			var user models.UserDetail
			if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName); err != nil {
				c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to scan user"})
				return
			}

			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.UserDetail
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{user_id} [get]
func GetUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		row := db.QueryRow("SELECT id, email, username, fullName FROM users WHERE id = ?", userId)

		var user models.UserDetail
		if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FullName); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{Error: "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to scan user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information by user ID
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param user body models.UserDetail true "User data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{user_id} [put]
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var updatedUser models.UserDetail

		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Error: "Invalid request body"})
			return
		}

		query := "UPDATE users SET email = ?, username = ?, fullName = ? WHERE id = ?"
		_, err := db.Exec(query, updatedUser.Email, updatedUser.Username, updatedUser.FullName, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "User updated successfully"})
	}
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the user's password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param password body models.ChangePasswordRequest true "Password data"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /users/{user_id}/change-password [put]
func ChangePassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

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

// DeleteUser godoc
// @Summary Soft delete user
// @Description Soft delete a user by user ID
// @Tags users
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} models.Message
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{user_id} [delete]
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

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

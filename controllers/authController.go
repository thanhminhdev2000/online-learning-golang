package controllers

import (
	"database/sql"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Log in
// @Description Log in using email or username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /auth/login [post]
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
			query = "SELECT id, email, username, fullName, password, gender, avatar, dateOfBirth, role FROM users WHERE email = ?"
		} else {
			query = "SELECT id, email, username, fullName, password, gender, avatar, dateOfBirth, role FROM users WHERE username = ?"
		}

		err := db.QueryRow(query, loginData.Identifier).
			Scan(&user.ID, &user.Email, &user.Username, &user.FullName, &storedPassword, &user.Gender, &user.Avatar, &user.DateOfBirth, &user.Role)
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

		accessToken, expiresIn, err := utils.CreateAccessToken(user.ID, string(user.Role))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to generate access token"})
			return
		}

		refreshToken, _, err := utils.CreateRefreshToken(user.ID, string(user.Role))
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
			ExpiresIn:   expiresIn,
		}

		c.JSON(http.StatusOK, response)
	}
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags Authentication
// @Produce json
// @Success 200 {object} models.AccessTokenResponse
// @Failure 401 {object} models.Error
// @Router /auth/refresh-token [post]
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{Error: "No refresh token found"})
			return
		}

		userId, role, err := utils.ValidToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		accessToken, expiresIn, err := utils.CreateAccessToken(userId, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to generate access token"})
			return
		}

		c.JSON(http.StatusOK, models.AccessTokenResponse{
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
		})
	}
}

// Logout godoc
// @Summary Log out
// @Description Log out by clearing the refresh token
// @Tags Authentication
// @Produce json
// @Success 200 {object} models.Message
// @Router /auth/logout [post]
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, models.Message{Message: "Logout successful"})
	}
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset link to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email body models.ForgotPasswordRequest true "User email"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /auth/forgot-password [post]
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
		_, err = db.Exec("INSERT INTO reset_pw_tokens (userId, token, expiry) VALUES (?, ?, ?)", user.ID, token, expiry)
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
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Reset token"
// @Param password body models.ResetPasswordRequest true "New password"
// @Success 200 {object} models.Message
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /auth/reset-password [post]
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

		var userId int
		var tokenExpiryRaw []uint8
		query := "SELECT userId, expiry FROM reset_pw_tokens WHERE token = ?"

		err := db.QueryRow(query, token).Scan(&userId, &tokenExpiryRaw)
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
		_, err = db.Exec(updateQuery, hashedPassword, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to reset password"})
			return
		}

		_, err = db.Exec("DELETE FROM reset_pw_tokens WHERE token = ?", token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Error: "Failed to invalidate token"})
			return
		}

		c.JSON(http.StatusOK, models.Message{Message: "Password reset successful"})
	}
}

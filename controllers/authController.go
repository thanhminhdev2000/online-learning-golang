package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"os"
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
// @Failure 400 {object} models.Error "Invalid request"
// @Failure 401 {object} models.Error "Authentication failed"
// @Failure 500 {object} models.Error "Server error"
// @Router /auth/login [post]
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData models.LoginRequest
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid request: " + err.Error(),
			})
			return
		}

		if err := validateLoginInput(loginData); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: err.Error(),
			})
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

		var user models.UserDetail
		var password string
		query := `
			SELECT id, email, username, fullName, password, 
				   gender, avatar, dateOfBirth, role 
			FROM users 
			WHERE (LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?)) 
			AND deletedAt IS NULL`

		err = tx.QueryRow(query, loginData.Identifier, loginData.Identifier).
			Scan(
				&user.ID,
				&user.Email,
				&user.Username,
				&user.FullName,
				&password,
				&user.Gender,
				&user.Avatar,
				&user.DateOfBirth,
				&user.Role,
			)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.Error{
					Error: "Invalid credentials",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to fetch user details",
			})
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginData.Password)); err != nil {

			c.JSON(http.StatusUnauthorized, models.Error{
				Error: "Invalid credentials",
			})
			return
		}

		accessToken, expiresIn, err := utils.CreateAccessToken(user.ID, string(user.Role))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to generate access token",
			})
			return
		}

		refreshToken, refreshExpiresIn, err := utils.CreateRefreshToken(user.ID, string(user.Role))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to generate refresh token",
			})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to complete login process",
			})
			return
		}

		setRefreshTokenCookie(c, refreshToken, int(refreshExpiresIn))

		c.JSON(http.StatusOK, models.LoginResponse{
			Message:     "Login successful",
			User:        user,
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
		})
	}
}

// Helper functions

func validateLoginInput(data models.LoginRequest) error {
	if data.Identifier == "" {
		return fmt.Errorf("email or username is required")
	}
	if data.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func setRefreshTokenCookie(c *gin.Context, token string, expiresIn int) {
	c.SetCookie(
		"refreshToken",
		token,
		expiresIn,
		"/",
		os.Getenv("COOKIE_DOMAIN"),
		os.Getenv("ENV") == "production", // Secure
		true,                             // HttpOnly
	)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh both access token and refresh token
// @Tags Authentication
// @Produce json
// @Success 200 {object} models.AccessTokenResponse "Returns new access token and sets new refresh token cookie"
// @Failure 400 {object} models.Error "Invalid user ID"
// @Failure 401 {object} models.Error "Invalid or missing refresh token"
// @Failure 500 {object} models.Error "Server error"
// @Router /auth/refresh-token [post]
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{
				Error: "No refresh token found",
			})
			return
		}

		userId, role, err := utils.ValidToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.Error{
				Error: "Invalid refresh token",
			})
			return
		}

		if userId <= 0 {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid user ID",
			})
			return
		}

		accessToken, expiresIn, err := utils.CreateAccessToken(userId, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to generate access token",
			})
			return
		}

		newRefreshToken, refreshExpiresIn, err := utils.CreateRefreshToken(userId, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to generate refresh token",
			})
			return
		}

		setRefreshTokenCookie(c, newRefreshToken, int(refreshExpiresIn))

		c.JSON(http.StatusOK, models.AccessTokenResponse{
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
		})
	}
}

// Logout godoc
// @Summary Log out
// @Description Log out by clearing the refresh token and invalidating the session
// @Tags Authentication
// @Produce json
// @Success 200 {object} models.Message "Logout successful"
// @Failure 401 {object} models.Error "No refresh token found"
// @Router /auth/logout [post]
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(
			"refreshToken",
			"",
			-1,
			"/",
			os.Getenv("COOKIE_DOMAIN"),
			os.Getenv("ENV") == "production",
			true,
		)

		c.JSON(http.StatusOK, models.Message{
			Message: "Logout successful",
		})
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
// @Failure 500 {object} models.Error "Server error"
// @Router /auth/forgot-password [post]
func ForgotPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid request format",
			})
			return
		}

		if req.Email == "" || !utils.IsValidEmail(req.Email) {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Valid email is required",
			})
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

		var user models.UserDetail
		err = tx.QueryRow(`
			SELECT id, email 
			FROM users 
			WHERE email = ? AND deletedAt IS NULL
		`, req.Email).Scan(&user.ID, &user.Email)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.Error{
					Error: "Email not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to query user",
			})
			return
		}

		_, err = tx.Exec("DELETE FROM reset_pw_tokens WHERE userId = ?", user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to clear existing tokens",
			})
			return
		}

		token, err := utils.GenerateResetToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to generate reset token",
			})
			return
		}

		expiry := time.Now().Add(1 * time.Hour)
		_, err = tx.Exec(`
			INSERT INTO reset_pw_tokens (userId, token, expiry) 
			VALUES (?, ?, ?)
		`, user.ID, token, expiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to store reset token",
			})
			return
		}

		err = utils.SendResetEmail(req.Email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to send reset email",
			})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to complete password reset request",
			})
			return
		}

		c.JSON(http.StatusOK, models.Message{
			Message: "Password reset instructions sent to your email",
		})
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
// @Failure 400 {object} models.Error "Invalid request or password"
// @Failure 401 {object} models.Error "Invalid or expired token"
// @Failure 500 {object} models.Error "Server error"
// @Router /auth/reset-password [post]
func ResetPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Token is required",
			})
			return
		}

		var req models.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: "Invalid request format",
			})
			return
		}

		if err := utils.ValidatePassword(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{
				Error: err.Error(),
			})
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

		var userId int
		var tokenExpiry time.Time
		err = tx.QueryRow(`
			SELECT userId, expiry 
			FROM reset_pw_tokens 
			WHERE token = ?
		`, token).Scan(&userId, &tokenExpiry)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.Error{
					Error: "Invalid or expired token",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to verify token",
			})
			return
		}

		if time.Now().After(tokenExpiry) {
			_, _ = tx.Exec("DELETE FROM reset_pw_tokens WHERE token = ?", token)
			c.JSON(http.StatusUnauthorized, models.Error{
				Error: "Token has expired",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to hash password",
			})
			return
		}

		_, err = tx.Exec(`
			UPDATE users 
			SET password = ?,
				updatedAt = NOW()
			WHERE id = ?
		`, hashedPassword, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to reset password",
			})
			return
		}

		_, err = tx.Exec("DELETE FROM reset_pw_tokens WHERE token = ?", token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to invalidate token",
			})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Error: "Failed to complete password reset",
			})
			return
		}

		c.JSON(http.StatusOK, models.Message{
			Message: "Password reset successful",
		})
	}
}

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

func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, email, username, fullName FROM users ORDER BY id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.User

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
				return
			}

			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		row := db.QueryRow("SELECT id, email, username, fullName FROM users WHERE id = ?", userId)

		var user models.User
		if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FullName); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.SignUp
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "err": err})
			return
		}

		var exists bool
		query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
		err := db.QueryRow(query, newUser.Email, newUser.Username).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing user"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email or username already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		query = "INSERT INTO users (email, username, fullName, password) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(query, newUser.Email, newUser.Username, newUser.FullName, hashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user", "err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData models.LoginRequest
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		var storedPassword string
		var user models.User
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
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginData.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		accessToken, err := utils.CreateAccessToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err := utils.CreateRefreshToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
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

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtKey = []byte(os.Getenv("JWT_KEY"))
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token found in cookies"})
			return
		}

		token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired or invalid refresh token"})
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID"})
			return
		}

		accessToken, err := utils.CreateAccessToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"accessToken": accessToken,
		})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}

func ForgotPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		// Check if the email exists in the database
		var user models.User
		err := db.QueryRow("SELECT id, email FROM users WHERE email = ?", req.Email).Scan(&user.ID, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
			return
		}

		// Generate a password reset token
		token, err := utils.GenerateResetToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}

		// Store the token in the password_reset_tokens table
		expiry := time.Now().Add(1 * time.Hour) // Token expires after 1 hour
		_, err = db.Exec("INSERT INTO password_reset_tokens (user_id, token, expiry) VALUES (?, ?, ?)", user.ID, token, expiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store reset token"})
			return
		}

		// Send an email with the password reset link
		err = utils.SendResetEmail(req.Email, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to your email"})
	}
}

func ResetPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Get the token from query parameters
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		// Step 2: Parse the new password from request body
		var req models.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
			return
		}

		// Step 3: Verify if the token is valid (e.g., check expiration and existence in DB)
		var userID int
		var tokenExpiryRaw []uint8
		query := "SELECT user_id, expiry FROM password_reset_tokens WHERE token = ?"

		err := db.QueryRow(query, token).Scan(&userID, &tokenExpiryRaw)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
			return
		}

		tokenExpiryString := string(tokenExpiryRaw)
		tokenExpiry, err := time.Parse("2006-01-02 15:04:05", tokenExpiryString)
		if err != nil {
			fmt.Printf("Error parsing expiry time: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token expiry time"})
			return
		}

		if time.Now().After(tokenExpiry) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		// Step 4: Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Step 5: Update the password in the users table
		updateQuery := "UPDATE users SET password = ? WHERE id = ?"
		_, err = db.Exec(updateQuery, hashedPassword, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
			return
		}

		// Step 6: Invalidate the token (delete it from the DB)
		_, err = db.Exec("DELETE FROM password_reset_tokens WHERE token = ?", token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate token"})
			return
		}

		// Step 7: Respond with success
		c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
	}
}

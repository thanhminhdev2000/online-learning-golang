package controllers

import (
	"database/sql"
	"net/http"
	"online-learning-golang/models"
	"online-learning-golang/utils"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, email, username, firstName, lastName FROM users ORDER BY id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.User

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName); err != nil {
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
		row := db.QueryRow("SELECT id, email, username, firstName, lastName FROM users WHERE id = ?", userId)

		var user models.User
		if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName); err != nil {
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

		query = "INSERT INTO users (email, username, firstName, lastName, password) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, newUser.Email, newUser.Username, newUser.FirstName, newUser.LastName, hashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
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
			query = "SELECT id, email, username, firstName, lastName, password FROM users WHERE email = ?"
		} else {
			query = "SELECT id, email, username, firstName, lastName, password FROM users WHERE username = ?"
		}

		err := db.QueryRow(query, loginData.Identifier).
			Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &storedPassword)
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

	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

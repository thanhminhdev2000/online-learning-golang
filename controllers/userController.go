package controllers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignUp(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

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

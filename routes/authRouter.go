package routes

import (
	"database/sql"
	"online-learning-golang/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.POST("/login", controllers.Login(db))
	router.POST("/logout", controllers.Logout())
	router.POST("/refresh-token", controllers.RefreshToken())
	router.POST("/forgot-password", controllers.ForgotPassword(db))
	router.POST("/reset-password", controllers.ResetPassword(db))
}

package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.POST("/", controllers.CreateUser(db))
	router.GET("/", middleware.AuthMiddleware(), controllers.GetUsers(db))
	router.GET("/:user_id", middleware.AuthMiddleware(), controllers.GetUserByID(db))
	router.PUT("/:user_id", middleware.AuthMiddleware(), controllers.UpdateUser(db))
	router.PUT("/:user_id/password", middleware.AuthMiddleware(), controllers.PasswordUpdate(db))
	router.DELETE("/:user_id", middleware.AuthMiddleware(), controllers.DeleteUser(db))
}

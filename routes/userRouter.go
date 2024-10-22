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
	router.GET("/:userId", middleware.AuthMiddleware(), controllers.GetUserByID(db))
	router.PUT("/:userId", middleware.AuthMiddleware(), controllers.UpdateUser(db))
	router.PUT("/:userId/password", middleware.AuthMiddleware(), controllers.UpdateUserPassword(db))
	router.PUT("/:userId/avatar", middleware.AuthMiddleware(), controllers.UpdateUserAvatar(db))
	router.DELETE("/:userId", middleware.AuthMiddleware(), controllers.DeleteUser(db))
}

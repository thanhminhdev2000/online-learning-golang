package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.POST("/", controllers.CreateUser(db))
	router.POST("/admin", middleware.AuthMiddleware(), controllers.CreateUserAdmin(db))
	router.GET("/", middleware.AuthMiddleware(), controllers.GetUsers(db))
	router.GET("/:id", middleware.AuthMiddleware(), controllers.GetUserByID(db))
	router.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateUser(db))
	router.PUT("/:id/password", middleware.AuthMiddleware(), controllers.UpdateUserPassword(db))
	router.PUT("/:id/avatar", middleware.AuthMiddleware(), controllers.UpdateUserAvatar(db))
	router.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteUser(db))
}

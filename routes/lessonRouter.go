package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func LessonRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetLessons(db))
	router.GET("/:id", controllers.GetLesson(db))
	router.POST("/", middleware.AuthMiddleware(), controllers.CreateLesson(db))
	router.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateLesson(db))
	router.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteLesson(db))
}

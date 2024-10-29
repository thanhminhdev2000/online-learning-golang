package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetCourses(db))
	router.GET("/:id", controllers.GetCourse(db))
	router.POST("/", middleware.AuthMiddleware(), controllers.CreateCourse(db))
	router.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateCourse(db))
	router.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteCourse(db))
}

package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func CourseRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetCourses(db))
	router.GET("/:id", middleware.AuthMiddleware(), controllers.GetCourse(db))
	router.POST("/", middleware.OnlyAdminMiddleware(), controllers.CreateCourse(db))
	router.PUT("/:id", middleware.OnlyAdminMiddleware(), controllers.UpdateCourse(db))
	router.DELETE("/:id", middleware.OnlyAdminMiddleware(), controllers.DeleteCourse(db))
}

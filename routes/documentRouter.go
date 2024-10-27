package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func DocumentRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetDocuments(db))
	router.GET("/subjects", controllers.GetListClassesWithSubjects(db))
	router.POST("/", middleware.AuthMiddleware(), controllers.CreateDocument(db))
	router.PUT("/:documentId", middleware.AuthMiddleware(), controllers.UpdateDocument(db))
	router.DELETE("/:documentId", middleware.AuthMiddleware(), controllers.DeleteDocument(db))
}

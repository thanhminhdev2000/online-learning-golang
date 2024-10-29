package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func DocumentRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetDocuments(db))
	router.GET("/classes", controllers.GetListClass(db))
	router.POST("/", middleware.AuthMiddleware(), controllers.CreateDocument(db))
	router.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateDocument(db))
	router.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteDocument(db))
}

package routes

import (
	"database/sql"
	"online-learning-golang/controllers"

	"github.com/gin-gonic/gin"
)

func DocumentRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetDocuments(db))
	router.GET("/subjects", controllers.GetListClassesWithSubjects(db))
	router.POST("/upload", controllers.UploadDocument(db))
}

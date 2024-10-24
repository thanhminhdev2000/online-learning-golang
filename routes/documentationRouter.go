package routes

import (
	"database/sql"
	"online-learning-golang/controllers"

	"github.com/gin-gonic/gin"
)

func DocumentationRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/", controllers.GetListClassesWithSubjects(db))
}

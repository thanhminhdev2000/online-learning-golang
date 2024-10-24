package routes

import (
	"database/sql"
	"online-learning-golang/controllers"

	"github.com/gin-gonic/gin"
)

func ContactRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.POST("/", controllers.Contact(db))
}

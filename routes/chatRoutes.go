package routes

import (
	"database/sql"

	"online-learning-golang/chat"
	"online-learning-golang/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.RouterGroup, db *sql.DB) {
	manager := chat.NewManager()
	go manager.Run()

	router.GET("/ws", controllers.HandleWebSocket(manager))
	router.GET("/history", controllers.GetChatHistory(db))
}

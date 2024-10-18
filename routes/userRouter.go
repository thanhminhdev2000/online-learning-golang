package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine, db *sql.DB) {
	userGroup := routes.Group("/users")
	{
		userGroup.GET("/", middleware.AuthMiddleware(), controllers.GetUsers(db))
		userGroup.GET("/:user_id", middleware.AuthMiddleware(), controllers.GetUser(db))
		userGroup.POST("/signup", controllers.SignUp(db))
		userGroup.POST("/login", controllers.Login(db))
		userGroup.POST("/refresh", controllers.RefreshToken())
		userGroup.POST("/logout", controllers.Logout())
	}

}

package routes

import (
	"database/sql"
	"online-learning-golang/controllers"
	"online-learning-golang/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.POST("/signup", controllers.SignUp(db))
	router.POST("/login", controllers.Login(db))
	router.POST("/refresh", controllers.RefreshToken())
	router.POST("/logout", controllers.Logout())
	router.POST("/forgot-password", controllers.ForgotPassword(db))
	router.POST("/reset-password", controllers.ResetPassword(db))

	router.GET("/", middleware.AuthMiddleware(), controllers.GetUsers(db))
	router.GET("/:user_id", middleware.AuthMiddleware(), controllers.GetUserByID(db))
	router.PUT("/:user_id", middleware.AuthMiddleware(), controllers.UpdateUser(db))
	router.PUT("/:user_id/change-password", middleware.AuthMiddleware(), controllers.ChangePassword(db))
	router.DELETE("/:user_id", middleware.AuthMiddleware(), controllers.DeleteUser(db))
}

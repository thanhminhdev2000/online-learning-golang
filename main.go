package main

import (
	"flag"
	"fmt"
	"log"
	"online-learning-golang/database"
	"online-learning-golang/routes"
	"os"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "online-learning-golang/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Online Learning API
// @version 1.0
// @description This is an online learning API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	reset := flag.Bool("reset", false, "Reset the database")
	flag.Parse()
	db, err := database.ConnectMySQL()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	if *reset {
		database.ResetDataBase(db)
		fmt.Println("Database reset successfully!")
		return
	}

	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	apiPrefix := os.Getenv("API_PREFIX")
	routes.UserRoutes(router.Group(apiPrefix+"/users"), db)
	routes.AuthRoutes(router.Group(apiPrefix+"/auth"), db)
	routes.DocumentationRoutes(router.Group(apiPrefix+"/documentations"), db)
	routes.ContactRoutes(router.Group(apiPrefix+"/contact"), db)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("localhost:8080")
}

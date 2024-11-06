package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"online-learning-golang/database"
	_ "online-learning-golang/docs"
	"online-learning-golang/routes"

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

	if err := database.CreateAllTablesIfNotExist(db); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	if *reset {
		if err := database.ResetDataBase(db); err != nil {
			log.Fatalf("Error resetting database: %v", err)
		}
		fmt.Println("Database reset successfully!")
		return
	}

	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://178.128.55.56"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	apiPrefix := os.Getenv("API_PREFIX")
	routes.UserRoutes(router.Group(apiPrefix+"/users"), db)
	routes.AuthRoutes(router.Group(apiPrefix+"/auth"), db)
	routes.ContactRoutes(router.Group(apiPrefix+"/contacts"), db)
	routes.DocumentRoutes(router.Group(apiPrefix+"/documents"), db)
	routes.CourseRoutes(router.Group(apiPrefix+"/courses"), db)
	routes.LessonRoutes(router.Group(apiPrefix+"/lessons"), db)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}

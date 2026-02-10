package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/config"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/database"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/handlers"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/middleware"
)

func main() {
	//Load environment variables from .env file if it exists
	env_configurations, err := config.LoadEnvironmentalVariables()
	if err != nil {
		log.Fatal("Failed to load environment variables: " + err.Error())
	}
	log.Println("Environment variables loaded successfully:")

	//connect to the database
	connectionPool, err := database.ConnectToDatabase(env_configurations.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database: " + err.Error())
	}
	log.Println("Database connection established successfully")

	defer connectionPool.Close()

	// Initialize Gin router
	var router *gin.Engine = gin.Default()

	router.SetTrustedProxies(nil)
	

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":"Todo api is running",
			"status":"success",
			"database": "connected",	
		})
	})

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(*env_configurations))

	// Setup protected routes
    {
	  protected.POST("", handlers.CreateTodoHandler(connectionPool))
	  protected.GET("/", handlers.GetAllTodosHandler(connectionPool))
	  protected.GET("/:id", handlers.GetTodoByIDHandler(connectionPool))
	  protected.PUT("/:id", handlers.UpdateTodoHandler(connectionPool))
	  protected.DELETE("/:id", handlers.DeleteTodoHandler(connectionPool))

	}
	router.POST("/auth/register", handlers.CreateUserHandler(connectionPool))
	router.POST("/auth/login", handlers.LoginUserHandler(connectionPool, env_configurations))

	

	// Protected route example
	router.GET("/protected-test", middleware.AuthMiddleware(*env_configurations), handlers.TestProtectedHandler())

	router.Run(":" + env_configurations.Port)
}